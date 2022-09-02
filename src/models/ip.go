package models

type IpType int

const (
	IPV4 IpType = 1
	IPV6 IpType = 2
)

type Ip struct {
	Ip   string
	Type IpType
}
