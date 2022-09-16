package services

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/cjlapao/common-go/guard"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/helper/linux_helper"
	"github.com/cjlapao/common-go/helper/linux_service"
	"github.com/cjlapao/common-go/helper/linux_user"
	"github.com/cjlapao/postfixcli-backend-api/ioc"
	"github.com/cjlapao/postfixcli-backend-api/models"
)

var globalOpenDKIMService *OpenDKIMService

const (
	OpenDKIMUserName  string = "opendkim"
	OpenDKIMGroupName string = "opendkim"
)

type OpenDKIMService struct {
	Context         context.Context
	SocketDirectory string
}

func GetOpenDKIMService() *OpenDKIMService {
	if globalOpenDKIMService != nil {
		return globalOpenDKIMService
	}

	return NewOpenDKIMService()
}

func NewOpenDKIMService() *OpenDKIMService {
	if helper.GetOperatingSystem() != helper.LinuxOs {
		ioc.Log.Fatal("This service only works in Linux, exiting")
	}

	globalOpenDKIMService = &OpenDKIMService{
		SocketDirectory: "/var/spool/postfix/opendmarc",
	}

	globalOpenDKIMService.Context = context.Background()

	return globalOpenDKIMService
}

func (svc *OpenDKIMService) Name() string {
	return "OpenDMARC"
}

func (svc *OpenDKIMService) Start() error {
	return linux_service.Start("opendmarc")
}

func (svc *OpenDKIMService) Stop() error {
	return linux_service.Stop("opendmarc")
}

func (svc *OpenDKIMService) Restart() error {
	return linux_service.Restart("opendmarc")
}

func (svc *OpenDKIMService) Status() linux_service.LinuxServiceState {
	return linux_service.Status("opendmarc")
}

func (svc *OpenDKIMService) Init() error {
	ioc.Log.Info("Starting initialization of OpenDMARC")
	if svc.Status() != linux_service.LinuxServiceRunning {
		if err := svc.Stop(); err != nil {
			return err
		}
	}

	if !helper.DirectoryExists(svc.SocketDirectory) {
		if !helper.CreateDirectory(svc.SocketDirectory, fs.ModePerm) {
			return fmt.Errorf("there was an error creating the folder %v", svc.SocketDirectory)
		} else {
			ioc.Log.Info("Created OpenDMARC default SOCKET folder")
		}
	}

	err := linux_helper.ChangeOwner(svc.SocketDirectory, OpenDMARCUserName, OpenDMARCGroupName, true)

	if err != nil {
		return err
	}

	err = linux_helper.ChangeFileMode(svc.SocketDirectory, "750", true)

	if err != nil {
		return err
	}
	ioc.Log.Info("Updated user permissions in the OpenDMARC default SOCKET folder")

	err = linux_user.AddToGroup(PostfixUserName, OpenDMARCGroupName)
	if err != nil {
		return err
	}
	ioc.Log.Info("Added %v to %v group", PostfixUserName, OpenDMARCGroupName)

	return nil
}

func (svc *OpenDKIMService) Configure(config models.MailServerConfig) error {
	if err := guard.EmptyOrNil(config); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(config.Domain); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(config.SubDomain); err != nil {
		return err
	}

	openDMARCConfig := models.ConfigFile{
		FileName:       "opendmarc.conf",
		DestinationDir: "/etc",
		TemplateName:   "opendmarc.conf.tpl",
	}

	return applyConfiguration(svc, config, openDMARCConfig)
}
