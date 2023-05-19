package db

import (
	"strings"
	"testing"
	"time"
)

type testOrg struct {
	EntityManaged `em:"cd_test_org, o"`
	ID            int64     `em:"-"`
	Updated       time.Time `em:"-"`
	Name          string    `em:"-"`
	LegalName     string    `em:"-"`
}

func TestFieldNameConversion(t *testing.T) {
	testStrings := [10]string{"ID", "HTML", "Data", "CashDash", "CashDASH", "BigC", "JSONString", "CDUserID", "RTimer", "XML10YAMLFormat5"}
	expected := [10]string{"id", "html", "data", "cash_dash", "cash_dash", "big_c", "json_string", "cd_user_id", "r_timer", "xml_10_yaml_format_5"}
	for i, str := range testStrings {
		result := fieldNameToColumnName(str)
		if result != expected[i] {
			t.Errorf("Expected '%v' but got '%v'", expected[i], result)
		}
	}
}

func TestColumnNameConversionFromFieldTags(t *testing.T) {
	expected := [4]string{"id", "updated", "name", "legal_name"}
	o := testOrg{}
	dtor := GetDescriptor(&o)
	names := dtor.columnNames
	if len(expected) != len(names) {
		t.Errorf("Expected '%d' fields but got '%d'", len(expected), len(names))
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != names[i] {
			t.Errorf("Expected field name '%v' but got '%v'", expected[i], names[i])
		}
	}
}

func TestColumnNameWithAliasConversionFromFieldTags(t *testing.T) {
	expected := [4]string{"o.id o_id", "o.updated o_updated", "o.name o_name", "o.legal_name o_legal_name"}
	o := testOrg{}
	names := columnNamesWithAlias(&o)
	if len(expected) != len(names) {
		t.Errorf("Expected '%d' fields but got '%d'", len(expected), len(names))
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != names[i] {
			t.Errorf("Expected field name '%v' but got '%v'", expected[i], names[i])
		}
	}
}

func TestScannableFields(t *testing.T) {
	u := new(testOrg)
	dtor := GetDescriptor(u)
	ns := dtor.columnNames
	expect := "id, updated, name, legal_name"
	if expect != strings.Join(ns, ", ") {
		t.Errorf("got '%v' expected '%v'", strings.Join(ns, ", "), expect)
	}
}

func TestScannableFieldsAndValues(t *testing.T) {
	o := filledNewOrg(t)
	f := FieldValues(&o)
	sc := []interface{}{
		&o.ID,
		&o.Updated,
		&o.Name,
		&o.LegalName,
	}
	if len(f) != len(sc) {
		t.Errorf("ScannableFields failed. Expected %d got %d", len(sc), len(f))
	}
	for i := 0; i < len(sc); i++ {
		if f[i] != sc[i] {
			t.Errorf("ID failed %v", f[i])
		}
	}
}

func filledNewOrg(t *testing.T) testOrg {
	id := int64(123)
	dateStr := "2020-01-18T08:15:00Z"
	dt, _ := time.Parse(time.RFC3339, dateStr)
	name := "VBC"
	lname := "Very Big Corp N.V."
	return newOrg(id, dt, name, lname)
}

func newOrg(id int64, dtu time.Time, name string, lname string) testOrg {
	return testOrg{
		ID:        id,
		Updated:   dtu,
		Name:      name,
		LegalName: lname,
	}
}
