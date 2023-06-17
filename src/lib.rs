//! This crate provides bindings to the [`mvdan.cc/sh` Golang library](https://pkg.go.dev/mvdan.cc/sh).
//!
//! # Note
//! As this library is a wrapper around a Golang library, you'll need the [Go
//! toolchain](https://go.dev/) installed on your system in order to build it.
pub mod syntax;
pub mod shell;
pub mod util;
mod ctypes;

pub use syntax::{QuoteError, ParseError};
use thiserror::Error as ThisError;
use std::{io, collections::HashMap};

#[allow(non_camel_case_types)]
#[allow(non_snake_case)]
mod bindings {
    include!(concat!(env!("OUT_DIR"), "/bindings.rs"));
}

/// The types of errors that may occur in this crate.
#[derive(ThisError, Debug)]
pub enum Error {
    /// An I/O error.
    #[error("{0}")]
    IO(io::Error),
    /// A quoting error.
    #[error("{0}")]
    Quote(syntax::QuoteError),
    /// A parsing error.
    #[error("{0}")]
    Parse(syntax::ParseError),
    /// An unknown error.
    ///
    /// Ideally there wouldn't need to be an unknown variant, but the Go library from which this
    /// one is a wrapper of returns the `error {}` interface* in a lot of code, so it's possible
    /// that an unknown kind of error could be returned.
    ///
    /// *See [here](https://go.dev/tour/methods/9) for a brushup on Go interfaces. TLDR is that
    /// they work similarly to Rust traits.
    #[error("{0}")]
    Unknown(String),
}

/// Alias for a [`std::result::Result`] with the error type always being a [`husk::Error`](crate::Error).
pub type Result<T> = std::result::Result<T, Error>;

/// Alias for the type that stored environment variables. You can create this yourself manually by
/// making a new [`HashMap`], or you can call [`util::env_map`] to get one
/// containing the system's environment varibales.
pub type EnvMap = std::collections::HashMap<String, String>;
