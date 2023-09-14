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
use types::{Kind, Status, StringError, Buffer};

use std::ffi::{c_char, CString};
use std::io::Cursor;
use polyglot_rs::{Decoder as PolyglotDecoder};

#[repr(C)]
#[derive(Debug)]
pub struct Decoder {
    cursor: Cursor<Vec<u8>>,
}

#[no_mangle]
pub extern "C" fn polyglot_new_decoder(status: *mut Status, buffer_pointer: *mut u8, buffer_size: u32) -> *mut Decoder {
    Status::check_not_null(status);

    if buffer_pointer.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return std::ptr::null_mut();
    }

    unsafe {
        let buffer = std::slice::from_raw_parts_mut(buffer_pointer, buffer_size as usize);
        *status = Status::Pass;
        Box::into_raw(Box::new(Decoder {
            cursor: Cursor::new(buffer.to_vec()),
        }))
    }
}

#[no_mangle]
pub extern "C" fn polyglot_free_decoder(decoder: *mut Decoder) {
    if !decoder.is_null() {
        unsafe {
            drop(Box::from_raw(decoder));
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_decode_none(status: *mut Status, decoder: *mut Decoder) -> bool {
    Status::check_not_null(status);

    if decoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
    }

    unsafe {
        (*decoder).cursor.decode_none()
    }
}

#[no_mangle]
pub extern "C" fn polyglot_decode_array(status: *mut Status, decoder: *mut Decoder, array_kind: Kind) -> u32 {
    Status::check_not_null(status);

    if decoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
    }

    unsafe {
        match (*decoder).cursor.decode_array(array_kind.into()) {
            Ok(size ) => {
                *status = Status::Pass;
                size as u32
            },
            Err(_) => {
                *status = Status::Fail;
                0
            }
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_decode_map(status: *mut Status, decoder: *mut Decoder, key_kind: Kind, value_kind: Kind) -> u32 {
    Status::check_not_null(status);

    if decoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
    }

    unsafe {
        match (*decoder).cursor.decode_map(key_kind.into(), value_kind.into()) {
            Ok(size ) => {
                *status = Status::Pass;
                size as u32
            },
            Err(_) => {
                *status = Status::Fail;
                0
            }
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_decode_bytes(status: *mut Status, decoder: *mut Decoder) -> *mut Buffer {
    Status::check_not_null(status);

    if decoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
    }

    unsafe {
        match (*decoder).cursor.decode_bytes() {
            Ok(value ) => {
                *status = Status::Pass;
                let mut boxed_value = value.into_boxed_slice();
                let buffer = Buffer::new_raw(boxed_value.as_mut_ptr(), boxed_value.len() as u32);
                std::mem::forget(boxed_value);
                buffer
            },
            Err(_) => {
                *status = Status::Fail;
                std::ptr::null_mut()
            }
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_free_decode_bytes(buffer: *mut Buffer) {
    if !buffer.is_null() {
        unsafe {
            let boxed_value = unsafe { std::slice::from_raw_parts_mut((*buffer).data, (*buffer).length as usize) };
            let value = boxed_value.as_mut_ptr();
            drop(Box::from_raw(value));
            drop(Box::from_raw(buffer));
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_decode_string(status: *mut Status, decoder: *mut Decoder) -> *mut c_char{
    Status::check_not_null(status);

    if decoder.is_null() {
        unsafe {
            *status = Status::NullPointer;
        }
        return std::ptr::null_mut();
    }

    unsafe {
        match (*decoder).cursor.decode_string() {
            Ok(value) => {
                return match CString::new(value) {
                    Ok(c_string) => {
                        *status = Status::Pass;
                        c_string.into_raw()
                    }
                    Err(_) => {
                        *status = Status::Fail;
                        std::ptr::null_mut()
                    }
                };
            },
            Err(_) => {
                *status = Status::Fail;
                std::ptr::null_mut()
            }
        }
    }
}

#[no_mangle]
pub extern "C" fn polyglot_free_decode_string(c_string: *mut c_char) {
    unsafe {
        if !c_string.is_null() {
            drop(CString::from_raw(c_string))
        }
    };
}

