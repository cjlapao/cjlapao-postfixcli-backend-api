package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/version"
	"github.com/cjlapao/postfixcli-backend-api/ioc"
	"github.com/cjlapao/postfixcli-backend-api/models"
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
		// ioc.Log.Info("Kubernetes Cluster Test Mode")
		// client := services.GetK8sService()
		// client.GetClusterIps()
		// r, _ := client.GetIngressIp("istio-ingress", "istio-system")
		// ioc.Log.Info("%v -> %v", r[0].Hostname, r[0].Ip)

		// linode := services.GetLinodeService()
		// _, err := linode.GetNodeBalancerDetails("185.3.92.171")
		// if err != nil {
		// 	ioc.Log.Exception(err, "Error")
		// }
		// linode.GetNodeInstances()

		// azure := services.GetAzureService()
		// azure.GetDnsZone("Infrastructure", "carloslapao.com")
		// azure.GetDnsZones("Infrastructure")
		// dnsr, _ := azure.GetDnsRecord("Infrastructure", "carloslapao.com", "AAAA", "2a01-7e00-f03c-93ff-fecc-6a6f")
		// ioc.Log.Info(dnsr.Name)
		// err = azure.UpsertADnsRecord("Infrastructure", "carloslapao.com", "test1", 50, "200.10.0.1", "200.10.0.2")
		// if err != nil {
		// 	ioc.Log.Exception(err, "error a dns")
		// }
		// err = azure.UpsertAAAADnsRecord("Infrastructure", "carloslapao.com", "test1aaaa", 30, "2a01:7e00::f03c:93ff:fecc:6a6f")
		// if err != nil {
		// 	ioc.Log.Exception(err, "error aaa dns")
		// }
		// err = azure.UpsertMXDnsRecord("Infrastructure", "carloslapao.com", 30, models.AzureDnsMXRecord{
		// 	Exchange:   "mail.carloslapao.com",
		// 	Preference: 10,
		// })
		// if err != nil {
		// 	ioc.Log.Exception(err, "error mx dns")
		// }
		// err = azure.UpsertCNAMEDnsRecord("Infrastructure", "carloslapao.com", "test_cname", 50, "mail.carloslapao.com")
		// if err != nil {
		// 	ioc.Log.Exception(err, "error CNAME dns")
		// }
		// err = azure.UpsertSRVDnsRecord("Infrastructure", "carloslapao.com", "_smtp", 50, models.AzureDnsSRVRecord{
		// 	Port:     25,
		// 	Priority: 10,
		// 	Target:   "mail.carloslapao.com",
		// 	Weight:   10,
		// })
		// if err != nil {
		// 	ioc.Log.Exception(err, "error srv dns")
		// }
		// err = azure.UpsertTXTDnsRecord("Infrastructure", "carloslapao.com", "_small", 50, "Hang on, my kittens are scratching at the bathtub and they'll upset by the lack of biscuits.Hang on, my kittens are scratching at the bathtub and they'll upset by the lack of biscuits.")
		// if err != nil {
		// 	ioc.Log.Exception(err, "error txt small dns")
		// }
		// err = azure.UpsertTXTDnsRecord("Infrastructure", "carloslapao.com", "_big", 50, "Hang on, my kittens are scratching at the bathtub and they'll upset by the lack of biscuits.Hang on, my kittens are scratching at the bathtub and they'll upset by the lack of biscuits.Hang on, my kittens are scratching at the bathtub and they'll upset by the lack of biscuits.")
		// if err != nil {
		// 	ioc.Log.Exception(err, "error txt small dns")
		// }
		sys := services.NewSystemService()
		err := sys.SetupDefaultVirtualMailFolder()
		if err != nil {
			ioc.Log.Exception(err, "error creating virtual mail folder")
		}

		mailConfig := models.MailServerConfig{
			Domain:    "carloslapao.com",
			SubDomain: "mail",
			LoadBalancer: models.MailServerLoadBalancer{
				Hostname: "ip.linode.com",
			},
		}

		postfix := services.GetPostfixService()
		err = postfix.Init()
		if err != nil {
			ioc.Log.Exception(err, "error init postfix")
		}
		err = postfix.Configure(mailConfig)
		if err != nil {
			ioc.Log.Exception(err, "error config postfix")
		}

		opendmarc := services.GetOpenDMARCService()
		err = opendmarc.Init()
		if err != nil {
			ioc.Log.Exception(err, "error init opendmarc")
		}
		err = opendmarc.Configure(mailConfig)
		if err != nil {
			ioc.Log.Exception(err, "error config opendmarc")
		}
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
