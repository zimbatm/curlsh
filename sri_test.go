package main

import (
	"io/ioutil"
	"testing"
)

func TestParseSRI(t *testing.T) {
	var err error

	content, err := ioutil.ReadFile("sri_test.js")
	if err != nil {
		t.Error(err)
	}

	sriList, err := parseSRIList("sha256-ySadHRVML1LfcwlPIxXx4CQpk64arq0Yv32cBpu9CFQ= sha384-z5SseNF3PzzKas2Pab3YN8u3SauawAs+a69rDYBiAFtJIlJxoWrg29HldeHVuv8g sha512-r6PRJqbVfzPNfeOP2OjtGmkuSAcXFR9Z22F1GiUicEycz0SHSD/gaMfvFKHwSqGTPf/O+WIGFhf48Vu1Xw8UhQ==")
	if err != nil {
		t.Fatal(err)
	}

	for _, sri := range sriList {
		if err != nil {
			t.Error(err)
		}
		if !sri.Check(content) {
			t.Error("check failed for", sri.Algo)
		}
	}

}
