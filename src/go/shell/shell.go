package shell

import (
    "runtime/cgo"
    "github.com/hwittenborn/husk/ctypes"
    "github.com/hwittenborn/husk/util"
    "mvdan.cc/sh/v3/shell"
)

// Wrapper around `shell.Expand`.
//
// # Arguments:
// - `shellString`: A C string to be expanded.
// - `envVarsArray`: A C array of environment variables (i.e. 'hi=me').
// - `envVarsArrayLength`: The length of the `envVarsArray` array.
//
// # Returns:
// - `outputString`: The quoted string.
// - `errorPtr`: A pointer to the error, if one was found.
//
//export HuskShellExpand
func Expand(shellString *ctypes.Char, envVarsArray **ctypes.Char, envVarsArrayLength ctypes.Int) (outputString *ctypes.Char, errorPtr ctypes.UintptrT) {
	goShellString := ctypes.GoString(shellString)
	goEnvVars := util.BuildStringArray(envVarsArray, envVarsArrayLength)
	goEnvMap := util.EnvListToEnvMap(goEnvVars)

	goQuotedString, err := shell.Expand(goShellString, func(envVar string) string {
		return goEnvMap[envVar]
	})

	if err != nil {
		errorPtr = ctypes.UintptrT(cgo.NewHandle(err))
	} else {
		outputString = ctypes.CString(goQuotedString)
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
// - `errorPtr`: The pointer to the error, if one was found.
// - `isError`: Whether this function returned an error, can be used to decide which of the above return values is a valid pointer.
// Only one of `goArray`/`errorString` will be set, the other will be a null pointer.
//
//export HuskShellFields
func Fields(shellString *ctypes.Char, envVarsArray **ctypes.Char, envVarsArrayLength ctypes.Int) (goArray ctypes.UintptrT, errorPtr ctypes.UintptrT, isError bool) {
	goShellString := ctypes.GoString(shellString)
	goEnvVars := util.BuildStringArray(envVarsArray, envVarsArrayLength)
	goEnvMap := util.EnvListToEnvMap(goEnvVars)

	goStrings, err := shell.Fields(goShellString, func(envVar string) string {
		return goEnvMap[envVar]
	})

	if err != nil {
		errorPtr = ctypes.UintptrT(cgo.NewHandle(err))
		isError = true
	} else {
		goArray = ctypes.UintptrT(cgo.NewHandle(goStrings))
		isError = false
	}

	return
}
