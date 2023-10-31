package main

import (
	"fmt"
	"math"
)

const (
	EPSILON = 1.0e-5
)

func VfromM(M [][]float64, j int) []float64 {
	m := len(M)
	x := make([]float64, m)
	for i := 0; i < m; i++ {
		x[i] = M[i][j]
	}
	return x
}

func makeM(m, n int) [][]float64 {
	arr := make([][]float64, m)
	for i := 0; i < m; i++ {
		arr[i] = make([]float64, n)
	}
	return arr
}

func prod(b []float64, c [][]float64, j int) float64 {
	sum := 0.0
	for k := 0; k < len(b); k++ {
		sum += b[k] * c[k][j]
	}
	return sum
}

func mult(b, c [][]float64) [][]float64 {
	m := len(b)
	n := len(c[0])
	a := makeM(m, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			a[i][j] = prod(b[i], c, j)
		}
	}
	return a
}

func testMult() {
	B := [][]float64{
		{1.0, 2.0},
		{2.0, 1.0},
	}
	C := [][]float64{
		{1.0, 1.0},
		{2.0, 3.0},
	}
	A := mult(B, C)
	displayM(A)
}

// return sum of jth column
func sumCol(A [][]float64, j int) float64 {
	m := len(A)
	sum := 0.0
	for i := 0; i < m; i++ {
		sum += A[i][j]
	}
	return sum
}

// rturn phase array from v with coeff E
// a(i,j) = a(i,j) * ( delta(i,j) - 1/E )
func phaseV(v []float64, E float64) []float64 {
	m := len(v)
	nv := make([]float64, m)
	sumE := -1.0 * float64(m) / E
	var sumV float64
	for i := 0; i < m; i++ {
		sumV = sumV + v[i]
	}
	scaleSumV := sumE * sumV
	for i := 0; i < m; i++ {
		nv[i] = v[i] - scaleSumV
	}
	return nv
}

// calculate | x(i,j) - x(i,k) |^2
func prodV(M [][]float64, j, k int) float64 {
	m := len(M)
	sum := 0.0
	for i := 0; i < m; i++ {
		v := M[i][j] - M[i][k]
		sum += v * v
	}
	return sum
}

// |x|
func Fabs(f float64) float64 {
	if f < 0.0 {
		return -f
	}
	return f
}

// |x| < epsilon
func Fsmall(f float64) float64 {
	if Fabs(f) < EPSILON {
		return 0.0
	}
	return f
}

func displayM(M [][]float64) {
	for i := range M {
		fmt.Printf("row %d : ", i)
		for j := range M[i] {
			fmt.Printf("%v ", M[i][j])
		}
		println()
	}
}

func displayV(V []float64) {
	for i := 0; i < len(V); i++ {
		fmt.Printf("%v ", V[i])
	}
	println()
}

// calculate  rbf K(X,X(j)
func EvalK(X [][]float64, j int, x []float64, sigma float64) float64 {
	m := len(x)
	gamma := 1.0 / (2.0 * sigma * sigma)
	var dprod float64
	for i := 0; i < m; i++ {
		v := X[i][j] - x[i]
		dprod += v * v
	}
	val := math.Exp(-gamma * dprod)
	//fmt.Printf("EvalK X%+v : %f %g\n", x, -gamma * dprod, val)
	return val
}

func EvalF(X [][]float64, x, alpha []float64, b, nu float64) float64 {
	var f float64
	f = 0.0
	for j := 0; j < len(alpha); j++ {
		k := EvalK(X, j, x, nu)
		f += alpha[j] * k
	}
	return f + b
}
