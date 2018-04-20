package misc

import "strconv"

const (
	l, r = 0, 1
)

func ParseInt64(v string) int64 {
	if len(v) >= 2 && v[:2] == "0x" {
		v = v[2:]
	}
	i, _ := strconv.ParseInt(v, 16, 32)
	return i
}

func ParseInt(v string) int {
	return int(ParseInt64(v))
}

func ParseUInt(v string) uint {

	return uint(ParseInt64(v))
}

func ConvertAddress(v1 []int) (int) {
	return v1[l]<<8 + v1[r]
}

func Convert1u(v1 []uint16) (uint) {
	return uint(v1[l])<<8 + uint(v1[r])
}

func ConvertWordAddress(v1 []uint16) (int) {
	return int(v1[1])<<8 + int(v1[0])
}

func ToBytes(sb string) []int {
	bytesCount := len(sb) / 2
	b := make([]int, bytesCount)
	for i := range b {
		si := i * 2
		b[i] = ParseInt(sb[si : si+2])
	}
	return b
}
