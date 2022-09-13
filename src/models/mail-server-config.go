package models

type ConfigFile struct {
	FileName       string
	DestinationDir string
	TemplateName   string
}

type MailServerConfig struct {
	Domain       string
	SubDomain    string
	LoadBalancer MailServerLoadBalancer
	Sql          SqlServer
}

type MailServerLoadBalancer struct {
	Hostname string
	IPV4     string
	IPV6     string
}

type SqlServer struct {
	ServerName   string
	DatabaseName string
	Username     string
	Password     string
}
