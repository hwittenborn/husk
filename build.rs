use std::path::PathBuf;
use std::{env, fs, process::Command};

fn main() {
    println!("cargo:rerun-if-changed=src/go");
    env::set_current_dir("src/go").unwrap();
    let out_dir = env::var("OUT_DIR").unwrap();
    let out_path = PathBuf::from(&out_dir);

    // The docs.rs builder blocks network access, which would require vendoring everything. I don't
    // want to go through that hassle right now, so just generate a dummy file to allow the docs to
    // build.
    if env::var("DOCS_RS").is_ok() {
        fs::write(out_path.join("bindings.rs"), "").unwrap();
        return;
    }

    // Copy our 'husk.h' file into `out_dir`, as `bindgen` needs it to generate stuff.
    fs::copy("husk.h", out_path.join("husk.h")).unwrap();

    let status = Command::new("go")
        .args(["build", "-buildmode=c-archive", "-o", "libhusk.a", "."])
        .status()
        .expect("`go build` failed. Is `go` installed and on the latest version?");

    if !status.success() {
        panic!("`go build` failed. Is `go` on the latest version?");
    }

    println!("cargo:rustc-link-search={out_dir}");
    // We need to remove this before merge, it's working to test Koca for some reason though ???
    println!("cargo:rustc-link-search=/home/hunter/Documents/Git/GitHub/hwittenborn/husk/src/go/");
    println!("cargo:rustc-link-lib=static=husk");

    bindgen::Builder::default()
        .header("libhusk.h")
        .allowlist_type("^Husk.*")
        .allowlist_function("^Husk.*")
        .blocklist_function("^HuskRust.*")
        .generate()
        .unwrap()
        .write_to_file(out_path.join("bindings.rs"))
        .unwrap();
}
