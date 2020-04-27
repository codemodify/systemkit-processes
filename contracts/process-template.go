package contracts

// ProcessOutputReader -
type ProcessOutputReader func([]byte)

// ProcessStoppedDelegate -
type ProcessStoppedDelegate func()

// ProcessTemplate -
type ProcessTemplate struct {
	Executable       string   `json:"executable"`
	Args             []string `json:"args"`
	WorkingDirectory string   `json:"workingDirectory"`
	Environment      []string `json:"environment"`
}
