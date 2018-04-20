package gpd

import (
	"io"
	"bufio"
	"strings"
	"fmt"
	"sort"
	"github.com/gosuri/uitable"
	"os"
	"github.com/lightAssemble/GoPicDissassemble/gpd/hexReader"
	"github.com/lightAssemble/GoPicDissassemble/gpd/misc"
)

const nextWordAddrAddValue = 2

type Disassembler struct {
	regMap      map[int]string
	debug       bool
	ListingMode bool
	HexStyle    bool
	TableStyle  bool
	Code        map[uint]*instruction
	EEPROM      map[uint]uint16
	Config      map[uint]uint16
	Processor   *ProcessorInfo
	instrMapper *InstructionMapper
}

func NewDisassembler(debug bool) *Disassembler {
	d := new(Disassembler)
	d.regMap = nil
	d.debug = debug
	d.ListingMode = false
	d.HexStyle = false
	d.Code = make(map[uint]*instruction)
	d.Config = make(map[uint]uint16)
	d.EEPROM = make(map[uint]uint16)
	d.instrMapper = NewInstructionMap()

	return d
}

func (d *Disassembler) ReadObjectCode(stream io.Reader) {
	records := hexReader.New().GetRecordsFrom(stream)
	for _, r := range records {
		data := r.Data
		address := r.Address
		if address < d.Processor.TopOfMemory {
			for len(data) != 0 {
				d.Code[address] = newInstruction(misc.ConvertWordAddress(data[:2]))
				data = data[2:]
				address += 2
			}
		} else {
			dest := d.EEPROM
			if address < d.Processor.TopOfConfig {
				dest = d.Config
			}
			for _, v := range data {
				dest[address] = v
				address++
			}

		}

	}
}

func (d *Disassembler) getSortedAddresses() []uint {
	sortedAddress := make([]uint, len(d.Code))
	index := 0
	for e := range d.Code {
		sortedAddress[index] = e
		index++
	}
	sort.Slice(sortedAddress, func(i, j int) bool { return sortedAddress[i] < sortedAddress[j] })
	return sortedAddress
}

func (d *Disassembler) Assemble() {
	for _, v := range d.getSortedAddresses() {
		d.assembleLine(v)
	}

	for _, v := range d.getSortedAddresses() {
		d.arrangeLine(v)
	}

}
func (d *Disassembler) arrangeLine(addr uint) {
	code := d.Code[addr]
	if _, exist := d.Code[addr-2]; !exist {
		code.PrefixLine += fmt.Sprintf("org 0x%x ", addr)
	}
	if len(code.Calls) > 0 {
		code.Label = d.makeLabel(addr) + ":"
		callsStr := make([]string, len(code.Calls))
		for i, v := range code.Calls {
			callsStr[i] = fmt.Sprintf("0x%X", v)
		}
		code.Comment += "entry from: " + strings.Join(callsStr, ",")
		if !d.ListingMode && len(callsStr) > 1 {
			code.PrefixLine += "\n"
		}
	}

	if len(code.Comment) < 3 {
		code.Comment = ""
	}

}
func (d *Disassembler) get(addr uint) {

}

func (d *Disassembler) assembleLine(addr uint) {
	code := d.Code[addr]
	cBin := code.Bin
	operand := d.instrMapper.matchingOpcode(code.Bin)
	line := make([]string, 0)
	if d.debug {
		fmt.Printf("0x%X %v\n", cBin, operand.AsmFunc)
	}
	for _, symb := range operand.AsmExtra {
		switch symb {
		case 'F':
			q := cBin & 0xff
			addLine := ""
			if ((cBin & 0x100) == 0) && q >= 0x80 {
				reg, ok := d.regMap[q|0xF00]
				if ok {
					addLine = reg
				}
			}
			if addLine == "" {
				addLine = fmt.Sprintf("0x%X", q|0xF00)
			}
			line = append(line, addLine)
		case 'D':
			if cBin&0x200 == 0 {
				line = append(line, "W")
			} else {
				line = append(line, "F")
			}
		case 'B':
			line = append(line, fmt.Sprintf("0x%x", (cBin>>9)&0x7))
		case 'K':
			line = append(line, fmt.Sprintf("0x%X", cBin&0xFF))
		case 'C':
			line = append(line, fmt.Sprintf("0x%X", cBin&0xF))
		case 'N':
			q := uint(cBin & 0xFF)
			dest := addr + nextWordAddrAddValue
			if q < 0x80 {
				dest += q * 2
			} else {
				dest -= (0x100 - q) * 2
			}

			line = append(line, d.makeLabel(dest))
			d.lookUpAddr(dest).CallsAddAddress(addr)
		case 'M':
			q := uint(cBin & 0x7FF)
			dest := addr + nextWordAddrAddValue
			if q < 0x400 {
				dest += q * 2
			} else {
				dest -= (0x800 - q) * 2
			}
			line = append(line, d.makeLabel(dest))
			d.lookUpAddr(dest).CallsAddAddress(addr)
		case 'A':
			if cBin&0x100 == 0 {
				line = append(line, "ACCESS")
			} else {
				line = append(line, "BANKED")
			}
		case 'S':
			if cBin&0x1 == 1 {
				line = append(line, "S")
			}
		case 'Y':
			c2Bin := d.lookUpAddr(addr + nextWordAddrAddValue).Bin
			r := ""
			if c, ok := d.regMap[cBin&0xFFF]; ok {
				r += c
			} else {
				r += fmt.Sprintf("0x%X", cBin&0xFFF)
			}

			r += ", "
			if c, ok := d.regMap[c2Bin&0xFFF]; ok {
				r += c
			} else {
				r += fmt.Sprintf("0x%X", c2Bin&0xFFF)
			}
			line = append(line, r)
		case 'W':
			lookUpAddr := d.lookUpAddr(addr + nextWordAddrAddValue)
			c2Bin := lookUpAddr.Bin
			dest := uint(((cBin & 0xFF) | (c2Bin&0xFFF)<<8) * 2)
			d.lookUpAddr(dest).CallsAddAddress(addr)
			line = append(line, d.makeLabel(dest))
			if ((cBin & 0x300) ^ 0x100) == 0 {
				line = append(line, "FAST")
			}
		case 'Z':
			c2Bin := d.lookUpAddr(addr + nextWordAddrAddValue).Bin
			line = append(line, fmt.Sprintf("%v, 0x%X", (cBin&0x30)>>4,
				((cBin&0xF)<<8)|(c2Bin&0xFF)))
		case 'X':
			line = append(line, fmt.Sprintf("DE 0x%X ;[%04b %04b] WARNING: unknown instruction!", cBin, (cBin>>4)&0xF, cBin&0xF))
		default:
			if symb != ',' {
				line = append(line, string(symb))
			}
		}
	}
	code.Asm = operand.AsmFunc + " " + strings.Join(line, ", ")

}
func (d *Disassembler) makeLabel(address uint) string {
	return strings.Replace(fmt.Sprintf("p%5v", fmt.Sprintf("%X", address)), " ", "_", -1)

}
func (d *Disassembler) lookUpAddr(address uint) *instruction {
	i, ok := d.Code[address]
	if ok {
		return i
	}
	i = newInstruction(0xffff).Dummy()
	d.Code[address] = i
	return i
}

func (d *Disassembler) WriteTo(stream io.Writer) {
	table := uitable.New()
	if d.TableStyle {
		table.Separator = " | "
	}

	writer := bufio.NewWriter(stream)
	writer.WriteString(";Generated by GoPicDissassebmle, RawLight (Yurii Gnevush) 2018\n")
	fmt.Fprintf(writer, "\tLIST %v\n\t#include \"%v\"\n",
		strings.ToUpper(d.Processor.Info["List"]), d.Processor.Info["Include"])

	configDirective := d.Processor.Info["ConfigDirective"]
	for k, c := range d.Processor.Config {
		table.AddRow("", fmt.Sprintf("%v  %v=%v", configDirective, k, c.GetNamed(d.Config)))

	}
	table.AddRow(";", " ===", "", "   ===")
	table.AddRow(";", "=======", "Code section", "=======")
	table.AddRow(";", " ===", "", "   ===")
	for _, addr := range d.getSortedAddresses() {
		code, _ := d.Code[uint(addr)]
		if d.ListingMode {
			prefix := fmt.Sprintf("%05X %04X", addr, code.Bin)
			table.AddRow(prefix, code.Label, code.Asm, code.Comment)
		} else {
			prefix := code.PrefixLine
			if len(prefix) > 3 && prefix[:3] == "org" {
				table.AddRow("",prefix,  "", "")
			}
			table.AddRow(code.Label, code.Asm, code.Comment)
		}

	}
	writer.WriteString(table.String())
	writer.WriteString("\n\tEND")
	writer.Flush()
}

func (d *Disassembler) ReadScheme(file *os.File) error {
	sr := NewSchemeReader(d.Processor)
	sr.ReadScheme(file)
	d.regMap = makeRegMap(sr.Registers)
	d.instrMapper.SetTable(makeOperandTables(sr.Opcodes))
	return nil
}
