package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewVar(t *testing.T) {
	assert.Equal(t, Var{
		"x", "", "2",
	}, NewVar(`var x = 2`))

	assert.Equal(t, Var{
		"x", "int", "2",
	}, NewVar(`var x int = 2`))

	assert.Equal(t, Var{
		"x", "", "2",
	}, NewVar(`x := 2`))
}
