package gpd

type instruction struct {
	Bin        int
	IsDummy    bool
	Calls      []uint
	Asm        string
	Label      string
	PrefixLine string
	Comment    string
}

func newInstruction(bin int) *instruction {
	i := new(instruction)
	i.Asm = ""
	i.Bin = bin
	i.IsDummy = false
	i.Calls = make([]uint, 0)
	i.Comment = ";"
	i.Label = " "
	i.PrefixLine = ""
	return i
}

func (i *instruction) CallsAddAddress(addr uint) {
	i.Calls = append(i.Calls, addr)
}

func (i *instruction) Dummy() *instruction {
	i.IsDummy = true
	return i
}
