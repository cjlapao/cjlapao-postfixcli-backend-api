package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/version"
	"github.com/cjlapao/postfixcli-backend-api/ioc"
	"github.com/cjlapao/postfixcli-backend-api/services"
)

var ver = "0.0.1"
var iocServices = execution_context.Get().Services

func main() {
	SetVersion()
	getVersion := helper.GetFlagSwitch("version", false)
	if getVersion {
		format := helper.GetFlagValue("o", "json")
		switch strings.ToLower(format) {
		case "json":
			fmt.Println(iocServices.Version.PrintVersion(int(version.JSON)))
		case "yaml":
			fmt.Println(iocServices.Version.PrintVersion(int(version.JSON)))
		default:
			fmt.Println("Please choose a valid format, this can be either json or yaml")
		}
		os.Exit(0)
	}

	iocServices.Version.PrintAnsiHeader()

	configFile := helper.GetFlagValue("config", "")
	if configFile != "" {
		iocServices.Logger.Command("Loading configuration from " + configFile)
		iocServices.Configuration.LoadFromFile(configFile)
	}

	defer func() {
	}()

	Init()
}

func Init() {
	mode := helper.GetFlagValue("mode", "none")
	if mode == "k8s" {
		ioc.Log.Info("Kubernetes Cluster Test Mode")
		client := services.GetK8sService()
		client.GetClusterIps()
		r, _ := client.GetIngressIp("istio-ingress", "istio-system")
		ioc.Log.Info("%v -> %v", r[0].Hostname, r[0].Ip)

		linode := services.GetLinodeService()
		_, err := linode.GetNodeBalancerDetails("185.3.92.171")
		if err != nil {
			ioc.Log.Exception(err, "Error")
		}
		linode.GetNodeInstances()
	}
}

func SetVersion() {
	iocServices.Version.Name = "Postfix Client Backend API Service"
	iocServices.Version.Author = "Carlos Lapao"
	iocServices.Version.License = "MIT"
	strVer, err := version.FromString(ver)
	if err == nil {
		iocServices.Version.Major = strVer.Major
		iocServices.Version.Minor = strVer.Minor
		iocServices.Version.Build = strVer.Build
		iocServices.Version.Rev = strVer.Rev
	}
}
