package printer

import (
	"io/ioutil"
	"syscall"
	"unsafe"
)

type docInfo struct {
	pDocName    []byte
	pOutputFile []byte
	pDatatype   []byte
}

var (
	dll               = syscall.MustLoadDLL("winspool.drv")
	getDefaultPrinter = dll.MustFindProc("GetDefaultPrinterW")
	openPrinter       = dll.MustFindProc("OpenPrinterW")
	startDocPrinter   = dll.MustFindProc("StartDocPrinterW")
	startPagePrinter  = dll.MustFindProc("StartPagePrinter")
	writePrinter      = dll.MustFindProc("WritePrinter")
	endPagePrinter    = dll.MustFindProc("EndPagePrinter")
	endDocPrinter     = dll.MustFindProc("EndDocPrinter")
	closePrinter      = dll.MustFindProc("ClosePrinter")
)

func main() {
	printerName, printerName16 := getDefaultPrinterName()
	printerHandle := openPrinterFunc(printerName, printerName16)
	startPrinter(printerHandle)
	startPagePrinter.Call(printerHandle)
	writePrinterFunc(printerHandle)
	endPagePrinter.Call(printerHandle)
	endDocPrinter.Call(printerHandle)
	closePrinter.Call(printerHandle)
}

func writePrinterFunc(printerHandle uintptr) {
	fileContents, _ := ioutil.ReadFile(FILENAME)
	var contentLen uintptr = uintptr(len(fileContents))
	var writtenLen int
	writePrinter.Call(printerHandle,
		uintptr(unsafe.Pointer(&fileContents[0])), contentLen,
		uintptr(unsafe.Pointer(&writtenLen)))
}

func startPrinter(printerHandle uintptr) {
	di := docInfo{[]byte(""), nil, []byte("RAW")}
	startDocPrinter.Call(printerHandle, 1, uintptr(unsafe.Pointer(&di)))
}

func openPrinterFunc(printerName string, printerName16 []uint16) uintptr {
	var printerHandle uintptr
	openPrinter.Call(uintptr(unsafe.Pointer(&printerName16[0])),
		uintptr(unsafe.Pointer(&printerHandle)), 0)
	return (printerHandle)
}

func getDefaultPrinterName() (string, []uint16) {
	var pn [256]uint16
	plen := len(pn)
	getDefaultPrinter.Call(uintptr(unsafe.Pointer(&pn)),
		uintptr(unsafe.Pointer(&plen)))
	printerName := syscall.UTF16ToString(pn[:])
	return printerName, syscall.StringToUTF16(printerName)
}
