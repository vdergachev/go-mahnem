package main

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

type stripArgs struct {
	val  string
	want string
}

func Test_strip(t *testing.T) {

	var tests = []stripArgs{
		{val: "", want: ""},
		{val: " ", want: ""},
		{val: "\n", want: ""},
		{val: "\r\n", want: ""},
		{val: "\r", want: ""},
		{val: " text ", want: "text"},
	}

	for _, tt := range tests {
		Equal(t, strip(tt.val), tt.want)
	}
}
