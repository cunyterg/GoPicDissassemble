package gpd

type InstructionMapper struct {
	operandsTable []*Operand
}

func NewInstructionMap() *InstructionMapper {
	return &InstructionMapper{nil}
}

func (im *InstructionMapper) matchingOpcode(bin int) *Operand {
	for _, operand := range im.operandsTable {
		if (bin & operand.CodeMask) == operand.CodeValue {
			return operand
		}
	}
	return &Operand{AsmExtra: "X", AsmFunc: "", CodeValue: 0, CodeMask: 0}

}
func (im *InstructionMapper) SetTable(tables []*Operand) {
	im.operandsTable = tables
}
