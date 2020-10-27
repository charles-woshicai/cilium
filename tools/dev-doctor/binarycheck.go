package main

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/blang/semver"
)

// A binaryCheck checks that a binary called name is installed and optionally at
// least version minVersion.
type binaryCheck struct {
	name          string
	ifNotFound    checkResult
	versionArgs   []string
	versionRegexp *regexp.Regexp
	minVersion    *semver.Version
}

func (c *binaryCheck) id() string {
	return c.name
}

func (c *binaryCheck) run() (checkResult, string) {
	path, err := exec.LookPath(c.name)
	if err != nil {
		return c.ifNotFound, err.Error()
	}

	if c.versionArgs == nil {
		return checkOK, fmt.Sprintf("found %s", path)
	}

	output, err := exec.Command(path, c.versionArgs...).CombinedOutput()
	if err != nil {
		return checkFailed, err.Error()
	}

	versionBytes := output
	if c.versionRegexp != nil {
		match := c.versionRegexp.FindSubmatch(versionBytes)
		if len(match) != 2 {
			return checkFailed, fmt.Sprintf("found %s, could not parse version from %s", path, versionBytes)
		}
		versionBytes = match[1]
	}
	version, err := semver.Parse(string(versionBytes))
	if err != nil {
		return checkFailed, err.Error()
	}

	if c.minVersion != nil && version.LT(*c.minVersion) {
		return checkError, fmt.Sprintf("found %s, version %s, need %s", path, version, c.minVersion)
	}

	return checkOK, fmt.Sprintf("found %s, version %s", path, version)
}
