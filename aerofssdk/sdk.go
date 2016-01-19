package aerofssdk

// A descriptor for a response from an AeroFS Appliance when an HTTP {4,5}XX
// code is received
type AeroError struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	HttpCode int
}

var SFPermissions []string = []string{"WRITE", "MANAGE"}
