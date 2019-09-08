package main

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestEndpointBuilder_String(t *testing.T) {
	// given
	var baseURL = "http://fake.mahnem.ru"

	tests := []struct {
		url  string
		path string
		want string
	}{
		{baseURL, "", "http://fake.mahnem.ru/"},
		{baseURL, "/", "http://fake.mahnem.ru/"},
		{"", "/", ""},
		{"", "", ""},
	}

	// then
	for _, tt := range tests {

		url := NewEndpointBuilder(tt.url).
			WithPath(tt.path).
			String()

		Equal(t, url, tt.want)
	}
}

func TestEndpointBuilderWithQueryParams_String(t *testing.T) {
	// given
	tests := []struct {
		builder *EndpointBuilder
		param   string
		value   string
		want    string
	}{
		{builder(), "k", "v ", "http://fake.mahnem.ru/path?k=v+"},
		{builder(), "k", " v", "http://fake.mahnem.ru/path?k=+v"},
		{builder(), "", "v", "http://fake.mahnem.ru/path"},
	}

	// then
	for _, tt := range tests {

		url := tt.builder.WithQueryParam(tt.param, tt.value).String()

		Equal(t, url, tt.want)
	}
}

func builder() *EndpointBuilder {
	return NewEndpointBuilder("http://fake.mahnem.ru").WithPath("path")
}
