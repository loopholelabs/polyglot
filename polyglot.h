// SPDX-License-Identifier: Apache-2.0

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

#ifndef _LIBPOLYGLOT_H_
#define _LIBPOLYGLOT_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include <stdbool.h>

#define POLYGLOT_VERSION_MAJOR        2
#define POLYGLOT_VERSION_MINOR        0
#define POLYGLOT_VERSION_MICRO        1

#define POLYGLOT_VERSION                \
    ((POLYGLOT_VERSION_MAJOR * 10000) + \
     (POLYGLOT_VERSION_MINOR * 100) +   \
     POLYGLOT_VERSION_MICRO)

#ifdef POLYGLOT_USE_C_ENUMS
# define _C_ENUM_THUNK(key, val) key = val,
# define DEFINE_ENUM(name, cb)			\
	typedef enum name : uint8_t {	    \
		cb(_C_ENUM_THUNK)				\
	} name ## _t;
#else
# define _C_CONST_THUNK(key, val) static const uint8_t key = val;
# define DEFINE_ENUM(name, cb)          \
	cb(_C_CONST_THUNK)                  \
		typedef uint8_t name ## _t;
#endif

#define STATUS_VALUES(_)                 \
    _(POLYGLOT_STATUS_PASS, 0)           \
    _(POLYGLOT_STATUS_FAILURE, 1)        \
    _(POLYGLOT_STATUS_NULL_PTR,  2)      \

#define KIND_VALUES(_)                 \
    _(POLYGLOT_KIND_NONE, 0)           \
    _(POLYGLOT_KIND_ARRAY, 1)          \
    _(POLYGLOT_KIND_MAP, 2)            \
    _(POLYGLOT_KIND_ANY, 3)            \
    _(POLYGLOT_KIND_BYTES, 4)          \
    _(POLYGLOT_KIND_STRING, 5)         \
    _(POLYGLOT_KIND_ERROR, 6)          \
    _(POLYGLOT_KIND_BOOL, 7)           \
    _(POLYGLOT_KIND_U8, 8)             \
    _(POLYGLOT_KIND_U16, 9)            \
    _(POLYGLOT_KIND_U32, 10)           \
    _(POLYGLOT_KIND_U64, 11)           \
    _(POLYGLOT_KIND_I32, 12)           \
    _(POLYGLOT_KIND_I64, 13)           \
    _(POLYGLOT_KIND_F32, 14)           \
    _(POLYGLOT_KIND_F64, 15)           \
    _(POLYGLOT_KIND_UNKNOWN, 16)       \

DEFINE_ENUM(polyglot_status, STATUS_VALUES)
DEFINE_ENUM(polyglot_kind, KIND_VALUES)

typedef struct polyglot_buffer {
    uint8_t *data;
    uint32_t length;
} polyglot_buffer_t;

typedef struct polyglot_encoder polyglot_encoder_t;

polyglot_encoder_t* polyglot_new_encoder(polyglot_status_t *status);
uint32_t polyglot_encoder_size(polyglot_status_t *status, polyglot_encoder_t *encoder);
void polyglot_encoder_buffer(polyglot_status_t *status, polyglot_encoder_t *encoder, uint8_t *buffer_pointer, uint32_t buffer_size);
void polyglot_free_encoder(polyglot_encoder_t *encoder);

void polyglot_encode_none(polyglot_status_t *status, polyglot_encoder_t *encoder);
void polyglot_encode_array(polyglot_status_t *status, polyglot_encoder_t *encoder, uint32_t array_size, polyglot_kind_t array_kind);
void polyglot_encode_map(polyglot_status_t *status, polyglot_encoder_t *encoder, uint32_t map_size, polyglot_kind_t key_kind, polyglot_kind_t value_kind);
void polyglot_encode_bytes(polyglot_status_t *status, polyglot_encoder_t *encoder, uint8_t *buffer_pointer, uint32_t buffer_size);
void polyglot_encode_string(polyglot_status_t *status, polyglot_encoder_t *encoder, char *string_pointer);
void polyglot_encode_error(polyglot_status_t *status, polyglot_encoder_t *encoder, char *string_pointer);
void polyglot_encode_bool(polyglot_status_t *status, polyglot_encoder_t *encoder, bool value);
void polyglot_encode_u8(polyglot_status_t *status, polyglot_encoder_t *encoder, uint8_t value);
void polyglot_encode_u16(polyglot_status_t *status, polyglot_encoder_t *encoder, uint16_t value);
void polyglot_encode_u32(polyglot_status_t *status, polyglot_encoder_t *encoder, uint32_t value);
void polyglot_encode_u64(polyglot_status_t *status, polyglot_encoder_t *encoder, uint64_t value);
void polyglot_encode_i32(polyglot_status_t *status, polyglot_encoder_t *encoder, int32_t value);
void polyglot_encode_i64(polyglot_status_t *status, polyglot_encoder_t *encoder, int64_t value);
void polyglot_encode_f32(polyglot_status_t *status, polyglot_encoder_t *encoder, float value);
void polyglot_encode_f64(polyglot_status_t *status, polyglot_encoder_t *encoder, double value);

typedef struct polyglot_decoder polyglot_decoder_t;

polyglot_decoder_t* polyglot_new_decoder(polyglot_status_t *status, uint8_t *buffer_pointer, uint32_t buffer_size);
void polyglot_free_decoder(polyglot_decoder_t *decoder);

bool polyglot_decode_none(polyglot_status_t *status, polyglot_decoder_t *decoder);
uint32_t polyglot_decode_array(polyglot_status_t *status, polyglot_decoder_t *decoder, polyglot_kind_t array_kind);
uint32_t polyglot_decode_map(polyglot_status_t *status, polyglot_decoder_t *decoder, polyglot_kind_t key_kind, polyglot_kind_t value_kind);

polyglot_buffer_t *polyglot_decode_bytes(polyglot_status_t *status, polyglot_decoder_t *decoder);
void polyglot_free_decode_bytes(polyglot_buffer_t *buffer);

char* polyglot_decode_string(polyglot_status_t *status, polyglot_decoder_t *decoder);
void polyglot_free_decode_string(char *c_string);

char* polyglot_decode_error(polyglot_status_t *status, polyglot_decoder_t *decoder);
void polyglot_free_decode_error(char *c_string);

bool polyglot_decode_bool(polyglot_status_t *status, polyglot_decoder_t *decoder);
uint8_t polyglot_decode_u8(polyglot_status_t *status, polyglot_decoder_t *decoder);
uint16_t polyglot_decode_u16(polyglot_status_t *status, polyglot_decoder_t *decoder);
uint32_t polyglot_decode_u32(polyglot_status_t *status, polyglot_decoder_t *decoder);
uint64_t polyglot_decode_u64(polyglot_status_t *status, polyglot_decoder_t *decoder);
int32_t polyglot_decode_i32(polyglot_status_t *status, polyglot_decoder_t *decoder);
int64_t polyglot_decode_i64(polyglot_status_t *status, polyglot_decoder_t *decoder);
float polyglot_decode_f32(polyglot_status_t *status, polyglot_decoder_t *decoder);
double polyglot_decode_f64(polyglot_status_t *status, polyglot_decoder_t *decoder);

#ifdef __cplusplus
} /* extern "C" */
#endif

#endif /* End of _LIBPOLYGLOT_H_ */
