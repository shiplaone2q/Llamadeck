package backends

import (
	"fmt"
	"llamadeck/pkg/validation"
	"strconv"
	"strings"
)

// CustomServerOptions runs a user-defined command configured as
// backends.custom.<name> in the config file. Args are literal strings,
// not a flag grammar; {port} and {model} are substituted at build time.
type CustomServerOptions struct {
	Name  string   `json:"name,omitempty"`
	Model string   `json:"model,omitempty"`
	Host  string   `json:"host,omitempty"`
	Port  int      `json:"port,omitempty"`
	Args  []string `json:"args,omitempty"`
}

func (o *CustomServerOptions) GetModel() string {
	if o == nil {
		return ""
	}
