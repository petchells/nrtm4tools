package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

const tagName = "em"
const entityName = "EntityManaged"

/*
EntityManaged Support for query execution and struct mapping. If your entity has an ID field, you can do CRUD
operations: Create, GetByID (Read), Update (not implemented: Delete).
*/
type EntityManaged interface {
}

var descriptors = map[string]Descriptor{}

var tableNameSet = map[string]struct{}{}
var tableAliasSet = map[string]struct{}{}

type taggedField struct {
	f   reflect.StructField
	tag string
}

// Descriptor extends EntityManaged with functions needed by query builders
type Descriptor struct {
	tableName   string
	tableAlias  string
	columnNames []string
	fields      []taggedField
}

// TableAlias return the table alias to be used in joined queries
func (d *Descriptor) TableAlias() string {
	return d.tableAlias
}

// TableName returns the table name to be used in queries
func (d *Descriptor) TableName() string {
	return d.tableName
}

// TableNameWithAlias is a shortcut for name  + " " + alias
func (d *Descriptor) TableNameWithAlias() string {
	return d.tableName + " " + d.tableAlias
}

// ColumnNames returns an array of column names
func (d *Descriptor) ColumnNames() []string {
	return d.columnNames
}

// ColumnNamesWithAlias returns an array of column names like "o.id o_id", "o.name o_name"
func (d *Descriptor) ColumnNamesWithAlias() []string {
	names := []string{}
	for _, name := range d.columnNames {
		names = append(names, d.tableAlias+"."+name+" "+d.tableAlias+"_"+name)
	}
	return names
}

// FieldValues returns field pointers so an EntityManaged row can be 'Scan(...)'ned
func FieldValues(e EntityManaged) []interface{} {
	sflds := []interface{}{}
	val := reflect.ValueOf(e).Elem()
	kind := val.Kind().String()
	if kind == "invalid" {
		return sflds
	}
	for _, tf := range GetDescriptor(e).fields {
		valueField := val.FieldByName(tf.f.Name)
		sflds = append(sflds, valueField.Addr().Interface())
	}
	return sflds
}

// GetDescriptor gives you a helper for building sql
func GetDescriptor(e EntityManaged) Descriptor {
	ty := reflect.TypeOf(e).Elem()
	if d, ok := descriptors[ty.Name()]; ok {
		return d
	}
	desc := Descriptor{}
	allFields := getTaggedFields(e)
	for _, field := range allFields {
		if field.f.Name == entityName {
			parts := strings.Split(field.tag, " ")
			if len(parts) != 2 {
				log.Fatalln("EntityManaged table name, alias is not defined.", ty.Name())
			}
			desc.tableName = strings.TrimSpace(parts[0])
			desc.tableAlias = strings.TrimSpace(parts[1])
			if len(desc.tableName) == 0 || len(desc.tableAlias) == 0 {
				log.Fatalln("EntityManaged table name or alias is empty.", ty.Name())
			}
			if _, ok := tableNameSet[desc.tableName]; ok {
				log.Fatalln("EntityManaged table name is not unique.", desc.tableName)
			}
			if _, ok := tableAliasSet[desc.tableAlias]; ok {
				log.Fatalln("EntityManaged table alias is not unique.", desc.tableName, desc.tableAlias)
			}
			tableNameSet[desc.tableName] = struct{}{}
			tableAliasSet[desc.tableAlias] = struct{}{}
		} else {
			desc.columnNames = append(desc.columnNames, fieldNameToColumnName(field.f.Name))
			desc.fields = append(desc.fields, field)
		}
	}
	descriptors[ty.Name()] = desc
	return desc
}

// GetAll performs a select and calls back a function with transformed entities. bit clumsy but it works
func GetAll[T EntityManaged](tx pgx.Tx, entity T, fn func(entity T)) ([]T, error) {
	dtor := GetDescriptor(entity)
	cols := strings.Join(dtor.columnNames, ", ")
	sql := fmt.Sprintf("SELECT %v FROM %v", cols, dtor.tableName)
	res := []T{}
	rows, err := tx.Query(context.Background(), sql)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(FieldValues(entity)...)
		if err != nil {
			return res, err
		}
		fn(entity)
		res = append(res, entity)
	}
	return res, nil
}

// GetAllByColumn performs a select and calls back a function with transformed entities. bit clumsy but it works
func GetAllByColumn[T EntityManaged](tx pgx.Tx, colname string, value any, entity T, fn func(entity T)) ([]T, error) {
	dtor := GetDescriptor(entity)
	cols := strings.Join(dtor.columnNames, ", ")
	sql := fmt.Sprintf("SELECT %v FROM %v WHERE %v=$1", cols, dtor.tableName, colname)
	res := []T{}
	rows, err := tx.Query(context.Background(), sql, value)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(FieldValues(entity)...)
		if err != nil {
			return res, err
		}
		fn(entity)
		res = append(res, entity)
	}
	return res, nil
}

// GetByID Fills the entityPtr with the row selected by ID
func GetByID(tx pgx.Tx, ID int64, entityPtr EntityManaged) error {
	dtor := GetDescriptor(entityPtr)
	cols := strings.Join(dtor.columnNames, ", ")
	tableName := dtor.tableName
	sql := fmt.Sprintf("SELECT %v FROM %v WHERE id=$1", cols, tableName)
	return tx.QueryRow(context.Background(), sql, ID).Scan(FieldValues(entityPtr)...)
}

// GetByColumn Fills the entityPtr with a single row matched by the value
func GetByColumn(tx pgx.Tx, colname string, value any, entityPtr EntityManaged) error {
	dtor := GetDescriptor(entityPtr)
	cols := strings.Join(dtor.columnNames, ", ")
	tableName := dtor.tableName
	sql := fmt.Sprintf("SELECT %v FROM %v WHERE %v=$1", cols, tableName, colname)
	return tx.QueryRow(context.Background(), sql, value).Scan(FieldValues(entityPtr)...)
}

// Create an entity -- entity must be a pointer
func Create(tx pgx.Tx, entityPtr EntityManaged) error {
	dtor := GetDescriptor(entityPtr)
	placeholders := []string{}
	cols := dtor.columnNames
	if len(cols) == 0 {
		return errors.New("Entity has no columns: " + dtor.tableName)
	}
	for i := range cols {
		placeholders = append(placeholders, "$"+strconv.Itoa(i+1))
	}
	values := FieldValues(entityPtr)
	sql := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)",
		dtor.tableName,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)
	tag, err := tx.Exec(context.Background(), sql, values...)
	if err != nil {
		return err
	}
	log.Println("DEBUG Insert", tag.RowsAffected(), "rows affected")
	return nil
}

// Update an entity
func Update(tx pgx.Tx, e EntityManaged) error {
	dtor := GetDescriptor(e)
	placeholders := []string{}
	for i, cn := range dtor.columnNames {
		placeholders = append(placeholders, cn+"=$"+strconv.Itoa(i+1))
	}
	values := FieldValues(e)
	_, err := tx.Exec(context.Background(),
		fmt.Sprintf("UPDATE %v SET %v WHERE id=$1",
			dtor.tableName,
			strings.Join(placeholders, ", "),
		), values...)
	return err
}

func columnNamesWithAlias(e EntityManaged) []string {
	dtor := GetDescriptor(e)
	names := []string{}
	for _, name := range dtor.columnNames {
		names = append(names, dtor.tableAlias+"."+name+" "+dtor.tableAlias+"_"+name)
	}
	return names
}

func getTaggedFields(t EntityManaged) []taggedField {
	var fields []taggedField
	ty := reflect.TypeOf(t).Elem()
	if len(ty.Name()) == 0 {
		return fields
	}
	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)
		tag := field.Tag.Get(tagName)
		if len(tag) > 0 {
			fields = append(fields, taggedField{f: field, tag: tag})
		}
	}
	return fields
}

var cre []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile("[0-9]+"),
	regexp.MustCompile("[A-Z][a-z]+"),
	regexp.MustCompile("[A-Z]+"),
}

func fieldNameToColumnName(fieldname string) string {
	s := fieldname
	for i := 0; i < len(cre); i++ {
		s = cre[i].ReplaceAllStringFunc(s, func(str string) string {
			return "_" + strings.ToLower(str)
		})
	}
	return s[1:]
}
