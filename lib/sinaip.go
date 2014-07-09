package sinaip

import (
	"encoding/binary"
	"io/ioutil"
	"os"
	"syscall"
)

// SINAIP struct
type SINAIP struct {
	Path        string
	Dat         []byte
	Size        uint32
	DataOffset  uint32
	IndexOffset uint32
	Count       int
	Preload     bool
}

// NewSINAIP returns a new SINAIP struct instance
func NewSINAIP(path string, preload bool) (*SINAIP, error) {
	var (
		dattmp []byte
		err    error
	)
	sinaip := &SINAIP{Path: path, Preload: preload}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if preload {
		dattmp, err = ioutil.ReadFile(path)
	} else {
		dattmp, err = syscall.Mmap(int(file.Fd()), 0, int(fi.Size()), syscall.PROT_READ, syscall.MAP_PRIVATE)
	}
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

// Query returns the location info of an IP
func (sinaip *SINAIP) Query(ipstr string) (*IP, error) {
	var (
		ip         uint32
		err        error
		ipStart    uint32
		ipEnd      uint32
		ipStartStr string
		ipEndStr   string
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
			ipStart = binary.LittleEndian.Uint32(sinaip.Dat[offset : offset+4])
			if offset+8 <= sinaip.Size {
				ipEnd = binary.LittleEndian.Uint32(sinaip.Dat[offset+4:offset+8]) - 1
			} else {
				ipEnd = 4294967295 // 255.255.255.255
			}
			ipStartStr, _ = LongToIP(ipStart)
			ipEndStr, _ = LongToIP(ipEnd)
			break
		}
		offset = sinaip.IndexOffset + uint32(middle*4)
		ipStart = binary.LittleEndian.Uint32(sinaip.Dat[offset : offset+4])
		if ip < ipStart {
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
	result.Start = ipStartStr
	result.End = ipEndStr

	return result, nil
}
