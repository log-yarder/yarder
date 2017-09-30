package index

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testDoc = "foo bar"

func TestEmptyIndex_MatchReturnsEmpty(t *testing.T) {
	var i MapIndex
	require.Empty(t, i.Match([]string{"bar"}))
	require.Empty(t, i.Match([]string{"baz"}))
}

func TestIndex_CanMatchAddedDoc(t *testing.T) {
	var i MapIndex
	i.Add(testDoc, strings.Split(testDoc, " "))
	require.Equal(t, []string{testDoc}, i.Match([]string{"foo"}))
	require.Equal(t, []string{testDoc}, i.Match([]string{"bar"}))
	require.Empty(t, i.Match([]string{"baz"}))
}

func TestUnmarshaledIndex_CanMatchAddedDoc(t *testing.T) {
	i := &MapIndex{}
	i.Add(testDoc, strings.Split(testDoc, " "))
	var buf bytes.Buffer
	err := i.WriteTo(&buf)
	require.NoError(t, err)
	require.NotEmpty(t, buf.String())
	i, err = ReadMapIndex(&buf)
	require.NoError(t, err)
	require.Equal(t, []string{testDoc}, i.Match([]string{"foo"}))
	require.Equal(t, []string{testDoc}, i.Match([]string{"bar"}))
	require.Empty(t, i.Match([]string{"baz"}))
}
