#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

enum PolyglotKind {
  None = 0,
  Array = 1,
  Map = 2,
  Any = 3,
  Bytes = 4,
  String = 5,
  Error = 6,
  Bool = 7,
  U8 = 8,
  U16 = 9,
  U32 = 10,
  U64 = 11,
  I32 = 12,
  I64 = 13,
  F32 = 14,
  F64 = 15,
  Unknown,
};
typedef uint8_t PolyglotKind;

enum PolyglotStatus {
  Pass,
  Fail,
  NullPointer,
};
typedef uint8_t PolyglotStatus;

typedef struct Encoder {
  Cursor<Vec<uint8_t>> cursor;
} Encoder;

typedef struct Decoder {
  Cursor<Vec<uint8_t>> cursor;
} Decoder;

typedef struct polyglot_decode_string_result {
  PolyglotStatus status;
  const uint8_t *value;
} polyglot_decode_string_result;

PolyglotStatus polyglot_new_encoder(struct Encoder **encoder);

unsigned int polyglot_encoder_size(struct Encoder *encoder);

PolyglotStatus polyglot_encoder_buffer(struct Encoder *encoder,
                                       char *buffer_pointer,
                                       unsigned int buffer_size);

void polyglot_free_encoder(struct Encoder *encoder);

PolyglotStatus polyglot_encode_none(struct Encoder *encoder);

PolyglotStatus polyglot_encode_array(struct Encoder *encoder,
                                     unsigned int array_size,
                                     PolyglotKind array_kind);

PolyglotStatus polyglot_encode_map(struct Encoder *encoder,
                                   unsigned int map_size,
                                   PolyglotKind key_kind,
                                   PolyglotKind value_kind);

PolyglotStatus polyglot_encode_bytes(struct Encoder *encoder,
                                     char *buffer_pointer,
                                     unsigned int buffer_size);

PolyglotStatus polyglot_encode_string(struct Encoder *encoder, const char *string_pointer);

PolyglotStatus polyglot_new_decoder(struct Decoder **decoder,
                                    char *buffer_pointer,
                                    unsigned int buffer_size);

void polyglot_decoder_free(struct Decoder *decoder);

struct polyglot_decode_string_result polyglot_decode_string(struct Decoder *decoder);
