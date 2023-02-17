package interp

import (
    _ "runtime/cgo"
    "github.com/hwittenborn/husk/ctypes"
    _ "github.com/hwittenborn/husk/util"
    _ "mvdan.cc/sh/v3/shell"
)

// Wrapper around `interp.New`
//
// # Arguments:
// - `callHandler`: A closure pointer to run as the call handler.
// - `execHandler`: A closure pointer to run as the exec handler.
// - `openHandler`: A closure pointer to run as the open handler.
// - `readDirHandler`: A closure pointer to run as the read directory handler.
// - `statHandler`: A closure pointer to run as the stat handler.
// - `stdinHandler`: A closure pointer to handle stdin.
// - `stdoutHandler`: A closure pointer to handle stdout.
// - `stderrHandler`: A closure pointer to handle stderr.
// - `dir`: The interpreter's working directory.
// - `env`: A C array of environment variables.
// - `envLen`: The length of `env`.
// - `params`: Sets the shell options and parameters
// - `paramsLen`: The length of `params`
func RunnerNew(
	callHandler, execHandler, openHandler, readDirHandler, statHandler, stdinHandler, stdoutHandler, stderrHandler *ctypes.Void,
	dir *ctypes.Char,
	env **ctypes.Char,
	envLen ctypes.Int,
	params **ctypes.Char,
	paramsLen ctypes.Int,
) {
}