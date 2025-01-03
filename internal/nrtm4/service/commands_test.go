package service

import (
	"errors"
	"testing"

	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
)

type ProcessorStub struct{}

func (ps ProcessorStub) Connect(url, label string) error {
	return errors.New("test error")
}

func (ps ProcessorStub) Update(srcName, label string) error {
	return nil
}

func (ps ProcessorStub) ListSources() ([]persist.NRTMSourceDetails, error) {
	return []persist.NRTMSourceDetails{}, nil
}

func (ps ProcessorStub) ReplaceLabel(src, fromLabel, toLabel string) (*persist.NRTMSource, error) {
	return nil, nil
}

func TestCommandExecutorConnect(t *testing.T) {
	ce := CommandExecutor{ProcessorStub{}}
	ce.Connect("url", "label")
}

func TestCommandExecutorUpdate(t *testing.T) {
	ce := CommandExecutor{ProcessorStub{}}
	ce.Update("srcName", "label")
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
