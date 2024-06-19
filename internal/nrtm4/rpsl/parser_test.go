package rpsl

import "testing"

/*

What about multi-line strings like address -- how is the second line distinguished from an attribute name?
- maybe indented? if so, then is it mandatory for the first char in a line to *always* be an attribute name, or comment?
- or is a colon disallowed?

*/

func TestErrorThrownForBadRPSL(t *testing.T) {

	var err error

	_, err = parseString("")

	if err != ErrCannotParseRPSL {
		t.Error("Parser should fail on empty object")
	}

	_, err = parseString("person")

	if err != ErrCannotParseRPSL {
		t.Error("Parser should fail when object type isn't defined")
	}

	_, err = parseString("person: fish")

	if err != ErrCannotParseRPSL {
		t.Error("Parser should fail on incomplete object definition")
	}

	str := `route:          37.37.37.0/24
	descr:          XYZ Network
	origin:         AS9876
	mnt-by:         XYZ-MNT
	created:        2021-02-22T04:30:00Z
	last-modified:  2021-02-22T14:30:00Z
	`
	_, err = parseString(str)

	if err != ErrCannotParseRPSL {
		t.Error("Parser should fail if source is missing")
	}

}

func TestRpslParsePerson(t *testing.T) {
	str := `
	person:      Daniel Karrenberg # is a person
    address:     RIPE Network Coordination Centre (NCC)
    address:     Singel 258
    address:     NL-1016 AB  Amsterdam
    address:     Netherlands
    phone:       +31 20 535 4444
    fax-no:      +31 20 535 4445
    e-mail:      Daniel.Karrenberg@ripe.net
    nic-hdl:     DK58 # xxxxxx
    changed:     Daniel.Karrenberg@ripe.net 19970616
    source:      RIPE`
	objectType := "PERSON"
	source := "RIPE"
	primaryKey := "DK58"

	obj, err := ParseString(str)

	if err != nil {
		t.Error("Parser doesn't work", err)
	}
	if obj.ObjectType != objectType {
		t.Error("Parser did not recognize object. Expected", objectType, "was", obj.ObjectType)
	}
	if obj.Source != source {
		t.Error("Parser did not parse a source. expected", source, "was", obj.Source)
	}
	if obj.PrimaryKey != primaryKey {
		t.Error("Parser did not parse a primary key. expected", primaryKey, "was", obj.PrimaryKey)
	}
}

func TestRpslParseRoute(t *testing.T) {
	str := `route:          37.37.37.0/24 # road to nowhere
	descr:          XYZ Network
	origin:         AS9876
	mnt-by:         XYZ-MNT
	created:        2021-02-22T04:30:00Z
	last-modified:  2021-02-22T14:30:00Z
	source:         RIPE # Filtered
	`
	objectType := "ROUTE"
	source := "RIPE"
	primaryKey := "37.37.37.0/24AS9876"

	obj, err := parseString(str)

	if err != nil {
		t.Error("Parser doesn't work", err)
	}
	if obj.ObjectType != objectType {
		t.Error("Parser did not recognize object. Expected", objectType, "was", obj.ObjectType)
	}
	if obj.Source != source {
		t.Error("Parser did not parse a source. expected", source, "was", obj.Source)
	}
	if obj.PrimaryKey != primaryKey {
		t.Error("Parser did not parse a primary key. expected", primaryKey, "was", obj.PrimaryKey)
	}
}

func TestRpslParseASBlock(t *testing.T) {
	str := `AS-BLOCK:       as3209 - as3353
	descr:          RIPE NCC ASN block
	remarks:        These AS Numbers are assigned to network operators in the RIPE NCC service region.
	mnt-by:         RIPE-NCC-HM-MNT
	created:        2018-11-22T15:27:19Z
	last-modified:  2018-11-22T15:27:19Z
	source:         RIPE
	`
	objectType := "AS-BLOCK"
	source := "RIPE"
	primaryKey := "AS3209 - AS3353"

	obj, err := parseString(str)

	if err != nil {
		t.Error("Parser doesn't work", err)
	}
	if obj.ObjectType != objectType {
		t.Error("Parser did not recognize object. Expected", objectType, "was", obj.ObjectType)
	}
	if obj.Source != source {
		t.Error("Parser did not parse a source. expected", source, "was", obj.Source)
	}
	if obj.PrimaryKey != primaryKey {
		t.Error("Parser did not parse a primary key. expected", primaryKey, "was", obj.PrimaryKey)
	}
}
