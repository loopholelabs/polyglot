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
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"testing"
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

	testData := make([]byte, b.Cap()*2)
	_, err = rand.Read(testData)
	assert.NoError(t, err)

	for {
		assert.Equal(t, NewBuffer().Bytes(), b.Bytes())
		assert.Equal(t, 0, b.Len())

		b.Write(testData)
		assert.Equal(t, len(testData), b.Len())
		assert.GreaterOrEqual(t, b.Cap(), len(testData))

		pool.Put(b)
		b = pool.Get()

		if b.Cap() < len(testData) {
			continue
		}
		assert.Equal(t, 0, b.Len())
		assert.GreaterOrEqual(t, b.Cap(), len(testData))
		break
	}

	pool.Put(b)
}
