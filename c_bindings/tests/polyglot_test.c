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
#include <polyglot.h>

int main(void) {
    polyglot_status_t status = POLYGLOT_STATUS_PASS;
    printf("initial status: %d\n", status);

    polyglot_encoder_t* encoder = polyglot_new_encoder(&status);
    printf("polyglot_new_encoder status: %d\n", status);

    polyglot_encode_none(&status, encoder);
    printf("polyglot_encode_none status: %d\n", status);

    char *input_string_pointer = "Hello, World!";
    polyglot_encode_string(&status, encoder, input_string_pointer);
    printf("polyglot_encode_string status: %d\n", status);

    uint8_t *input_buffer_pointer = malloc(32);
    uint8_t *current_input_buffer_pointer = input_buffer_pointer;
    for (uint8_t i = 0; i < 32; i++) {
        *current_input_buffer_pointer = i;
        current_input_buffer_pointer++;
    }

    polyglot_encode_bytes(&status, encoder, input_buffer_pointer, 32);
    printf("polyglot_encode_bytes status: %d\n", status);
    free(input_buffer_pointer);

    polyglot_kind_t polyglot_kind = POLYGLOT_KIND_ARRAY;
    polyglot_encode_array(&status, encoder, 8, polyglot_kind);
    printf("polyglot_encode_array status: %d\n", status);

    uint32_t buffer_size = polyglot_encoder_size(&status, encoder);
    printf("polyglot_encoder_size status: %d\n", status);
    printf("polyglot_encoder_size buffer_size: %d\n", buffer_size);

    uint8_t *buffer_pointer = malloc(buffer_size);
    polyglot_encoder_buffer(&status, encoder, buffer_pointer, buffer_size);
    printf("polyglot_encoder_buffer status: %d\n", status);
    polyglot_free_encoder(encoder);

    uint8_t *current_buffer_pointer = buffer_pointer;
    printf("polyglot_encoder_buffer contents: ");
    for(uint32_t i = 0; i < buffer_size; i++) {
        printf("%d ", *current_buffer_pointer);
        current_buffer_pointer++;
    }
    printf("\n");

    polyglot_decoder_t *polyglot_decoder = polyglot_new_decoder(&status, buffer_pointer, buffer_size);
    printf("polyglot_new_decoder status: %d\n", status);

    bool decode_none_success = polyglot_decode_none(&status, polyglot_decoder);
    printf("polyglot_decode_none status: %d\n", status);
    printf("polyglot_decode_none success: %s\n", decode_none_success ? "true" : "false");

    char *output_string_pointer = polyglot_decode_string(&status, polyglot_decoder);
    printf("polyglot_decode_string status: %d\n", status);
    printf("polyglot_decode_string value: %s\n", output_string_pointer);
    polyglot_free_decode_string(output_string_pointer);

    polyglot_free_decoder(polyglot_decoder);
    free(buffer_pointer);
}