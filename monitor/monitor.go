package monitor

import (
	"fmt"
	"sync"

	helpersJSON "github.com/codemodify/systemkit-helpers"
	helpersReflect "github.com/codemodify/systemkit-helpers"
	logging "github.com/codemodify/systemkit-logging"
	loggingC "github.com/codemodify/systemkit-logging/contracts"
	"github.com/codemodify/systemkit-processes/contracts"
	"github.com/codemodify/systemkit-processes/internal"
)

const logID = "SKIT-PROCESS-MONITOR"

// processMonitor - Represents Windows service
type processMonitor struct {
	procs     map[string]contracts.RuningProcess
	procsSync *sync.RWMutex
}

// New -
func New() contracts.Monitor {
	return &processMonitor{
		procs:     map[string]contracts.RuningProcess{},
		procsSync: &sync.RWMutex{},
	}
}

// Spawn -
func (thisRef *processMonitor) Spawn(id string, processTemplate contracts.ProcessTemplate) error {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: preparing to spawn [%s], details [%s]", logID, id, helpersJSON.AsJSONString(processTemplate)),
	})

	thisRef.procsSync.Lock()

	thisRef.procs[id] = internal.NewRuningProcess(processTemplate)
	thisRef.procsSync.Unlock()

	return thisRef.Start(id)
}

// Start -
func (thisRef *processMonitor) Start(id string) error {
	if thisRef.GetRuningProcess(id).IsRunning() {
		return nil
	}

	thisRef.procsSync.Lock()
	defer thisRef.procsSync.Unlock()

	// CHECK-IF-EXISTS
	if _, ok := thisRef.procs[id]; !ok {
		return fmt.Errorf("ID %s, CHECK-IF-EXISTS failed", id)
	}

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: requesting to start [%s]", logID, id),
	})

	err := thisRef.procs[id].Start()
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: error starting [%s], details [%s]", logID, thisRef.procs[id], err.Error()),
		})

		return err
	}

	return nil
}

// Stop -
func (thisRef *processMonitor) Stop(id string) error {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: requesting to stop [%s]", logID, id),
	})

	if !thisRef.GetRuningProcess(id).IsRunning() {
		return nil
	}

	thisRef.procsSync.RLock()
	defer thisRef.procsSync.RUnlock()

	return thisRef.procs[id].Stop()
}

// Restart -
func (thisRef processMonitor) Restart(id string) error {
	err := thisRef.Stop(id)
	if err != nil {
		return err
	}

	return thisRef.Start(id)
}

// StopAll -
func (thisRef processMonitor) StopAll() []error {
	thisRef.procsSync.RLock()
	defer thisRef.procsSync.RUnlock()

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: requesting to stop all", logID),
	})

	allErrors := []error{}
	for k := range thisRef.procs {
		allErrors = append(allErrors, thisRef.Stop(k))
	}

	return allErrors
}

// GetRuningProcess -
func (thisRef processMonitor) GetRuningProcess(id string) contracts.RuningProcess {
	thisRef.procsSync.RLock()
	defer thisRef.procsSync.RUnlock()

	// CHECK-IF-EXISTS
	if _, ok := thisRef.procs[id]; !ok {
		return internal.NewEmptyRuningProcess()
	}

	return thisRef.procs[id]
}

// RemoveFromMonitor -
func (thisRef *processMonitor) RemoveFromMonitor(id string) {
	thisRef.procsSync.Lock()
	defer thisRef.procsSync.Unlock()

	if _, ok := thisRef.procs[id]; ok {
		delete(thisRef.procs, id) // delete
	}
}

// GetAllIDs -
func (thisRef processMonitor) GetAllIDs() []string {
	thisRef.procsSync.RLock()
	defer thisRef.procsSync.RUnlock()

	allIds := []string{}
	for k := range thisRef.procs {
		allIds = append(allIds, k)
	}

	return allIds
}
