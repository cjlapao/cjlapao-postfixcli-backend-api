package services

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"text/template"

	"github.com/cjlapao/common-go/guard"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/helper/linux_helper"
	"github.com/cjlapao/common-go/helper/linux_service"
	"github.com/cjlapao/common-go/helper/linux_user"
	"github.com/cjlapao/postfixcli-backend-api/ioc"
	"github.com/cjlapao/postfixcli-backend-api/models"
)

var globalOpenDMARCService *OpenDMARCService

const (
	SocketDirectory string = "/var/spool/postfix/opendmarc"
	UserName        string = "opendmarc"
	GroupName       string = "opendmarc"
	ConfigFileName  string = "opendmarc.conf"
	ConfigFolder    string = "/etc"
)

type OpenDMARCService struct {
	Context context.Context
}

func GetOpenDMARCService() *OpenDMARCService {
	if globalOpenDMARCService != nil {
		return globalOpenDMARCService
	}

	return NewOpenDMARCService()
}

func NewOpenDMARCService() *OpenDMARCService {
	if helper.GetOperatingSystem() != helper.LinuxOs {
		ioc.Log.Fatal("This service only works in Linux, exiting")
	}

	globalOpenDMARCService = &OpenDMARCService{}

	globalOpenDMARCService.Context = context.Background()

	return globalOpenDMARCService
}

func (svc *OpenDMARCService) Start() error {
	return linux_service.Start("opendmarc")
}

func (svc *OpenDMARCService) Stop() error {
	return linux_service.Stop("opendmarc")
}

func (svc *OpenDMARCService) Restart() error {
	return linux_service.Restart("opendmarc")
}

func (svc *OpenDMARCService) Status() linux_service.LinuxServiceState {
	return linux_service.Status("opendmarc")
}

func (svc *OpenDMARCService) Init() error {
	ioc.Log.Info("Starting initialization of OpenDMARC")
	if svc.Status() != linux_service.LinuxServiceRunning {
		svc.Stop()
	}

	if !helper.DirectoryExists(SocketDirectory) {
		if !helper.CreateDirectory(SocketDirectory, fs.ModePerm) {
			return fmt.Errorf("there was an error creating the folder %v", SocketDirectory)
		} else {
			ioc.Log.Info("Created OpenDMARC default SOCKET folder")
		}
	}

	err := linux_helper.ChangeOwner(SocketDirectory, UserName, GroupName, true)

	if err != nil {
		return err
	}

	err = linux_helper.ChangeFileMode(SocketDirectory, "750", true)

	if err != nil {
		return err
	}
	ioc.Log.Info("Updated user permissions in the OpenDMARC default SOCKET folder")

	err = linux_user.AddToGroup(PostfixUserName, GroupName)
	if err != nil {
		return err
	}
	ioc.Log.Info("Added %v to %v group", PostfixUserName, GroupName)

	return nil
}

func (svc *OpenDMARCService) Config(config models.MailServerConfig) error {
	if err := guard.EmptyOrNil(config); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(config.Domain); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(config.SubDomain); err != nil {
		return err
	}

	templateContent, err := helper.ReadFromFile("./templates/opendmarc.tpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("openDMARC").Parse(string(templateContent))
	if err != nil {
		return err
	}

	var configResult bytes.Buffer

	if err = tmpl.Execute(&configResult, config); err != nil {
		return err
	}

	if err := svc.saveConfig(configResult.String()); err != nil {
		return err
	}

	return nil
}

func (svc *OpenDMARCService) saveConfig(content string) error {
	configPath := helper.JoinPath(ConfigFolder, ConfigFileName)
	tempPath := "./template/opendmarc.tmp"

	// Comparing checksums to see if the file is different
	if !helper.FileExists(configPath) {
		ioc.Log.Info("OpenDMARC configuration file was not found, writing new one")
		if err := helper.WriteToFile(content, configPath); err != nil {
			return nil
		}

		if err := svc.Restart(); err != nil {
			return err
		}
	} else {
		ioc.Log.Info("OpenDMARC configuration file was found, making sure it is up to date")
		existingChecksum, err := helper.Checksum(configPath)
		if err != nil {
			return err
		}

		if err := helper.WriteToFile(content, tempPath); err != nil {
			return err
		}

		currentChecksum, err := helper.Checksum(tempPath)
		if err != nil {
			return err
		}

		ioc.Log.Info("currentCheckSum: %v  |  existingChecksum: %v", currentChecksum, existingChecksum)
		if currentChecksum != existingChecksum {
			ioc.Log.Info("OpenDMARC configuration files differ, updating configuration")
			if err := helper.WriteToFile(content, configPath); err != nil {
				return err
			}

			if err := svc.Restart(); err != nil {
				return err
			}
		} else {
			ioc.Log.Info("OpenDMARC configuration file is up to date")
		}
	}

	return nil
}
