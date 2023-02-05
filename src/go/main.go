// Sadly we have to include all of our code inside of one file, because cgo doesn't allow using C types across different files :P.
package main

// int HUSK_LANG_VARIANT_BASH = 0;
// int HUSK_LANG_VARIANT_POSIX = 1;
// int HUSK_LANG_VARIANT_MKSH = 2;
// int HUSK_LANG_VARIANT_BATS = 3;
// int HUSK_LANG_VARIANT_AUTO = 4;
// #include <stdint.h>
// #include <stdlib.h>
import "C"

import (
	"mvdan.cc/sh/v3/shell"
	"mvdan.cc/sh/v3/syntax"
	"runtime/cgo"
	"strings"
	"unsafe"
)

func main() {}

/*********/
/* UTILS */
/*********/
// Return an error. These should be handled in the `husk` Rust crate.
func huskError(errorString string, huskDefined bool) string {
	var huskErrorString string

	if huskDefined {
		huskErrorString = "husk-err:"
	} else {
		huskErrorString = "err:"
	}

	huskErrorString += errorString

	return huskErrorString
}

// Build a Go string array from a C string array.
func buildStringArray(stringArray **C.char, stringArraySize C.int) []string {
	var goCStringArray []*C.char = unsafe.Slice(stringArray, stringArraySize)
	var goStringArray []string

	for _, item := range goCStringArray {
		goStringArray = append(goStringArray, C.GoString(item))
	}

	return goStringArray
}

// Convert an array of environment variables into a map of key-value pairs.
func envListToEnvMap(envList []string) map[string]string {
	envMap := make(map[string]string)

	for _, keyValue := range envList {
		parts := strings.Split(keyValue, "=")
		key := parts[0]
		value := strings.Join(parts[1:], "=")

		envMap[key] = value
	}

	return envMap
}

// Free the memory for an item from Go.
//
// # Arguments:
// - `ptr`: The pointer to the Go item.
//
//export HuskDeleteGoItem
func HuskDeleteGoItem(ptr C.uintptr_t) {
	handle := cgo.Handle(ptr)
	handle.Delete()
}

// Get a C string out of a Go string array.
//
// # Arguments:
// - `goArray`: The pointer to the Go array.
// - `itemPosition`: The position of the item to return.
//
// # Returns:
// - `cString`: A pointer to the C string. Is null if `itemPosition` isn't a valid index in the list.
//
//export HuskGetCStringFromArray
func HuskGetCStringFromArray(goArray C.uintptr_t, itemPosition C.int) (cString *C.char) {
	arrayHandle := cgo.Handle(goArray)
	array := arrayHandle.Value().([]string)
	goItemPosition := int(itemPosition)

	if !(goItemPosition > (len(array) - 1)) {
		cString = C.CString(array[goItemPosition])
	}

	return
}

/*********/
/* SHELL */
/*********/
// Wrapper around `shell.Expand`.
//
// # Arguments:
// - `shellString`: A C string to be expanded.
// - `envVarsArray`: A C array of environment variables (i.e. 'hi=me').
// - `envVarsArrayLength`: The length of the `envVarsArray` array.
//
// # Returns:
// - `outputString`: The quoted string/error string.
// - `isError`: If `outputString` is an error string, or the quoted string.
//
//export HuskShellExpand
func HuskShellExpand(shellString *C.char, envVarsArray **C.char, envVarsArrayLength C.int) (outputString *C.char, isError bool) {
	goShellString := C.GoString(shellString)
	goEnvVars := buildStringArray(envVarsArray, envVarsArrayLength)
	goEnvMap := envListToEnvMap(goEnvVars)

	goQuotedString, err := shell.Expand(goShellString, func(envVar string) string {
		return goEnvMap[envVar]
	})

	if err != nil {
		outputString = C.CString(huskError(err.Error(), false))
		isError = true
	} else {
		outputString = C.CString(goQuotedString)
		isError = false
	}

	return
}

// Wrapper around `shell.Fields`.
//
// # Arguments:
// - `shellString`: A C string to be expanded.
// - `envVarsArray`: A C array of environment variables (i.e. `hi=me`).
// - `envVarsArrayLength`: The length of the `envVarsArray` array.
//
// # Returns:
// - `goArray`: The Go array of strings, passed under a pointer.
// - `errorString`: The error from parsing `shellString`, passed under a pointer.
// Only one of `goArray`/`errorString` will be set, the other will be a null pointer.
//
//export HuskShellFields
func HuskShellFields(shellString *C.char, envVarsArray **C.char, envVarsArrayLength C.int) (goArray C.uintptr_t, errorString *C.char) {
	goShellString := C.GoString(shellString)
	goEnvVars := buildStringArray(envVarsArray, envVarsArrayLength)
	goEnvMap := envListToEnvMap(goEnvVars)

	goStrings, err := shell.Fields(goShellString, func(envVar string) string {
		return goEnvMap[envVar]
	})

	if err != nil {
		errorString = C.CString(err.Error())
	} else {
		goArray = C.uintptr_t(cgo.NewHandle(goStrings))
	}

	return
}

/**********/
/* SYNTAX */
/**********/
// Wrapper around `syntax.IsKeyword`.
//
// # Arguments:
// - `word`: The keyword to check.
//
// # Returns:
// - A boolean: `true` if the given word is part of the language keywords.
//
//export HuskSyntaxIsKeyword
func HuskSyntaxIsKeyword(word *C.char) bool {
	return syntax.IsKeyword(C.GoString(word))
}

// Wrapper around `syntax.Quote`.
//
// # Arguments:
// - `inputString`: The string to quote.
// - `langVariant`: The language variant, from one of the `HUSK_LANG_VARIANT_*` constants.
//
// # Returns:
// - `outputString`: The quoted string/error string.
// - `isError`: Whether `outputString` is the quoted string or an error string.
//
//export HuskSyntaxQuote
func HuskSyntaxQuote(inputString *C.char, langVariant C.int) (outputString *C.char, isError bool) {
	var goLangVariant syntax.LangVariant

	switch langVariant {
	case C.HUSK_LANG_VARIANT_BASH:
		goLangVariant = syntax.LangBash
	case C.HUSK_LANG_VARIANT_POSIX:
		goLangVariant = syntax.LangPOSIX
	case C.HUSK_LANG_VARIANT_MKSH:
		goLangVariant = syntax.LangMirBSDKorn
	case C.HUSK_LANG_VARIANT_BATS:
		goLangVariant = syntax.LangBats
	case C.HUSK_LANG_VARIANT_AUTO:
		goLangVariant = syntax.LangAuto
	default:
		panic("Invalid language variant supplied: " + string(langVariant))
	}

	quotedString, err := syntax.Quote(C.GoString(inputString), goLangVariant)

	if err == nil {
		outputString = C.CString(quotedString)
		isError = false
	} else {
		outputString = C.CString(err.Error())
		isError = true
	}

	return
}

// Wrapper around `syntax.ValidName`.
//
// # Arguments:
// - `value`: The value to check.
//
// # Returns:
// - A boolean, `true` if `value` is a valid name.
//
//export HuskSyntaxValidName
func HuskSyntaxValidName(value *C.char) bool {
	return syntax.ValidName(C.GoString(value))
}

// vim: set sw=4 expandtab:
