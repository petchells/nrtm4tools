package cli

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

func (ps ProcessorStub) RemoveSource(src, label string) error {
	return nil
}

func TestCommandExecutorConnect(t *testing.T) {
	ce := CommandExecutor{ProcessorStub{}}
	ce.Connect("url", "label")
}

func TestCommandExecutorUpdate(t *testing.T) {
	ce := CommandExecutor{ProcessorStub{}}
	ce.Update("srcName", "label")
}
