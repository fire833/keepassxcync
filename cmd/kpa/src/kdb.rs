/*
*	Copyright (C) 2022  Kendall Tauser
*
*	This program is free software; you can redistribute it and/or modify
*	it under the terms of the GNU General Public License as published by
*	the Free Software Foundation; either version 2 of the License, or
*	(at your option) any later version.
*
*	This program is distributed in the hope that it will be u
    flags: [u8; 4],

    version: [u8; 4],seful,
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
const KDBFRONTHEADER: usize = 4 + 4;

/// Information stored within keepass (kdb) header file. This will be a constant set of bytes,
/// so the struct is a lot simpler than the TLV parser setup within KDBXHeader struct.
pub struct KDBHeader {
    flags: [u8; 4],
    version: [u8; 4],
    master_seed: [u8; 16],
    encr_iv: [u8; 16],
    num_groups: u32,
    num_entries: u32,
    sha256: String,
    transform_seed: String,
    encr_rounds: u32,
}

pub fn new() -> KDBHeader {
    KDBHeader {
        flags: [0; 4],
        version: [0; 4],
        master_seed: [0; 16],
        encr_iv: [0; 16],
        num_groups: 0,
        num_entries: 0,
        sha256: String::from(""),
        transform_seed: String::from(""),
        encr_rounds: 0,
    }
}

impl KDBXInfo for KDBHeader {
    fn parse(&mut self, data: &[u8]) -> Result<(), String> {
        if data.len() < (KDB1HEADERSIZE + KDBFRONTHEADER) {
            return Err(format!(
                "header size is not greater than {} bytes, exiting",
                KDB1HEADERSIZE + KDBFRONTHEADER
            ));
        }

        // I know this is inefficient, will figure out another way to do this eventually.
        self.flags = [data[0], data[1], data[2], data[3]];
        self.version = [data[4], data[5], data[6], data[7]];
        self.master_seed = [
            data[8], data[9], data[10], data[12], data[13], data[14], data[15], data[16], data[17],
            data[18], data[19], data[20], data[21], data[22], data[23], data[24],
        ];
        self.encr_iv = [
            data[25], data[26], data[27], data[28], data[29], data[30], data[31], data[32],
            data[33], data[34], data[35], data[36], data[37], data[38], data[39], data[40],
        ];

        self.num_groups = u32::from_ne_bytes([data[41], data[42], data[43], data[44]]);
        self.num_entries = u32::from_ne_bytes([data[45], data[46], data[47], data[48]]);

        self.sha256 = format!(
            "{:x?}",
            [
                // Yes, pretty inefficient, I know.
                data[49], data[50], data[51], data[52], data[53], data[54], data[55], data[56],
                data[57], data[58], data[59], data[60], data[61], data[62], data[63], data[64],
                data[65], data[66], data[67], data[68], data[69], data[70], data[71], data[72],
                data[73], data[74], data[75], data[76], data[77], data[78], data[79], data[80],
            ]
        );

        self.transform_seed = format!(
            "{:x?}",
            [
                // Yes, pretty inefficient, I know.
                data[81], data[82], data[83], data[84], data[85], data[86], data[87], data[88],
                data[89], data[90], data[91], data[92], data[93], data[94], data[95], data[96],
                data[97], data[98], data[99], data[101], data[102], data[103], data[104],
                data[105], data[106], data[107], data[108], data[109], data[110], data[111],
                data[112], data[113],
            ]
        );

        self.encr_rounds = u32::from_ne_bytes([data[114], data[115], data[116], data[117]]);

        Ok(())
    }

    fn format_info(&self) -> String {
        format!(
            "
Flags: {:x?}
Version: {:?}
Master Seed: {:x?}
Encryption IV: {:x?}

Number of groups: {}
Number of entries: {}
Number of encryption rounds: {}

Database XML Hash: {}

Transform Key: {}

",
            self.flags,
            self.version,
            self.master_seed,
            self.encr_iv,
            self.num_groups,
            self.num_entries,
            self.encr_rounds,
            self.sha256,
            self.transform_seed,
        )
    }
}
