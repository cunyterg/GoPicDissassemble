package hexReader

import (
	"io"
	"bufio"
	"strings"
	"fmt"
	"github.com/qiniu/log"
	"github.com/lightAssemble/GoPicDissassemble/gpd/misc"
)

type Hex32Reader struct {
	records       []*Record
	addressOffset uint
}

func New() *Hex32Reader {
	reader := new(Hex32Reader)
	reader.records = make([]*Record, 0)
	reader.addressOffset = 0
	return reader
}

func (h *Hex32Reader) GetRecordsFrom(stream io.Reader) []*Record {
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line[0] != ':' { // Intel form
			fmt.Println("Ignoring  hex line > " + line)
			continue
		}
		record := h.decodeLine(line[1:])
		if record != nil {
			h.records = append(h.records, record)
		}
	}
	return h.records

}

func (h *Hex32Reader) decodeLine(line string) *Record {
	record := h.tokenize(line)
	if record == nil {
		return nil
	}
	switch record.Type {
	case RT_ExtendedSegmentedAddress:
		h.addressOffset = misc.Convert1u(record.Data[0:2]) << 4
		//log.Println("Set ExtendedSegmentedAddress : " + string(h.addressOffset)) // FIXMe: string
	case RT_ExtendedLinearAddress:
		h.addressOffset = misc.Convert1u(record.Data[0:2]) << 16
		//log.Println("Set ExtendedLinearAddress : " + string(h.addressOffset))
	case RT_EOF:
		return record
	case RT_Data:
		record.Address += h.addressOffset
		if record.Address > 0x10000 {
			log.Warn(fmt.Sprintf("Data at record `0x%X : %v` wraps over 0xFFFF.", record.Address, line))
		}
		//log.Println(fmt.Sprintf("Write data record on [0x%X:0x%X] : size %v",
		//	record.Address, record.Address+uint(record.ByteCount), record.ByteCount))
		return record

	}

	return nil
}

func (h *Hex32Reader) tokenize(line string) *Record {

	lineBytes := misc.ToBytes(line)
	byteCount := uint8(lineBytes[0])
	address := uint(misc.ConvertAddress(lineBytes[1:3]))
	rType := RecordType(lineBytes[3])
	checkSum := uint8(lineBytes[len(lineBytes)-1])
	data := make([]uint16, byteCount)
	const dataOffset = 4
	for w := uint8(0); w < byteCount; w++ {
		data[w] = uint16(lineBytes[dataOffset+w])
	}
	c := 0
	for _, v := range lineBytes {
		c += v
	}
	c &= 0xff
	if c != 0x0 {
		log.Warn("Record checksum error > " + line)
	}
	return &Record{Data: data, Type: rType, Address: address,
		CheckSum: checkSum, ByteCount: byteCount}
}
