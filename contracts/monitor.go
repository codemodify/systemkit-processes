package contracts

// Monitor - process monitor
type Monitor interface {
	Spawn(process ProcessTemplate) (string, error)
	SpawnWithTag(process ProcessTemplate, tag string) error
	Start(tag string) error
	Stop(tag string) error
	Restart(tag string) error
	StopAll() []error
	GetProcess(tag string) RuningProcess
	RemoveFromMonitor(tag string)
	GetAllTags() []string
}
