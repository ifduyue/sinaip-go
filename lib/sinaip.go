package sinaip

import (
	"encoding/binary"
	"os"
	"syscall"
)

type SINAIP struct {
	Path        string
	Dat         []byte
	Size        uint32
	DataOffset  uint32
	IndexOffset uint32
	Count       int
}

func NewSINAIP(path string) (*SINAIP, error) {
	sinaip := &SINAIP{Path: path}

	file, err := os.Open(sinaip.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	dattmp, err := syscall.Mmap(int(file.Fd()), 0, int(fi.Size()), syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		return nil, err
	}

	sinaip.Path = path
	sinaip.Dat = dattmp
	sinaip.Size = uint32(fi.Size())
	sinaip.DataOffset = binary.LittleEndian.Uint32(sinaip.Dat[0:4])
	sinaip.IndexOffset = binary.LittleEndian.Uint32(sinaip.Dat[4:8])
	sinaip.Count = int((sinaip.Size - sinaip.IndexOffset) / 4)

	return sinaip, nil
}

func (sinaip *SINAIP) gettext(offset uint32) string {
	i := offset
	for sinaip.Dat[i] != 0 {
		i++
	}
	return string(sinaip.Dat[offset:i])
}

func (sinaip *SINAIP) Query(ipstr string) (*IP, error) {
	var (
		ip         uint32
		err        error
		ip_start   uint32
		ip_end     uint32
		ip_start_s string
		ip_end_s   string
		start      int
		end        int
		middle     int
		offset     uint32
		result     *IP
	)
	ip, err = IPToLong(ipstr)
	if err != nil {
		return nil, err
	}

	start = 0
	end = sinaip.Count
	for {
		middle = (end-start)/2 + start
		if start == middle {
			offset = sinaip.IndexOffset + uint32(middle*4)
			ip_start = binary.LittleEndian.Uint32(sinaip.Dat[offset : offset+4])
			if offset+8 <= sinaip.Size {
				ip_end = binary.LittleEndian.Uint32(sinaip.Dat[offset+4 : offset+8])
				ip_end -= 1
			} else {
				ip_end = 4294967295 // 255.255.255.255
			}
			ip_start_s, _ = LongToIP(ip_start)
			ip_end_s, _ = LongToIP(ip_end)
			break
		}
		offset = sinaip.IndexOffset + uint32(middle*4)
		ip_start = binary.LittleEndian.Uint32(sinaip.Dat[offset : offset+4])
		if ip < ip_start {
			end = middle
		} else {
			start = middle
		}
	}

	offset = sinaip.DataOffset + uint32(start*16)
	countryOffset := binary.LittleEndian.Uint32(sinaip.Dat[offset : offset+4])
	provinceOffset := binary.LittleEndian.Uint32(sinaip.Dat[offset+4 : offset+8])
	cityOffset := binary.LittleEndian.Uint32(sinaip.Dat[offset+8 : offset+12])
	ispOffset := binary.LittleEndian.Uint32(sinaip.Dat[offset+12 : offset+16])

	result = &IP{}
	result.Country = sinaip.gettext(countryOffset)
	result.Province = sinaip.gettext(provinceOffset)
	result.City = sinaip.gettext(cityOffset)
	result.ISP = sinaip.gettext(ispOffset)
	result.IP = ipstr
	result.Start = ip_start_s
	result.End = ip_end_s

	return result, nil
}
