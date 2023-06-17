//! Parsing and formatting of shell programs. Supports POSIX Shell, Bash, and mksh.
use crate::{bindings, ctypes, util, Error};
use std::{ptr, io, ffi::{self, CString}, fmt, ops::Deref};

// One-off functions.
/// Check if the given word is a language keyword.
#[husk_proc::unsafe_wrapper]
pub fn is_keyword(word: &str) -> bool {
    let word_ffi = CString::new(word).unwrap();
    util::to_bool(bindings::HuskSyntaxIsKeyword(word_ffi.as_ptr() as *mut libc::c_char))
}

/// Quote `input_string` so that the quoted version is expanded or interpreted as the original
/// string in the language variant set by `lang_variant`.
#[husk_proc::unsafe_wrapper]
pub fn quote(input_string: &str, lang_variant: LangVariant) -> crate::Result<String> {
    let input_ffi = CString::new(input_string).unwrap();
    let resp = bindings::HuskSyntaxQuote(
        input_ffi.as_ptr() as *mut libc::c_char,
        lang_variant.to_int(),
    );

    if resp.r1 == ctypes::HUSK_SYNTAX_QUOTE_STRING {
        let chars = bindings::HuskMiscStringerToString(resp.r0);
        bindings::HuskMiscDeleteGoObj(resp.r0);
        Ok(CString::from_raw(chars).into_string().unwrap())
    } else {
        Err(util::error_obj_to_rust(resp.r0, resp.r1))
    }
}

/// Check if a value is a valid name as per the POSIX spec.
#[husk_proc::unsafe_wrapper]
pub fn valid_name(value: &str) -> bool {
    let value_ffi = CString::new(value).unwrap();
    util::to_bool(bindings::HuskSyntaxValidName(value_ffi.as_ptr() as *mut libc::c_char))
}

// Errors
pub struct QuoteError {
    ptr: usize,
    /// Where the error occurred at, represented as the number of bytes from the beginning of the
    /// string.
    pub byte_offset: u32,
    /// The error message.
    pub msg: String
}

impl QuoteError {
    #[husk_proc::unsafe_wrapper]
    pub(crate) fn new(ptr: usize) -> Self {
        let quote_error = bindings::HuskSyntaxQuoteErrorData(ptr);
        let quote_error_string = CString::from_raw(quote_error.r1).into_string().unwrap();

        Self {
            ptr,
            byte_offset: quote_error.r0.try_into().unwrap(),
            msg: quote_error_string
        }
    }
}

impl Drop for QuoteError {
    #[husk_proc::unsafe_wrapper]
    fn drop(&mut self) {
        bindings::HuskMiscDeleteGoObj(self.ptr);
    }
}

impl fmt::Debug for QuoteError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.debug_struct("QuoteError")
            .field("byte_offset", &self.byte_offset)
            .field("msg", &self.msg)
            .finish()
    }
}

impl fmt::Display for QuoteError {
    #[husk_proc::unsafe_wrapper]
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let error_str = bindings::HuskMiscErrorToString(self.ptr);
        let string = CString::from_raw(error_str).into_string().unwrap();
        write!(f, "{string}")
    }
}

pub struct ParseError {
    ptr: usize,
    pub filename: Option<String>,
    pos: Pos,
    pub text: String,
    pub incomplete: bool
}

impl ParseError {
    #[husk_proc::unsafe_wrapper]
    pub(crate) fn new(ptr: usize) -> Self {
        let data = bindings::HuskSyntaxParseErrorData(ptr);
        let filename = {
            let string = CString::from_raw(data.r0).into_string().unwrap();
            if string.is_empty() {
                None
            } else {
                Some(string)
            }
        };
        let pos = Pos { ptr: data.r1 };
        let text = CString::from_raw(data.r2).into_string().unwrap();
        let incomplete = util::to_bool(data.r3);

        Self { ptr, filename, pos, text, incomplete }
    }
}

impl Drop for ParseError {
    #[husk_proc::unsafe_wrapper]
    fn drop(&mut self) {
        bindings::HuskMiscDeleteGoObj(self.ptr);
    }
}

impl fmt::Debug for ParseError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.debug_struct("ParseError")
            .field("filename", &self.filename)
            .field("pos", &self.pos)
            .field("text", &self.text)
            .field("incomplete", &self.incomplete)
            .finish()
    }
}

impl fmt::Display for ParseError {
    #[husk_proc::unsafe_wrapper]
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let error_str = bindings::HuskMiscErrorToString(self.ptr);
        let string = CString::from_raw(error_str).into_string().unwrap();
        write!(f, "{string}")
    }
}

impl Deref for ParseError {
    type Target = Pos;

    fn deref(&self) -> &Self::Target {
        &self.pos
    }
}

// Enums
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

    /// Automatically detect the language.
    Auto,
}

impl LangVariant {
    fn to_int(&self) -> i32 {
        match self {
            Self::Bash => ctypes::HUSK_LANG_BASH,
            Self::Posix => ctypes::HUSK_LANG_POSIX,
            Self::Mksh => ctypes::HUSK_LANG_MKSH,
            Self::Bats => ctypes::HUSK_LANG_BATS,
            Self::Auto => ctypes::HUSK_LANG_AUTO,
        }
    }
}

impl fmt::Display for LangVariant {
    #[husk_proc::unsafe_wrapper]
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let error_str = bindings::HuskSyntaxLangVariantString(self.to_int());
        let string = CString::from_raw(error_str).into_string().unwrap();
        write!(f, "{string}")
    }
}

// General types.
/// A position within a shell source file.
pub struct Pos {
    ptr: usize
}

impl Pos {
    /// Ceate a new position with the given `offset`, `line`, and `column`.
    #[husk_proc::unsafe_wrapper]
    pub fn new(offset: u32, line: u32, column: u32) -> Self {
        let ptr = bindings::HuskSyntaxNewPos(offset, line, column);
        Self { ptr }
    }

    /// Check whether `self` is after `p2`. It is a more expressive version of
    /// `self.offset() > p2.offset`.
    #[husk_proc::unsafe_wrapper]
    pub fn after(&self, p2: &Self) -> bool {
        let res = bindings::HuskSyntaxPosAfter(self.ptr, p2.ptr);
        util::to_bool(res)
    }

    /// Get the offset of this position, starting at 0. It counts in bytes.
    ///
    /// The returned number is protected against overflows; if an input is larger than 4GiB, the
    /// offset will wrap around to 0.
    #[husk_proc::unsafe_wrapper]
    pub fn offset(&self) -> u32 {
        bindings::HuskSyntaxPosOffset(self.ptr)
    }

    /// Get the line number of this position, starting at 1. It counts in bytes.
    ///
    /// The returned number is protected against overflows; if an input line has too many columns,
    /// extra columns will have a number of 0, rendered as `?` by `Self::to_string`.
    #[husk_proc::unsafe_wrapper]
    pub fn line(&self) -> u32 {
        bindings::HuskSyntaxPosLine(self.ptr)
    }

    /// Get the column number of this position, starting at 1. It counts in bytes.
    ///
    /// The returned number is protected against overflows; if an input line has too many columns,
    /// extra columns will have a number of 0, rendered as `?` by `Self::to_string`.
    #[husk_proc::unsafe_wrapper]
    pub fn col(&self) -> u32 {
        bindings::HuskSyntaxPosCol(self.ptr)
    }

    /// Check whether this [`Pos`] contains useful position information. Some positions returned
    /// via [`Parser::parse`] may be invalid: for example, [`Stmt::semicolon`] will only be valid
    /// if a statement contained a closing token suck as `;`.
    #[husk_proc::unsafe_wrapper]
    pub fn is_valid(&self) -> bool {
        let res = bindings::HuskSyntaxPosIsValid(self.ptr);
        util::to_bool(res)
    }
}

impl Drop for Pos {
    #[husk_proc::unsafe_wrapper]
    fn drop(&mut self) {
        bindings::HuskMiscDeleteGoObj(self.ptr);
    }
}

impl fmt::Debug for Pos {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        f.debug_struct("Pos")
            .field("offset", &self.offset())
            .field("line", &self.line())
            .field("col", &self.col())
            .finish()
    }
}

impl fmt::Display for Pos {
    #[husk_proc::unsafe_wrapper]
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let error_str = bindings::HuskMiscStringerToString(self.ptr);
        let string = CString::from_raw(error_str).into_string().unwrap();
        write!(f, "{string}")
    }
}

/// A builder to configure a new [`Parser`].
#[derive(Clone)]
pub struct ParserBuilder {
    keep_comments: Option<bool>,
    stop_at: Option<String>,
    lang_variant: Option<LangVariant>
}

impl ParserBuilder {
    fn new() -> Self {
        Self {
            keep_comments: None,
            stop_at: None,
            lang_variant: None,
        }
    }

    /// When set to `true`, makes the parse parse comments and attach them to [`Node`]s, as opposed to discarding them.
    pub fn keep_comments(mut self, enabled: bool) -> Self {
        self.keep_comments = Some(enabled);
        self
    }

    /// Configure the lexer to stop at an arbitrary word, treating it as if it were the end of the
    /// input. It can contain any characters except whitespace, and cannot be over four bytes in
    /// size.
    ///
    /// This can be useful to embed shell code within another language, as one can use a special
    /// word to mark the delimiters betwee the two.
    ///
    /// As a word, it will only apply when following whitespace or a separating token. For example,
    /// `$$` will act on the inputs `foo $$` and `foo;$$`, but not on `foo '$$'`.
    ///
    /// The match is done by prefix, so the example about will also act on `foo $$bar`.
    pub fn stop_at(mut self, word: &str) -> Self {
        self.stop_at = Some(word.to_owned());
        self
    }

    /// The shell language variant that the parser will accept.
    pub fn variant(mut self, lang: LangVariant) -> Self {
        self.lang_variant = Some(lang);
        self
    }

    /// Build this [`ParserBuilder`] into a [`Parser`].
    #[husk_proc::unsafe_wrapper]
    pub fn build(self) -> Parser {
        let mut keep_comments = self.keep_comments.map(|keep| util::to_u8(keep));
        let mut stop_at = self.stop_at.map(|stop| CString::new(stop).unwrap());
        let mut lang_int = self.lang_variant.map(|lang| lang.to_int());

        let keep_comments_ptr = match &mut keep_comments {
            Some(keep) => keep as *mut u8,
            None => ptr::null_mut(),
        };
        let stop_at_ptr = match &stop_at {
            Some(stop) => stop.as_ptr() as *mut ffi::c_char,
            None => ptr::null_mut(),
        };
        let lang_ptr = match &mut lang_int {
            Some(lang) => lang as *mut i32,
            None => ptr::null_mut()
        };

        let dptr = bindings::HuskSyntaxNewParser(keep_comments_ptr, stop_at_ptr, lang_ptr);
        Parser { dptr }
    }
}

/// The internal state of the parsing mechanism of a program.
pub struct Parser {
    dptr: usize,
}

impl Parser {
    /// Create a new [`Parser`] with the default configuration options.
    pub fn new() -> Self {
        Self::builder().build()
    }

    /// Create and configure a new [`Parser`].
    pub fn builder() -> ParserBuilder {
        ParserBuilder::new()
    }

    /// Read and parse a shell program, with an optional name. [`Error::Parser`] is returned if
    /// there was an issue parsing the program.
    #[husk_proc::unsafe_wrapper]
    pub fn parse<R: io::Read>(&self, content: &mut R, name: Option<&str>) -> crate::Result<()> {
        let bcontainer = Box::new(ctypes::ReadContainer(content));
        let bcontainer_ptr = Box::into_raw(bcontainer) as *mut ffi::c_void;
        let c_name = name.map(|name| CString::new(name).unwrap());
        let c_name_ptr = match &c_name {
            Some(name) => name.as_ptr() as *mut ffi::c_char,
            None => ptr::null_mut()
        };
        let data = bindings::HuskSyntaxParserParse(self.dptr, bcontainer_ptr, c_name_ptr);

        if data.r1 == ctypes::HUSK_SYNTAX_PARSER_FILE {
            Ok(())
        } else {
            Err(util::error_obj_to_rust(data.r0, data.r1))
        }
    }
}

impl Drop for Parser {
    #[husk_proc::unsafe_wrapper]
    fn drop(&mut self) {
        bindings::HuskMiscDeleteGoObj(self.dptr);
    }
}
