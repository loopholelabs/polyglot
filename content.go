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

package polyglot

const (
	defaultSize = 512
)

type C []byte

func (c *C) Reset() {
	*c = (*c)[:0]
}

func (c *C) Write(b []byte) int {
	if cap(*c)-len(*c) < len(b) {
		*c = append((*c)[:len(*c)], b...)
	} else {
		*c = (*c)[:len(*c)+copy((*c)[len(*c):cap(*c)], b)]
	}
	return len(b)
}

func CNew() *C {
	c := make(C, 0, defaultSize)
	return &c
}
