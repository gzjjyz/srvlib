package utils

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gzjjyz/srvlib/logger"
	"github.com/huandu/go-clone"
	jsoniter "github.com/json-iterator/go"
)

type PointSt struct {
	X, Y int32
}

// 加载json配置
func LoadJson(mgr interface{}, module string) bool {
	file := GetCurrentDir() + "config/" + strings.ToLower(module) + ".json"
	data, err := ioutil.ReadFile(file)
	if err != nil {
		if IsDev() {
			logger.Errorf("未找到配置文件数据:[%s] %s", file, err)
		} else {
			logger.Fatalf("未找到配置文件数据:[%s] %s", file, err)
		}
		return false
	}

	if err := jsoniter.Unmarshal(data, mgr); err != nil {
		logger.Fatalf("load %s Unmarshal json error:%s", file, err)
		return false
	}
	return true
}

var (
	isDev bool
)

func ParseCmdInput() {
	flag.BoolVar(&isDev, "dev", false, "is it dev now?")
	flag.Parse()
}

// IsDev 是否开发版本
func IsDev() bool {
	return isDev
}

// GetMd5 获取md5
func GetMd5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GetCurrentDir 获取当前目录
func GetCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1) + "/"
}

// Int2ip Convert uint to net.IP
func Int2ip(ipnr int32) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// Ip2int Convert net.IP to int64
func Ip2int(ipnr string) int64 {
	ipnr = strings.Split(ipnr, ":")[0]
	bits := strings.Split(ipnr, ".")

	if len(bits) < 4 {
		return 0
	}
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

// IsSetBit64 是否设置指定位
func IsSetBit64(value uint64, bit uint32) bool {
	if bit > 63 {
		return false
	}
	return (value & (1 << bit)) > 0
}

// SetBit64 设置位
func SetBit64(value uint64, bit uint32) uint64 {
	if bit > 63 {
		return value
	}
	return value | (1 << bit)
}

// ClearBit64 清空位
func ClearBit64(value uint64, bit uint32) uint64 {
	if bit > 63 {
		return value
	}
	return value & ^(1 << bit)
}

// IsSetBit 是否设置指定位
func IsSetBit(value uint32, bit uint32) bool {
	if bit > 31 {
		return false
	}
	return (value & (1 << bit)) > 0
}

// SetBit 设置位
func SetBit(value uint32, bit uint32) uint32 {
	if bit > 31 {
		return value
	}
	return value | (1 << bit)
}

// ClearBit 清空位
func ClearBit(value uint32, bit uint32) uint32 {
	if bit > 31 {
		return value
	}
	return value & ^(1 << bit)
}

// High32 64位高32位值
func High32(value uint64) uint32 {
	return uint32((value & 0xFFFFFFFF00000000) >> 32)
}

// Low32 64位低32位值
func Low32(value uint64) uint32 {
	return uint32(value & 0x00000000FFFFFFFF)
}

// Make64 32+32组装64
func Make64(low, high uint32) uint64 {
	return uint64(low) | uint64(high)<<32
}

func High16(value uint32) uint16 {
	return uint16(value >> 16)
}

func Low16(value uint32) uint16 {
	return uint16(value & 0xFFFF)
}

func Make32(low, high uint16) uint32 {
	return uint32(low) | uint32(high)<<16
}

// RemoveSpace 去除空格类字符
func RemoveSpace(str *string) string {
	s := strings.Replace(*str, "\t", "", -1)
	s = strings.Replace(s, "\r\n", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, " ", "", -1)
	return s
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func MinUInt32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func MaxUInt32(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func MinFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func MaxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func MaxFloat64InSlice(args ...float64) float64 {
	max := float64(0)
	for _, v := range args {
		if v > max {
			max = v
		}
	}
	return max
}

func RoundFloat64(f float64) float64 {
	return math.Floor(f + 0.5)
}

// DeepCopy  结构体深拷贝(谨慎调用)
func DeepCopy(i interface{}) interface{} {
	return clone.Clone(i)
}

func GetSrcServerByActorId(actorId uint64) uint32 {
	tmp := High32(actorId)
	return uint32(Low16(tmp))
}

// CalcMillionRate 万分比加成计算
func CalcMillionRate(base, rate uint32) uint32 {
	return (10000 + rate) * base / 10000
}

func CalcMillionRate64(base, rate int64) int64 {
	return (10000 + rate) * base / 10000
}

func CalcMillionRateBoth64(base, rate int64) int64 {
	return (10000 + rate) * base / 10000
}

// CalcMillionRateRevert 万分比加成计算, 越加越小
func CalcMillionRateRevert(base, rate int64) int64 {
	return ((10000 - rate) * base) / 10000
}

// CalcBillionRate 百万分比加成计算
func CalcBillionRate(base, rate uint32) uint32 {
	return (1000000 + rate) * base / 1000000
}

func PanicDebugStack(tag string) {
	if err := recover(); err != nil {
		logger.Errorf("tag[%s] %+v: %s", tag, err, debug.Stack())
	}
}

func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func BindNum(flag bool) uint32 {
	if flag {
		return 1
	}
	return 0
}

func Get1Num(x uint32) uint32 {
	x = (x & 0x55555555) + ((x & 0xaaaaaaaa) >> 1)
	x = (x & 0x33333333) + ((x & 0xcccccccc) >> 2)
	x = (x & 0x0f0f0f0f) + ((x & 0xf0f0f0f0) >> 4)
	x = (x & 0x00ff00ff) + ((x & 0xff00ff00) >> 8)
	x = (x & 0x0000ffff) + ((x & 0xffff0000) >> 16)
	return x
}

func Get1Num64(x uint64) uint32 {
	return Get1Num(uint32(x>>32)) + Get1Num(uint32(x))
}

// 1-10星为1阶
// 11-20星为2阶
func GetStageByLevel(level, step uint32) (stage uint32) {
	if level%step != 0 {
		stage = 1
	}
	stage += level / step
	return
}

// 切片翻转
func SliceReverseUint32(nums []uint32) {
	for i, n := 0, len(nums)-1; i <= n/2; i++ {
		nums[i], nums[n] = nums[n], nums[i]
		n--
	}
}

func GetUint32SliceFromString(str string) []uint32 {
	ret := make([]uint32, 0)
	splitStr := strings.Split(str, ",")
	for _, numStr := range splitStr {
		ret = append(ret, AtoUint32(numStr))
	}
	return ret
}

func Ternary(flag bool, i, j interface{}) interface{} {
	if flag {
		return i
	}
	return j
}

func IsRobot(actorId uint64) bool {
	return IsSetBit64(actorId, 31)
}
