//! This crate provides bindings to the [`mvdan.cc/sh` Golang library](https://pkg.go.dev/mvdan.cc/sh).
//!
//! # Note
//! As this library is a wrapper around a Golang library, you'll need the [Go
//! toolchain](https://go.dev/) installed on your system in order to build it.
pub mod shell;
pub mod syntax;
pub mod util;

mod bindings {
    include!(concat!(env!("OUT_DIR"), "/bindings.rs"));
}
