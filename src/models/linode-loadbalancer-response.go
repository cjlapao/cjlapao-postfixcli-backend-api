package models

type LinodeLoadBalancerResponse struct {
	Data    []LinodeLoadBalancer `json:"data"`
	Page    int64                `json:"page"`
	Pages   int64                `json:"pages"`
	Results int64                `json:"results"`
}

type LinodeLoadBalancer struct {
	ID                 int64         `json:"id"`
	Label              string        `json:"label"`
	Region             string        `json:"region"`
	Hostname           string        `json:"hostname"`
	Ipv4               string        `json:"ipv4"`
	Ipv6               string        `json:"ipv6"`
	Created            string        `json:"created"`
	Updated            string        `json:"updated"`
	ClientConnThrottle int64         `json:"client_conn_throttle"`
	Tags               []interface{} `json:"tags"`
	Transfer           Transfer      `json:"transfer"`
}

type Transfer struct {
	In    float64 `json:"in"`
	Out   float64 `json:"out"`
	Total float64 `json:"total"`
}
