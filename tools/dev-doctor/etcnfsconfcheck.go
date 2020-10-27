package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/ini.v1"
)

// An etcdNFSConfCheck checks /etc/nfs.conf.
type etcNFSConfCheck struct{}

func (etcNFSConfCheck) id() string {
	return "/etc/nfs.conf"
}

func (etcNFSConfCheck) run() (checkResult, string) {
	data, err := ioutil.ReadFile("/etc/nfs.conf")
	switch {
	case os.IsNotExist(err):
		return checkError, "/etc/nfs.conf does not exist"
	case err != nil:
		return checkFailed, err.Error()
	}

	var nfsConf struct {
		NFSD struct {
			TCP string `ini:"tcp"`
		} `ini:"nfsd"`
	}
	if err := ini.MapTo(&nfsConf, data); err != nil {
		return checkError, err.Error()
	}

	switch {
	case nfsConf.NFSD.TCP == "":
		return checkError, `nfsd.tcp is not set, want "y"`
	case nfsConf.NFSD.TCP != "y":
		return checkError, fmt.Sprintf(`nfsd.tcp is %q, want "y"`, nfsConf.NFSD.TCP)
	}

	return checkOK, `nfsd.tcp is "y"`
}
