//! Various helper utilities.
use crate::{ctypes, bindings, syntax::{QuoteError, ParseError}};
use std::{collections::HashMap, env, ffi::{self, CString}, io};

/// Convert a [`u8`] into a [`bool`].
pub(crate) fn to_bool(num: u8) -> bool {
    match num {
        0 => false,
        1 => true,
        _ => panic!("Got invalid value for 'bool' conversion: {num}"),
    }
}

/// Convert a [`bool`] into a [`u8`].
pub(crate) fn to_u8(boolean: bool) -> u8 {
    match boolean {
        true => 1,
        false => 0,
    }
}

/// Convert a Go object pointer and integer error type into a [`crate::Error`].
///
/// After this function is called, `obj_ptr` belongs to the types returned from this function, and
/// the'll handle cleaning up memory upon being dropped. I.e. do **not** call
/// [`bindings::HuskMiscDeleteGoObj`] on the `obj_ptr` after this function is called.
pub(crate) unsafe fn error_obj_to_rust(obj_ptr: usize, err_int: ffi::c_int) -> crate::Error {
    if err_int == ctypes::HUSK_ERROR_IO {
        let ptr = bindings::HuskMiscGetReadContainerPtr(obj_ptr) as *mut io::Error;
        bindings::HuskMiscDeleteGoObj(obj_ptr);
        let berror = Box::from_raw(ptr);
        crate::Error::IO(*berror)
    } else if err_int == ctypes::HUSK_ERROR_UNEXPECTED_COMMAND {
        todo!()
    } else if err_int == ctypes::HUSK_ERROR_UNSET_PARAMETER {
        todo!()
    } else if err_int == ctypes::HUSK_ERROR_EXIT_STATUS {
        todo!()
    } else if err_int == ctypes::HUSK_ERROR_LANG {
        todo!()
    } else if err_int == ctypes::HUSK_ERROR_PARSE {
        crate::Error::Parse(ParseError::new(obj_ptr))
    } else if err_int == ctypes::HUSK_ERROR_QUOTE {
        crate::Error::Quote(QuoteError::new(obj_ptr))
    } else if err_int == ctypes::HUSK_ERROR_UNKNOWN {
        let error_str = bindings::HuskMiscErrorToString(obj_ptr);
        crate::Error::Unknown(CString::from_raw(error_str).into_string().unwrap())
    } else {
        unreachable!()
    }
}

/// Get a [`Vec<String>`] from a pointer of a Go string array. The pointer doesn't point to valid
/// data after this function is ran.
pub(crate) unsafe fn go_to_rust_string_vec(array_ptr: usize) -> Vec<String> {
    let mut index = 0;
    let mut rust_vec = vec![];

    loop {
        let char_ptr = bindings::HuskMiscGetStringFromArray(array_ptr, index);

        if char_ptr.is_null() {
            break;
        }

        rust_vec.push(CString::from_raw(char_ptr).into_string().unwrap());
        index += 1;
    }

    bindings::HuskMiscDeleteGoObj(array_ptr);
    rust_vec
}

/// Get a [`HashMap`] containing the system's environment variables.
pub fn env_map() -> HashMap<String, String> {
    let mut env_hashmap = HashMap::new();

    for (key, value) in env::vars() {
        env_hashmap.insert(key, value);
    }

    env_hashmap
}
