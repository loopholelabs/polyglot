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

use std::io::Cursor;
use crate::PolyglotStatus::{NullPointer, Pass};

#[repr(u8)]
pub enum PolyglotStatus {
    Pass,
    Fail,
    NullPointer,
}

pub struct Encoder {
    cursor: Cursor<Vec<u8>>,
}

impl Encoder {
    pub fn new() -> Self {
        Self {
            cursor: Cursor::new(Vec::new()),
        }
    }
}

#[allow(clippy::not_unsafe_ptr_arg_deref)]
#[no_mangle]
pub extern "C" fn encoder_new(encoder: *mut *mut Encoder) -> PolyglotStatus {
    if encoder.is_null() {
        return PolyglotStatus::NullPointer;
    }

    unsafe {
        *encoder = std::ptr::null_mut();
    }

    unsafe {
        *encoder = Box::into_raw(Box::new(Encoder::new()));
    }

    PolyglotStatus::Pass
}

#[no_mangle]
pub extern "C" fn encoder_free(encoder: *mut Encoder) {
    if !encoder.is_null() {
        unsafe {
            drop(Box::from_raw(encoder));
        }
    }
}