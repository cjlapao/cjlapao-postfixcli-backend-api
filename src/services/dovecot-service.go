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
	DovecotUserName  string = "postfix"
	DovecotGroupName string = "postfix"
)

var globalDovecotService *DovecotService

type DovecotService struct {
	Context context.Context
}

func GetDovecotService() *DovecotService {
	if globalDovecotService != nil {
		return globalDovecotService
	}

	return NewDovecotService()
}

func NewDovecotService() *DovecotService {
	if helper.GetOperatingSystem() != helper.LinuxOs {
		ioc.Log.Fatal("This service only works in Linux, exiting")
	}

	globalDovecotService = &DovecotService{}

	globalDovecotService.Context = context.Background()

	return globalDovecotService
}

func (svc *DovecotService) Name() string {
	return "Dovecot"
}

func (svc *DovecotService) Start() error {
	return linux_service.Start("dovecot")
}

func (svc *DovecotService) Stop() error {
	return linux_service.Stop("dovecot")
}

func (svc *DovecotService) Restart() error {
	return linux_service.Restart("dovecot")
}

func (svc *DovecotService) Status() linux_service.LinuxServiceState {
	return linux_service.Status("dovecot")
}

func (svc *DovecotService) Init() error {
	ioc.Log.Info("Starting initialization of Dovecot")
	if svc.Status() != linux_service.LinuxServiceRunning {
		if err := svc.Stop(); err != nil {
			return err
		}
	}

	return nil
}

func (svc *DovecotService) Configure(config models.MailServerConfig) error {
	if err := guard.EmptyOrNil(config); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(config.Domain); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(config.SubDomain); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(config.Sql); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(config.Sql.ServerName); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(config.Sql.DatabaseName); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(config.Sql.Username); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(config.Sql.Password); err != nil {
		return err
	}

	dovecotConf := models.ConfigFile{
		FileName:       "dovecot.conf",
		DestinationDir: "/etc/dovecot",
		TemplateName:   "dovecot.conf.tpl",
	}
	if err := applyConfiguration(svc, config, dovecotConf); err != nil {
		return err
	}

	authConf := models.ConfigFile{
		FileName:       "10-auth.conf",
		DestinationDir: "/etc/dovecot/conf.d",
		TemplateName:   "dovecot_conf.d_10-auth.conf.tpl",
	}
	if err := applyConfiguration(svc, config, authConf); err != nil {
		return err
	}

	mailConf := models.ConfigFile{
		FileName:       "10-mail.conf",
		DestinationDir: "/etc/dovecot/conf.d",
		TemplateName:   "dovecot_conf.d_10-mail.conf.tpl",
	}
	if err := applyConfiguration(svc, config, mailConf); err != nil {
		return err
	}

	masterConf := models.ConfigFile{
		FileName:       "10-master.conf",
		DestinationDir: "/etc/dovecot/conf.d",
		TemplateName:   "dovecot_conf.d_10-master.conf.tpl",
	}
	if err := applyConfiguration(svc, config, masterConf); err != nil {
		return err
	}

	sslConf := models.ConfigFile{
		FileName:       "10-ssl.conf",
		DestinationDir: "/etc/dovecot/conf.d",
		TemplateName:   "dovecot_conf.d_10-ssl.conf.tpl",
	}
	if err := applyConfiguration(svc, config, sslConf); err != nil {
		return err
	}

	authSqlConf := models.ConfigFile{
		FileName:       "auth-sql.conf.ext",
		DestinationDir: "/etc/dovecot/conf.d",
		TemplateName:   "dovecot_conf.d_auth-sql.conf.ext.tpl",
	}
	if err := applyConfiguration(svc, config, authSqlConf); err != nil {
		return err
	}

	sqlConf := models.ConfigFile{
		FileName:       "sql.conf.ext",
		DestinationDir: "/etc/dovecot/conf.d",
		TemplateName:   "dovecot_conf.d_sql.conf.ext.tpl",
	}
	if err := applyConfiguration(svc, config, sqlConf); err != nil {
		return err
	}

	return nil
}
