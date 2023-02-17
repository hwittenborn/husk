package ctypes

// #include <stdint.h>
// #include <stdlib.h>
import "C"

type Char = C.char
type UintptrT = C.uintptr_t
type Int = C.int

func CString(str string) *Char {
	return C.CString(str)
}

func GoString(str *Char) string {
	return C.GoString(str)
}
