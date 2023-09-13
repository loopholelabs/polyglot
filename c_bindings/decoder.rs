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

use std::ffi::{c_char, c_uint, CString};
use std::io::Cursor;
use polyglot_rs::{Decoder as PolyglotDecoder};

#[repr(C)]
#[derive(Debug)]
pub struct Decoder {
    cursor: Cursor<Vec<u8>>,
}

#[no_mangle]
pub extern "C" fn polyglot_new_decoder(status: *mut PolyglotStatus, buffer_pointer: *mut c_char, buffer_size: c_uint) -> *mut Decoder {
    PolyglotStatus::check_not_null(status);

    if buffer_pointer.is_null() {
        unsafe {
            *status = PolyglotStatus::NullPointer;
        }
        return std::ptr::null_mut();
    }

    unsafe {
        let buffer = std::slice::from_raw_parts_mut(buffer_pointer as *mut u8, buffer_size as usize);
        *status = PolyglotStatus::Pass;
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
pub extern "C" fn polyglot_decode_none(status: *mut PolyglotStatus, decoder: *mut Decoder) -> bool {
    PolyglotStatus::check_not_null(status);

    if decoder.is_null() {
        unsafe {
            *status = PolyglotStatus::NullPointer;
        }
    }

    unsafe {
        (*decoder).cursor.decode_none()
    }
}

#[no_mangle]
pub extern "C" fn polyglot_decode_array(status: *mut PolyglotStatus, decoder: *mut Decoder, array_kind: PolyglotKind) -> c_uint {
    PolyglotStatus::check_not_null(status);

    if decoder.is_null() {
        unsafe {
            *status = PolyglotStatus::NullPointer;
        }
    }

    unsafe {
        match (*decoder).cursor.decode_array(array_kind.into()) {
            Ok(size ) => {
                *status = PolyglotStatus::Pass;
                size as c_uint
            },
            Err(_) => {
                *status = PolyglotStatus::Fail;
                0
            }
        }
    }
}





#[no_mangle]
pub extern "C" fn polyglot_decode_string(status: *mut PolyglotStatus, decoder: *mut Decoder) -> *mut c_char{
    PolyglotStatus::check_not_null(status);

    if decoder.is_null() {
        unsafe {
            *status = PolyglotStatus::NullPointer;
        }
        return std::ptr::null_mut();
    }

    unsafe {
        match (*decoder).cursor.decode_string() {
            Ok(value) => {
                return match CString::new(value) {
                    Ok(c_string) => {
                        *status = PolyglotStatus::Pass;
                        c_string.into_raw()
                    }
                    Err(_) => {
                        *status = PolyglotStatus::Fail;
                        std::ptr::null_mut()
                    }
                };
            },
            Err(_) => {
                *status = PolyglotStatus::Fail;
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

