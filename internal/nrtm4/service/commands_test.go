package service

import (
	"errors"
	"testing"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

type ProcessorStub struct{}

func (ps ProcessorStub) Connect(url, label string) error {
	return errors.New("test error")
}

func (ps ProcessorStub) Update(srcName, label string) error {
	return nil
}

func (ps ProcessorStub) ListSources() ([]persist.NRTMSource, error) {
	return []persist.NRTMSource{}, nil
}

func TestCommandExecutorConnect(t *testing.T) {
	ce := CommandExecutor{ProcessorStub{}}
	ce.Connect("url", "label")
}

func TestCommandExecutorUpdate(t *testing.T) {
	ce := CommandExecutor{ProcessorStub{}}
	ce.Connect("srcName", "label")
}

type labelExpectation struct {
	label  string
	expect bool
}

func TestLabelRegex(t *testing.T) {
	lbls := [...]labelExpectation{{
		"This_one_is-100.OK", true},
		{"1_is_ok", true},
		{"YES$nowerky", false},
		{"F", true},
		{"1970-01-01", true},
		{"This one is not OK", false},
		{"-------", false},
		{"------1", true},
	}
	for _, lbl := range lbls {
		match := labelRegex.MatchString(lbl.label)
		if match != lbl.expect {
			if lbl.expect {
				t.Error("Label regex should succeed", lbl.label)
			} else {
				t.Error("Label regex should fail", lbl.label)
			}
		}
	}
}
