package hexReader_test

import (
	"testing"
	"GoPicDissassemble/gpd/hexReader"
	"strings"
)

func TestHex32Reader_ReadFrom(t *testing.T) {
	source := ":10246200464C5549442050524F46494C4500464C33"
	reader := hexReader.New()
	records := reader.GetRecordsFrom(strings.NewReader(source))
	if len(records) != 1 {
		t.Error("Line not parsed")
	}
	record := records[0]
	if record.ByteCount != 0x10 {
		t.Error("Byte count parse error")
	} else if record.Address != 0x2462 {
		t.Error("Address parse error")
	} else if record.CheckSum != 0x33 {
		t.Error("CheckSum parse error")
	}
}
