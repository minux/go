// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sys

const (
	TheChar       = 'u'
	BigEndian     = 1
	CacheLineSize = 32
	PhysPageSize  = 8192
	PCQuantum     = 4
	Int64Align    = 8
	HugePageSize  = 0
	MinFrameSize  = 176
	SpAlign       = 16
)

type Uintreg uint64
