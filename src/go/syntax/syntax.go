package syntax

import (
    "github.com/hwittenborn/husk/ctypes"
    "mvdan.cc/sh/v3/syntax"
)

// Wrapper around `syntax.IsKeyword`.
//
// # Arguments:
// - `word`: The keyword to check.
//
// # Returns:
// - A boolean: `true` if the given word is part of the language keywords.
func IsKeyword(word *ctypes.Char) bool {
	return syntax.IsKeyword(ctypes.GoString(word))
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
func Quote(inputString *ctypes.Char, langVariant ctypes.Int) (outputString *ctypes.Char, isError bool) {
	var goLangVariant syntax.LangVariant

	switch langVariant {
	case 0:
		goLangVariant = syntax.LangBash
	case 1:
		goLangVariant = syntax.LangPOSIX
	case 2:
		goLangVariant = syntax.LangMirBSDKorn
	case 3:
		goLangVariant = syntax.LangBats
	case 4:
		goLangVariant = syntax.LangAuto
	default:
		panic("Invalid language variant supplied: " + string(langVariant))
	}

	quotedString, err := syntax.Quote(ctypes.GoString(inputString), goLangVariant)

	if err == nil {
		outputString = ctypes.CString(quotedString)
		isError = false
	} else {
		outputString = ctypes.CString(err.Error())
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
func ValidName(value *ctypes.Char) bool {
	return syntax.ValidName(ctypes.GoString(value))
}
