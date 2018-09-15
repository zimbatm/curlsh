package main

import (
	"flag"
)

var _ flag.Value = &SRIFlag{}

// Implements the flag.Value interface
type SRIFlag struct {
	List []*SRI
}

func (x *SRIFlag) String() string {
	return ""
}

func (x *SRIFlag) Set(rawurl string) (err error) {
	x.List, err = parseSRIList(rawurl)
	return err
}
