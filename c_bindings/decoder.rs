/*
    Copyright 2023 Loophole Labs

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

           http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

use crate::kind;
use kind::PolyglotStatus;

use std::ffi::{c_char, c_uint, CString};
use std::io::{Cursor, Read};
use polyglot_rs::{Decoder as PolyglotDecoder};

#[repr(C)]
#[derive(Debug)]
pub struct Decoder {
    cursor: Cursor<Vec<u8>>,
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn polyglot_new_decoder(decoder: *mut *mut Decoder, buffer_pointer: *mut c_char, buffer_size: c_uint) -> PolyglotStatus {
    if decoder.is_null() {
        return PolyglotStatus::NullPointer;
    }

    if buffer_pointer.is_null() {
        return PolyglotStatus::NullPointer;
    }

    unsafe {
        let buffer = std::slice::from_raw_parts_mut(buffer_pointer as *mut u8, buffer_size as usize);
        *decoder = Box::into_raw(Box::new(Decoder {
            cursor: Cursor::new(buffer.to_vec()),
        }));
    }

    PolyglotStatus::Pass
}

impl Decoder {
    fn decode_none(&mut self) -> bool {
        self.cursor.decode_none()
    }

    fn decode_array(&mut self, val_kind: polyglot_rs::Kind) -> Result<usize, polyglot_rs::DecodingError> {
        self.cursor.decode_array(val_kind)
    }
}