package main

// #include "husk.h"
import "C"

import (
	"github.com/hwittenborn/husk/util"
	"github.com/hwittenborn/husk/ctypes"
	"mvdan.cc/sh/v3/syntax"
	"runtime/cgo"
	"unsafe"
	"io"
)

// A struct to wrap Rust's `io::Read` trait into Go's `io.Reader` interface.
type rustReader struct {
	// A pointer to the `ReadContainer` struct from Rust.
	container *C.void
}

func (reader rustReader) Read(bytes []byte) (int, error) {
	data := C.HuskRustGetByteFromReadTrait(unsafe.Pointer(reader.container))

	if data.isOk {
		bytesRead := int(data.bytesRead)
		if bytesRead == 0 {
			return 0, io.EOF
		}

		byteArray := make([]byte, 1)
		byteArray[0] = byte(data.byte)
		copy(bytes, byteArray)
		return bytesRead, nil
	} else {
		return 0, util.RustReaderError { ErrObj: data.errObj }
	}
}

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
// - `langVariantInt`: The language variant.
//
// # Returns:
// - 0: A pointer to the quoted string, or the error object.
// - 1: What kind of object `outputString` is.
//
//export HuskSyntaxQuote
func HuskSyntaxQuote(inputString *C.char, langVariantInt C.int) (C.ptr, C.int) {
	ctypesLangVariantInt := ctypes.Int(langVariantInt)

	langVariant := util.GetLangVariant(ctypesLangVariantInt)
	quotedString, err := syntax.Quote(C.GoString(inputString), langVariant)

	if err == nil {
		return C.ptr(cgo.NewHandle(quotedString)), C.HUSK_SYNTAX_QUOTE_STRING
	} else {
		handle, errType := util.RustError(err)
		return C.ptr(handle), C.int(errType)
	}
}

// Get the data out of a `syntax.QuoteError` pointer.
//
// # Returns
// 0. The byte offset of the error.
// 1. The error message.
//
//export HuskSyntaxQuoteErrorData
func HuskSyntaxQuoteErrorData(obj C.ptr) (C.int, *C.char) {
	handle := cgo.Handle(obj)
	quoteError := handle.Value().(syntax.QuoteError)
	cByteOffset := C.int(quoteError.ByteOffset)
	cMessage := C.CString(quoteError.Message)

	return cByteOffset, cMessage
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

// Wrapper around `syntax.LangVariant.String`.
//
// # Arguments
// - `variant`: The integer for the language variant.
//
// # Returns
// - The string
//
//export HuskSyntaxLangVariantString
func HuskSyntaxLangVariantString(variant C.int) *C.char {
	langVariant := util.GetLangVariant(ctypes.Int(variant))
	return C.CString(langVariant.String())
}

// Wrapper around `syntax.NewPos`.
//
// # Arguments
// - `offset`: The offset position.
// - `line`: The line position.
// - `column`: The column position.
//
// # Returns
// A pointer to the `syntax.Pos` object.
//
//export HuskSyntaxNewPos
func HuskSyntaxNewPos(offset, line, column C.uint) C.ptr {
	goOffset := uint(offset)
	goLine := uint(line)
	goColumn := uint(column)

	handle := cgo.NewHandle(syntax.NewPos(goOffset, goLine, goColumn))
	return C.ptr(handle)
}

// Wrapper around `syntax.Pos.After`.
//
// # Arguments
// - `pos1`: A pointer to the first `syntax.Pos` object.
// - `pos2`: A pointer to the second `syntax.Pos` object.
//
// # Returns
// A boolean, indicating if `pos1` is after `pos2`.
//
//export HuskSyntaxPosAfter
func HuskSyntaxPosAfter(pos1, pos2 C.ptr) bool {
	goPos1 := cgo.Handle(pos1).Value().(syntax.Pos)
	goPos2 := cgo.Handle(pos2).Value().(syntax.Pos)
	return goPos1.After(goPos2)
}

// Wrapper around `syntax.Pos.Col`.
//
// # Arguments
// - `obj`: A pointer to the `syntax.Pos` object.
//
// # Returns
// The column number of the position.
//
//export HuskSyntaxPosCol
func HuskSyntaxPosCol(obj C.ptr) C.uint {
	pos := cgo.Handle(obj).Value().(syntax.Pos)
	return C.uint(pos.Col())
}

// Wrapper around `syntax.Pos.Line`.
//
// # Arguments
// - `obj`: A pointer to the `syntax.Pos` object.
//
// # Returns
// The line number of the position.
//
//export HuskSyntaxPosLine
func HuskSyntaxPosLine(obj C.ptr) C.uint {
	pos := cgo.Handle(obj).Value().(syntax.Pos)
	return C.uint(pos.Line())
}

// Wrapper around `syntax.Pos.Offset`.
//
// # Arguments
// - `obj`: A pointer to the `syntax.Pos` object.
//
// # Returns
// The offset number of the position.
//
//export HuskSyntaxPosOffset
func HuskSyntaxPosOffset(obj C.ptr) C.uint {
	pos := cgo.Handle(obj).Value().(syntax.Pos)
	return C.uint(pos.Offset())
}

// Wrapper around `syntax.Pos.IsValid`.
//
// # Arguments
// - `obj`: A pointer to the `syntax.Pos` object.
//
// # Returns
// A boolean, representing if the position is valid.
//
//export HuskSyntaxPosIsValid
func HuskSyntaxPosIsValid(obj C.ptr) bool {
	pos := cgo.Handle(obj).Value().(syntax.Pos)
	return pos.IsValid()
}

// Wrapper around `syntax.NewParser`.
//
// # Arguments
// - `keepComments`: A pointer pointing to whether to keep comments.
// - `stopAt`: A pointer to the character to stop at.
// - `langVariant`: A pointer to the int of the language variant to use.
//
// Make any of the above pointers null to avoid using them.
//
// # Returns
// A pointer to the `*syntax.Parser` object.
//
//export HuskSyntaxNewParser
func HuskSyntaxNewParser(keepComments *bool, stopAt *C.char, langVariant *C.int) C.ptr {
	var args []syntax.ParserOption

	if keepComments != nil {
		args = append(args, syntax.KeepComments(*keepComments))
	}
	if stopAt != nil {
		args = append(args, syntax.StopAt(C.GoString(stopAt)))
	}
	if langVariant != nil {
		shLangVariant := util.GetLangVariant(ctypes.Int(*langVariant))
		args = append(args, syntax.Variant(shLangVariant))
	}

	parser := syntax.NewParser(args...)
	return C.ptr(cgo.NewHandle(parser))
}

// Wrapper around `syntax.Parser.parse`.
//
// # Arguments
// - `obj`: A pointer to the `*syntax.Parser` object.
// - `container`: A pointer to the rust `ReadContainer` object to read the bytes from.
// - `name`: The name of the shell program.
//
// # Returns
// - 0: A pointer to the Go object, containing either the `*File` or the error.
// - 1: The return type.
//
//export HuskSyntaxParserParse
func HuskSyntaxParserParse(obj C.ptr, container *C.void, name *C.char) (C.ptr, C.int) {
	parser := cgo.Handle(obj).Value().(*syntax.Parser)
	reader := rustReader { container }
	file, err := parser.Parse(reader, C.GoString(name))

	if err == nil {
		return C.ptr(cgo.NewHandle(file)), C.HUSK_SYNTAX_PARSER_FILE
	} else {
		cErr, cErrInt := util.RustError(err)
		return C.ptr(cErr), C.int(cErrInt)
	}
}

// Get the data out of a `syntax.ParseError`.
//
// # Arguments
// - `obj`: A pointer to the `syntax.ParseError` object.
//
// # Returns
// - 0: The filename
// - 1: A pointer to the `syntax.Pos` object.
// - 2: The text
// - 3: Whether the error is from being incomplete.
//
//export HuskSyntaxParseErrorData
func HuskSyntaxParseErrorData(obj C.ptr) (*C.char, C.ptr, *C.char, bool) {
	parseError := cgo.Handle(obj).Value().(syntax.ParseError)
	posHandle := cgo.NewHandle(parseError.Pos)

	return C.CString(parseError.Filename), C.ptr(posHandle), C.CString(parseError.Text), parseError.Incomplete
}
