package models

type NodeIp struct {
	Hostname   string
	ExternalIp Ip
	InternalIp Ip
}
