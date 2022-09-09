package services

import (
	"context"

	"github.com/cjlapao/common-go/guard"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/helper/linux_service"
	"github.com/cjlapao/postfixcli-backend-api/ioc"
	"github.com/cjlapao/postfixcli-backend-api/models"
)

const (
	PostfixUserName  string = "postfix"
	PostfixGroupName string = "postfix"
)

var globalPostfixService *PostfixService

type PostfixService struct {
	Context context.Context
}

func GetPostfixService() *PostfixService {
	if globalPostfixService != nil {
		return globalPostfixService
	}

	return NewPostfixService()
}

func NewPostfixService() *PostfixService {
	if helper.GetOperatingSystem() != helper.LinuxOs {
		ioc.Log.Fatal("This service only works in Linux, exiting")
	}

	globalPostfixService = &PostfixService{}

	globalPostfixService.Context = context.Background()

	return globalPostfixService
}

func (svc *PostfixService) Name() string {
	return "Postfix"
}

func (svc *PostfixService) Start() error {
	return linux_service.Start("postfix")
}

func (svc *PostfixService) Stop() error {
	return linux_service.Stop("postfix")
}

func (svc *PostfixService) Restart() error {
	return linux_service.Restart("postfix")
}

func (svc *PostfixService) Status() linux_service.LinuxServiceState {
	return linux_service.Status("postfix")
}

func (svc *PostfixService) Init() error {
	ioc.Log.Info("Starting initialization of postfix")
	if svc.Status() != linux_service.LinuxServiceRunning {
		if err := svc.Stop(); err != nil {
			return err
		}
	}

	return nil
}

func (svc *PostfixService) Configure(config models.MailServerConfig) error {
	if err := guard.EmptyOrNil(config); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(config.Domain); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(config.SubDomain); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(config.LoadBalancer); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(config.LoadBalancer.Hostname); err != nil {
		return err
	}

	headerChecks := models.ConfigFile{
		FileName:       "header_checks",
		DestinationDir: "/etc/postfix",
		TemplateName:   "postfix_header_checks.tpl",
	}
	if err := applyConfiguration(svc, config, headerChecks); err != nil {
		return err
	}

	mainConfig := models.ConfigFile{
		FileName:       "main.cf",
		DestinationDir: "/etc/postfix",
		TemplateName:   "postfix_main.cf.tpl",
	}
	if err := applyConfiguration(svc, config, mainConfig); err != nil {
		return err
	}

	masterConfig := models.ConfigFile{
		FileName:       "master.cf",
		DestinationDir: "/etc/postfix",
		TemplateName:   "postfix_master.cf.tpl",
	}

	if err := applyConfiguration(svc, config, masterConfig); err != nil {
		return err
	}
	return nil
}
