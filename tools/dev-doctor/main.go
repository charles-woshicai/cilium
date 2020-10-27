// Command doctor checks the development setup for common problems.
package main

import (
	"fmt"
	"os"
	"regexp"
	"text/tabwriter"

	"github.com/blang/semver"
)

// A checkResult is the result of a check.
type checkResult int

const (
	checkSkipped checkResult = -1 // The check was skipped.
	checkOK      checkResult = 0  // The check completed and did not find any problems.
	checkWarning checkResult = 1  // The check completed and found something that might indicate a problem.
	checkError   checkResult = 2  // The check completed and found a problem.
	checkFailed  checkResult = 3  // The check could not be completed.
)

// A check is an individual check.
type check interface {
	id() string
	run() (checkResult, string)
}

var checkResultStr = map[checkResult]string{
	checkSkipped: "skipped",
	checkOK:      "ok",
	checkWarning: "warning",
	checkError:   "error",
	checkFailed:  "failed",
}

func main() {
	checks := []check{
		osArchCheck{},
		&binaryCheck{
			name:          "go",
			versionArgs:   []string{"version"},
			versionRegexp: regexp.MustCompile(`go version go(\d+\.\d+\.\d+)`),
			minVersion:    &semver.Version{Major: 1, Minor: 15, Patch: 0},
		},
		&binaryCheck{
			name:          "clang",
			versionArgs:   []string{"--version"},
			versionRegexp: regexp.MustCompile(`clang version (\d+\.\d+\.\d+)`),
			minVersion:    &semver.Version{Major: 10, Minor: 0, Patch: 0},
		},
		// FIXME add llvm check?
		// FIXME add libelf-devel check?
		&binaryCheck{
			name:          "ginkgo",
			versionArgs:   []string{"version"},
			versionRegexp: regexp.MustCompile(`Ginkgo Version (\d+\.\d+\.\d+)`),
			minVersion:    &semver.Version{Major: 1, Minor: 4, Patch: 0},
		},
		// FIXME add gomega check?
		&binaryCheck{
			name:          "golangci-lint",
			versionArgs:   []string{"--version"},
			versionRegexp: regexp.MustCompile(`(\d+\.\d+\.\d+\S*)`),
			minVersion:    &semver.Version{Major: 1, Minor: 27, Patch: 0},
		},
		&binaryCheck{
			name: "docker",
		},
		&binaryCheck{
			name:       "docker-compose",
			ifNotFound: checkWarning,
		},
		&binaryCheck{
			name:          "vagrant",
			versionArgs:   []string{"--version"},
			versionRegexp: regexp.MustCompile(`Vagrant (\d+\.\d+\.\d+)`),
			minVersion:    &semver.Version{Major: 2, Minor: 0, Patch: 0},
		},
		&binaryCheck{
			name:       "virtualbox",
			ifNotFound: checkWarning,
		},
		&binaryCheck{
			name:          "pip3",
			versionArgs:   []string{"--version"},
			versionRegexp: regexp.MustCompile(`pip (\d+\.\d+\.\d+)`),
		},
		dockerGroupCheck{},
		etcNFSConfCheck{},
		&iptablesRuleCheck{
			rule: []string{"INPUT", "-p", "tcp", "-s", "192.168.34.0/24", "--dport", "111", "-j", "ACCEPT"},
		},
		&iptablesRuleCheck{
			rule: []string{"INPUT", "-p", "tcp", "-s", "192.168.34.0/24", "--dport", "2049", "-j", "ACCEPT"},
		},
		&iptablesRuleCheck{
			rule: []string{"INPUT", "-p", "tcp", "-s", "192.168.34.0/24", "--dport", "20048", "-j", "ACCEPT"},
		},
	}

	worstResult := checkOK
	w := tabwriter.NewWriter(os.Stdout, 3, 0, 3, ' ', 0)
	fmt.Fprintf(w, "RESULT\tCHECK\tMESSAGE\n")
	for _, check := range checks {
		checkResult, message := check.run()
		fmt.Fprintf(w, "%s\t%s\t%s\n", checkResultStr[checkResult], check.id(), message)
		if checkResult > worstResult {
			worstResult = checkResult
		}
	}
	w.Flush()

	if worstResult > checkOK {
		fmt.Println("\nSee https://docs.cilium.io/en/latest/contributing/development/dev_setup/")
	}
	if worstResult > checkWarning {
		os.Exit(1)
	}
}
