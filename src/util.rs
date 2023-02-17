//! Various helper utilities.
use crate::bindings;
use std::{collections::HashMap, env, ffi::CString};

/// Get a [`Vec<String>`] from a pointer to a Go string array. The pointer doesn't point to valid
/// data after this function is ran.
pub(crate) unsafe fn go_to_rust_string_vec(array_ptr: usize) -> Vec<String> {
    let mut index = 0;
    let mut rust_vec = vec![];

    loop {
        let char_ptr = bindings::HuskUtilGetCStringFromArray(array_ptr, index);

        if char_ptr.is_null() {
            break;
        }

        rust_vec.push(CString::from_raw(char_ptr).into_string().unwrap());
        index += 1;
    }

    bindings::HuskUtilDeleteGoItem(array_ptr);
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
