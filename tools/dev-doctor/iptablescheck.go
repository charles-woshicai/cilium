package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// An iptablesRuleCheck checks that the give iptables rule is present.
type iptablesRuleCheck struct {
	rule []string
}

func (c *iptablesRuleCheck) id() string {
	return "iptables-rule"
}

func (c *iptablesRuleCheck) run() (checkResult, string) {
	if runtime.GOOS != "linux" {
		return checkSkipped, "iptables only used on linux"
	}

	var cmd *exec.Cmd
	if os.Getuid() == 0 {
		cmd = exec.Command("iptables", append([]string{"-C"}, []string(c.rule)...)...)
		cmd.Stdin = nil
		cmd.Stdout = nil
		cmd.Stderr = nil

	} else {
		cmd = exec.Command("sudo", append([]string{"iptables", "-C"}, []string(c.rule)...)...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = nil
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return checkFailed, err.Error()
	}
	if cmd.ProcessState.ExitCode() != 0 {
		return checkError, fmt.Sprintf("%s rule not found", strings.Join(c.rule, " "))
	}

	return checkOK, strings.Join(c.rule, " ")
}
