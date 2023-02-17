package syntax

import (
    "github.com/hwittenborn/husk/ctypes"
    "github.com/hwittenborn/husk/util"
    "mvdan.cc/sh/v3/syntax"
    "runtime/cgo"
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
// - `langVariantInt`: The language variant.
//
// # Returns:
// - `outputString`: The quoted string/error string.
// - `isError`: Whether `outputString` is the quoted string or an error string.
func Quote(inputString *ctypes.Char, langVariantInt ctypes.Int) (outputString *ctypes.Char, isError bool) {
	langVariant := util.GetLangVariant(langVariantInt)
	quotedString, err := syntax.Quote(ctypes.GoString(inputString), langVariant)

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

// Wrapper around `syntax.Parser.NewParser`.
//
// # Arguments:
// - `keepComments`: Whether to keep comments.
// - `stopAt`: A string pointer for the word to stop at. Set to a null pointer to avoid stopping.
// - `variantInt`: The language variant.
//
// # Returns
// - A pointer to the `syntax.Parser`.
func NewParser(keepComments bool, stopAt *ctypes.Char, variantInt ctypes.Int) ctypes.UintptrT {
	parser := syntax.NewParser()
	syntax.KeepComments(keepComments)(parser)
	syntax.Variant(util.GetLangVariant(variantInt))

	if stopAt != nil {
		syntax.StopAt(ctypes.GoString(stopAt))(parser)
	}

	return ctypes.UintptrT(cgo.NewHandle(parser))
}

// Wrapper around `syntax.NewPos`.
//
// # Arguments:
// - `offset`: The offset.
// - `line`: The line.
// - `column`: The column.
//
// # Returns:
// - A pointer to the `syntax.Pos`.
func NewPos(offset, line, column ctypes.Uint) ctypes.UintptrT {
	pos := syntax.NewPos(uint(offset), uint(line), uint(column))
	return ctypes.UintptrT(cgo.NewHandle(pos))
}

// Wrapper around `syntax.Pos.After`.
//
// # Arguments:
// - `pos1`: A pointer to the first `syntax.Pos` object.
// - `pos2`: A pointer to the second `syntax.Pos` object.
//
// # Returns:
// - Whether `pos1` is after `p2`.
func PosAfter(pos1, pos2 ctypes.UintptrT) bool {
	pos1Handle := cgo.Handle(pos1)
	pos2Handle := cgo.Handle(pos2)
	goPos1 := pos1Handle.Value().(syntax.Pos)
	goPos2 := pos2Handle.Value().(syntax.Pos)

	return goPos1.After(goPos2)
}

// Wrapper around `syntax.Pos.Col`.
//
// # Arguments:
// - `pos`: A pointer to the `syntax.Pos` object.
//
// # Returns:
// - The column number of the position.
func PosCol(pos ctypes.UintptrT) ctypes.Uint {
	posHandle := cgo.Handle(pos)
	goPos := posHandle.Value().(syntax.Pos)
	return ctypes.Uint(goPos.Col())
}

// Wrapper around `syntax.Pos.IsValid`.
//
// # Arguments:
// - `pos`: A pointer to the `syntax.Pos` object.
//
// # Returns:
// - Whether the position contains useful position information.
func PosIsValid(pos ctypes.UintptrT) bool {
	posHandle := cgo.Handle(pos)
	goPos := posHandle.Value().(syntax.Pos)
	return goPos.IsValid()
}

// Wrapper around `syntax.Pos.Line`.
//
// # Arguments:
// - `pos`: A pointer to the `syntax.Pos` object.
//
// # Returns:
// - The line number of the position.
func PosLine(pos ctypes.UintptrT) ctypes.Uint {
	posHandle := cgo.Handle(pos)
	goPos := posHandle.Value().(syntax.Pos)
	return ctypes.Uint(goPos.Line())
}

// Wrapper around `syntax.Pos.Offset`.
//
// # Arguments:
// - `pos`: A pointer to the `syntax.Pos` object.
//
// # Returns:
// - The offset of the position.
func PosOffset(pos ctypes.UintptrT) ctypes.Uint {
	posHandle := cgo.Handle(pos)
	goPos := posHandle.Value().(syntax.Pos)
	return ctypes.Uint(goPos.Offset())
}
