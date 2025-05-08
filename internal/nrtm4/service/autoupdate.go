package service

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/util"
)

var (
	updaterPool    = util.NewExecutionPool(2)
	updateRegister = make(map[uint64]*AutoUpdater)
)

// GetAutoUpdaterInstance should be started in a goroutine
func GetAutoUpdaterInstance(p NRTMProcessor, sourceID uint64) *AutoUpdater {
	au, ok := updateRegister[sourceID]
	if ok && au != nil {
		return au
	}
	au = &AutoUpdater{
		sourceID: sourceID,
		p:        p,
	}
	updateRegister[sourceID] = au
	au.initialize()
	return au
}

// AutoUpdater handles the timing of updates to a repo
type AutoUpdater struct {
	sourceID uint64
	p        NRTMProcessor
	mu       sync.Mutex
	t        *time.Ticker
	ch       chan error
	interval int
	running  bool
}

// IsRunning returns true if the autoupdater is running
func (u *AutoUpdater) IsRunning() bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.running
}

// ResetTimer restarts the auto update timer with the latest duration
// func (u *AutoUpdater) ResetTimer() {
// 	src := latestSource(u.p, u.sourceID)
// 	u.mu.Lock()
// 	defer u.mu.Unlock()
// 	if src == nil {
// 		if u.running {
// 			u.ch <- nil
// 		}
// 		return
// 	}
// 	dur, err := time.ParseDuration(fmt.Sprintf("%dm", interval))
// 	if err != nil {
// 		return err
// 	}

// }

func (u *AutoUpdater) initialize() error {
	src := latestSource(u.p, u.sourceID)
	if src == nil {
		return errors.New("source not found")
	}
	u.interval = src.Properties.AutoUpdateInterval
	if u.interval < 0 {
		logger.Error("Bad value in Properties", "AutoUpdateInterval", u.interval)
		return nil
	}
	u.ch = make(chan error)
	var err error

	if u.interval == 0 {
		dur, _ := time.ParseDuration("9h")
		if u.t == nil {
			u.t = time.NewTicker(dur)
		}
		u.t.Stop()
	} else {
		dur, _ := time.ParseDuration(fmt.Sprintf("%dm", u.interval))
		logger.Info("setting ticker", "dur", dur, "u.interval", u.interval)
		u.t = time.NewTicker(dur)
		u.running = true
	}
	go func() {
		for {
			select {
			case <-u.t.C:
				go u.doUpdateWithReconnectOnFail()
			case err = <-u.ch:
				u.t.Stop()
				u.mu.Lock()
				defer u.mu.Unlock()
				u.running = false
				delete(updateRegister, u.sourceID)
				if err != nil {
					logger.Error("Failed to synchronize with source", "source", src.Source, "label", src.Label)
				}
				return
			}
		}
	}()
	return nil
}

// Start sets up the timer and performs updates.
func (u *AutoUpdater) Start(preDelay bool) error {
	src := latestSource(u.p, u.sourceID)
	if src == nil {
		return errors.New("source not found")
	}
	if u.interval == src.Properties.AutoUpdateInterval {
		return nil
	}
	u.interval = src.Properties.AutoUpdateInterval
	if u.interval < 0 {
		logger.Error("Bad value in Properties", "AutoUpdateInterval", u.interval)
		return nil
	}
	var err error
	if u.interval == 0 {
		u.t.Stop()
		u.running = false
	} else {
		// if !preDelay {
		// 	go u.doUpdateWithReconnectOnFail()
		// }
		dur, _ := time.ParseDuration(fmt.Sprintf("%dm", u.interval))
		logger.Info("setting ticker", "dur", dur, "u.interval", u.interval)
		u.t.Reset(dur)
		u.running = true
	}
	return err
}

// Destroy stops this autoupdater and removes it from the registry
func (u *AutoUpdater) Destroy() {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.running {
		u.ch <- nil
	}
}

func (u *AutoUpdater) doUpdateWithReconnectOnFail() {
	if !u.mu.TryLock() {
		return
	}
	defer u.mu.Unlock()

	src := latestSource(u.p, u.sourceID)
	if src == nil {
		u.Destroy()
		return
	}
	if src.Status != "ok" {
		return
	}

	lock := updaterPool.Acquire()
	defer func() {
		updaterPool.Release(lock)
	}()

	if src.Properties.AutoUpdateInterval != u.interval {
		// interval changed
		u.interval = src.Properties.AutoUpdateInterval
		if u.interval == 0 {
			u.t.Stop()
			return
		} else {
			d, _ := time.ParseDuration(fmt.Sprintf("%dm", u.interval))
			u.t.Reset(d)
		}
	}
	// Do an update
	logger.Info("AutoUpdater will update", "src.Source", src.Source, "src.Label", src.Label)
	_, err := u.p.Update(src.Source, src.Label)
	if err == nil {
		logger.Info("AutoUpdater finished updating", "src.Source", src.Source, "src.Label", src.Label)
		return
	}
	// There was en error with the update; delete or rename repo
	logger.Info("AutoUpdater failed to update, reconnecting...", "src.Source", src.Source, "src.Label", src.Label)
	label := src.Label // preserve original label in case of rename
	switch src.Properties.UpdateMode {
	case persist.UpdateModeReplace:
		err = u.p.RemoveSource(src.Source, src.Label)
	case persist.UpdateModePreserve:
		err = relabelAutoUpdateFailure(u.p, src)
	}
	if err != nil {
		logger.Error("AutoUpdater got an error when removing or renaming", "src.Source", src.Source, "src.Label", src.Label)
		u.ch <- err
		return
	}
	u.t.Stop()
	// Connect to source
	err = u.p.Connect(src.NotificationURL, label)
	u.ch <- err
	// srcs, err := u.p.ListSources()
	// if err != nil {
	// 	return
	// }
	// logger.Debug("Looking for latest source...", "src.Source", src.Source, "label", label)
	// for _, s := range srcs {
	// 	if s.Source == src.Source && s.Label == label {
	// 		logger.Debug("Found it")
	// 		u.sourceID = s.ID
	// 		if s.Properties.AutoUpdateInterval != u.interval {
	// 			u.interval = s.Properties.AutoUpdateInterval
	// 			if u.interval == 0 {
	// 				u.t.Stop()
	// 			} else {
	// 				d, _ := time.ParseDuration(fmt.Sprintf("%dm", u.interval))
	// 				u.t.Reset(d)
	// 			}
	// 		}
	// 		break
	// 	}
	// }
}

func latestSource(p NRTMProcessor, id uint64) *persist.NRTMSourceDetails {
	srcs, err := p.ListSources()
	if err != nil {
		return nil
	}
	for _, s := range srcs {
		if s.ID == id {
			return &s
		}
	}
	return nil
}

var autoLabelRe = regexp.MustCompile(` :AUTO:([1-9][0-9])*$`)

func relabelAutoUpdateFailure(p NRTMProcessor, src *persist.NRTMSourceDetails) error {
	srcs, err := p.ListSources()
	if err != nil {
		return nil
	}
	idx := 0
	for _, s := range srcs {
		if s.ID != src.ID && s.Source == src.Source && strings.HasPrefix(s.Label, src.Label) {
			ext := s.Label[len(src.Label):]
			m := autoLabelRe.FindSubmatch([]byte(ext))
			if len(m) == 2 {
				n, err := strconv.Atoi(string(m[1]))
				if err != nil {
					logger.Debug("Error converting string to int", "string", string(m[1]), "error", err)
				}
				if n > idx {
					idx = n
				}
			}
		}
	}
	idx++
	_, err = p.ReplaceLabel(src.Source, src.Label, fmt.Sprintf("%v :AUTO:%d", src.Label, idx))
	return err
}
