package gpd

import (
	"strings"
	"strconv"
)

type Operand struct {
	CodeValue int
	CodeMask  int
	AsmFunc   string
	AsmExtra  string
}

func update(old, newCode map[uint]*instruction) {
	for k, v := range newCode {
		old[k] = v
	}
}

func makeOperandTables(opcodes []string) []*Operand {
	operands := make([]*Operand, len(opcodes))

	for index := range opcodes {
		splits := strings.Fields(opcodes[index])
		asm := splits[0]
		template := strings.Join(splits[1:len(splits)-4], " ")
		code := strings.Join(splits[len(splits)-4:], "")
		codeValue, codeMask := 0, 0
		for _, v := range code {
			codeValue = codeValue << 1
			codeMask = codeMask << 1
			if v == '0' {
				codeMask |= 1
			} else if v == '1' {
				codeMask |= 1
				codeValue |= 1
			}
		}
		operands[index] = &Operand{CodeValue: codeValue,
			AsmExtra: template, CodeMask: codeMask, AsmFunc: asm}
	}

	return operands
}
func makeRegMap(reg []string) map[int]string {
	regMap := make(map[int]string)

	for index := range reg {
		splits := strings.Fields(reg[index])
		address, _ := strconv.ParseInt(splits[0], 16, 31)
		regMap[int(address)] = splits[1]
	}
	return regMap
}

