package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"
)

var (
	port         int
	ipdatPath    string
	ipdat        []byte
	ipdatSize    uint32
	dataOffset   uint32
	indexOffset  uint32
	ipRangeCount int
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

func ip2long(ipstr string) (uint32, error) {
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

func long2ip(iplong uint32) (string, error) {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, iplong)
	return ip.String(), nil
}

func gettext(offset uint32) string {
	i := offset
	for ipdat[i] != 0 {
		i++
	}
	return string(ipdat[offset:i])
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ipstr      string
		ip_start   uint32
		ip_end     uint32
		ip_start_s string
		ip_end_s   string
		start      int
		end        int
		middle     int
		offset     uint32
	)
	ipstr = r.URL.Path[1:]
	ip, err := ip2long(ipstr)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	start = 0
	end = ipRangeCount
	for {
		middle = (end-start)/2 + start
		if start == middle {
			offset = indexOffset + uint32(middle*4)
			ip_start = binary.LittleEndian.Uint32(ipdat[offset : offset+4])
			if offset+8 <= ipdatSize {
				ip_end = binary.LittleEndian.Uint32(ipdat[offset+4 : offset+8])
				ip_end -= 1
			} else {
				ip_end = 4294967295 // 255.255.255.255
			}
			ip_start_s, _ = long2ip(ip_start)
			ip_end_s, _ = long2ip(ip_end)
			break
		}
		offset = indexOffset + uint32(middle*4)
		ip_start = binary.LittleEndian.Uint32(ipdat[offset : offset+4])
		if ip < ip_start {
			end = middle
		} else {
			start = middle
		}
	}

	offset = dataOffset + uint32(start*16)
	countryOffset := binary.LittleEndian.Uint32(ipdat[offset : offset+4])
	provinceOffset := binary.LittleEndian.Uint32(ipdat[offset+4 : offset+8])
	cityOffset := binary.LittleEndian.Uint32(ipdat[offset+8 : offset+12])
	ispOffset := binary.LittleEndian.Uint32(ipdat[offset+12 : offset+16])

	result := IP{}
	result.Country = gettext(countryOffset)
	result.Province = gettext(provinceOffset)
	result.City = gettext(cityOffset)
	result.ISP = gettext(ispOffset)
	result.IP = ipstr
	result.Start = ip_start_s
	result.End = ip_end_s

	json, _ := json.Marshal(result)
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.Write(json)
}

func init() {
	flag.IntVar(&port, "port", 8080, "HTTP Server Port")
	flag.StringVar(&ipdatPath, "ipdat", "ip.dat", "Path to ip.dat")
	flag.Parse()
}

func main() {
	httpAddr := fmt.Sprintf(":%v", port)

	file, err := os.Open(ipdatPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err.Error())
	}
	ipdatSize = uint32(fi.Size())

	ipdat_tmp, err := syscall.Mmap(int(file.Fd()), 0, int(fi.Size()), syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		log.Fatal(err.Error())
	}
	ipdat = ipdat_tmp
	dataOffset = binary.LittleEndian.Uint32(ipdat[0:4])
	indexOffset = binary.LittleEndian.Uint32(ipdat[4:8])
	ipRangeCount = int((uint32(fi.Size()) - indexOffset) / 4)
	log.Printf("%10v", fi.Size())
	log.Printf("%10v", indexOffset)
	log.Printf("%10v", ipRangeCount)
	log.Printf("Listening to %v", httpAddr)
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(httpAddr, nil))
}
