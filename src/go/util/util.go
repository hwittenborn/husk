package util

// #include "../husk.h"
import "C"

import (
    "github.com/hwittenborn/husk/ctypes"
    "mvdan.cc/sh/v3/expand"
    "mvdan.cc/sh/v3/interp"
    "mvdan.cc/sh/v3/syntax"
    "unsafe"
    "runtime/cgo"
)

// An error struct for use by `rustReader.Read`.
type RustReaderError struct {
        // A pointer to the error object.
        ErrObj unsafe.Pointer
}

func (err RustReaderError) Error() string {
        panic("`github.com/hwittenborn/husk/syntax/rustReaderError` doesn't have a Display implementation")
}

// Convert a Rust `HashMap<String, String>` into a go `map[string]string`.
func RustStringMapToGo(hmap unsafe.Pointer) map[string]string {
	index := C.uint(0)
	var keys []string

	for {
		key := C.HuskRustGetHashMapKey(hmap, index)
		if key == nil {
			break
		}

		keys = append(keys, C.GoString(key))
		index += 1
	}

	stringMap := make(map[string]string)

	for _, key := range keys {
		value := C.HuskRustGetHashMapValue(hmap, C.CString(key))
		stringMap[key] = C.GoString(value)
	}

	return stringMap
}

// Get the `syntax.LangVariant` from an integer.
func GetLangVariant(langVariant ctypes.Int) syntax.LangVariant {
        var shLangVariant syntax.LangVariant
	cLangVariant := C.int(langVariant)

        switch cLangVariant {
        case C.HUSK_LANG_BASH:
                shLangVariant = syntax.LangBash
        case C.HUSK_LANG_POSIX:
                shLangVariant = syntax.LangPOSIX
        case C.HUSK_LANG_MKSH:
                shLangVariant = syntax.LangMirBSDKorn
        case C.HUSK_LANG_BATS:
                shLangVariant = syntax.LangBats
        case C.HUSK_LANG_AUTO:
                shLangVariant = syntax.LangAuto
        default:
                panic("Invalid language variant supplied: " + string(cLangVariant))
        }

	return shLangVariant
}

// Convert an `error` interface into a Go pointer, and the error type as an integer.
func RustError(err error) (cgo.Handle, int) {
	var typeInt C.int

	switch errType := err.(type) {
	case RustReaderError:
		typeInt = C.HUSK_ERROR_IO
	case expand.UnexpectedCommandError:
		typeInt = C.HUSK_ERROR_UNEXPECTED_COMMAND
	case expand.UnsetParameterError:
		typeInt = C.HUSK_ERROR_UNSET_PARAMETER
	case syntax.LangError:
		typeInt = C.HUSK_ERROR_LANG
	case syntax.ParseError:
		typeInt = C.HUSK_ERROR_PARSE
	case *syntax.QuoteError:
		// Dereference the QuoteError pointer so that we aren't returning a handle to a pointer.
		err = *errType
		typeInt = C.HUSK_ERROR_QUOTE
	default:
		_, isExitStatus := interp.IsExitStatus(err)

		if isExitStatus {
			typeInt = C.HUSK_ERROR_EXIT_STATUS
		} else {
			typeInt = C.HUSK_ERROR_UNKNOWN
		}
	}

	return cgo.NewHandle(err), int(typeInt)
}
