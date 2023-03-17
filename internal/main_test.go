package main_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	m "github.com/tsanton/gohub-action"
)

func Test_main(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	m.Greet("foo", "bar")

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	assert.Equal(t, string(out), "foo, bar!")
}
