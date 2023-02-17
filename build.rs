use std::{env, process::Command};

fn main() {
    println!("cargo:rerun-if-changed=src/go");
    env::set_current_dir("src/go").unwrap();
    let out_dir = env::var("OUT_DIR").unwrap();

    let status = Command::new("go")
        .args(["build", "-buildmode=c-archive", "-o"])
        .arg(&format!("{out_dir}/libhusk.a"))
        .arg("main.go")
        .status()
        .expect("`go build` failed. Is `go` installed and on the latest version?");

    if !status.success() {
        panic!("`go build` failed. Is `go` on the latest version?");
    }

    env::set_current_dir("../../").unwrap();
    println!("cargo:rustc-link-search=native={out_dir}");
    println!("cargo:rustc-link-lib=static=husk");

    let bindings = bindgen::Builder::default()
        .header(format!("{out_dir}/libhusk.h"))
        .allowlist_function("HuskDeleteGoItem")
        .allowlist_function("HuskGetCStringFromArray")
        .allowlist_function("HuskShellExpand")
        .allowlist_function("HuskShellFields")
        .allowlist_function("HuskSyntaxIsKeyword")
        .allowlist_function("HuskSyntaxQuote")
        .allowlist_function("HuskSyntaxValidName")
        /*.allowlist_var("HUSK_LANG_VARIANT_BASH")
        .allowlist_var("HUSK_LANG_VARIANT_POSIX")
        .allowlist_var("HUSK_LANG_VARIANT_MKSH")
        .allowlist_var("HUSK_LANG_VARIANT_BATS")*/
        .generate()
        .unwrap();
    bindings.write_to_file("src/bindings.rs").unwrap();
}
