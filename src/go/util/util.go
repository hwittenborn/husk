package util

import (
    "github.com/hwittenborn/husk/ctypes"
    "strings"
    "unsafe"
    "runtime/cgo"
    shSyntax "mvdan.cc/sh/v3/syntax"
)

// Return an error. These should be handled in the `husk` Rust crate.
func HuskError(errorString string, huskDefined bool) *ctypes.Char {
        var huskErrorString string

        if huskDefined {
                huskErrorString = "husk-err:"
        } else {
                huskErrorString = "err:"
        }

        huskErrorString += errorString

        return ctypes.CString(huskErrorString)
}

// Get the error string from an error.
func ErrorString(errorPtr ctypes.UintptrT) *ctypes.Char {
        errorPtrHandle := cgo.Handle(errorPtr)
        goErrorPtr := errorPtrHandle.Value().(error)
        return ctypes.CString(goErrorPtr.Error())
}

// Build a Go string array from a C string array.
func BuildStringArray(stringArray **ctypes.Char, stringArraySize ctypes.Int) []string {
        var goCStringArray []*ctypes.Char = unsafe.Slice(stringArray, stringArraySize)
        var goStringArray []string                   
         
        for _, item := range goCStringArray {
                goStringArray = append(goStringArray, ctypes.GoString(item))
        }

        return goStringArray
}

// Build a go byte array from a C byte array.
func BuildByteArray(byteArray *ctypes.Uint8, byteArraySize ctypes.Int) []byte {
        var goCByteArray []ctypes.Uint8 = unsafe.Slice(byteArray, byteArraySize)
        var goByteArray []byte

        for _, item := range goCByteArray {
                goByteArray = append(goByteArray, byte(item))
        }

        return goByteArray
}
// Convert an array of environment variables into a map of key-value pairs.
func EnvListToEnvMap(envList []string) map[string]string {
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
func HuskDeleteGoItem(ptr ctypes.UintptrT) {
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
func HuskGetCStringFromArray(goArray ctypes.UintptrT, itemPosition ctypes.Int) (cString *ctypes.Char) {
        arrayHandle := cgo.Handle(goArray)
        array := arrayHandle.Value().([]string)
        goItemPosition := int(itemPosition)

        if !(goItemPosition > (len(array) - 1)) {
                cString = ctypes.CString(array[goItemPosition])
        }

        return
}

// Get the `syntax.LangVariant` from an integer.
func GetLangVariant(langVariantInt ctypes.Int) shSyntax.LangVariant {
        var langVariant shSyntax.LangVariant

        switch langVariant {
        case 0:
                langVariant = shSyntax.LangBash
        case 1:
                langVariant = shSyntax.LangPOSIX
        case 2:
                langVariant = shSyntax.LangMirBSDKorn
        case 3:
                langVariant = shSyntax.LangBats
        case 4:
                langVariant = shSyntax.LangAuto
        default:
                panic("Invalid language variant supplied: " + string(langVariant))
        }

	return langVariant
}
