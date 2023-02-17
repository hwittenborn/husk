//! This crate provides bindings to the [`mvdan.cc/sh` Golang library](https://pkg.go.dev/mvdan.cc/sh).
mod bindings;
pub mod shell;
pub mod syntax;
pub mod util;
use std::{ffi::CString, fmt};

/// An error for several functions throughout this library.
///
/// This contains errors that are returned from the Go library. You can access the string representation of the error via the [`ToString`] implementation.
pub struct Error {
    ptr: usize,
}

impl Error {
    fn new(ptr: usize) -> Self {
        Self { ptr }
    }
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        unsafe {
            let str_ptr = bindings::HuskUtilErrorString(self.ptr);
            let err_string = CString::from_raw(str_ptr).into_string().unwrap();
            write!(f, "{err_string}")
        }
    }
}

impl fmt::Debug for Error {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let err_string = self.to_string();
        write!(f, "{err_string:?}")
    }
}
