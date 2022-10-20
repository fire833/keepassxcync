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

mod kdb;
mod kdbx;

// use std::fs::File;
use clap::{App, Arg, ArgMatches};
use sha1::{Digest, Sha1};
use sha2::{Sha256, Sha512};
use std::{
    fs::{self},
    io::Write,
    path::Path,
    process::exit,
};

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
                exit(1);
            }
            Ok(file) => {
                print_data(file, path);
            }
        },
        None => {
            println!("database file path not provided");
            exit(1);
        }
    }
}

/// Header signature for keepass databases. All files will start with this value.
const KEEPASSSIGNATURE1: [u8; 4] = [3, 217, 162, 154];

/// for .kdb files (KeePass 1.x file format)
const KDBSIGNATURE: [u8; 4] = [101, 251, 75, 181];

/// for kdbx file of KeePass 2.x pre-release (alpha & beta)
const KDBXSIGNATURE1: [u8; 4] = [102, 251, 75, 181];

/// for kdbx file of KeePass post-release
const KDBXSIGNATURE2: [u8; 4] = [103, 251, 75, 181];

trait KDBXInfo {
    fn parse(&mut self, data: &[u8]) -> Result<(), String>;

    fn format_info(&self) -> String;
}

fn print_data(file_data: Vec<u8>, file_name: &str) {
    let file_bytes = file_data.as_slice();
    let count: usize = file_bytes.len();
    let file_data: String;

    // Check to make sure this is a keepass database.
    if file_bytes[0..4] != KEEPASSSIGNATURE1 {
        println!(
            "file {} does not appear to be a keepass database, exiting",
            file_name
        );
        exit(1);
    }

    let sig2 = &file_bytes[4..8];

    if sig2 == KDBSIGNATURE {
        let mut kdbheader = kdb::new();
        match kdbheader.parse(&file_bytes[8..]) {
            Ok(()) => {
                file_data = kdbheader.format_info();
            }
            Err(error) => {
                println!("{}", error);
                exit(1);
            }
        }
    } else if sig2 == KDBXSIGNATURE1 || sig2 == KDBXSIGNATURE2 {
        let mut kdbxheader = kdbx::new();
        match kdbxheader.parse(&file_bytes[8..]) {
            Ok(()) => {
                file_data = kdbxheader.format_info();
            }
            Err(error) => {
                println!("{}", error);
                exit(1);
            }
        }
    } else {
        println!(
            "error in processing file: signature2 ({:?}) does not match any kdb|x signatures",
            sig2
        );
        exit(1);
    }

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

    print!(
        "Database file {} info:

File Size: {} bytes, or {} kilobytes
SHA1 Hash: {:x}
SHA256 Hash: {:x}
SHA512 Hash: {:x}

Metadata:
{}
",
        file_name,
        &count,
        &count / 1000,
        sha1.finalize(),
        sha256.finalize(),
        sha512.finalize(),
        file_data,
    );

    let _ = file_bytes[0..3];
}
