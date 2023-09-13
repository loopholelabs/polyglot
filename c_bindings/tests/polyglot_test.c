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
#include <stdlib.h>
#include <polyglot.h>

int main(void) {
    polyglot_status_t polyglot_status = POLYGLOT_STATUS_PASS;
    printf("initial status: %d\n", polyglot_status);

    struct polyglot_encoder *encoder = NULL;

    polyglot_status = polyglot_new_encoder(&encoder);
    printf("polyglot_new_encoder status: %d\n", polyglot_status);

    uint8_t *input_buffer_pointer = malloc(32);
    uint8_t *current_input_buffer_pointer = input_buffer_pointer;
    for (uint8_t i = 0; i < 32; i++) {
        *current_input_buffer_pointer = i;
        current_input_buffer_pointer++;
    }

    polyglot_status = polyglot_encode_bytes(encoder, input_buffer_pointer, 32);
    printf("polyglot_encode_bytes status: %d\n", polyglot_status);
    free(input_buffer_pointer);

    polyglot_kind_t polyglot_kind = POLYGLOT_KIND_ARRAY;
    polyglot_status = polyglot_encode_array(encoder, 8, polyglot_kind);
    printf("polyglot_encode_array status: %d\n", polyglot_status);

    polyglot_status = polyglot_encode_none(encoder);
    printf("polyglot_encode_none status: %d\n", polyglot_status);

    char *input_string_pointer = "Hello, World!";
    polyglot_status = polyglot_encode_string(encoder, input_string_pointer);
    printf("polyglot_encode_string status: %d\n", polyglot_status);

    uint32_t buffer_size = 0;
    buffer_size = polyglot_encoder_size(encoder);
    printf("polyglot_encoder_size: %d\n", buffer_size);

    uint8_t *buffer_pointer = malloc(buffer_size);
    polyglot_status = polyglot_encoder_buffer(encoder, buffer_pointer, buffer_size);
    printf("polyglot_encoder_buffer status: %d\n", polyglot_status);
    polyglot_free_encoder(encoder);

    uint8_t *current_buffer_pointer = buffer_pointer;
    printf("polyglot_encoder_buffer contents: ");
    for(uint32_t i = 0; i < buffer_size; i++) {
        printf("%d ", *current_buffer_pointer);
        current_buffer_pointer++;
    }
    printf("\n");

    struct polyglot_decoder *decoder = NULL;
    polyglot_status = polyglot_new_decoder(&decoder, buffer_pointer, buffer_size);
    printf("polyglot_new_decoder status: %d\n", polyglot_status);
    
    polyglot_
    

    free(buffer_pointer);
}