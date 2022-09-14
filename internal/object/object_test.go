package object

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	require.Equal(t, hello1.HashKey(), hello2.HashKey())
	require.Equal(t, diff1.HashKey(), diff2.HashKey())
	require.NotEqual(t, hello1.HashKey(), diff1.HashKey())
}
