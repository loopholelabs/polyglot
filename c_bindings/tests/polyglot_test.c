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

#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>

#include <assert.h>
#include <string.h>

#include <polyglot.h>

int main(void) {
    polyglot_status_t status = POLYGLOT_STATUS_PASS;

    polyglot_encoder_t* encoder = polyglot_new_encoder(&status);
    assert(status == POLYGLOT_STATUS_PASS);
    assert(encoder != NULL);

    uint32_t buffer_size = polyglot_encoder_size(&status, encoder);
    assert(status == POLYGLOT_STATUS_PASS);
    assert(buffer_size == 0);

    polyglot_encode_none(&status, encoder);
    assert(status == POLYGLOT_STATUS_PASS);

    buffer_size = polyglot_encoder_size(&status, encoder);
    assert(status == POLYGLOT_STATUS_PASS);
    assert(buffer_size == 1);

    char *input_string_pointer = "Hello, World!";
    polyglot_encode_string(&status, encoder, input_string_pointer);
    assert(status == POLYGLOT_STATUS_PASS);

    buffer_size = polyglot_encoder_size(&status, encoder);
    assert(status == POLYGLOT_STATUS_PASS);
    assert(buffer_size == 1 + (1 + 15));

    polyglot_encode_array(&status, encoder, 8, POLYGLOT_KIND_STRING);
    assert(status == POLYGLOT_STATUS_PASS);

    buffer_size = polyglot_encoder_size(&status, encoder);
    assert(status == POLYGLOT_STATUS_PASS);
    assert(buffer_size == 1 + (1 + 15) + (1 + 1 + 1 + 1));

    uint8_t *input_buffer_pointer = malloc(32);
    uint8_t *current_input_buffer_pointer = input_buffer_pointer;
    for (uint8_t i = 0; i < 32; i++) {
        *current_input_buffer_pointer = i;
        current_input_buffer_pointer++;
    }

    polyglot_encode_bytes(&status, encoder, input_buffer_pointer, 32);
    assert(status == POLYGLOT_STATUS_PASS);
    free(input_buffer_pointer);

    buffer_size = polyglot_encoder_size(&status, encoder);
    assert(status == POLYGLOT_STATUS_PASS);
    assert(buffer_size == 1 + (1 + 15) + (1 + 1 + 1 + 1 ) + (1 + 1 + 1 + 32));

    uint8_t *buffer_pointer = malloc(buffer_size);
    polyglot_encoder_buffer(&status, encoder, buffer_pointer, buffer_size);
    assert(status == POLYGLOT_STATUS_PASS);
    polyglot_free_encoder(encoder);

//    uint8_t *current_buffer_pointer = buffer_pointer;
//    printf("polyglot_encoder_buffer contents: ");
//    for(uint32_t i = 0; i < buffer_size; i++) {
//        printf("%d ", *current_buffer_pointer);
//        current_buffer_pointer++;
//    }
//    printf("\n");

    polyglot_decoder_t *decoder = polyglot_new_decoder(&status, buffer_pointer, buffer_size);
    assert(status == POLYGLOT_STATUS_PASS);
    assert(decoder != NULL);

    bool decode_none_success = polyglot_decode_none(&status, decoder);
    assert(status == POLYGLOT_STATUS_PASS);
    assert(decode_none_success == true);

    char *output_string_pointer = polyglot_decode_string(&status, decoder);
    assert(status == POLYGLOT_STATUS_PASS);
    assert(output_string_pointer != NULL);
    assert(strcmp(input_string_pointer, output_string_pointer) == 0);
    polyglot_free_decode_string(output_string_pointer);

    uint32_t output_array_size = polyglot_decode_array(&status, decoder, POLYGLOT_KIND_STRING);
    assert(status == POLYGLOT_STATUS_PASS);
    assert(output_array_size == 8);

    polyglot_free_decoder(decoder);
    free(buffer_pointer);
}