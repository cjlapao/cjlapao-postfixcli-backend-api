package services

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/cjlapao/common-go/guard"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/helper/linux_service"
	"github.com/cjlapao/postfixcli-backend-api/ioc"
	"github.com/cjlapao/postfixcli-backend-api/models"
)

type MailService interface {
	Name() string
	Restart() error
	Start() error
	Stop() error
	Status() linux_service.LinuxServiceState
	Init() error
	Configure(config models.MailServerConfig) error
}

func applyConfiguration(service MailService, config models.MailServerConfig, templateFile models.ConfigFile) error {
	if err := guard.EmptyOrNil(config); err != nil {
		return err
	}

	if err := guard.EmptyOrNil(templateFile); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(templateFile.DestinationDir); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(templateFile.FileName); err != nil {
		return err
	}
	if err := guard.EmptyOrNil(templateFile.TemplateName); err != nil {
		return err
	}

	baseTemplateDir := "./templates"
	filePath := helper.ToOsPath(helper.JoinPath(baseTemplateDir, templateFile.TemplateName))

	ioc.Log.Info("Applying template %v to %v for %v service", templateFile.TemplateName, templateFile.FileName, service.Name())
	templateContent, err := helper.ReadFromFile(filePath)
	if err != nil {
		return err
	}

	tmpl, err := template.New(service.Name()).Parse(string(templateContent))
	if err != nil {
		return err
	}

	var configResult bytes.Buffer

	if err = tmpl.Execute(&configResult, config); err != nil {
		return err
	}

	if err := saveConfig(service, templateFile.FileName, templateFile.DestinationDir, configResult.String()); err != nil {
		return err
	}

	return nil
}

func saveConfig(service MailService, filename string, destinationFolder string, content string) error {
	serviceName := service.Name()
	configPath := helper.ToOsPath(helper.JoinPath(destinationFolder, filename))
	tempPath := fmt.Sprintf("%v.tmp", filename)

	// Comparing checksums to see if the file is different
	if !helper.FileExists(configPath) {
		ioc.Log.Info("%v service configuration file %v was not found, writing new one", serviceName, filename)
		if err := helper.WriteToFile(content, configPath); err != nil {
			return nil
		}

		if err := service.Restart(); err != nil {
			return err
		}
	} else {
		ioc.Log.Info("%v service configuration file %v was found, making sure it is up to date", serviceName, filename)
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

		if err := helper.DeleteFile(tempPath); err != nil {
			ioc.Log.Exception(err, "error deleting template temporary file %v", tempPath)
		}

		if currentChecksum != existingChecksum {
			ioc.Log.Info("%v service configuration file %v differ, updating configuration", serviceName, filename)
			if err := helper.WriteToFile(content, configPath); err != nil {
				return err
			}

			if err := service.Restart(); err != nil {
				return err
			}
		} else {
			ioc.Log.Info("%v service configuration file %v is up to date", serviceName, filename)
		}
	}

	return nil
}
