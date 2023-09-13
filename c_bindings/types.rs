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

#[repr(u8)]
#[derive(Debug)]
pub enum PolyglotStatus {
    Pass,
    Fail,
    NullPointer,
}

impl PolyglotStatus {
    pub fn check_not_null(status: *mut PolyglotStatus) {
        if status.is_null() {
            panic!("status pointer is null");
        }
    }
}

#[repr(u8)]
#[derive(Debug)]
pub enum PolyglotKind {
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

impl PolyglotKind {
    pub fn into(self) -> polyglot_rs::Kind {
        match self {
            PolyglotKind::None => polyglot_rs::Kind::None,
            PolyglotKind::Array => polyglot_rs::Kind::Array,
            PolyglotKind::Map => polyglot_rs::Kind::Map,
            PolyglotKind::Any => polyglot_rs::Kind::Any,
            PolyglotKind::Bytes => polyglot_rs::Kind::Bytes,
            PolyglotKind::String => polyglot_rs::Kind::String,
            PolyglotKind::Error => polyglot_rs::Kind::Error,
            PolyglotKind::Bool => polyglot_rs::Kind::Bool,
            PolyglotKind::U8 => polyglot_rs::Kind::U8,
            PolyglotKind::U16 => polyglot_rs::Kind::U16,
            PolyglotKind::U32 => polyglot_rs::Kind::U32,
            PolyglotKind::U64 => polyglot_rs::Kind::U64,
            PolyglotKind::I32 => polyglot_rs::Kind::I32,
            PolyglotKind::I64 => polyglot_rs::Kind::I64,
            PolyglotKind::F32 => polyglot_rs::Kind::F32,
            PolyglotKind::F64 => polyglot_rs::Kind::F64,
            PolyglotKind::Unknown => polyglot_rs::Kind::Unknown,
        }
    }
}