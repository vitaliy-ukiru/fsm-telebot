package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_generateKey(t *testing.T) {
	s := Storage{prefix: "test"}

	assert.Equal(t,
		"test:1:1:data:myKey",
		s.generateKey(1, 1, stateDataKey, "myKey"),
	)
	assert.Equal(t,
		"test:1:1:data:multiple:keys",
		s.generateKey(1, 1, stateDataKey, "multiple", "keys"),
	)
	assert.Equal(t,
		"test:1:1:data:*",
		s.generateKey(1, 1, stateDataKey, "*"),
	)
	assert.Equal(t,
		"test:1:1:state",
		s.generateKey(1, 1, stateKey),
	)
}
