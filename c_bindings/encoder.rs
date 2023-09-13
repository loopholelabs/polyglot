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

#![allow(clippy::not_unsafe_ptr_arg_deref)]

use crate::types;
use types::PolyglotKind;
use types::PolyglotStatus;

use std::ffi::{c_char, c_uint, CStr};
use std::io::{Cursor, Write};
use polyglot_rs::{Encoder as PolyglotEncoder};

#[repr(C)]
#[derive(Debug)]
pub struct Encoder {
    cursor: Cursor<Vec<u8>>,
}

#[no_mangle]
pub extern "C" fn polyglot_new_encoder(status: *mut PolyglotStatus) -> *mut Encoder {
    PolyglotStatus::check_not_null(status);

    unsafe {
        *status = PolyglotStatus::Pass;
    }

    Box::into_raw(Box::new(Encoder {
        cursor: Cursor::new(Vec::new()),
    }))
}

#[no_mangle]
pub extern "C" fn polyglot_encoder_size(status: *mut PolyglotStatus, encoder: *mut Encoder) -> c_uint {
    PolyglotStatus::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = PolyglotStatus::NullPointer;
        }
        return 0;
    }

    unsafe {
        (*encoder).cursor.get_ref().len() as c_uint
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encoder_buffer(status: *mut PolyglotStatus, encoder: *mut Encoder, buffer_pointer: *mut c_char, buffer_size: c_uint) {
    PolyglotStatus::check_not_null(status);

    if encoder.is_null() || buffer_pointer.is_null() {
        unsafe {
            *status = PolyglotStatus::NullPointer;
        }
        return;
    }

    unsafe {
        if buffer_size < (*encoder).cursor.get_ref().len() as c_uint {
            *status = PolyglotStatus::Fail;
            return;
        }

        let mut buffer = std::slice::from_raw_parts_mut(buffer_pointer as *mut u8, buffer_size as usize);
        match buffer.write((*encoder).cursor.get_ref().as_slice()) {
            Ok(_) => {
                *status = PolyglotStatus::Pass;
            },
            Err(_) => {
                *status = PolyglotStatus::Fail;
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_free_encoder(encoder: *mut Encoder) {
    if !encoder.is_null() {
        unsafe {
            drop(Box::from_raw(encoder));
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_none(status: *mut PolyglotStatus, encoder: *mut Encoder) {
    PolyglotStatus::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = PolyglotStatus::NullPointer
        }
        return;
    }

    unsafe {
        match (*encoder).cursor.encode_none() {
            Ok(_) => {
                *status = PolyglotStatus::Pass
            },
            Err(_) => {
                *status = PolyglotStatus::Fail
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_array(status: *mut PolyglotStatus, encoder: *mut Encoder, array_size: c_uint, array_kind: PolyglotKind) {
    PolyglotStatus::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = PolyglotStatus::NullPointer;
        }
        return;
    }

    unsafe {
        match (*encoder).cursor.encode_array(array_size as usize, array_kind.into()) {
            Ok(_) => {
                *status = PolyglotStatus::Pass
            },
            Err(_) => {
                *status = PolyglotStatus::Fail
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_map(status: *mut PolyglotStatus, encoder: *mut Encoder, map_size: c_uint, key_kind: PolyglotKind, value_kind: PolyglotKind) {
    PolyglotStatus::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = PolyglotStatus::NullPointer;
        }
        return;
    }

    unsafe {
        match (*encoder).cursor.encode_map(map_size as usize, key_kind.into(), value_kind.into()) {
            Ok(_) => {
                *status = PolyglotStatus::Pass
            },
            Err(_) =>  {
                *status = PolyglotStatus::Fail;
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_bytes(status: *mut PolyglotStatus, encoder: *mut Encoder, buffer_pointer: *mut c_char, buffer_size: c_uint) {
    PolyglotStatus::check_not_null(status);

    if encoder.is_null() || buffer_pointer.is_null() {
        unsafe {
            *status = PolyglotStatus::NullPointer;
        }
        return;
    }

    unsafe {
        let buffer = std::slice::from_raw_parts_mut(buffer_pointer as *mut u8, buffer_size as usize);
        match (*encoder).cursor.encode_bytes(buffer) {
            Ok(_) => {
                *status = PolyglotStatus::Pass
            },
            Err(_) =>  {
                *status = PolyglotStatus::Fail;
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_string(status: *mut PolyglotStatus, encoder: *mut Encoder, string_pointer: *const c_char) {
    PolyglotStatus::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = PolyglotStatus::NullPointer;
        }
        return;
    }

    unsafe {
        let c_str = match CStr::from_ptr(string_pointer).to_str() {
            Ok(c_str) => c_str,
            Err(_) => {
                *status = PolyglotStatus::Fail;
                return;
            },
        };
        match (*encoder).cursor.encode_str(c_str) {
            Ok(_) => {
                *status = PolyglotStatus::Pass
            },
            Err(_) =>  {
                *status = PolyglotStatus::Fail;
            },
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