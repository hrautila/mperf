
// Copyright (c) Harri Rautila, 2013

// This file is part of github.com/hrautila/mperf package. It is free software,
// distributed under the terms of GNU Lesser General Public License Version 3, or
// any later version. See the COPYING tile included in this archive.

package mperf

import (
    "github.com/hrautila/matrix"
    "github.com/hrautila/linalg/blas"
    "time"
    "fmt"
    "os"
)

var nullData []float64

// Write and read to a large array to force matrix data out of cache
func FlushCache() {
    if len(nullData) == 0 {
        nullData = make([]float64, 1500*1500)
    }
    for i, _ := range nullData {
        nullData[i] = 1e-10
    }
    zero := 0.0
    for i, _ := range nullData {
        zero += nullData[i]
    }
    if zero < 10 {
        zero = 0.0
    }
}


func MakeData(M, N, P int, randomData, diagonal bool) (A, B, C *matrix.FloatMatrix) {

    if diagonal && M != N {
        diagonal = false
        fmt.Printf("cannot make diagonal if B.rows != B.cols\n")
    }
	
	if randomData {
        A = matrix.FloatNormal(M, P)
        if diagonal {
            d := matrix.FloatNormal(P, 1)
            B := matrix.FloatDiagonal(P, 0.0)
            B.SetIndexesFromArray(d.FloatArray(), matrix.DiagonalIndexes(B)...)
        } else {
            B = matrix.FloatNormal(P, N)
        }
    } else {
        A = matrix.FloatWithValue(M, P, 1.0)
        if diagonal {
            B = matrix.FloatDiagonal(P, 1.0)
        } else {
            B = matrix.FloatWithValue(P, N, 1.0)
        }
    }
    C = matrix.FloatZeros(M, N)
    return
}

func Timeit(fn func ()) time.Duration {
    t1 := time.Now()
    fn()
    t2 := time.Now()
    return t2.Sub(t1)
}

func Check(A, B, C0 *matrix.FloatMatrix) (dt time.Duration, result bool) {
    C := matrix.FloatZeros(A.Rows(), B.Cols())
    fnc := func() {
        blas.GemmFloat(A, B, C, 1.0, 1.0)
    }
    FlushCache()
    dt = Timeit(fnc)
    result = C0.AllClose(C)
    return
}

func CheckWithFunc(A, B, C0 *matrix.FloatMatrix, cfunc MatrixCheckFunc) (dt time.Duration, result bool) {
    C := matrix.FloatZeros(C0.Rows(), C0.Cols())
    fnc := func() {
        cfunc(A, B, C)
    }
    FlushCache()
    dt = Timeit(fnc)
    result = C0.AllClose(C)
    return
}


// executions times for matrix sizes
type Timings map[int]float64

// Check function takes 3 matrix arguments (A, B and C) and return in the 3rd argument
// matrix product of the two first matrices. (Blas GEMM matrix arguments)
type MatrixCheckFunc func(*matrix.FloatMatrix,*matrix.FloatMatrix,*matrix.FloatMatrix)

// function that returns testable function and its data matrices A, B and C
type MatrixTestFunc func(int, int, int)(func(), *matrix.FloatMatrix,*matrix.FloatMatrix,*matrix.FloatMatrix,)

// Run single test and return elapsed seconds
func SingleTest(name string, testAndData MatrixTestFunc, m, n, p int, check, verbose bool) (float64, bool) {
    result := true
    fnc, A, B, C := testAndData(m, n, p)
    FlushCache()
    tm := Timeit(fnc)
    if check {
        reftime, ok := Check(A, B, C)
        if verbose {
            fmt.Fprintf(os.Stderr, "%s: %v\n", name, tm)
            fmt.Fprintf(os.Stderr, "Reference: [%v] %v (%.2f) \n",
                ok, reftime, tm.Seconds()/reftime.Seconds())
        }
        result = ok
    } 
    return tm.Seconds(), result
}


// Run tests for multiple matrix sizes
func MultipleSizeTests(testAndData MatrixTestFunc, sizes []int, testCount int, verbose bool) Timings {
    times := make(Timings, len(sizes))
    for _, sz := range sizes {
        fnc, _, _, _ := testAndData(sz, sz, sz)
        minTime := 0.0
        for i := 0; i < testCount; i++ {
            FlushCache()
            tm := Timeit(fnc)
            if minTime == 0.0 {
                minTime = tm.Seconds()
            } else {
                if tm.Seconds() < minTime {
                    minTime = tm.Seconds()
                }
            }
            if verbose {
                fmt.Fprintf(os.Stderr, "%4d: %v\n", sz, tm)
            }
        }
        times[sz] = minTime
    }
    return times
}


// Local Variables:
// tab-width: 4
// indent-tabs-mode: nil
// End:
