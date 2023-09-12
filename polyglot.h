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

#define POLYGLOT_VERSION_MAJOR        1
#define POLYGLOT_VERSION_MINOR        1
#define POLYGLOT_VERSION_MICRO        3

#define POLYGLOT_VERSION              \
    ((POLYGLOT_VERSION_MAJOR * 10000) + \
     (POLYGLOT_VERSION_MINOR * 100) + \
     POLYGLOT_VERSION_MICRO)

#define POLYGLOT_PASS           0
#define POLYGLOT_FAIL           1
#define POLYGLOT_NULL_POINTER   2

struct encoder;

int encoder_new(struct encoder **encoder);
void encoder_free(struct encoder *encoder);

#ifdef __cplusplus
} /* extern "C" */
#endif

#endif /* End of _LIBPOLYGLOT_H_ */
