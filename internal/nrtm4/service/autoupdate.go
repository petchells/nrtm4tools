package service

import (
	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/util"
)

var pool = util.NewExecutionPool(2)

func NewAutoUpdater(p NRTMProcessor, source persist.NRTMSource) autoupdater {
	return autoupdater{
		source: source,
		p:      p,
	}
}

type autoupdater struct {
	source persist.NRTMSource
	p      NRTMProcessor
}

func (u *autoupdater) Start(preDelay bool) error {

	exe := pool.Acquire()
	defer pool.Release(exe)
	_, err := u.p.Update(u.source.Source, u.source.Label)
	return err
}
