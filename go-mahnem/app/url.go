package main

import (
	"log"
	"net/url"
)

// EndpointBuilder defenition of site URL builder
type EndpointBuilder struct {
	url    string
	path   string
	params map[string]string
}

// NewEndpointBuilder creates new endpoint builder
func NewEndpointBuilder(baseURL string) *EndpointBuilder {
	return &EndpointBuilder{url: baseURL, params: make(map[string]string)}
}

// WithPath setter for url path
func (builder *EndpointBuilder) WithPath(path string) *EndpointBuilder {
	if len(path) > 0 {
		builder.path = path
	} else {
		builder.path = "/"
	}
	return builder
}

// WithQueryParam setter for url param
func (builder *EndpointBuilder) WithQueryParam(param, value string) *EndpointBuilder {
	if len(param) > 0 {
		builder.params[param] = value
	}
	return builder
}

func (builder *EndpointBuilder) String() string {

	reqURL, err := url.ParseRequestURI(builder.url)
	if err != nil {
		log.Printf("Can't parse base site url '%s', error '%s'\n", builder.url, err.Error())
		return ""
	}

	reqURL.Path = builder.path

	if builder.params != nil {
		reqParams := url.Values{}
		for k, v := range builder.params {
			if len(k) > 0 && len(v) > 0 {
				reqParams.Add(k, v)
			}
		}
		if len(reqParams) > 0 {
			reqURL.RawQuery = reqParams.Encode()
		}
	}

	return reqURL.String()
}
