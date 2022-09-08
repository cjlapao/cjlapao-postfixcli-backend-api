package services

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/helper/linux_helper"
	"github.com/cjlapao/common-go/helper/linux_service"
	"github.com/cjlapao/postfixcli-backend-api/ioc"
)

var globalOpenDMarc *OpenDMarc

const (
	OpenDMarcSocketDirectory string = "/var/spool/postfix/opendmarc"
	OpenDMarcUserName        string = "opendmarc"
	OpenDMarcGroupName       string = "opendmarc"
)

type OpenDMarc struct {
	Context context.Context
}

func GetOpenDMarc() *OpenDMarc {
	if globalOpenDMarc != nil {
		return globalOpenDMarc
	}

	return NewOpenDMarc()
}

func NewOpenDMarc() *OpenDMarc {
	if helper.GetOperatingSystem() != helper.LinuxOs {
		ioc.Log.Fatal("This service only works in Linux, exiting")
	}

	globalOpenDMarc = &OpenDMarc{}

	globalOpenDMarc.Context = context.Background()

	return globalOpenDMarc
}

func (svc *OpenDMarc) Start() error {
	return linux_service.Start("opendmarc")
}

func (svc *OpenDMarc) Stop() error {
	return linux_service.Stop("opendmarc")
}

func (svc *OpenDMarc) Restart() error {
	return linux_service.Restart("opendmarc")
}

func (svc *OpenDMarc) Status() linux_service.LinuxServiceState {
	return linux_service.Status("opendmarc")
}

func (svc OpenDMarc) Init() error {
	ioc.Log.Info("Starting initialization of OpenDMarc")
	if svc.Status() != linux_service.LinuxServiceRunning {
		svc.Stop()
	}

	if !helper.DirectoryExists(OpenDMarcSocketDirectory) {
		if !helper.CreateDirectory(OpenDMarcSocketDirectory, fs.ModePerm) {
			return fmt.Errorf("there was an error creating the folder %v", OpenDMarcSocketDirectory)
		} else {
			ioc.Log.Info("Created OpenDMarc default SOCKET folder")
		}
	}

	err := linux_helper.ChangeOwner(OpenDMarcSocketDirectory, OpenDMarcUserName, OpenDMarcGroupName, true)

	if err != nil {
		return err
	}

	err = linux_helper.ChangeFileMode(OpenDMarcSocketDirectory, "750", true)

	if err != nil {
		return err
	}
	ioc.Log.Info("Updated user permissions in the OpenDMarc default SOCKET folder")

	return nil
}
