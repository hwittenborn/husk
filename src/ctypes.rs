//! Functions and types that are used between here and the Go side.
use std::{collections::HashMap, ptr, ffi::{self, CString}, mem, io, ops::Deref};
use crate::bindings;

husk_proc::gen_ints!(
    HUSK_ERROR_IO,
    HUSK_ERROR_UNEXPECTED_COMMAND,
    HUSK_ERROR_UNSET_PARAMETER,
    HUSK_ERROR_EXIT_STATUS,
    HUSK_ERROR_LANG,
    HUSK_ERROR_PARSE,
    HUSK_ERROR_QUOTE,
    HUSK_ERROR_UNKNOWN,

    HUSK_SHELL_EXPAND_STRING,
    HUSK_SHELL_FIELDS_STRINGS,
    HUSK_SYNTAX_QUOTE_STRING,
    HUSK_SYNTAX_PARSER_FILE,

    HUSK_LANG_BASH,
    HUSK_LANG_POSIX,
    HUSK_LANG_MKSH,
    HUSK_LANG_BATS,
    HUSK_LANG_AUTO,
);

/// A function to get a key from a [`Hashmap<String, String>`].
#[no_mangle]
pub unsafe extern "C" fn HuskRustGetHashMapKey(hmap: *mut ffi::c_void, index: ffi::c_uint) -> *mut ffi::c_char {
    let bmap: Box<HashMap<String, String>> = Box::from_raw(hmap as *mut HashMap<String, String>);

    let value = if let Some(key) = bmap.keys().nth(index.try_into().unwrap()) {
        CString::new(key.to_owned()).unwrap().into_raw()
    } else {
        ptr::null_mut()
    };

    mem::forget(bmap);
    value
}

/// A function to call from the Go side, to get a string out of a [`HashMap<String, String>`].
#[no_mangle]
pub unsafe extern "C" fn HuskRustGetHashMapValue(hmap: *mut ffi::c_void, key: *mut ffi::c_char) -> *mut ffi::c_char {
    let bmap: Box<HashMap<String, String>> = Box::from_raw(hmap as *mut HashMap<String, String>);
    let rust_key = CString::from_raw(key).into_string().unwrap();

    let value = if let Some(value) = bmap.get(&rust_key) {
        CString::new(value.as_str()).unwrap().into_raw()
    } else {
        ptr::null_mut()
    };

    mem::forget(bmap);
    value
}

/// A struct to pass an `io::Read` trait object in between Go and Rust.
pub struct ReadContainer<'a>(pub &'a mut dyn io::Read);

/// A function to call from the Go side, to get a byte out of a [`ReadContainer`].
///
/// # Returns
/// - A pointer to the [`io::Result<usize>`] object.
#[no_mangle]
pub unsafe extern "C" fn HuskRustGetByteFromReadTrait(container: *mut ffi::c_void) -> bindings::HuskRustRead {
    let mut bcontainer: Box<ReadContainer> = Box::from_raw(container as *mut ReadContainer);
    let mut byte_container = [0; 1];

    let res = bcontainer.0.read(&mut byte_container);

    let is_ok = res.is_ok();
    let bytes_read = match &res {
        Ok(num_bytes) => *num_bytes,
        Err(_) => 0
    };
    let err_obj = match res {
        Ok(_) => ptr::null_mut(),
        Err(err) => {
            let berr = Box::new(err);
            Box::into_raw(berr) as *mut ffi::c_void
        }
    };

    mem::forget(bcontainer);
    bindings::HuskRustRead {
        isOk: is_ok,
        byte: byte_container[0],
        bytesRead: bytes_read.try_into().unwrap(),
        errObj: err_obj
    }
}
