package qutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/masterstat/qutil"
)

func TestUniqueStrings(t *testing.T) {
	values := []string{"alpha", "beta", "beta", "gamma"}
	result := qutil.UniqueStrings(values)
	expect := []string{"alpha", "beta", "gamma"}
	assert.Equal(t, expect, result)
}
