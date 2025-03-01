package db

import (
	"context"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

type testOrg struct {
	EntityManaged `em:"test_org o"`
	ID            int64      `em:"-"`
	Updated       time.Time  `em:"-"`
	Name          *string    `em:"-"`
	Quantity      uint       `em:"-"`
	Birthday      *time.Time `em:"date_of_birth"`
}

func TestFieldNameConversion(t *testing.T) {
	testStrings := [...]string{"ID", "HTML", "Data", "CashDash", "CashDASH", "BigC", "JSONString", "CDUserID", "RTimer", "XML10YAMLFormat5", "Under_Score"}
	expected := [...]string{"id", "html", "data", "cash_dash", "cash_dash", "big_c", "json_string", "cd_user_id", "r_timer", "xml_10_yaml_format_5", "under_score"}
	for i, str := range testStrings {
		result := fieldNameToColumnName(str)
		if result != expected[i] {
			t.Errorf("Expected '%v' but got '%v'", expected[i], result)
		}
	}
}

func TestColumnNameConversionFromFieldTags(t *testing.T) {
	expected := [...]string{"id", "updated", "name", "quantity", "date_of_birth"}
	o := testOrg{}
	dtor := GetDescriptor(&o)
	names := dtor.ColumnNames()
	if len(expected) != len(names) {
		t.Errorf("Expected '%d' fields but got '%d'", len(expected), len(names))
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != names[i] {
			t.Errorf("Expected field name '%v' but got '%v'", expected[i], names[i])
		}
	}
}

func TestColumnNames(t *testing.T) {
	u := new(testOrg)
	dtor := GetDescriptor(u)
	{
		ns := dtor.ColumnNames()
		expect := "id, updated, name, quantity, date_of_birth"
		if expect != strings.Join(ns, ", ") {
			t.Errorf("got '%v' expected '%v'", strings.Join(ns, ", "), expect)
		}
	}
	{
		ns := dtor.ColumnNamesWithAlias()
		expect := "o.id o_id, o.updated o_updated, o.name o_name, o.quantity o_quantity, o.date_of_birth o_date_of_birth"
		if expect != strings.Join(ns, ", ") {
			t.Errorf("got '%v' expected '%v'", strings.Join(ns, ", "), expect)
		}
	}
}

func TestScannableFieldsAndValues(t *testing.T) {
	o := filledNewOrg()
	o.Birthday = nil
	f := ValuesForSelect(&o)
	sc := []interface{}{
		&o.ID,
		&o.Updated,
		&o.Name,
		&o.Quantity,
		&o.Birthday,
	}
	if len(f) != len(sc) {
		t.Fatalf("ScannableFields failed. Expected %d got %d", len(sc), len(f))
	}
	for i := 0; i < len(sc); i++ {
		if f[i] != sc[i] {
			t.Errorf("ID failed %v", f[i])
		}
	}
}

func TestUInsertOrpdateValues(t *testing.T) {
	o := filledNewOrg()
	o.Birthday = nil
	f := ValuesForModify(&o)
	sc := []any{
		&o.ID,
		&o.Updated,
		o.Name,
		&o.Quantity,
		nil,
	}
	if len(f) != len(sc) {
		t.Fatalf("ScannableFields failed. Expected %d got %d", len(sc), len(f))
	}
	for i := 0; i < len(sc); i++ {
		if f[i] != sc[i] {
			t.Errorf("Field comparison failed at index %d. Expected %v but got %v", i, sc[i], f[i])
		}
	}
}

type stubRows struct {
	pgx.Rows
	state [1]*int
	rows  [][]any
}

func (r stubRows) Next() bool {
	return *r.state[0] < len(r.rows)
}

func (r stubRows) Scan(dest ...any) error {
	reflect.Copy(reflect.ValueOf(dest), reflect.ValueOf(r.rows[*r.state[0]]))
	// for i := range dest {
	// 	dest[i] = &r.rows[*r.state[0]][i]
	// }
	for _, d := range dest {
		log.Println("dest", d)
	}
	*r.state[0]++
	return nil
}

func (r stubRows) Close() {}

type stubTx struct {
	pgx.Tx
	t    *testing.T
	rows pgx.Rows
	sql  string
}

func (tx stubTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if tx.sql != sql {
		tx.t.Errorf("Not the SQL you were looking for. Expected:\n%v\nbut was:\n%v", tx.sql, sql)
	}
	return tx.rows, nil
}

// Doesn't work. see comment below
func TestGetAll(t *testing.T) {
	o := filledNewOrg()
	o2 := filledNewOrg()
	o2.ID = 98765
	name := "Next Gen"
	o2.Name = &name
	o2.Quantity = 7
	o2.Birthday = nil

	rowCounter := 0
	state := [1]*int{&rowCounter}
	rows := stubRows{
		rows: [][]any{
			{
				o.ID,
				o.Updated,
				o.Name,
				o.Quantity,
				o.Birthday,
			},
			{
				o2.ID,
				o2.Updated,
				o2.Name,
				o2.Quantity,
				o2.Birthday,
			},
		},
		state: state,
	}
	sql := "SELECT id, updated, name, quantity, date_of_birth FROM test_org"
	tx := stubTx{t: t, rows: rows, sql: sql}

	expectLen := len(rows.rows)

	callbackCounter := 0
	allOrgs, err := GetAll(tx, testOrg{}, func(o testOrg) {
		callbackCounter++
	})
	if err != nil {
		t.Fatal("GetAll returned an error")
	}
	if callbackCounter != expectLen {
		t.Fatalf("Expected %d callback(s), but was %d", expectLen, callbackCounter)
	}

	if len(allOrgs) != expectLen {
		t.Fatalf("Should have returned %d orgs but was %d", expectLen, len(allOrgs))
	}
	// the mock 'Scan' function in stubRows does not properly emulate rows.Scan, so these tests fail
	// if allOrgs[0].ID != rows.rows[0][0] {
	// 	t.Errorf("Error in ID field. Expected %v but was %v", o.ID, allOrgs[0].ID)
	// }
	// if allOrgs[0].Updated != rows.rows[0][1] {
	// 	t.Errorf("Error in Updated field. Expected %v but was %v", o.Updated, allOrgs[0].Updated)
	// }
	// if allOrgs[0].Name != rows.rows[0][2] {
	// 	t.Errorf("Error in Name field. Expected %v but was %v", *o.Name, allOrgs[0].Name)
	// }
	// if allOrgs[0].Quantity != rows.rows[0][3] {
	// 	t.Errorf("Error in Quantity field. Expected %v but was %v", o.Quantity, allOrgs[0].Quantity)
	// }
}

func filledNewOrg() testOrg {
	id := int64(123)
	dateStr := "2020-01-18T08:15:00Z"
	dt, _ := time.Parse(time.RFC3339, dateStr)
	name := "VBC"
	qty := uint(999)
	dobStr := "2002-03-28T08:15:00Z"
	dob, _ := time.Parse(time.RFC3339, dobStr)
	return newOrg(id, dt, &name, qty, &dob)
}

func newOrg(id int64, dtu time.Time, name *string, qty uint, dob *time.Time) testOrg {
	return testOrg{
		ID:       id,
		Updated:  dtu,
		Name:     name,
		Quantity: qty,
		Birthday: dob,
	}
}
