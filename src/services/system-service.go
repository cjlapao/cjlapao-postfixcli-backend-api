package services

import (
	"context"
	"io/fs"

	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/helper/linux_helper"
	"github.com/cjlapao/common-go/helper/linux_user"
	"github.com/cjlapao/common-go/helper/linux_user_group"
	"github.com/cjlapao/postfixcli-backend-api/ioc"
)

const (
	VirtualHostFolderName      string = "vhosts"
	VirtualHostDefaultMailPath string = "/var/mail"
	VirtualHostsGroupName      string = "vmail"
	VirtualHostsGroupId        int    = 5000
	VirtualHostsUserName       string = "vmail"
	VirtualHostsUserId         int    = 5000
)

var globalSystemService *SystemService

type SystemService struct {
	Context context.Context
}

func GetSystemService() *SystemService {
	if globalSystemService != nil {
		return globalSystemService
	}

	return NewSystemService()
}

func NewSystemService() *SystemService {
	if helper.GetOperatingSystem() != helper.LinuxOs {
		ioc.Log.Fatal("This service only works in Linux, exiting")
	}

	globalSystemService = &SystemService{}

	globalSystemService.Context = context.Background()

	return globalSystemService
}
func (svc *SystemService) SetupDefaultVirtualMailFolder() error {
	return svc.SetupVirtualMailFolder(VirtualHostDefaultMailPath)
}

func (svc *SystemService) SetupVirtualMailFolder(basePath string) error {
	ioc.Log.Info("Checking folders")
	basePath = helper.ToOsPath(basePath)
	if !helper.DirectoryExists(basePath) {
		helper.CreateDirectory(basePath, fs.ModePerm)
		ioc.Log.Info("Folder %v was created", basePath)
	} else {
		ioc.Log.Info("Folder %v already exists, skipping", basePath)
	}

	fullPath := helper.JoinPath(basePath, VirtualHostFolderName)
	if !helper.DirectoryExists(fullPath) {
		helper.CreateDirectory(fullPath, fs.ModePerm)
		ioc.Log.Info("Folder %v was created", fullPath)
	} else {
		ioc.Log.Info("Folder %v already exists, skipping", fullPath)
	}

	ioc.Log.Info("Checking user and groups")

	groupExists := linux_user_group.Exists(VirtualHostsGroupName)

	if !groupExists {
		err := linux_user_group.Create(VirtualHostsGroupName, VirtualHostsGroupId)
		if err != nil {
			return err
		}
		ioc.Log.Info("Group %v was created successfully.", VirtualHostsGroupName)
	} else {
		ioc.Log.Info("Group %v already exists, skipping", VirtualHostsGroupName)
	}

	userExists := linux_user.Exists(VirtualHostsUserName)
	if !userExists {
		err := linux_user.Create(VirtualHostsUserName, VirtualHostsUserId, linux_user.UserGroupCreateOption(VirtualHostsGroupName), linux_user.UserHomeDirectoryCreateOption(fullPath))
		if err != nil {
			return err
		}
		ioc.Log.Info("User %v was created in group %v successfully.", VirtualHostsGroupName, VirtualHostsGroupName)
	} else {
		ioc.Log.Info("User %v already exists, skipping", VirtualHostsGroupName)
	}

	err := linux_helper.ChangeOwner(basePath, VirtualHostsUserName, VirtualHostsGroupName, true)
	if err != nil {
		return err
	}

	ioc.Log.Info("Changed the owner of %v to %v", basePath, VirtualHostsUserName)
	return nil
}
