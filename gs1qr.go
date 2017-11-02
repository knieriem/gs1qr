// Copyright 2017 The gs1qr Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gs1qr encodes GS1 QR codes.
//
// The GS1 specification is at https://www.gs1.org/docs/barcodes/GS1_General_Specifications.pdf. See also https://2016archive.gs1us.org/gs1-us-library/command/core_download/entryid/768/method/attachment
//
// QR code generation is based on Russ Cox' qr package: https://github.com/rsc/qr
//
package gs1qr

import (
	"strings"

	"github.com/knieriem/gs1qr/ai"
	"rsc.io/qr/coding"
)

// An Elem is an ai.Elem with additional encoding specific information.
type Elem struct {
	preGS    string
	postGS   string
	preFNC1  bool
	postFNC1 bool
	hideID   bool
	*ai.Elem
	isAlpha bool
	isNum   bool
	c       coding.Encoding
}

func (e *Elem) String() string {
	s := ""
	if !e.hideID {
		s += e.preGS + e.AI.String()
	}
	return s + e.Data + e.postGS
}

func (e *Elem) setupEncoding() {
	text := e.String()
	e.isNum = false
	e.isAlpha = false
	// Pick data encoding, smallest first.
	// We could split the string and use different encodings
	// but that seems like overkill for now.
	switch {
	case coding.Num(text).Check() == nil:
		e.c = coding.Num(text)
		e.isNum = true
	case coding.Alpha(text).Check() == nil:
		e.c = coding.Alpha(text)
		e.isAlpha = true
	default:
		e.c = coding.String(text)
	}
}

const (
	ascGS = "\x1d"
)

// ElemList is a slice of elements.
type ElemList []Elem

// ConvertElements initializes a slice of Elem from a slice of ai.Elem.
func ConvertElements(list []ai.Elem) ElemList {
	var prev *Elem
	el := make(ElemList, len(list))
	for i := range list {
		e := &el[i]
		e.Elem = &list[i]
		e.setupEncoding()
		if i == 0 {
			prev = e
			continue
		}
		if prev.AI.Variable {
			if prev.isNum {
				if e.isNum {
					if e.AI.Variable {
						// make e Alpha encoded
						//	e.preFNC1 = true
						e.preGS = "%"
					} else {
						prev.postGS = "%"
						//	e.postFNC1 = true
					}
				} else if e.isAlpha {
					e.preGS = "%"
				} else {
					//	e.preFNC1 = true
					e.preGS = ascGS
				}
				e.setupEncoding()
			} else if prev.isAlpha {
				prev.postGS = "%"
			} else {
				prev.postGS = ascGS
				//	prev.postFNC1 = true
			}
			prev.setupEncoding()
		}
		prev = e
	}
	return el
}

func (list ElemList) Strings() []string {
	sym := make([]string, 1, 1+len(list))
	sym[0] = "<FNC1>"
	r := strings.NewReplacer(ascGS, "<GS>")
	for i := range list {
		sym = append(sym, r.Replace(list[i].String()))
	}
	return sym
}

// Compile translates a list of elements into a corresponding list of
// QR data encoding schemes.
func (list ElemList) Compile() []coding.Encoding {
	var clist = make([]coding.Encoding, 0, 1+len(list))
	clist = append(clist, FNC1{})
	text := ""
	for i := range list {
		e := &list[i]
		if e.preFNC1 {
			clist = append(clist, coding.Alpha("%"))
		}
		if i == 0 || e.preFNC1 {
			clist = append(clist, e.c)
			text = e.String()
			if e.postFNC1 {
				clist = append(clist, coding.Alpha("%"))
				text = ""
			}
			continue
		}
		prev := clist[len(clist)-1]
		equal := false
		switch prev.(type) {
		case coding.Num:
			equal = e.isNum
			if !equal && e.preGS == "" {
				text += e.AI.String()
				clist[len(clist)-1] = pickEncoding(text)
				e.hideID = true
				e.setupEncoding()
			}
		case coding.Alpha:
			equal = e.isAlpha
		case coding.String:
			_, equal = e.c.(coding.String)
		}
		if equal {
			text += e.String()
			clist[len(clist)-1] = pickEncoding(text)
		} else {
			clist = append(clist, e.c)
			text = e.String()
		}
		if e.postFNC1 {
			clist = append(clist, coding.Alpha("%"))
		}
	}
	return clist
}

// FNC1 defines an QR encoding scheme for the Function Code 1
type FNC1 struct{}

func (s FNC1) String() string {
	return "FNC1"
}

func (s FNC1) Check() error {
	return nil
}

func (s FNC1) Bits(v coding.Version) int {
	return 4
}

func (s FNC1) Encode(b *coding.Bits, v coding.Version) {
	b.Write(5, 4)
}

// FNC12nd defines an QR encoding scheme for the Function Code 1
// in second position.
type FNC12nd struct{}

func (s FNC12nd) String() string {
	return "fnc1"
}

func (s FNC12nd) Check() error {
	return nil
}

func (s FNC12nd) Bits(v coding.Version) int {
	return 4
}

func (s FNC12nd) Encode(b *coding.Bits, v coding.Version) {
	b.Write(9, 4)
}
