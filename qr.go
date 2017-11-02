// Copyright 2017 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gs1qr

import (
	"errors"

	"rsc.io/qr"
	"rsc.io/qr/coding"
)

// Encode returns an encoding of the list of QR data encoding schemes
// at the given error correction level.
func Encode(list []coding.Encoding, level coding.Level, minVersion coding.Version) (*qr.Code, *coding.Plan, error) {
	// Pick size.
	l := coding.Level(level)
	v := minVersion
	if v == 0 {
		v = coding.MinVersion
	}
	for ; ; v++ {
		if v > coding.MaxVersion {
			return nil, nil, errors.New("text too long to encode as QR")
		}
		nBits := 0
		for _, enc := range list {
			nBits += enc.Bits(v)
		}
		if nBits <= v.DataBytes(l)*8 {
			break
		}
	}

	// Build and execute plan.
	p, err := coding.NewPlan(v, l, 0)
	if err != nil {
		return nil, nil, err
	}

	cc, err := p.Encode(list...)
	if err != nil {
		return nil, p, err
	}

	// TODO: Pick appropriate mask.

	return &qr.Code{
		Bitmap: cc.Bitmap,
		Size:   cc.Size,
		Stride: cc.Stride,
		Scale:  8}, p, nil
}

func pickEncoding(text string) coding.Encoding {
	// Pick data encoding, smallest first.
	// We could split the string and use different encodings
	// but that seems like overkill for now.
	switch {
	case coding.Num(text).Check() == nil:
		return coding.Num(text)
	case coding.Alpha(text).Check() == nil:
		return coding.Alpha(text)
	}
	return coding.String(text)
}

// Bytes returns, for debugging purposes, a binary representation of the
// list of QR data encoding schemes at the given error correction level.
func Bytes(list []coding.Encoding, l coding.Level, v coding.Version) ([]byte, error) {
	var b coding.Bits
	for _, t := range list {
		if err := t.Check(); err != nil {
			return nil, err
		}
		t.Encode(&b, v)
	}
	b.AddCheckBytes(v, l)
	if n := b.Bits() % 8; n != 0 {
		b.Write(0, 8-n)
	}
	return b.Bytes(), nil
}
