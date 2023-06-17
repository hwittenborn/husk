package main

// #include "husk.h"
import "C"

import (
	"github.com/hwittenborn/husk/util"
	"mvdan.cc/sh/v3/shell"
	"unsafe"
	"runtime/cgo"
)

// Wrapper around `shell.Expand`.
//
// # Arguments:
// - `shellString`: The string to be expanded.
// - `hmap`: A pointer to the Rust `HashMap` to read environment variables from.
//
// # Returns:
// - 0: A object pointing to the string/error.
// - 1: The type of the returned object.
//
//export HuskShellExpand
func HuskShellExpand(shellString *C.char, hmap *C.void) (C.ptr, C.int) {
	hmapPtr := unsafe.Pointer(hmap)
	envMap := util.RustStringMapToGo(hmapPtr)

	expandedString, err := shell.Expand(C.GoString(shellString), func(name string) string {
		return envMap[name]
	})

	if err == nil {
		return C.ptr(cgo.NewHandle(expandedString)), C.HUSK_SHELL_EXPAND_STRING
	} else {
                handle, errType := util.RustError(err)
                return C.ptr(handle), C.int(errType)
	}
}

// Wrapper around `shell.Fields`.
//
// # Arguments:
// - `shellString`: The string to be expanded.
// - `hmap`: A pointer to the Rust `HashMap` to read environment variables from.
//
// # Returns:
// - 0: A pointer to the Go object of the string array, or the error.
// - 1: Whether the returned string was the string array, or the error.
//
//export HuskShellFields
func HuskShellFields(shellString *C.char, hmap *C.void) (C.ptr, C.int) {
	hmapPtr := unsafe.Pointer(hmap)
	envMap := util.RustStringMapToGo(hmapPtr)

	expandedStrings, err := shell.Fields(C.GoString(shellString), func(name string) string {
		return envMap[name]
	})

	if err == nil {
		return C.ptr(cgo.NewHandle(expandedStrings)), C.HUSK_SHELL_FIELDS_STRINGS
	} else {
                handle, errType := util.RustError(err)
                return C.ptr(handle), C.int(errType)
	}
}
