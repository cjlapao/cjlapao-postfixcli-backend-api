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
}

type MailServerLoadBalancer struct {
	Hostname string
	IPV4     string
	IPV6     string
}
