/*
*	Copyright (C) 2022  Kendall Tauser
*
*	This program is free software; you can redistribute it and/or modify
*	it under the terms of the GNU General Public License as published by
*	the Free Software Foundation; either version 2 of the License, or
*	(at your option) any later version.
*
*	This program is distributed in the hope that it will be useful,
*	but WITHOUT ANY WARRANTY; without even the implied warranty of
*	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*	GNU General Public License for more details.
*
*	You should have received a copy of the GNU General Public License along
*	with this program; if not, write to the Free Software Foundation, Inc.,
*	51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

// use std::fs::File;
use clap::{App, Arg, ArgMatches};
use std::{path::Path, fs::{self}, io::{Write}};
use sha1::{Sha1, Digest};
use sha2::{Sha256, Sha512};

fn main() {
    let cli: App = App::new("kpa")
        .about("A simple CLI tool for analysing the headers of a keepass database.")
        .author("Kendall Tauser <kttpsy@gmail.com>")
        .arg(
            Arg::new("file")
                .default_value("passwords.kdbx")
                .alias("f")
                // .required(true)
                .help("Specify the file name to be read in and analysed."),
        );

    let args: ArgMatches = cli.get_matches();

    match args.value_of("file") {
        Some(path) => match fs::read(Path::new(path)) {
            Err(error) => {
                println!("Unable to read file, error: {}\n", error);
                std::process::exit(1);
            }
            Ok(file) => {
                print_data(file, path);
            }
        },
        None => {
            println!("database file path not provided");
            std::process::exit(1);
        }
    }
}

fn print_data(file_data: Vec<u8>, file_name: &str) {
    let file_bytes = file_data.as_slice();
    let count: usize = file_bytes.len();

    let mut sha1 = Sha1::new();
    let mut sha256 = Sha256::new();
    let mut sha512 = Sha512::new();
    
    match sha1.write(&file_bytes) {
        Ok(_) => {}
        Err(_) => {}
    }
    match sha256.write(&file_bytes) {
        Ok(_) => {}
        Err(_) => {}
    }
    match sha512.write(&file_bytes) {
        Ok(_) => {}
        Err(_) => {}
    }

    print!("Database file {} info:

File Size: {} bytes, or {} kilobytes
SHA1 Hash: {:x}
SHA256 Hash: {:x}
SHA512 Hash: {:x}", file_name, 
    &count, &count / 1000, 
    sha1.finalize(), 
    sha256.finalize(),
    sha512.finalize(),
);

    // let mut i = 0;
    // for byte in file_bytes {

    // }

}
