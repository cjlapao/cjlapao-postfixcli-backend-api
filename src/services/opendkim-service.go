package services

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/cjlapao/common-go/commands"
	"github.com/cjlapao/common-go/guard"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/helper/linux_helper"
	"github.com/cjlapao/common-go/helper/linux_service"
	"github.com/cjlapao/postfixcli-backend-api/ioc"
	"github.com/cjlapao/postfixcli-backend-api/models"
)

var globalOpenDKIMService *OpenDKIMService

const (
	OpenDKIMUserName  string = "opendkim"
	OpenDKIMGroupName string = "opendkim"
)

type OpenDKIMService struct {
	Context          context.Context
	ConfigFilePath   string
	ConfigFolderPath string
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
		ConfigFilePath:   "/etc/opendkim.conf",
		ConfigFolderPath: "/etc/opendkim",
	}

	globalOpenDKIMService.Context = context.Background()

	return globalOpenDKIMService
}

func (svc *OpenDKIMService) Name() string {
	return "OpenDKIM"
}

func (svc *OpenDKIMService) Start() error {
	return linux_service.Start("opendkim")
}

func (svc *OpenDKIMService) Stop() error {
	return linux_service.Stop("opendkim")
}

func (svc *OpenDKIMService) Restart() error {
	return linux_service.Restart("opendkim")
}

func (svc *OpenDKIMService) Status() linux_service.LinuxServiceState {
	return linux_service.Status("opendkim")
}

func (svc *OpenDKIMService) ConfigKeyFolderPath() string {
	return helper.JoinPath(svc.ConfigFolderPath, "keys")
}

func (svc *OpenDKIMService) SigningTableFilePath() string {
	return helper.JoinPath(svc.ConfigFolderPath, "signing.table")
}

func (svc *OpenDKIMService) KeysTableFilePath() string {
	return helper.JoinPath(svc.ConfigFolderPath, "key.table")
}

func (svc *OpenDKIMService) TrustedHostsFilePath() string {
	return helper.JoinPath(svc.ConfigFolderPath, "trusted.hosts")
}

func (svc *OpenDKIMService) Init() error {
	ioc.Log.Info("Starting initialization of OpenDKIM")
	if svc.Status() != linux_service.LinuxServiceRunning {
		if err := svc.Stop(); err != nil {
			return err
		}
	}

	// Changing the config file mode
	if !helper.FileExists(svc.ConfigFilePath) {
		linux_helper.ChangeFileMode(svc.ConfigFilePath, "u=rw,go=r", false)
	}

	// Setting new config files variables
	// Creating the config folders
	if !helper.DirectoryExists(svc.ConfigFolderPath) {
		if !helper.CreateDirectory(svc.ConfigKeyFolderPath(), fs.ModePerm) {
			return fmt.Errorf("there was an error creating the folder %v", svc.ConfigFolderPath)
		} else {
			ioc.Log.Info("Created OpenDKIM default config folder")
			// Creating the config keys folder
			if !helper.CreateDirectory(svc.ConfigKeyFolderPath(), fs.ModePerm) {
				return fmt.Errorf("there was an error creating the folder %v", svc.ConfigFolderPath)
			} else {
				ioc.Log.Info("Created OpenDKIM default config keys folder")
			}
		}
	}

	// Changing config folder owner
	err := linux_helper.ChangeOwner(svc.ConfigFolderPath, OpenDKIMUserName, OpenDKIMGroupName, true)

	if err != nil {
		return err
	}
	ioc.Log.Info("Updated user permissions in the OpenDKIM default config folder")

	// Changing config keys folder mode
	err = linux_helper.ChangeFileMode(svc.ConfigFolderPath, "go-rw", true)

	if err != nil {
		return err
	}

	ioc.Log.Info("Updated folder mode in the OpenDKIM default config keys folder")

	// Creating default configuration files
	if !helper.FileExists(svc.SigningTableFilePath()) {
		_, err = commands.Execute("touch", svc.SigningTableFilePath())
		if err != nil {
			return err
		} else {
			ioc.Log.Info("Created OpenDKIM signing table configuration file")
		}
	}

	if !helper.FileExists(svc.KeysTableFilePath()) {
		_, err = commands.Execute("touch", svc.KeysTableFilePath())
		if err != nil {
			return err
		} else {
			ioc.Log.Info("Created OpenDKIM keys table configuration file")
		}
	}

  if !helper.FileExists(svc.TrustedHostsFilePath()) {
		_, err = commands.Execute("touch", svc.TrustedHostsFilePath())
		if err != nil {
			return err
		} else {
			ioc.Log.Info("Created OpenDKIM trusted host configuration file")
		}
	}

  hostname, err := commands.Execute("hostname", "-s")
  
  svc.TrustHost("127.0.0.1")
  svc.TrustHost("::1")
  svc.TrustHost("localhost")
  svc.TrustHost(hostname)

	return nil
}

func (svc *OpenDKIMService) TrustHost(host string) error {
  rawFileContent, err := helper.ReadFromFile(svc.TrustedHostsFilePath())
  if err != nil {
    return err
  }

  fileContent = 

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
