package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/version"
)

var ver = "0.0.0"
var services = execution_context.Get().Services

func main() {
	SetVersion()
	getVersion := helper.GetFlagSwitch("version", false)
	if getVersion {
		format := helper.GetFlagValue("o", "json")
		switch strings.ToLower(format) {
		case "json":
			fmt.Println(services.Version.PrintVersion(int(version.JSON)))
		case "yaml":
			fmt.Println(services.Version.PrintVersion(int(version.JSON)))
		default:
			fmt.Println("Please choose a valid format, this can be either json or yaml")
		}
		os.Exit(0)
	}

	services.Version.PrintAnsiHeader()

	configFile := helper.GetFlagValue("config", "")
	if configFile != "" {
		services.Logger.Command("Loading configuration from " + configFile)
		services.Configuration.LoadFromFile(configFile)
	}

	defer func() {
	}()

	Init()
}

func Init() {
}

func SetVersion() {
	services.Version.Name = "GoLang Template"
	services.Version.Author = "Carlos Lapao"
	services.Version.License = "MIT"
	strVer, err := version.FromString(ver)
	if err == nil {
		services.Version.Major = strVer.Major
		services.Version.Minor = strVer.Minor
		services.Version.Build = strVer.Build
		services.Version.Rev = strVer.Rev
	}
}
