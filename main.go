package main

import (
	flag "github.com/spf13/pflag"
	"fmt"
	"os"
	"github.com/lightAssemble/GoPicDissassemble/gpd"
)

func main() {
	hexFilePtr := flag.String("hex", "", "Hex file for disassemble (require) [disassemble --hex <hexFile>]")
	outAsmFilePtr := flag.String("out", "", "Disassembled asm out file (default <hex>_.asm) ")
	picNamePtr := flag.String("pic", "", "Target pic processor")
	//hexFile    = kingpin.Arg("hex", "Hex file for disassemble").Required().String()
	listingModePtr := flag.Bool("l", false, "Listing mode")
	debugPtr := flag.Bool("d", false, "debug")
	//tableAsm := flag.Bool("t", false, "write human friendly asm output (not for compiling)")

	flag.Parse()

	hexFile := *hexFilePtr
	asmFile := *outAsmFilePtr
	if hexFile == "" {
		fmt.Printf("GoPicDissassemble v%v\n", "0.0.1a")
		flag.PrintDefaults()
		fmt.Printf("(C) 2018 RawLight\n")
		os.Exit(0)
	}
	if asmFile == "" {
		asmFile = hexFile[:len(hexFile)-4] + "_.asm"
	}
	picName := *picNamePtr
	if picName == "" {
		picName = "p18f97j60" // TODO Make select processor
	}

	dis := gpd.NewDisassembler(*debugPtr)
	dis.ListingMode = *listingModePtr
	dis.Processor = gpd.NewProcessorInfo(picName)
	dis.TableStyle = dis.ListingMode

	fmt.Println("Building opcode tables for `" + picName + "`")

	file, err := os.Open(picName + ".hsch")
	gpd.CheckErr(err)
	defer file.Close()

	err = dis.ReadScheme(file)
	gpd.CheckErr(err)

	fmt.Println("Reading object file..." + hexFile)
	hexFileStream, err := os.Open(hexFile)
	gpd.CheckErr(err)
	dis.ReadObjectCode(hexFileStream)

	fmt.Println("Disassemble...")
	dis.Assemble()
	fmt.Println("Arranging...")
	fmt.Println("Writing...  " + asmFile)
	asmFileStream, err := os.OpenFile(asmFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	gpd.CheckErr(err)
	defer asmFileStream.Close()
	dis.WriteTo(asmFileStream)
	asmFileStream.Sync()
	fmt.Println("Done...")

}
