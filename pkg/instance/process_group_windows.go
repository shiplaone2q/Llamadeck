//go:build windows

package instance

import "os/exec"

func setProcAttrs(cmd *exec.Cmd) {
