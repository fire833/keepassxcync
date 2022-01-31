use std::fs;
use std::fs::File;
use std::path::Path;

fn main() {

    let f = fs::read_to_string(Path::new(""));
    match f {
        Err(error) => {
            println!("Unable to read file {}", error);
            std::process::exit(1);
        }
        Ok(file) => {
            println!("it works! {}", file);
        }
    }
}

struct KDB {
    keepass_flags: u32,
    keepass_version: u32,
    master_seed: Vec<u8>, // 16 bytes
    encryption_i: Vec<u8>, // 16 bytes
    num_groups: u32,
    num_entries: u32,
    hash: Vec<u8>, // 32 byte long SHA hash.
    transform_seed: Vec<u8> // 32 bytes
}

impl KDB {

fn isFileKDB(input: std::string::String) -> bool {
    return false;
}

}

struct KDBX {
    
}

impl KDBX {
    
fn new() -> KDBX {
    return KDBX{};
}

}