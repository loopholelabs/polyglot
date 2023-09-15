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

use std::fmt::{Debug, Display, Formatter};

#[repr(u8)]
#[derive(Debug)]
pub enum Status {
    Pass,
    Fail,
    NullPointer,
}

impl Status {
    pub fn check_not_null(status: *mut Status) {
        if status.is_null() {
            panic!("status pointer is null");
        }
    }
}

#[repr(u8)]
#[derive(Debug)]
pub enum Kind {
    None = 0x00,
    Array = 0x01,
    Map = 0x02,
    Any = 0x03,
    Bytes = 0x04,
    String = 0x05,
    Error = 0x06,
    Bool = 0x07,
    U8 = 0x08,
    U16 = 0x09,
    U32 = 0x0a,
    U64 = 0x0b,
    I32 = 0x0c,
    I64 = 0x0d,
    F32 = 0x0e,
    F64 = 0x0f,
    Unknown,
}

impl Kind {
    pub fn into(self) -> polyglot_rs::Kind {
        match self {
            Kind::None => polyglot_rs::Kind::None,
            Kind::Array => polyglot_rs::Kind::Array,
            Kind::Map => polyglot_rs::Kind::Map,
            Kind::Any => polyglot_rs::Kind::Any,
            Kind::Bytes => polyglot_rs::Kind::Bytes,
            Kind::String => polyglot_rs::Kind::String,
            Kind::Error => polyglot_rs::Kind::Error,
            Kind::Bool => polyglot_rs::Kind::Bool,
            Kind::U8 => polyglot_rs::Kind::U8,
            Kind::U16 => polyglot_rs::Kind::U16,
            Kind::U32 => polyglot_rs::Kind::U32,
            Kind::U64 => polyglot_rs::Kind::U64,
            Kind::I32 => polyglot_rs::Kind::I32,
            Kind::I64 => polyglot_rs::Kind::I64,
            Kind::F32 => polyglot_rs::Kind::F32,
            Kind::F64 => polyglot_rs::Kind::F64,
            Kind::Unknown => polyglot_rs::Kind::Unknown,
        }
    }
}

#[repr(C)]
#[derive(Debug)]
pub struct Buffer {
    pub(crate) data: *mut u8,
    pub(crate) length: u32,
}

impl Buffer {
    pub fn new_raw(data: *mut u8, length: u32) -> *mut Self {
        Box::into_raw(Box::new(Buffer {
            data,
            length,
        }))
    }
}

pub(crate) struct StringError(pub(crate) String);

impl Debug for StringError {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        Debug::fmt(&self.0, f)
    }
}

impl Display for StringError {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        Display::fmt(&self.0, f)
    }
}

impl std::error::Error for StringError {
    fn description(&self) -> &str {
        &self.0
    }
}