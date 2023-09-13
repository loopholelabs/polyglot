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
use kind::PolyglotKind;
use kind::PolyglotStatus;

use std::ffi::{c_char, c_uint, CStr};
use std::io::{Cursor, Write};
use polyglot_rs::{Encoder as PolyglotEncoder};

#[repr(C)]
#[derive(Debug)]
pub struct Encoder {
    cursor: Cursor<Vec<u8>>,
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn polyglot_new_encoder(encoder: *mut *mut Encoder) -> PolyglotStatus {
    if encoder.is_null() {
        return PolyglotStatus::NullPointer;
    }

    unsafe {
        *encoder = Box::into_raw(Box::new(Encoder {
            cursor: Cursor::new(Vec::new()),
        }));
    }

    PolyglotStatus::Pass
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn polyglot_encoder_size(encoder: *mut Encoder) -> c_uint {
    if encoder.is_null() {
        return 0;
    }

    unsafe {
        (*encoder).cursor.get_ref().len() as c_uint
    }
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn polyglot_encoder_buffer(encoder: *mut Encoder, buffer_pointer: *mut c_char, buffer_size: c_uint) -> PolyglotStatus {
    if encoder.is_null() {
        return PolyglotStatus::NullPointer;
    }

    if buffer_pointer.is_null() {
        return PolyglotStatus::NullPointer;
    }

    unsafe {
        if buffer_size < (*encoder).cursor.get_ref().len() as c_uint {
            return PolyglotStatus::Fail;
        }

        let mut buffer = std::slice::from_raw_parts_mut(buffer_pointer as *mut u8, buffer_size as usize);
        match buffer.write((*encoder).cursor.get_ref().as_slice()) {
            Ok(_) => PolyglotStatus::Pass,
            Err(_) => PolyglotStatus::Fail,
        }
    }
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn polyglot_free_encoder(encoder: *mut Encoder) {
    if !encoder.is_null() {
        unsafe {
            drop(Box::from_raw(encoder));
        }
    }
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn polyglot_encode_none(encoder: *mut Encoder) -> PolyglotStatus {
    if encoder.is_null() {
        return PolyglotStatus::NullPointer;
    }

    unsafe {
        match (*encoder).cursor.encode_none() {
            Ok(_) => PolyglotStatus::Pass,
            Err(_) => PolyglotStatus::Fail,
        }
    }
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn polyglot_encode_array(encoder: *mut Encoder, array_size: c_uint, array_kind: PolyglotKind) -> PolyglotStatus {
    if encoder.is_null() {
        return PolyglotStatus::NullPointer;
    }

    unsafe {
        match (*encoder).cursor.encode_array(array_size as usize, array_kind.into_polyglot()) {
            Ok(_) => PolyglotStatus::Pass,
            Err(_) => PolyglotStatus::Fail,
        }
    }
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn polyglot_encode_map(encoder: *mut Encoder, map_size: c_uint, key_kind: PolyglotKind, value_kind: PolyglotKind) -> PolyglotStatus {
    if encoder.is_null() {
        return PolyglotStatus::NullPointer;
    }

    unsafe {
        match (*encoder).cursor.encode_map(map_size as usize, key_kind.into_polyglot(), value_kind.into_polyglot()) {
            Ok(_) => PolyglotStatus::Pass,
            Err(_) => PolyglotStatus::Fail,
        }
    }
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn polyglot_encode_bytes(encoder: *mut Encoder, buffer_pointer: *mut c_char, buffer_size: c_uint) -> PolyglotStatus {
    if encoder.is_null() {
        return PolyglotStatus::NullPointer;
    }

    if buffer_pointer.is_null() {
        return PolyglotStatus::NullPointer;
    }

    unsafe {
        let buffer = std::slice::from_raw_parts_mut(buffer_pointer as *mut u8, buffer_size as usize);
        match (*encoder).cursor.encode_bytes(buffer) {
            Ok(_) => PolyglotStatus::Pass,
            Err(_) => PolyglotStatus::Fail,
        }
    }
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn polyglot_encode_string(encoder: *mut Encoder, string_pointer: *const c_char) -> PolyglotStatus {
    if encoder.is_null() {
        return PolyglotStatus::NullPointer;
    }

    unsafe {
        let string = match CStr::from_ptr(string_pointer).to_str() {
            Ok(string) => string,
            Err(_) => return PolyglotStatus::Fail,
        };
        match (*encoder).cursor.encode_str(string) {
            Ok(_) => PolyglotStatus::Pass,
            Err(_) => PolyglotStatus::Fail,
        }
    }
}

/*
pub trait Encoder {
    fn encode_str(self, val: &str) -> Result<Self, EncodingError>
    where
        Self: Sized;
    fn encode_error(self, val: Box<dyn Error>) -> Result<Self, EncodingError>
    where
        Self: Sized;
    fn encode_bool(self, val: bool) -> Result<Self, EncodingError>
    where
        Self: Sized;
    fn encode_u8(self, val: u8) -> Result<Self, EncodingError>
    where
        Self: Sized;
    fn encode_u16(self, val: u16) -> Result<Self, EncodingError>
    where
        Self: Sized;
    fn encode_u32(self, val: u32) -> Result<Self, EncodingError>
    where
        Self: Sized;
    fn encode_u64(self, val: u64) -> Result<Self, EncodingError>
    where
        Self: Sized;
    fn encode_i32(self, val: i32) -> Result<Self, EncodingError>
    where
        Self: Sized;
    fn encode_i64(self, val: i64) -> Result<Self, EncodingError>
    where
        Self: Sized;
    fn encode_f32(self, val: f32) -> Result<Self, EncodingError>
    where
        Self: Sized;
    fn encode_f64(self, val: f64) -> Result<Self, EncodingError>
    where
        Self: Sized;
}
 */