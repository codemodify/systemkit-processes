package monitor

import (
	"fmt"
	"sync"
	"time"

	logging "github.com/codemodify/systemkit-logging"
	"github.com/codemodify/systemkit-processes/contracts"
	"github.com/codemodify/systemkit-processes/helpers"
	"github.com/codemodify/systemkit-processes/internal"
)

const logID = "PROCESS-MONITOR"

// processMonitor - Represents Windows service
type processMonitor struct {
	procs        map[string]contracts.RuningProcess
	procsSync    *sync.Mutex
	procTagIndex int64
}

// New -
func New() contracts.Monitor {
	return &processMonitor{
		procs:        map[string]contracts.RuningProcess{},
		procsSync:    &sync.Mutex{},
		procTagIndex: 0,
	}
}

// Spawn -
func (thisRef *processMonitor) Spawn(processTemplate contracts.ProcessTemplate) (string, error) {
	thisRef.procsSync.Lock()
	defer thisRef.procsSync.Unlock()

	tag := fmt.Sprintf("gen-tag-%d", thisRef.procTagIndex)
	thisRef.procTagIndex++

	return tag, thisRef.SpawnWithTag(processTemplate, tag)
}

// SpawnWithID -
func (thisRef *processMonitor) SpawnWithTag(processTemplate contracts.ProcessTemplate, tag string) error {
	logging.Debugf("%s: spawn %s, %s", logID, tag, helpers.AsJSONString(processTemplate))

	thisRef.procsSync.Lock()
	thisRef.procs[tag] = internal.NewRuningProcess(processTemplate)
	thisRef.procsSync.Unlock()

	return thisRef.Start(tag)
}

// Start -
func (thisRef *processMonitor) Start(tag string) error {
	if thisRef.GetProcess(tag).IsRunning() {
		return nil
	}

	thisRef.procsSync.Lock()
	defer thisRef.procsSync.Unlock()

	// CHECK-IF-EXISTS
	if _, ok := thisRef.procs[tag]; !ok {
		return fmt.Errorf("ID %s, CHECK-IF-EXISTS failed", tag)
	}

	logging.Debugf("%s: start %s", logID, tag)

	err := thisRef.procs[tag].Start()
	if err != nil {
		logging.Errorf("%s: start-FAIL %s, %s", logID, thisRef.procs[tag], err.Error())
		return err
	}

	return nil
}

// Stop -
func (thisRef *processMonitor) Stop(tag string) error {
	return thisRef.StopWithTimeout(tag, 3, 0*time.Millisecond)
}

func (thisRef *processMonitor) StopWithTimeout(tag string, attempts int, waitTimeout time.Duration) error {
	if !thisRef.GetProcess(tag).IsRunning() {
		return nil
	}

	thisRef.procsSync.Lock()

	// CHECK-IF-EXISTS
	var rp contracts.RuningProcess
	var keyExists bool

	if rp, keyExists = thisRef.procs[tag]; !keyExists {
		thisRef.procsSync.Unlock()
		return nil
	}

	thisRef.procsSync.Unlock()

	return rp.Stop(tag, attempts, waitTimeout)
}

// Restart -
func (thisRef *processMonitor) Restart(tag string) error {
	err := thisRef.Stop(tag)
	if err != nil {
		return err
	}

	return thisRef.Start(tag)
}

// StopAll -
func (thisRef *processMonitor) StopAllInParallel() {
	thisRef.procsSync.Lock()
	defer thisRef.procsSync.Unlock()

	for k := range thisRef.procs {
		go func(tag string) {
			thisRef.Stop(tag)
		}(k)
	}
}

// GetRuningProcess -
func (thisRef *processMonitor) GetProcess(tag string) contracts.RuningProcess {
	thisRef.procsSync.Lock()
	defer thisRef.procsSync.Unlock()

	// CHECK-IF-EXISTS
	if _, ok := thisRef.procs[tag]; !ok {
		return internal.NewEmptyRuningProcess()
	}

	return thisRef.procs[tag]
}

// RemoveFromMonitor -
func (thisRef *processMonitor) RemoveFromMonitor(tag string) {
	thisRef.procsSync.Lock()
	defer thisRef.procsSync.Unlock()

	if _, ok := thisRef.procs[tag]; ok {
		delete(thisRef.procs, tag) // delete
	}
}

// GetAllTags -
func (thisRef *processMonitor) GetAllTags() []string {
	thisRef.procsSync.Lock()
	defer thisRef.procsSync.Unlock()

	allTags := []string{}
	for k := range thisRef.procs {
		allTags = append(allTags, k)
	}

	return allTags
}
