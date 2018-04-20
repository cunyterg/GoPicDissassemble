package gpd

import (
	"bufio"
	"strings"
	"io"
	"fmt"
)

const (
	SectionNone    = iota
	SectionOpcodes
	SectionReg
	SectionInfo
	SectionConfig
)

type SchemeReader struct {
	Opcodes   []string
	Registers []string
	Processor *ProcessorInfo
}

func NewSchemeReader(proc *ProcessorInfo) *SchemeReader {
	s := new(SchemeReader)
	s.Opcodes = make([]string, 0)
	s.Registers = make([]string, 0)
	s.Processor = proc
	return s
}

func (s *SchemeReader) ReadScheme(stream io.Reader) {

	currentSection := SectionNone
	section := map[string]int{"REG": SectionReg, "OPCODE": SectionOpcodes,
		"INFO": SectionInfo, "CONFIG": SectionConfig}

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		ll := len(line)
		if 1 > ll {
			continue
		}
		if line[0] == '[' && line[ll-1] == ']' {
			newSectionName := line[1 : ll-1]
			if newSection, ok := section[newSectionName]; ok {
				currentSection = newSection
			} else {
				fmt.Printf("[Warn] Unknown section `%v` in schematic file\n", newSectionName)
			}
		} else {
			switch currentSection {
			case SectionOpcodes:
				s.Opcodes = append(s.Opcodes, line)
			case SectionReg:
				s.Registers = append(s.Registers, line)
			case SectionInfo:
				s.Processor.Scan(line)
			case SectionConfig:
				s.Processor.ScanConfig(line)

			}
		}

	}
}
