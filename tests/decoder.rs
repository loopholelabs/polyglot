/*
    Copyright 2022 Loophole Labs

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

extern crate polyglot_rs;

use polyglot_rs::Decoder;
use polyglot_rs::DecodingError;
use polyglot_rs::Encoder;
use polyglot_rs::Kind;
use std::collections::HashMap;
use std::error::Error;
use std::io::Cursor;

#[test]
fn test_decode_nil() {
    let mut encoder = Cursor::new(Vec::with_capacity(512));
    encoder.encode_none().unwrap();

    let mut decoder = Cursor::new(encoder.get_mut());
    let val = decoder.decode_none();
    assert_eq!(val, true);
    assert_eq!(decoder.get_ref().len() - decoder.position() as usize, 0);
    let next_val = decoder.decode_none();
    assert_eq!(next_val, false);
}

#[test]
fn test_decode_array() {
    let mut encoder = Cursor::new(Vec::with_capacity(512));
    let m = ["1".to_string(), "2".to_string(), "3".to_string()];

    encoder.encode_array(m.len(), Kind::String).unwrap();
    for i in m.clone() {
        encoder.encode_string(&i).unwrap();
    }

    let mut decoder = Cursor::new(encoder.get_mut());
    let size = decoder.decode_array(Kind::String).unwrap() as usize;
    assert_eq!(size, m.len());

    let mut mv: Vec<String> = Vec::with_capacity(size);
    for i in 0..size {
        let val = decoder.decode_string().unwrap();
        mv.push(val.to_string());
        assert_eq!(mv[i], m[i]);
    }
    assert_eq!(mv, m);

    let error = decoder.decode_array(Kind::String).unwrap_err();
    assert_eq!(error, DecodingError::InvalidArray);
}

#[test]
fn test_decode_map() {
    let mut encoder = Cursor::new(Vec::with_capacity(512));
    let mut m: HashMap<String, u32> = HashMap::new();
    m.insert(String::from("1"), 1);
    m.insert(String::from("2"), 2);
    m.insert(String::from("3"), 3);

    encoder
        .encode_map(m.len(), Kind::String, Kind::U32)
        .unwrap();
    for (k, v) in m.clone() {
        encoder.encode_string(&k).unwrap().encode_u32(v).unwrap();
    }

    let mut decoder = Cursor::new(encoder.get_mut());
    let size = decoder.decode_map(Kind::String, Kind::U32).unwrap() as usize;
    assert_eq!(size, m.len());

    let mut mv = HashMap::new();
    for _ in 0..size {
        let k = decoder.decode_string().unwrap();
        let v = decoder.decode_u32().unwrap();
        mv.insert(k, v);
    }
    assert_eq!(mv, m);

    let error = decoder.decode_map(Kind::String, Kind::U32).unwrap_err();
    assert_eq!(error, DecodingError::InvalidMap);
}

#[test]
fn test_decode_bytes() {
    let mut encoder = Cursor::new(Vec::with_capacity(512));
    let v = "Test String".as_bytes();
    encoder.encode_bytes(v).unwrap();

    let mut decoder = Cursor::new(encoder.get_mut());
    let bytes = decoder.decode_bytes().unwrap();
    assert_eq!(bytes, v);
}

#[test]
fn test_decode_string() {
    let mut encoder = Cursor::new(Vec::with_capacity(512));
    let v = "Test String".to_string();
    encoder.encode_string(&v).unwrap();

    let mut decoder = Cursor::new(encoder.get_mut());
    let val = decoder.decode_string().unwrap();
    assert_eq!(val, v);

    let error = decoder.decode_string().unwrap_err();
    assert_eq!(error, DecodingError::InvalidString);
}

#[test]
fn test_decode_error() {
    let mut encoder = Cursor::new(Vec::with_capacity(512));
    let v = "Test String";
    encoder.encode_error(Box::<dyn Error>::from(v)).unwrap();

    let mut decoder = Cursor::new(encoder.get_mut());
    let val = decoder.decode_error().unwrap();
    assert_eq!(val.to_string(), v);

    let error = decoder.decode_error().unwrap_err();
    assert_eq!(error, DecodingError::InvalidError);
}

#[test]
fn test_decode_bool() {
    let mut encoder = Cursor::new(Vec::with_capacity(512));
    encoder.encode_bool(true).unwrap();

    let mut decoder = Cursor::new(encoder.get_mut());
    let val = decoder.decode_bool().unwrap();
    assert_eq!(val, true);

    let error = decoder.decode_bool().unwrap_err();
    assert_eq!(error, DecodingError::InvalidBool);
}

#[test]
fn test_decode_u8() {
    let mut encoder = Cursor::new(Vec::with_capacity(512));
    let v = 32 as u8;
    encoder.encode_u8(v).unwrap();

    let mut decoder = Cursor::new(encoder.get_mut());
    let val = decoder.decode_u8().unwrap();
    assert_eq!(val, v);

    let error = decoder.decode_u8().unwrap_err();
    assert_eq!(error, DecodingError::InvalidU8);
}

#[test]
fn test_decode_u16() {
    let mut encoder = Cursor::new(Vec::with_capacity(512));
    let v = 1024 as u16;
    encoder.encode_u16(v).unwrap();

    let mut decoder = Cursor::new(encoder.get_mut());
    let val = decoder.decode_u16().unwrap();
    assert_eq!(val, v);

    let error = decoder.decode_u16().unwrap_err();
    assert_eq!(error, DecodingError::InvalidU16);
}

#[test]
fn test_decode_u32() {
    let mut encoder = Cursor::new(Vec::with_capacity(512));
    let v = 4294967290 as u32;
    encoder.encode_u32(v).unwrap();

    let mut decoder = Cursor::new(encoder.get_mut());
    let val = decoder.decode_u32().unwrap();
    assert_eq!(val, v);

    let error = decoder.decode_u32().unwrap_err();
    assert_eq!(error, DecodingError::InvalidU32);
}

#[test]
fn test_decode_u64() {
    let mut encoder = Cursor::new(Vec::with_capacity(512));
    let v = 18446744073709551610 as u64;
    encoder.encode_u64(v).unwrap();

    let mut decoder = Cursor::new(encoder.get_mut());
    let val = decoder.decode_u64().unwrap();
    assert_eq!(val, v);

    let error = decoder.decode_u64().unwrap_err();
    assert_eq!(error, DecodingError::InvalidU64);
}

#[test]
fn test_decode_i32() {
    let mut encoder = Cursor::new(Vec::with_capacity(512));
    let v = 2147483647 as i32;
    let vneg = -32 as i32;
    encoder.encode_i32(v).unwrap();
    encoder.encode_i32(vneg).unwrap();

    let mut decoder = Cursor::new(encoder.get_mut());
    let val = decoder.decode_i32().unwrap();
    let valneg = decoder.decode_i32().unwrap();
    assert_eq!(val, v);
    assert_eq!(valneg, vneg);

    let error = decoder.decode_i32().unwrap_err();
    assert_eq!(error, DecodingError::InvalidI32);
}

#[test]
fn test_decode_i64() {
    let mut encoder = Cursor::new(Vec::with_capacity(512));
    let v = 9223372036854775807 as i64;
    let vneg = -32 as i64;
    encoder.encode_i64(v).unwrap();
    encoder.encode_i64(vneg).unwrap();

    let mut decoder = Cursor::new(encoder.get_mut());
    let val = decoder.decode_i64().unwrap();
    let valneg = decoder.decode_i64().unwrap();
    assert_eq!(val, v);
    assert_eq!(valneg, vneg);

    let error = decoder.decode_i64().unwrap_err();
    assert_eq!(error, DecodingError::InvalidI64);
}

#[test]
fn test_decode_f32() {
    let mut encoder = Cursor::new(Vec::with_capacity(512));
    let v = -2147483.648 as f32;
    encoder.encode_f32(v).unwrap();

    let mut decoder = Cursor::new(encoder.get_mut());
    let val = decoder.decode_f32().unwrap();
    assert_eq!(val, v);

    let error = decoder.decode_f32().unwrap_err();
    assert_eq!(error, DecodingError::InvalidF32);
}

#[test]
fn test_decode_f64() {
    let mut encoder = Cursor::new(Vec::with_capacity(512));
    let v = -922337203.477580 as f64;
    encoder.encode_f64(v).unwrap();

    let mut decoder = Cursor::new(encoder.get_mut());
    let val = decoder.decode_f64().unwrap();
    assert_eq!(val, v);

    let error = decoder.decode_f64().unwrap_err();
    assert_eq!(error, DecodingError::InvalidF64);
}
