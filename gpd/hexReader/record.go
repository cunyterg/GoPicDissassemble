package hexReader

type RecordType uint8

const (
	RT_Data                     RecordType = iota
	RT_EOF
	RT_ExtendedSegmentedAddress
	RT_StartSegmentedAddress
	RT_ExtendedLinearAddress
	RT_StartLinearAddress
)

type Record struct {
	ByteCount uint8
	Address   uint
	Type      RecordType
	Data      []uint16
	CheckSum  uint8
}
