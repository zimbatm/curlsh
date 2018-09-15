package main

// https://w3c.github.io/webappsec-subresource-integrity/

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"fmt"
	"strings"
)

var b64 = base64.StdEncoding

type ErrAlgoNotSupported struct{ Algo string }

func (x ErrAlgoNotSupported) Error() string {
	return fmt.Sprintf("unknown algo '%s'", x.Algo)
}

type SRI struct {
	Algo string
	Hash crypto.Hash
	Sum  []byte
}

// Check validates that the recorded Sum actually matches the content
func (sri SRI) Check(content []byte) bool {
	hash := sri.Hash.New()
	hash.Write(content)
	newSum := hash.Sum(nil)
	return bytes.Equal(newSum, sri.Sum)
}

// parseSRI conversts a single SRI entry to a Hash
func parseSRI(sriStr string) (*SRI, error) {
	args := strings.SplitN(sriStr, "-", 2)
	if len(args) < 2 {
		return nil, fmt.Errorf("SRI parsing error")
	}

	algo := args[0]
	base64sum := args[1]

	var hash crypto.Hash
	switch algo {
	case "sha256":
		hash = crypto.SHA256
	case "sha384":
		hash = crypto.SHA384
	case "sha512":
		hash = crypto.SHA512
	default:
		return nil, ErrAlgoNotSupported{algo}
	}

	sri := &SRI{
		Algo: algo,
		Hash: hash,
		// FIXME: why +2 is necessary?
		Sum: make([]byte, hash.Size()+2),
	}

	size, err := b64.Decode(sri.Sum, []byte(base64sum))
	if err != nil {
		return nil, err
	}

	// FIXME: see above
	sri.Sum = sri.Sum[0:hash.Size()]

	if size != sri.Hash.Size() {
		return nil, fmt.Errorf("SRI parsing error: expected hash to be of size %d, not %d", sri.Hash.Size(), size)
	}

	return sri, nil
}

func parseSRIList(sriStrList string) (ret []*SRI, err error) {
	list := strings.Split(sriStrList, " ")
	ret = make([]*SRI, 0, len(list))

	for _, sriStr := range list {
		sri, err := parseSRI(sriStr)
		ret = append(ret, sri)
		if err != nil {
			switch err.(type) {
			case ErrAlgoNotSupported:
				continue
			default:
				return nil, err
			}
		}
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("no supported SRI found")
	}

	return ret, nil
}
