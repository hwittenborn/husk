//! Parsing and formatting of shell programs. Supports POSIX Shell, Bash, and mksh.
use crate::bindings;
use std::ffi::CString;

/// A shell language variant to use when tokenizing and parsing shell code.
#[derive(Clone, Debug)]
pub enum LangVariant {
    /// The GNU Bash language, as described in its manual at
    /// https://www.gnu.org/software/bash/manual/bash.html.
    ///
    /// Currently Bash version 5.1 is followed.
    Bash,

    /// The POSIX Shell language, as described at
    /// https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html.
    Posix,

    /// The MirBSD Korn Shell, also known as mksh, as described at
    /// http://www.mirbsd.org/htman/i386/man1/mksh.htm. Note that it shares some features with
    /// Bash, due to the shared ancestry that is ksh.
    ///
    /// Currently mksh version 59 is followed.
    Mksh,

    /// The Bash Automated Testing System language, as described at
    /// https://github.com/bats-core/bats-core. Note that it's just a small extension of the Bash
    /// language.
    Bats,
}

impl LangVariant {
    fn to_ffi_int(&self) -> i32 {
        match self {
            Self::Bash => 0,
            Self::Posix => 1,
            Self::Mksh => 2,
            Self::Bats => 3,
        }
    }
}

/// Check if the given word is a language keyword.
pub fn is_keyword(word: &str) -> bool {
    let word_ffi = CString::new(word).unwrap();
    unsafe { bindings::HuskSyntaxIsKeyword(word_ffi.as_ptr() as *mut libc::c_char) != 0 }
}

/// Quote `input_string` so that the quoted version is expanded or interpreted as the original
/// string in the language variant set by `lang_variant`.
pub fn quote(input_string: &str, lang_variant: LangVariant) -> Result<String, String> {
    let input_ffi = CString::new(input_string).unwrap();
    let resp = unsafe {
        bindings::HuskSyntaxQuote(
            input_ffi.as_ptr() as *mut libc::c_char,
            lang_variant.to_ffi_int(),
        )
    };
    let res = unsafe { CString::from_raw(resp.r0).into_string().unwrap() };

    if resp.r1 != 0 {
        Ok(res)
    } else {
        Err(res)
    }
}

/// Check if a value is a valid name as per the POSIX spec.
pub fn valid_name(value: &str) -> bool {
    let value_ffi = CString::new(value).unwrap();
    unsafe { bindings::HuskSyntaxValidName(value_ffi.as_ptr() as *mut libc::c_char) != 0 }
}
