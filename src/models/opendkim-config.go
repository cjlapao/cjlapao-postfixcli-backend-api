package models

type OpenDKIMConfig struct {
	FQDN      string
	KeyName   string
	ServerIps []ServiceIp
}
