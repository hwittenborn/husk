//! This crate provides bindings to the [`mvdan.cc/sh` Golang library](https://pkg.go.dev/mvdan.cc/sh).
pub mod shell;
pub mod syntax;
pub mod util;

mod bindings {
    include!(concat!(env!("OUT_DIR"), "/bindings.rs"));
}
