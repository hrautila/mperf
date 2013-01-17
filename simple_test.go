
// Copyright (c) Harri Rautila, 2013

// This file is part of github.com/hrautila/mperf package. It is free software,
// distributed under the terms of GNU Lesser General Public License Version 3, or
// any later version. See the COPYING tile included in this archive.


package mperf

import (
    "github.com/hrautila/matrix"
	"testing"
)

func testAndData(m, n, p int)(fnc func(), A, B, C *matrix.FloatMatrix) {
    A, B, C = MakeData(m, n, p, true, false)
    fnc = func() {
        C.Plus(A.Times(B))
    }
    return
}

func TestSimple(t *testing.T) {
    A, B, C := MakeData(1000, 1000, 1000, true, false)
    fnc := func() {
        C.Plus(A.Times(B))
    }
    FlushCache()
    tm := Timeit(fnc)
    t.Logf("execution time: %v\n", tm)
}

func TestSingle(t *testing.T) {
    tm, _ := SingleTest("times", testAndData, 1000, 1000, 1000, true, true)
    t.Logf("execution time: %v\n", tm)
}

func TestMultiple(t *testing.T) {
    sizes := []int{400, 600, 800}
    data := MultipleSizeTests(testAndData, sizes, 3, true)
    for _, sz := range sizes {
        t.Logf("%d: %.3fmsec\n", sz, data[sz]*1000.0)
    }
}
    
// Local Variables:
// tab-width: 4
// indent-tabs-mode: nil
// End:
