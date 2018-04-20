# GoPicDissassemble
Microchip Pic microcontroller HEX code disassembler

The code was written with go1.10

## Usage
```
./GoPicDissassemble.exe --hex "someFile.hex"
```
#### Arguments

- target hex file (require)
```
--hex  myFile.hex
```
-  target out asm file (optional, default <hex>_.asm)
```
--out outFile.asm
```

-    listing mode (optional)
```
--l 
```

-   target microcontroller (optional) (Only p18f97j60 available and tested)
```
--pic "p18f97j60"
```
