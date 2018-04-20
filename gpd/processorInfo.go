package gpd

import (
	"strings"
	"github.com/lightAssemble/GoPicDissassemble/gpd/misc"
)

type ProcessorInfo struct {
	Name        string
	Info        map[string]string
	Config      map[string]*configer
	TopOfMemory uint
	TopOfConfig uint
}

func NewProcessorInfo(name string) *ProcessorInfo {
	p := new(ProcessorInfo)
	p.Name = name
	p.Info = make(map[string]string)
	p.Config = make(map[string]*configer)
	return p
}

func (p *ProcessorInfo) Scan(line string) {
	items := strings.SplitN(line, "=", 2)
	switch items[0] {
	case "TopOfMemory":
		p.TopOfConfig = misc.ParseUInt(items[1])
		p.TopOfMemory = p.TopOfConfig - 7
	default:
		p.Info[items[0]] = items[1]
	}
}

func (p *ProcessorInfo) ConfigValue(directive string, memory map[uint]uint) {

}
func (p *ProcessorInfo) ScanConfig(line string) {
	sep := strings.Index(line, "|")
	items := strings.Fields(line[sep+1:])
	line = strings.TrimSpace(line[:sep])
	ll := len(line)
	shift := uint16(misc.ParseUInt(line[ll-2 : ll-1]))
	nameValue := strings.SplitN(line[:ll-3], "=", 2)
	values := strings.SplitN(nameValue[1], "&", 2)
	configName := nameValue[0]
	address := misc.ParseUInt(strings.TrimSpace(values[0]))
	mask := uint16(misc.ParseInt64(strings.TrimSpace(values[1])))

	p.Config[configName] = &configer{mask: mask, address: address, shift: shift, items: items}

}

type configer struct {
	address uint
	mask    uint16
	shift   uint16
	items   []string
}

func (c *configer) Set(m map[uint]uint16, v uint16) {
	mv := m[c.address] & ^c.mask
	m[c.address] = mv | ((v << c.shift) & c.mask)
}
func (c *configer) SetNamed(m map[uint]uint16, name string) {
	for v, vName := range c.items {
		if name == vName {
			mv := m[c.address] & ^c.mask
			m[c.address] = mv | ((uint16(v) << c.shift) & c.mask)
			return 
		}
	}

}

func (c *configer) Get(m map[uint]uint16) uint16 {
	return (m[c.address] & c.mask) >> c.shift
}

func (c *configer) GetNamed(m map[uint]uint16) string {
	return c.items[(m[c.address]&c.mask)>>c.shift ]
}
