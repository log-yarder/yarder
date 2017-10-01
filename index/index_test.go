package index

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testDoc = "foo bar"
const testDoc2 = "bar baz"

func TestEmptyIndex(t *testing.T) {
	var i MapIndex

	require.Empty(t, i.Match(nil))
	require.Empty(t, i.Match([]string{}))
	require.Empty(t, i.Match([]string{"bar"}))
	require.Empty(t, i.Match([]string{"baz"}))
	require.Empty(t, i.Match([]string{"bar", "baz"}))
}

func TestIndexWithOneDoc(t *testing.T) {
	var i MapIndex
	i.Add(testDoc, strings.Split(testDoc, " "))

	require.Equal(t, []string{testDoc}, i.Match([]string{"foo"}))
	require.Equal(t, []string{testDoc}, i.Match([]string{"bar"}))
	require.Equal(t, []string{testDoc}, i.Match([]string{"foo", "bar"}))
	require.Empty(t, i.Match(nil))
	require.Empty(t, i.Match([]string{}))
	require.Empty(t, i.Match([]string{"baz"}))
}

func TestUnmarshaledIndexWithOneDoc(t *testing.T) {
	var i Index = &MapIndex{}
	i.Add(testDoc, strings.Split(testDoc, " "))
	var buf bytes.Buffer
	err := i.MarshalTo(&buf)
	require.NoError(t, err)
	require.NotEmpty(t, buf.String())
	i, err = UnmarshalFrom(&buf)
	require.NoError(t, err)

	require.Equal(t, []string{testDoc}, i.Match([]string{"foo"}))
	require.Equal(t, []string{testDoc}, i.Match([]string{"bar"}))
	require.Equal(t, []string{testDoc}, i.Match([]string{"foo", "bar"}))
	require.Empty(t, i.Match(nil))
	require.Empty(t, i.Match([]string{}))
	require.Empty(t, i.Match([]string{"baz"}))
}

func TestIndexWithMultipleDocs(t *testing.T) {
	var i MapIndex
	i.Add(testDoc, strings.Split(testDoc, " "))
	i.Add(testDoc2, strings.Split(testDoc2, " "))

	require.Equal(t, []string{testDoc}, i.Match([]string{"foo"}))
	require.Equal(t, []string{testDoc, testDoc2}, i.Match([]string{"bar"}))
	require.Equal(t, []string{testDoc2}, i.Match([]string{"baz"}))
	require.Empty(t, i.Match(nil))
	require.Empty(t, i.Match([]string{}))
	require.Empty(t, i.Match([]string{"quux"}))
}
