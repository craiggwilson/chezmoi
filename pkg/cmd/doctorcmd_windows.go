package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/twpayne/chezmoi/v2/pkg/chezmoi"
	"github.com/twpayne/chezmoi/v2/pkg/chezmoilog"
)

type (
	systeminfoCheck struct{}
	umaskCheck      struct{ skippedCheck }
	unameCheck      struct{ skippedCheck }
)

func (systeminfoCheck) Name() string {
	return "systeminfo"
}

func (systeminfoCheck) Run(system chezmoi.System, homeDirAbsPath chezmoi.AbsPath) (checkResult, string) {
	cmd := exec.Command("systeminfo")
	data, err := chezmoilog.LogCmdOutput(cmd)
	if err != nil {
		return checkResultFailed, err.Error()
	}

	var osName, osVersion string
	s := bufio.NewScanner(bytes.NewReader(data))
	for s.Scan() {
		switch key, value, found := chezmoi.CutString(s.Text(), ":"); {
		case !found:
			// Do nothing.
		case key == "OS Name":
			osName = strings.TrimSpace(value)
		case key == "OS Version":
			osVersion = strings.TrimSpace(value)
		}
	}
	return checkResultOK, fmt.Sprintf("%s (%s)", osName, osVersion)
}
