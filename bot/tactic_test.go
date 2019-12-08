package bot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefineRole(t *testing.T) {
	for i := uint32(2); i <= 11; i++ {
		assert.NotEqual(t, "", DefineRole(i))
	}
}
