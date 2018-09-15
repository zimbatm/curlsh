package main

import (
	"flag"
	"net/url"
)

var _ flag.Value = URLFlag{}

// Implements the flag.Value interface
type URLFlag struct {
	*url.URL
}

func (x URLFlag) String() string {
	if x.URL == nil {
		return ""
	}
	return x.URL.String()
}

func (x URLFlag) Set(rawurl string) error {
	url, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	(*x.URL) = (*url)
	return nil
}
