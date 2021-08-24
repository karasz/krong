package krong

import (
	"fmt"
)

type Agent struct {
	Endpoint string
	Command  string
	Secret   string
}

func (a *Agent) String() string {
	return fmt.Sprintf("Endpoint: %s, Command: %s, Secret: %s", a.Endpoint, a.Command, a.Secret)
}
