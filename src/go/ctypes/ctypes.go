package ctypes

/*
#include <stdint.h>
#include <stdlib.h>

// These would be defined in the `interp` package (which is where they're used), but then we wouldn't be able to use our `ctypes` in it.
struct husk_interp_call_return { char** strings; char* error };
struct husk_interp_call_return husk_interp_call_handler(void* closure_ptr, uintptr_t strings);

uint8_t husk_interp_exec_handler(void* closure_ptr, uintptr_t strings);

husk_interp_open_handler(void* closure_ptr, char* path, uint8_t flag, uintptr_t perm)
*/
import "C"

type Char = C.char
type UintptrT = C.uintptr_t
type Int = C.int
type Uint = C.uint
type Uint8 = C.uint8_t
type Void = C.void

func CString(str string) *Char {
	return C.CString(str)
}

func GoString(str *Char) string {
	return C.GoString(str)
}
