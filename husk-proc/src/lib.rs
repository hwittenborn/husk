extern crate proc_macro;
use proc_macro::TokenStream;
use quote::{quote, ToTokens};
use syn::{parse_macro_input, Expr, ExprUnsafe, ItemFn, Stmt, Token};

#[proc_macro_attribute]
pub fn unsafe_wrapper(_args: TokenStream, input: TokenStream) -> TokenStream {
    let mut input = parse_macro_input!(input as ItemFn);
    let mut fn_block = input.block.as_ref().to_owned();
    let unsafe_block = Expr::Unsafe(ExprUnsafe {
        attrs: vec![],
        unsafe_token: Default::default(),
        block: fn_block.clone(),
    });

    fn_block.stmts.clear();
    fn_block.stmts.push(Stmt::Expr(unsafe_block, None));
    input.block = Box::new(fn_block);

    TokenStream::from(quote!(#input))
}

#[proc_macro]
pub fn gen_ints(items: TokenStream) -> TokenStream {
    let mut num = 0;
    let mut code = String::new();

    for item in items {
        if item.to_string() == "," {
            continue;
        }

        let line = format!("#[no_mangle]\npub static {item}: ::std::ffi::c_int = {num};\n");
        code += &line;
        num += 1;
    }

    code.parse().unwrap()
}
