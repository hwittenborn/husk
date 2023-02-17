package main

// #include <stdint.h>
// #include <stdlib.h>
// void my_func();
import "C"

import (
	"github.com/hwittenborn/husk/ctypes"
	"github.com/hwittenborn/husk/shell"
	"github.com/hwittenborn/husk/syntax"
	"github.com/hwittenborn/husk/util"
	_ "github.com/hwittenborn/husk/interp"
	"unsafe"
)

func main() {}

/****************/
/* C TYPE UTILS */
/****************/
// Go doesn't allow using a C type across packages, so we have to manually convert between one that other packages in this use and one that this package uses directly.

// Convert a `*C.char` into `*ctypes.Char`.
func convertToCtypesChar(char *C.char) *ctypes.Char {
	return *(**ctypes.Char)(unsafe.Pointer(&char))
}

// Convert a `*ctypes.Char` into `*C.char`.
func convertToRawChar(char *ctypes.Char) *C.char {
	return *(**C.char)(unsafe.Pointer(&char))
}

// Convert a `**C.char` into `**ctypes.Char`.
func convertToCtypesCharArray(char **C.char) **ctypes.Char {
	return *(***ctypes.Char)(unsafe.Pointer(&char))
}

// Convert a `**ctypes.Char` into `**C.char`.
func convertToRawCharArray(char **ctypes.Char) **C.char {
	return *(***C.char)(unsafe.Pointer(&char))
}

// Convert a `C.int` into `ctypes.Int`.
func convertToCtypesInt(num C.int) ctypes.Int {
	return *(*ctypes.Int)(unsafe.Pointer(&num))
}

// Convert a `ctypes.Int` into `C.int`.
func convertToRawInt(num ctypes.Int) C.int {
	return *(*C.int)(unsafe.Pointer(&num))
}

// Convert a `C.uint` into `ctypes.Uint`.
func convertToCtypesUint(num C.uint) ctypes.Uint {
	return *(*ctypes.Uint)(unsafe.Pointer(&num))
}

// Convert a `ctypes.Uint` into `C.uint`.
func convertToRawUint(num ctypes.Uint) C.uint {
	return *(*C.uint)(unsafe.Pointer(&num))
}

// Convert a `C.uint8_t` into `ctypes.Uint8`.
func convertToCtypesUint8(num C.uint8_t) ctypes.Uint8 {
	return *(*ctypes.Uint8)(unsafe.Pointer(&num))
}

// Convert a `ctypes.Uint8` into `C.uint8_t`.
func convertToRawUint8(num ctypes.Uint8) C.uint8_t {
	return *(*C.uint8_t)(unsafe.Pointer(&num))
}

// Convert a `*C.uint8_t` into `*ctypes.Uint8`.
func convertToCtypesUint8Array(num *C.uint8_t) *ctypes.Uint8 {
	return *(**ctypes.Uint8)(unsafe.Pointer(&num))
}

// Convert a `*ctypes.Uint8` into `*C.uint8_t`.
func convertToRawUint8Array(num *ctypes.Uint8) *C.uint8_t {
	return *(**C.uint8_t)(unsafe.Pointer(&num))
}

// Convert a `C.uintptr_t` into `ctypes.UintptrT`.
func convertToCtypesUintptrT(num C.uintptr_t) ctypes.UintptrT {
	return *(*ctypes.UintptrT)(unsafe.Pointer(&num))
}

// Convert a `ctypes.UintptrT` into `C.uintptr_t`.
func convertToRawUintptrT(num ctypes.UintptrT) C.uintptr_t {
	return *(*C.uintptr_t)(unsafe.Pointer(&num))
}

/********/
/* UTIL */
/********/
//export HuskUtilDeleteGoItem
func HuskUtilDeleteGoItem(ptr C.uintptr_t) {
	ctypesPtr := convertToCtypesUintptrT(ptr)
	util.HuskDeleteGoItem(ctypesPtr)
}

//export HuskUtilGetCStringFromArray
func HuskUtilGetCStringFromArray(goArray C.uintptr_t, itemPosition C.int) (cString *C.char) {
	ctypesGoArray := convertToCtypesUintptrT(goArray)
	ctypesItemPosition := convertToCtypesInt(itemPosition)
	cString = convertToRawChar(util.HuskGetCStringFromArray(ctypesGoArray, ctypesItemPosition))
	return
}

//export HuskUtilErrorString
func HuskUtilErrorString(errorPtr C.uintptr_t) *C.char {
	ctypesErrorPtr := convertToCtypesUintptrT(errorPtr)
	goErrorString := util.ErrorString(ctypesErrorPtr)
	return convertToRawChar(goErrorString)
}

/*********/
/* SHELL */
/*********/
//export HuskShellExpand
func HuskShellExpand(shellString *C.char, envVarsArray **C.char, envVarsArrayLength C.int) (outputString *C.char, errorPtr C.uintptr_t) {
	ctypesShellString := convertToCtypesChar(shellString)
	ctypesEnvVarsArray := convertToCtypesCharArray(envVarsArray)
	ctypesEnvVarsArrayLength := convertToCtypesInt(envVarsArrayLength)
	ctypesOutputString, ctypesErrorPtr := shell.Expand(ctypesShellString, ctypesEnvVarsArray, ctypesEnvVarsArrayLength)

	outputString = convertToRawChar(ctypesOutputString)
	errorPtr = convertToRawUintptrT(ctypesErrorPtr)
	return
}

//export HuskShellFields
func HuskShellFields(shellString *C.char, envVarsArray **C.char, envVarsArrayLength C.int) (goArray C.uintptr_t, errorPtr C.uintptr_t, isError bool) {
	ctypesShellString := convertToCtypesChar(shellString)
	ctypesEnvVarsArray := convertToCtypesCharArray(envVarsArray)
	ctypesEnvVarsArrayLength := convertToCtypesInt(envVarsArrayLength)
	ctypesGoArray, ctypesErrorString, ctypesIsError := shell.Fields(ctypesShellString, ctypesEnvVarsArray, ctypesEnvVarsArrayLength)

	goArray = convertToRawUintptrT(ctypesGoArray)
	errorPtr = convertToRawUintptrT(ctypesErrorString)
	isError = ctypesIsError
	return
}

/**********/
/* SYNTAX */
/**********/
//export HuskSyntaxIsKeyword
func HuskSyntaxIsKeyword(word *C.char) bool {
	return syntax.IsKeyword(convertToCtypesChar(word))
}

//export HuskSyntaxQuote
func HuskSyntaxQuote(inputString *C.char, langVariant C.int) (outputString *C.char, errorPtr C.uintptr_t, isError bool) {
	ctypesInputString := convertToCtypesChar(inputString)
	ctypesInt := convertToCtypesInt(langVariant)
	ctypesOutputString, ctypesErrorPtr, ctypesIsError := syntax.Quote(ctypesInputString, ctypesInt)

	outputString = convertToRawChar(ctypesOutputString)
	errorPtr = convertToRawUintptrT(ctypesErrorPtr)
	isError = ctypesIsError
	return
}

//export HuskSyntaxValidName
func HuskSyntaxValidName(value *C.char) bool {
	return syntax.ValidName(convertToCtypesChar(value))
}

//export HuskSyntaxNewParser
func HuskSyntaxNewParser(keepComments bool, stopAt *C.char, variantInt C.int) C.uintptr_t {
	ctypesStopAt := convertToCtypesChar(stopAt)
	ctypesVariantInt := convertToCtypesInt(variantInt)

	return convertToRawUintptrT(syntax.NewParser(keepComments, ctypesStopAt, ctypesVariantInt))
}

//export HuskSyntaxParserParse
func HuskSyntaxParserParse(parser C.uintptr_t, shellCode *C.uint8_t, shellCodeLen C.int, name *C.char) (file C.uintptr_t, errorPtr C.uintptr_t, isError bool) {
	ctypesParser := convertToCtypesUintptrT(parser)
	ctypesShellCode := convertToCtypesUint8Array(shellCode)
	ctypesShellCodeLen := convertToCtypesInt(shellCodeLen)
	ctypesName := convertToCtypesChar(name)

	ctypesFile, ctypesErrorPtr, ctypesIsError := syntax.ParserParse(ctypesParser, ctypesShellCode, ctypesShellCodeLen, ctypesName)
	file = convertToRawUintptrT(ctypesFile)
	errorPtr = convertToRawUintptrT(ctypesErrorPtr)
	isError = ctypesIsError
	return
}

//export HuskSyntaxNewPos
func HuskSyntaxNewPos(offset, line, column C.uint) C.uintptr_t {
	ctypesOffset := convertToCtypesUint(offset)
	ctypesLine := convertToCtypesUint(line)
	ctypesColumn := convertToCtypesUint(column)

	return convertToRawUintptrT(syntax.NewPos(ctypesOffset, ctypesLine, ctypesColumn))
}

//export HuskSyntaxPosAfter
func HuskSyntaxPosAfter(pos1, pos2 C.uintptr_t) bool {
	ctypesPos1 := convertToCtypesUintptrT(pos1)
	ctypesPos2 := convertToCtypesUintptrT(pos2)

	return syntax.PosAfter(ctypesPos1, ctypesPos2)
}

//export HuskSyntaxPosCol
func HuskSyntaxPosCol(pos C.uintptr_t) C.uint {
	ctypesPos := convertToCtypesUintptrT(pos)
	return convertToRawUint(syntax.PosCol(ctypesPos))
}

//export HuskSyntaxPosIsValid
func HuskSyntaxPosIsValid(pos C.uintptr_t) bool {
	ctypesPos := convertToCtypesUintptrT(pos)
	return syntax.PosIsValid(ctypesPos)
}

//export HuskSyntaxPosLine
func HuskSyntaxPosLine(pos C.uintptr_t) C.uint {
	ctypesPos := convertToCtypesUintptrT(pos)
	return convertToRawUint(syntax.PosLine(ctypesPos))
}

//export HuskSyntaxPosOffset
func HuskSyntaxPosOffset(pos C.uintptr_t) C.uint {
	ctypesPos := convertToCtypesUintptrT(pos)
	return convertToRawUint(syntax.PosOffset(ctypesPos))
}
