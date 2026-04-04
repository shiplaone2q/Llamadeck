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
	return o.Model
}

func (o *CustomServerOptions) GetPort() int {
	if o == nil {
		return 0
	}
	return o.Port
}

func (o *CustomServerOptions) SetPort(port int) {
	if o == nil {
		return
	}
	o.Port = port
}

func (o *CustomServerOptions) GetHost() string {
	if o == nil {
		return "localhost"
	}
	return o.Host
}

func (o *CustomServerOptions) Validate() error {
	if o == nil {
		return validation.ValidationError(fmt.Errorf("custom server options cannot be nil for custom backend"))
	}

	if o.Name == "" {
		return validation.ValidationError(fmt.Errorf("custom backend requires a name referencing a backends.custom.<name> config entry"))
	}

	if o.Port < 0 || o.Port > 65535 {
		return validation.ValidationError(fmt.Errorf("invalid port range: %d", o.Port))
