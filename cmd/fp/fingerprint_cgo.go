package main

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

//export CGenerateFingerprint
func CGenerateFingerprint() *C.char {
	s := GenerateFingerprint()
	return C.CString(s)
}

//export CGetHostname
func CGetHostname() *C.char {
	s := GetHostname()
	return C.CString(s)
}

//export FreeCStr
func FreeCStr(s *C.char) {
	C.free(unsafe.Pointer(s))
}

func main() {

}
