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
use types::{Kind, Status, StringError};

use std::ffi::{c_char, CStr};
use std::io::{Cursor, Write};
use polyglot_rs::{Encoder as PolyglotEncoder};

#[repr(C)]
#[derive(Debug)]
pub struct Encoder {
    cursor: Cursor<Vec<u8>>,
}

#[no_mangle]
pub extern "C" fn polyglot_new_encoder(status: *mut Status) -> *mut Encoder {
    Status::check_not_null(status);

    unsafe {
        *status = Status::Pass;
    }

    Box::into_raw(Box::new(Encoder {
        cursor: Cursor::new(Vec::new()),
    }))
}

#[no_mangle]
pub extern "C" fn polyglot_encoder_size(status: *mut Status, encoder: *mut Encoder) -> u32 {
    Status::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return 0;
    }

    unsafe {
        (*encoder).cursor.get_ref().len() as u32
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encoder_buffer(status: *mut Status, encoder: *mut Encoder, buffer_pointer: *mut u8, buffer_size: u32) {
    Status::check_not_null(status);

    if encoder.is_null() || buffer_pointer.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return;
    }

    unsafe {
        if buffer_size < (*encoder).cursor.get_ref().len() as u32 {
            *status = Status::Fail;
            return;
        }

        let mut buffer = std::slice::from_raw_parts_mut(buffer_pointer, buffer_size as usize);
        match buffer.write((*encoder).cursor.get_ref().as_slice()) {
            Ok(_) => {
                *status = Status::Pass;
            },
            Err(_) => {
                *status = Status::Fail;
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
pub extern "C" fn polyglot_encode_none(status: *mut Status, encoder: *mut Encoder) {
    Status::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = Status::NullPointer
        }
        return;
    }

    unsafe {
        match (*encoder).cursor.encode_none() {
            Ok(_) => {
                *status = Status::Pass
            },
            Err(_) => {
                *status = Status::Fail
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_array(status: *mut Status, encoder: *mut Encoder, array_size: u32, array_kind: Kind) {
    Status::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return;
    }

    unsafe {
        match (*encoder).cursor.encode_array(array_size as usize, array_kind.into()) {
            Ok(_) => {
                *status = Status::Pass
            },
            Err(_) => {
                *status = Status::Fail
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_map(status: *mut Status, encoder: *mut Encoder, map_size: u32, key_kind: Kind, value_kind: Kind) {
    Status::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return;
    }

    unsafe {
        match (*encoder).cursor.encode_map(map_size as usize, key_kind.into(), value_kind.into()) {
            Ok(_) => {
                *status = Status::Pass
            },
            Err(_) =>  {
                *status = Status::Fail;
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_bytes(status: *mut Status, encoder: *mut Encoder, buffer_pointer: *mut u8, buffer_size: u32) {
    Status::check_not_null(status);

    if encoder.is_null() || buffer_pointer.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return;
    }

    unsafe {
        let buffer = std::slice::from_raw_parts_mut(buffer_pointer, buffer_size as usize);
        match (*encoder).cursor.encode_bytes(buffer) {
            Ok(_) => {
                *status = Status::Pass
            },
            Err(_) =>  {
                *status = Status::Fail;
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_string(status: *mut Status, encoder: *mut Encoder, string_pointer: *const c_char) {
    Status::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return;
    }

    unsafe {
        let c_str = match CStr::from_ptr(string_pointer).to_str() {
            Ok(c_str) => c_str,
            Err(_) => {
                *status = Status::Fail;
                return;
            },
        };
        match (*encoder).cursor.encode_str(c_str) {
            Ok(_) => {
                *status = Status::Pass
            },
            Err(_) =>  {
                *status = Status::Fail;
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_error(status: *mut Status, encoder: *mut Encoder, string_pointer: *mut c_char) {
    Status::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return;
    }

    unsafe {
        let error_string = match CStr::from_ptr(string_pointer).to_str() {
            Ok(c_str) => Box::new(StringError(c_str.to_string())),
            Err(_) => {
                *status = Status::Fail;
                return;
            },
        };
        match (*encoder).cursor.encode_error(error_string) {
            Ok(_) => {
                *status = Status::Pass
            },
            Err(_) =>  {
                *status = Status::Fail;
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_bool(status: *mut Status, encoder: *mut Encoder, val: bool) {
    Status::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return;
    }

    unsafe {
        match (*encoder).cursor.encode_bool(val) {
            Ok(_) => {
                *status = Status::Pass
            },
            Err(_) =>  {
                *status = Status::Fail;
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_u8(status: *mut Status, encoder: *mut Encoder, value: u8) {
    Status::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return;
    }

    unsafe {
        match (*encoder).cursor.encode_u8(value) {
            Ok(_) => {
                *status = Status::Pass
            },
            Err(_) =>  {
                *status = Status::Fail;
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_u16(status: *mut Status, encoder: *mut Encoder, value: u16) {
    Status::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return;
    }

    unsafe {
        match (*encoder).cursor.encode_u16(value) {
            Ok(_) => {
                *status = Status::Pass
            },
            Err(_) =>  {
                *status = Status::Fail;
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_u32(status: *mut Status, encoder: *mut Encoder, value: u32) {
    Status::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return;
    }

    unsafe {
        match (*encoder).cursor.encode_u32(value) {
            Ok(_) => {
                *status = Status::Pass
            },
            Err(_) =>  {
                *status = Status::Fail;
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_u64(status: *mut Status, encoder: *mut Encoder, value: u64) {
    Status::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return;
    }

    unsafe {
        match (*encoder).cursor.encode_u64(value) {
            Ok(_) => {
                *status = Status::Pass
            },
            Err(_) =>  {
                *status = Status::Fail;
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_i32(status: *mut Status, encoder: *mut Encoder, value: i32) {
    Status::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return;
    }

    unsafe {
        match (*encoder).cursor.encode_i32(value) {
            Ok(_) => {
                *status = Status::Pass
            },
            Err(_) =>  {
                *status = Status::Fail;
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_i64(status: *mut Status, encoder: *mut Encoder, value: i64) {
    Status::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return;
    }

    unsafe {
        match (*encoder).cursor.encode_i64(value) {
            Ok(_) => {
                *status = Status::Pass
            },
            Err(_) =>  {
                *status = Status::Fail;
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_f32(status: *mut Status, encoder: *mut Encoder, value: f32) {
    Status::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return;
    }

    unsafe {
        match (*encoder).cursor.encode_f32(value) {
            Ok(_) => {
                *status = Status::Pass
            },
            Err(_) =>  {
                *status = Status::Fail;
            },
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_encode_f64(status: *mut Status, encoder: *mut Encoder, value: f64) {
    Status::check_not_null(status);

    if encoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return;
    }

    unsafe {
        match (*encoder).cursor.encode_f64(value) {
            Ok(_) => {
                *status = Status::Pass
            },
            Err(_) =>  {
                *status = Status::Fail;
            },
        }
    }
}