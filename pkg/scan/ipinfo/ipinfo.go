package ipinfo

import (
	log2 "BeeScan-scan/pkg/log"
	"embed"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"strings"
	"time"
)

/*
创建人员：云深不知处
创建时间：2022/1/6
程序功能：获取ip详细信息
*/

const (
	INDEX_BLOCK_LENGTH  = 12
	TOTAL_HEADER_LENGTH = 8192
)

var err error
var ipInfo IpInfo

type Ip2Region struct {
	// db file handler
	dbFileHandler embed.FS

	//header block info

	headerSip []int64
	headerPtr []int64
	headerLen int64

	// super block index info
	firstIndexPtr int64
	lastIndexPtr  int64
	totalBlocks   int64

	// for memory mode only
	// the original db binary string

	dbBinStr []byte
	dbFile   string
}

type IpInfo struct {
	CityId   int64
	Country  string
	Region   string
	Province string
	City     string
	ISP      string
}

func (ip IpInfo) String() string {
	return strconv.FormatInt(ip.CityId, 10) + "|" + ip.Country + "|" + ip.Region + "|" + ip.Province + "|" + ip.City + "|" + ip.ISP
}

func getIpInfo(cityId int64, line []byte) IpInfo {

	lineSlice := strings.Split(string(line), "|")
	ipInfo := IpInfo{}
	length := len(lineSlice)
	ipInfo.CityId = cityId
	if length < 5 {
		for i := 0; i <= 5-length; i++ {
			lineSlice = append(lineSlice, "")
		}
	}

	ipInfo.Country = lineSlice[0]
	ipInfo.Region = lineSlice[1]
	ipInfo.Province = lineSlice[2]
	ipInfo.City = lineSlice[3]
	ipInfo.ISP = lineSlice[4]
	return ipInfo
}

func IpInfoInit(f embed.FS) *Ip2Region {
	ip2region, err := New(f)
	if err != nil {
		log2.Error(err)
		fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[IpInfoInit]:", err)
	}
	ip2region.dbBinStr, err = ip2region.dbFileHandler.ReadFile("ip2region.db")
	if err != nil {
		log2.Error(err)
		fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[IpInfoInit]:", err)
	}
	return ip2region
}

func New(f embed.FS) (*Ip2Region, error) {

	return &Ip2Region{
		dbFileHandler: f,
	}, nil
}

func (this *Ip2Region) MemorySearch(ipStr string) (ipInfo IpInfo, err error) {
	ipInfo = IpInfo{}

	if this.totalBlocks == 0 {

		this.firstIndexPtr = getLong(this.dbBinStr, 0)
		this.lastIndexPtr = getLong(this.dbBinStr, 4)
		this.totalBlocks = (this.lastIndexPtr-this.firstIndexPtr)/INDEX_BLOCK_LENGTH + 1
	}

	ip, err := ip2long(ipStr)
	if err != nil {
		return ipInfo, err
	}

	h := this.totalBlocks
	var dataPtr, l int64
	for l <= h {

		m := (l + h) >> 1
		p := this.firstIndexPtr + m*INDEX_BLOCK_LENGTH
		sip := getLong(this.dbBinStr, p)
		if ip < sip {
			h = m - 1
		} else {
			eip := getLong(this.dbBinStr, p+4)
			if ip > eip {
				l = m + 1
			} else {
				dataPtr = getLong(this.dbBinStr, p+8)
				break
			}
		}
	}
	if dataPtr == 0 {
		return ipInfo, errors.New("not found")
	}

	dataLen := ((dataPtr >> 24) & 0xFF)
	dataPtr = (dataPtr & 0x00FFFFFF)
	ipInfo = getIpInfo(getLong(this.dbBinStr, dataPtr), this.dbBinStr[(dataPtr)+4:dataPtr+dataLen])
	return ipInfo, nil

}

func getLong(b []byte, offset int64) int64 {

	val := (int64(b[offset]) |
		int64(b[offset+1])<<8 |
		int64(b[offset+2])<<16 |
		int64(b[offset+3])<<24)

	return val

}

func ip2long(IpStr string) (int64, error) {
	bits := strings.Split(IpStr, ".")
	if len(bits) != 4 {
		return 0, errors.New("ip format error")
	}

	var sum int64
	for i, n := range bits {
		bit, _ := strconv.ParseInt(n, 10, 64)
		sum += bit << uint(24-8*i)
	}

	return sum, nil
}

// GetIpinfo 获取IP详细信息
func GetIpinfo(r *Ip2Region, ip string) (IpInfo, error) {
	log2.Info("[GetIPInfo]:", ip)
	fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[GetIPInfo]:", ip)
	ipinfo, err := r.MemorySearch(ip)
	return ipinfo, err
}
