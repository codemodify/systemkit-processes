package contracts

// ProcessOutputReader -
type ProcessOutputReader func(params interface{}, outputData []byte)

// ProcessStoppedDelegate -
type ProcessStoppedDelegate func(params interface{})

// ProcessTemplate -
type ProcessTemplate struct {
	Executable       string   `json:"executable"`
	Args             []string `json:"args"`
	WorkingDirectory string   `json:"workingDirectory"`
	Environment      []string `json:"environment"`
}
