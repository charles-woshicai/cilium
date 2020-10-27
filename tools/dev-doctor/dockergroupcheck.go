package main

import (
	"fmt"
	"os/user"
	"runtime"
)

// A dockerGroupCheck checks that the current user is in the docker group.
type dockerGroupCheck struct{}

func (dockerGroupCheck) id() string {
	return "docker-group"
}

func (dockerGroupCheck) run() (checkResult, string) {
	if runtime.GOOS != "linux" {
		return checkSkipped, "docker group only used on linux"
	}

	currentUser, err := user.Current()
	if err != nil {
		return checkFailed, err.Error()
	}

	groupIDs, err := currentUser.GroupIds()
	if err != nil {
		return checkFailed, err.Error()
	}

	dockerGroup, err := user.LookupGroup("docker")
	if err != nil {
		return checkFailed, err.Error()
	}

	for _, groupID := range groupIDs {
		if groupID == dockerGroup.Gid {
			return checkOK, fmt.Sprintf("%s in docker group", currentUser.Username)
		}
	}

	return checkError, fmt.Sprintf("%s not in docker group", currentUser.Username)
}
