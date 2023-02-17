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
// - `outputString`: The quoted string/error string.
// - `isError`: If `outputString` is an error string, or the quoted string.
//
//export HuskShellExpand
func Expand(shellString *ctypes.Char, envVarsArray **ctypes.Char, envVarsArrayLength ctypes.Int) (outputString *ctypes.Char, isError bool) {
	goShellString := ctypes.GoString(shellString)
	goEnvVars := util.BuildStringArray(envVarsArray, envVarsArrayLength)
	goEnvMap := util.EnvListToEnvMap(goEnvVars)

	goQuotedString, err := shell.Expand(goShellString, func(envVar string) string {
		return goEnvMap[envVar]
	})

	if err != nil {
		outputString = util.HuskError(err.Error(), false)
		isError = true
	} else {
		outputString = ctypes.CString(goQuotedString)
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
func Fields(shellString *ctypes.Char, envVarsArray **ctypes.Char, envVarsArrayLength ctypes.Int) (goArray ctypes.UintptrT, errorString *ctypes.Char) {
	goShellString := ctypes.GoString(shellString)
	goEnvVars := util.BuildStringArray(envVarsArray, envVarsArrayLength)
	goEnvMap := util.EnvListToEnvMap(goEnvVars)

	goStrings, err := shell.Fields(goShellString, func(envVar string) string {
		return goEnvMap[envVar]
	})

	if err != nil {
		errorString = ctypes.CString(err.Error())
	} else {
		goArray = ctypes.UintptrT(cgo.NewHandle(goStrings))
	}

	return
}
