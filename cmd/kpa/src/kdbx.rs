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

pub struct KDBXHeader {}

pub fn new() -> KDBXHeader {
    KDBXHeader {}
}

impl KDBXInfo for KDBXHeader {
    fn parse(&mut self, data: &[u8]) -> Result<(), String> {
        loop {
            let typeBytes: u16 = i16::from_ne_bytes([data[0], data[1]]) as u16;
            let lenBytes: u16 = i16::from_ne_bytes([data[1], data[2]]) as u16;
            // vec!();
            match typeBytes {
                // (? bytes) Comment : this field doesn’t seem to be used in KeePass source code, but you may encounter it.
                // Maybe its value is an ASCII string.
                1 => {}
                // (16 bytes) Cipher ID : the UUID identifying the cipher. KeePass 2 has an implementation allowing the use
                // of other cipher, but only AES with a UUID of [0x31, 0xC1, 0xF2, 0xE6, 0xBF, 0x71, 0x43, 0x50, 0xBE, 0x58,
                // 0x05, 0x21, 0x6A, 0xFC, 0x5A, 0xFF] is implemented in KeePass 2 for now.
                2 => {}
                // (4 bytes) Compression Flags : give the compression algorithm used to compress the database. For now only
                // 2 possible value can be set : None (Value : 0) and GZip (Value : 1). For now, the Value of this field
                // should not be greater or equal to 2. The 4 bytes of the Value should be convert to a 32 bit signed integer
                // before comparing it to known values,
                3 => {}
                // (16 bytes) Master Seed : The salt that will be concatenated to the transformed master key then hashed, 
                // to create the final master key,
                4 => {}
                // (32 bytes) Transform Seed : The key used by AES as seed to generate the transformed master key,
                5 => {}
                // (8 bytes) Transform Rounds : The number of Rounds you have to compute the transformed master key,
                6 => {}
                // (? bytes ) Encryption IV : The IV used by the cipher that encrypted the database,
                7 => {}
                // (? bytes ) Protected Stream Key : The key/seed for the cipher used to encrypt the password of an entry 
                // in the database (see later),
                8 => {}
                // (32 bytes) Stream Start Bytes, indicates the first 32 unencrypted bytes of the database part of the file 
                // (to check if the file is corrupt, or the key correct,etc). These 32 bytes should have been randomly 
                // generated when the file was saved. Length should be 32 bytes ,
                9 => {}
                // (4 bytes) Inner Random Stream ID : the ID of the cipher used to encrypted the password of an entry in the 
                // database (see later), for now you can expect to have : 0 for nothing, 1 for ARC4, 2 for Salsa20.
                10 => {}
                _ => {
                    return Ok(());
                }
            }
        }

        Ok(())
    }

    fn format_info(&self) -> String {
        format!("")
    }
}
