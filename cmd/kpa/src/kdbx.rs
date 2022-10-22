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

use std::io::BufReader;

use crate::KDBXInfo;

const AESCIPHERID: [u8; 16] = [
    0x31, 0xC1, 0xF2, 0xE6, 0xBF, 0x71, 0x43, 0x50, 0xBE, 0x58, 0x05, 0x21, 0x6A, 0xFC, 0x5A, 0xFF,
];

pub struct KDBXHeader {
    db_comment: String,
    transform_rounds: u64,
}

pub fn new() -> KDBXHeader {
    KDBXHeader {
        db_comment: String::from(""),
        transform_rounds: 0,
    }
}

impl KDBXInfo for KDBXHeader {
    fn parse(&mut self, data: &[u8]) -> Result<(), String> {
        let mut counter: usize = 0;

        loop {
            let type_bytes: u8 = u8::from_ne_bytes([data[counter]]);
            let len_bytes: u16 = u16::from_le_bytes([
                data[counter + 2],
                data[counter + 3],
                /*data[counter+4], data[counter+5]*/
            ]);
            counter += 3;
            // vec!();
            match type_bytes {
                0x00 => {
                    if len_bytes != 2 {}

                    println!("found eoh: {}", len_bytes);
                    return Ok(());
                }
                // (? bytes) Comment : this field doesn’t seem to be used in KeePass source code, but you may encounter it.
                // Maybe its value is an ASCII string.`
                0x01 => {
                    println!("found comment: {}", len_bytes);
                }
                // (16 bytes) Cipher ID : the UUID identifying the cipher. KeePass 2 has an implementation allowing the use
                // of other cipher, but only AES with a UUID of [0x31, 0xC1, 0xF2, 0xE6, 0xBF, 0x71, 0x43, 0x50, 0xBE, 0x58,
                // 0x05, 0x21, 0x6A, 0xFC, 0x5A, 0xFF] is implemented in KeePass 2 for now.
                0x02 => {
                    if len_bytes != 16 {}
                    println!("found cipher ID: {}", len_bytes);
                }
                // (4 bytes) Compression Flags : give the compression algorithm used to compress the database. For now only
                // 2 possible value can be set : None (Value : 0) and GZip (Value : 1). For now, the Value of this field
                // should not be greater or equal to 2. The 4 bytes of the Value should be convert to a 32 bit signed integer
                // before comparing it to known values,
                0x03 => {
                    if len_bytes != 4 {}
                    println!("found compression flags: {}", len_bytes);
                }
                // (16 bytes) Master Seed : The salt that will be concatenated to the transformed master key then hashed,
                // to create the final master key,
                0x04 => {
                    if len_bytes != 16 {}

                    println!("found master seed: {}", len_bytes);
                }
                // (32 bytes) Transform Seed : The key used by AES as seed to generate the transformed master key,
                0x05 => {
                    if len_bytes != 32 {}

                    println!("found transform seed: {}", len_bytes);
                }
                // (8 bytes) Transform Rounds : The number of Rounds you have to compute the transformed master key,
                0x06 => {
                    if len_bytes != 8 {}

                    println!("found transform rounds: {}", len_bytes);

                    self.transform_rounds = u64::from_ne_bytes([
                        data[counter],
                        data[counter + 1],
                        data[counter + 2],
                        data[counter + 3],
                        data[counter + 4],
                        data[counter + 5],
                        data[counter + 6],
                        data[counter + 7],
                    ]);
                    continue;
                }
                // (? bytes ) Encryption IV : The IV used by the cipher that encrypted the database,
                0x07 => {
                    println!("found encryption IV: {}", len_bytes);
                }
                // (? bytes ) Protected Stream Key : The key/seed for the cipher used to encrypt the password of an entry
                // in the database (see later),
                0x08 => {
                    println!("found stream key: {}", len_bytes);
                }
                // (32 bytes) Stream Start Bytes, indicates the first 32 unencrypted bytes of the database part of the file
                // (to check if the file is corrupt, or the key correct,etc). These 32 bytes should have been randomly
                // generated when the file was saved. Length should be 32 bytes ,
                0x09 => {
                    if len_bytes != 32 {}
                    println!("found stream start bytes: {}", len_bytes);
                }
                // (4 bytes) Inner Random Stream ID : the ID of the cipher used to encrypted the password of an entry in the
                // database (see later), for now you can expect to have : 0 for nothing, 1 for ARC4, 2 for Salsa20.
                0x10 => {
                    if len_bytes != 4 {}
                    println!("found Inner Random Stream ID: {}", len_bytes);
                }
                (g) => {
                    println!("found end {}, bytes: {}", g, len_bytes);
                }
            }

            let bytes_i32 = len_bytes as u16;
            counter += usize::from(bytes_i32);
        }
    }

    fn format_info(&self) -> String {
        format!("")
    }
}
