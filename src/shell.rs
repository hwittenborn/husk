//! Parsing and formatting of shell programs. Supports POSIX Shell, Bash, and mksh.
use crate::{bindings, ctypes, util, EnvMap, Error};
use std::ffi::{self, CString};

/// Perform shell expansion on `shell_string` as if it were within double quotes, using `env_vars`
/// to resolve variables. This includes parameter expansion, arithmetic expansion, and quote
/// removal.
///
/// Command substitutions like `$(echo foo)` aren't supported to avoid running arbitrary code. To
/// support those, use an interpreter from the `expand` module.
///
/// If the input string has invalid syntax, [`Error::Expansion`] is returned.
#[husk_proc::unsafe_wrapper]
pub fn expand(shell_string: &str, env_vars: EnvMap) -> crate::Result<String> {
    let c_shell_string = CString::new(shell_string).unwrap();
    let boxed_env_vars = Box::new(env_vars);
    let env_vars_ptr = Box::into_raw(boxed_env_vars) as *mut ffi::c_void;
    let resp = bindings::HuskShellExpand(c_shell_string.as_ptr() as *mut ffi::c_char, env_vars_ptr);

    if resp.r1 == ctypes::HUSK_SHELL_EXPAND_STRING {
        let resp_str = bindings::HuskMiscStringerToString(resp.r0);
        Ok(CString::from_raw(resp_str).into_string().unwrap())
    } else {
        Err(util::error_obj_to_rust(resp.r0, resp.r1))
    }
}

/// Perform shell expansion on `shell_string` as if it were a command's arguments, using `env_vars`
/// to resolve variables. It is similar to [`expand`], but includes brace expansion, tilde
/// expansion, and globbing.
///
/// If the input string has invalid syntax, [`Error::Expansion`] is returned.
#[husk_proc::unsafe_wrapper]
pub fn fields(shell_string: &str, env_vars: EnvMap) -> crate::Result<Vec<String>> {
    let c_shell_string = CString::new(shell_string).unwrap();
    let boxed_env_vars = Box::new(env_vars);
    let env_vars_ptr = Box::into_raw(boxed_env_vars) as *mut ffi::c_void;
    
    let resp = bindings::HuskShellFields(c_shell_string.as_ptr() as *mut ffi::c_char, env_vars_ptr);

    if resp.r1 == ctypes::HUSK_SHELL_FIELDS_STRINGS {
        Ok(util::go_to_rust_string_vec(resp.r0))
    } else {
        Err(util::error_obj_to_rust(resp.r0, resp.r1))
    }
}
