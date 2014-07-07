package sinaip

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
)

type IP struct {
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	County   string `json:"county"`
	ISP      string `json:"isp"`
	IP       string `json:"ip"`
	Start    string `json:"start"`
	End      string `json:"end"`
}

func (ip *IP) Json() []byte {
	js, _ := json.Marshal(ip)
	return js
}

func (ip *IP) Long() uint32 {
	iplong, _ := IPToLong(ip.IP)
	return iplong
}

func IPToLong(ipstr string) (uint32, error) {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0, errors.New("Invalid IP")
	}
	ip = ip.To4()
	if ip == nil {
		return 0, errors.New("Not IPv4")
	}
	return binary.BigEndian.Uint32(ip), nil
}

func LongToIP(iplong uint32) (string, error) {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, iplong)
	return ip.String(), nil
}
