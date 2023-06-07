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

package polyglot

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecycle(t *testing.T) {
	pool := NewPool()

	data := make([]byte, 512)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	b := pool.Get()
	b.Write(data)

	pool.Put(b)
	b = pool.Get()

	testData := make([]byte, cap(*b)*2)
	_, err = rand.Read(testData)
	assert.NoError(t, err)

	for {
		assert.Equal(t, Buffer([]byte{}), *b)
		assert.Equal(t, 0, len(*b))

		b.Write(testData)
		assert.Equal(t, len(testData), len(*b))
		assert.GreaterOrEqual(t, cap(*b), len(testData))

		pool.Put(b)
		b = pool.Get()

		if cap(*b) < len(testData) {
			continue
		}
		assert.Equal(t, 0, len(*b))
		assert.GreaterOrEqual(t, cap(*b), len(testData))
		break
	}

	pool.Put(b)
}
