// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/compile/internal/gc"
	"cmd/internal/obj"
	"cmd/internal/obj/sparc64"
)

func blockcopy(n, res *gc.Node, osrc, odst, w int64) {
	// determine alignment.
	// want to avoid unaligned access, so have to use
	// smaller operations for less aligned types.
	// for example moving [4]byte must use 4 MOVB not 1 MOVW.
	align := int(n.Type.Align)

	var op int
	switch align {
	default:
		gc.Fatalf("sgen: invalid alignment %d for %v", align, n.Type)

	case 1:
		op = sparc64.AMOVB

	case 2:
		op = sparc64.AMOVH

	case 4:
		op = sparc64.AMOVW

	case 8:
		op = sparc64.AMOVD
	}

	if w%int64(align) != 0 {
		gc.Fatalf("sgen: unaligned size %d (align=%d) for %v", w, align, n.Type)
	}
	c := int32(w / int64(align))

	if osrc%int64(align) != 0 || odst%int64(align) != 0 {
		gc.Fatalf("sgen: unaligned offset src %d or dst %d (align %d)", osrc, odst, align)
	}

	// if we are copying forward on the stack and
	// the src and dst overlap, then reverse direction
	dir := align

	if osrc < odst && odst < osrc+w {
		dir = -dir
	}

	var dst gc.Node
	var src gc.Node
	if n.Ullman >= res.Ullman {
		gc.Agenr(n, &dst, res) // temporarily use dst
		gc.Regalloc(&src, gc.Types[gc.Tptr], nil)
		gins(sparc64.AMOVD, &dst, &src)
		if res.Op == gc.ONAME {
			gc.Gvardef(res)
		}
		gc.Agen(res, &dst)
	} else {
		if res.Op == gc.ONAME {
			gc.Gvardef(res)
		}
		gc.Agenr(res, &dst, res)
		gc.Agenr(n, &src, nil)
	}

	var tmp, tmp1 gc.Node
	gc.Regalloc(&tmp, gc.Types[gc.Tptr], nil)
	gc.Regalloc(&tmp1, gc.Types[gc.Tptr], nil)

	// set up end marker
	var nend gc.Node

	// move src and dest to the end of block if necessary
	if dir < 0 {
		if c >= 4 {
			gc.Regalloc(&nend, gc.Types[gc.Tptr], nil)
			gins(sparc64.AMOVD, &src, &nend)
		}

		p := gins(sparc64.AADD, nil, &src)
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = w

		p = gins(sparc64.AADD, nil, &dst)
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = w
	} else {
		p := gins(sparc64.AADD, nil, &src)
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = int64(-dir)

		p = gins(sparc64.AADD, nil, &dst)
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = int64(-dir)

		if c >= 4 {
			gc.Regalloc(&nend, gc.Types[gc.Tptr], nil)
			p := gins(sparc64.AMOVD, &src, &nend)
			p.From.Type = obj.TYPE_ADDR
			p.From.Offset = w
		}
	}

	// move
	// TODO: enable duffcopy for larger copies.
	if c >= 4 {
		// TODO(aram): instead of manually updating both src and dst, update
		// only the index register and change the comparison.
		ginscon(sparc64.AMOVD, int64(dir), &tmp1)

		p := gins(op, &src, &tmp)
		p.From.Type = obj.TYPE_MEM
		p.From.Index = tmp1.Reg
		p.From.Scale = 1
		ploop := p

		p = gins(sparc64.AADD, &tmp1, &src)

		p = gins(op, &tmp, &dst)
		p.To.Type = obj.TYPE_MEM
		p.To.Index = tmp1.Reg
		p.To.Scale = 1

		p = gins(sparc64.AADD, &tmp1, &dst)

		p = gcmp(sparc64.ACMP, &src, &nend)

		gc.Patch(gc.Gbranch(sparc64.ABNED, nil, 0), ploop)
		gc.Regfree(&nend)
	} else {
		// TODO(aram): instead of manually updating both src and dst, update
		// only the index register.
		ginscon(sparc64.AMOVD, int64(dir), &tmp1)
		var p *obj.Prog
		for ; c > 0; c-- {
			p = gins(op, &src, &tmp)
			p.From.Type = obj.TYPE_MEM
			p.From.Index = tmp1.Reg
			p.From.Scale = 1

			p = gins(sparc64.AADD, &tmp1, &src)

			p = gins(op, &tmp, &dst)
			p.To.Type = obj.TYPE_MEM
			p.To.Index = tmp1.Reg
			p.To.Scale = 1

			p = gins(sparc64.AADD, &tmp1, &dst)
		}
	}

	gc.Regfree(&dst)
	gc.Regfree(&src)
	gc.Regfree(&tmp)
	gc.Regfree(&tmp1)
}
