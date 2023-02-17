//! Higher-level abstractions to the [`syntax`](mod@crate::syntax), [`expand`](mod@crate::expand), and [`interp`](mod@crate::interp) modules.
use crate::{bindings, util};
use std::{collections::HashMap, ffi::CString};

/// Perform shell expansion on `shell_string` as if it were within double quotes, using `env_vars`
/// to resolve variables. This includes parameter expansion, arithmetic expansion, and quote
/// removal.
///
/// System environment variables can be used for expansion by using the [`util::env_map`]
/// function.
///
/// Command substitutions like `$(echo foo)` aren't supported to avoid running arbitrary code. To
/// support those, use an interpreter from the [`expand`](mod@crate::expand) module.
///
/// If the input string has invalid syntax, [`Err`] is returned.
///
/// # Example
/// ```rust
/// use std::env;
/// use husk::{shell, util};
///
/// env::set_var("NAME", "Foo Bar");
/// let env_map = util::env_map();
/// let expanded_string = shell::expand("Hello, ${NAME}!", env_map).unwrap();
/// assert_eq!(expanded_string, "Hello, Foo Bar!");
/// ```
pub fn expand<S: AsRef<str>>(
    shell_string: &str,
    env_vars: HashMap<S, S>,
) -> Result<String, crate::Error> {
    let shell_string_ffi = CString::new(shell_string).unwrap();
    let mut env_vec: Vec<CString> = vec![];

    for (key, value) in env_vars {
        let env_var_ffi =
            CString::new(format!("{}={}", key.as_ref(), value.as_ref()).as_str()).unwrap();
        env_vec.push(env_var_ffi);
    }

    let env_vec_ptrs: Vec<*mut libc::c_char> = env_vec
        .iter()
        .map(|string| string.as_ptr() as *mut libc::c_char)
        .collect();
    let resp = unsafe {
        bindings::HuskShellExpand(
            shell_string_ffi.as_ptr() as *mut libc::c_char,
            env_vec_ptrs.as_ptr() as *mut *mut libc::c_char,
            env_vec_ptrs.len().try_into().unwrap(),
        )
    };

    if !resp.r0.is_null() {
        Ok(unsafe { CString::from_raw(resp.r0).into_string().unwrap() })
    } else {
        Err(crate::Error::new(resp.r1))
    }
}

/// Perform shell expansion on `shell_string` as if it were command arguments, using `env_vars` to
/// resolve variables. It is similar to [`expand`], but includes brace expansion, tilde expansion,
/// and globbing.
///
/// System environment variables can be used for expansion by using the [`util::env_map`]
/// function.
///
/// If the input string has invalid syntax, [`Err`] is returned.
pub fn fields<S: AsRef<str>>(
    shell_string: &str,
    env_vars: HashMap<S, S>,
) -> Result<Vec<String>, crate::Error> {
    let shell_string_ffi = CString::new(shell_string).unwrap();

    let mut env_vec: Vec<CString> = vec![];

    for (key, value) in env_vars {
        let env_var_ffi =
            CString::new(format!("{}={}", key.as_ref(), value.as_ref()).as_str()).unwrap();
        env_vec.push(env_var_ffi);
    }

    let env_vec_ptrs: Vec<*mut libc::c_char> = env_vec
        .iter()
        .map(|string| string.as_ptr() as *mut libc::c_char)
        .collect();
    let resp = unsafe {
        bindings::HuskShellFields(
            shell_string_ffi.as_ptr() as *mut libc::c_char,
            env_vec_ptrs.as_ptr() as *mut *mut libc::c_char,
            env_vec_ptrs.len().try_into().unwrap(),
        )
    };

    if resp.r2 != 0 {
        Ok(unsafe { util::go_to_rust_string_vec(resp.r0) })
    } else {
        Err(crate::Error::new(resp.r1))
    }
}
