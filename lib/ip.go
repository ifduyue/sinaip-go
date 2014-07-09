package sinaip

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
)

// IP struct, including fields from SINA IP API.
// All fields, except county is from district, mapped to SINA IP API.
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

// JSON marshals IP struct to JSON
func (ip *IP) JSON() []byte {
	js, _ := json.Marshal(ip)
	return js
}

// Long returns the corresponding uint32 value of an IP
func (ip *IP) Long() uint32 {
	iplong, _ := IPToLong(ip.IP)
	return iplong
}

// IPToLong converts a dot notation IP string to uint32
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

// LongToIP converts an IP in uint32 form to dot notation string.
func LongToIP(iplong uint32) (string, error) {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, iplong)
	return ip.String(), nil
}
