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

use crate::KDBXInfo;

const KDB1HEADERSIZE: usize = 4 + 4 + 16 + 16 + 4 + 4 + 32 + 32 + 4;

/// Information stored within keepass (kdb) header file. This will be a constant set of bytes,
/// so the struct is a lot simpler than the TLV parser setup within KDBXHeader struct.
struct KDBHeader {
    flags: [u8; 4],

    version: [u8; 4],
}

pub fn new() -> KDBHeader {
    KDBHeader {
        flags: [0; 4],
        version: [0; 4],
    };
}

impl KDBXInfo for KDBHeader {
    fn parse(&mut self, data: &[u8]) -> Result<u8, std::error::Error> {}

    fn format_info(&self) -> String {}
}
