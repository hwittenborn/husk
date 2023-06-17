package main

// #include "husk.h"
import "C"

import (
	"github.com/hwittenborn/husk/util"
	"runtime/cgo"
	"unsafe"
)

// Delete the object connected to a pointer.
//
//export HuskMiscDeleteGoObj
func HuskMiscDeleteGoObj(obj C.ptr) {
	handle := cgo.Handle(obj)
	handle.Delete()
}

// Convert a `Stringer` pointer into a string.
//
//export HuskMiscStringerToString
func HuskMiscStringerToString(obj C.ptr) *C.char {
	handle := cgo.Handle(obj)
	value := handle.Value().(string)
	return C.CString(value)
}

// Convert an `error` interface into a string.
//
//export HuskMiscErrorToString
func HuskMiscErrorToString(obj C.ptr) *C.char {
	handle := cgo.Handle(obj)
	value := handle.Value().(error)
	return C.CString(value.Error())
}

// Get a C string out of a Go string array.
//
// # Arguments
// - `array`: A pointer to the Go array.
// - `index`: The position of the item in the index to return.
//
// # Returns
// A pointer to the C string. It will be null if `index` isn't a valid index in the array.
//export HuskMiscGetStringFromArray
func HuskMiscGetStringFromArray(array C.ptr, index C.int) *C.char {
	arrayHandle := cgo.Handle(array)
	goArray := arrayHandle.Value().([]string)
	goIndex := int(index)

	if goIndex < len(goArray) {
		return C.CString(goArray[goIndex])
	}

	return nil
}

// Get the Rust pointer for a `ReadContainer` from the Go `RustReaderError` object.
//
// # Arguments
// `obj`: The Go pointer
//
// # Returns
// The Rust pointer
//
//export HuskMiscGetReadContainerPtr
func HuskMiscGetReadContainerPtr(obj C.ptr) unsafe.Pointer {
	readerError := cgo.Handle(obj).Value().(util.RustReaderError)
	return readerError.ErrObj
}
