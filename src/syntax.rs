//! Parsing and formatting of shell programs. Supports POSIX Shell, Bash, and mksh.
use crate::bindings;
use std::ffi::CString;

/// A shell language variant to use when tokenizing and parsing shell code.
#[derive(Clone, Debug)]
pub enum LangVariant {
    /// The GNU Bash language, as described in its manual at
    /// <https://www.gnu.org/software/bash/manual/bash.html>.
    ///
    /// Currently Bash version 5.1 is followed.
    Bash,

    /// The POSIX Shell language, as described at
    /// <https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html>.
    Posix,

    /// The MirBSD Korn Shell, also known as mksh, as described at
    /// <http://www.mirbsd.org/htman/i386/man1/mksh.htm>. Note that it shares some features with
    /// Bash, due to the shared ancestry that is ksh.
    ///
    /// Currently mksh version 59 is followed.
    Mksh,

    /// The Bash Automated Testing System language, as described at
    /// <https://github.com/bats-core/bats-core>. Note that it's just a small extension of the Bash
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

/// A position within a shell source file.
pub struct Pos {
    ptr: usize,
    /// The column number of the position, starting at 1. It counts in bytes.
    ///
    /// This field is protected against overflows; if an input line has too many columns, extra
    /// columns will have a column number of 0.
    pub col: u32,
    /// The line number of the position, starting at 1.
    ///
    /// This field is protected against overflows; if an input line has too many columns, extra
    /// columns will have a column number of 0.
    pub line: u32,
    /// The byte offset of the position in the original source file. Byte offsets start at 0.
    ///
    /// Note that this field is not protected against overflows; if an input is larger than 4GiB,
    /// the offset will wrap around to 0.
    pub offset: u32,
}

impl Pos {
    /// Create a position with the given `offset`, `line`, and `column`.
    ///
    /// Note that [`Pos`] uses a limited number of bits to store these numbers. If `line` or
    /// `column` overflow their allocated space, they are replaced with 0.
    pub fn new(offset: u32, line: u32, column: u32) -> Self {
        let ptr = unsafe { bindings::HuskSyntaxNewPos(offset, line, column) };
        let offset = unsafe { bindings::HuskSyntaxPosOffset(ptr) };
        let line = unsafe { bindings::HuskSyntaxPosLine(ptr) };
        let col = unsafe { bindings::HuskSyntaxPosCol(ptr) };
        Self {
            ptr,
            col,
            line,
            offset,
        }
    }

    /// Report whether [`self`] is after `pos2`. It is a more expressive version of
    /// `self.offset() > pos2.offset()`.
    pub fn after(&self, pos2: &Self) -> bool {
        unsafe { bindings::HuskSyntaxPosAfter(self.ptr, pos2.ptr) != 0 }
    }

    /// Report whether the position contains useful position information. Some positions returned
    /// via [`parse()`] may be invalid: for example, [`SemiColon`] will only be valid if a
    /// statement contained a closing token such as `;`.
    pub fn is_valid(&self) -> bool {
        unsafe { bindings::HuskSyntaxPosIsValid(self.ptr) != 0 }
    }
}

impl Drop for Pos {
    fn drop(&mut self) {
        unsafe { bindings::HuskDeleteGoItem(self.ptr) }
    }
}

/// A struct that holds the internal state of the parsing mechanism of a program.
pub struct Parser {
    ptr: usize,
}

impl Parser {
    /// Create a new parser.
    ///
    /// # Arguments
    /// ## `keep_comments`
    /// Make the parser parse comments and attach them to nodes, as opposed to discarding them.
    ///
    /// ## `stop_at`
    /// Configures the lexer to stop at an arbitrary word, treating it as if it were the end of the
    /// input. It can contain any characters except whitespace, and cannot be over four bytes in
    /// size.
    ///
    /// This can be useful to embed shell code within another language, as one can use a special
    /// word to mark the delimiters between the two.
    ///
    /// As a word, it will only apply when following whitespace or a separating token. For example,
    /// `$$` will act on the inputs `foo $$` and `foo;$$`, but not on `foo '$$'`.
    ///
    /// The match is done by prefix, so the example above will also act on `foo $$bar`.
    ///
    /// ## `lang_variant`
    /// The language variant to use.
    pub fn new(keep_comments: bool, stop_at: Option<&str>, lang_variant: LangVariant) -> Self {
        let stop_at_ffi_ptr: *mut libc::c_char;
        let stop_at_cstring = stop_at.map(|stop_at| CString::new(stop_at).unwrap());

        if let Some(string) = stop_at_cstring {
            stop_at_ffi_ptr = string.as_ptr() as *mut libc::c_char;
        } else {
            stop_at_ffi_ptr = std::ptr::null_mut();
        }

        let ptr = unsafe {
            bindings::HuskSyntaxNewParser(
                keep_comments as u8,
                stop_at_ffi_ptr,
                lang_variant.to_ffi_int(),
            )
        };

        Self { ptr }
    }
}

impl Drop for Parser {
    fn drop(&mut self) {
        unsafe { bindings::HuskDeleteGoItem(self.ptr) }
    }
}
