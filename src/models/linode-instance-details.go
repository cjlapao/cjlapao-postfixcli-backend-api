package models

type LinodeInstanceDetail struct {
	ID   int64
	Name string
	Ips  []Ip
}
