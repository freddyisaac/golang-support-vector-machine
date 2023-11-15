package main

import (
	"fmt"
	"github.com/lukpank/go-glpk/glpk"
	"log"
	"math"
)

func createPhase(d []float64, E int) [][]float64 {
	N := len(d)
	L := N - E + 1
	T := make([][]float64, E)
	for i := 0; i < E; i++ {
		T[i] = make([]float64, L)
	}
	for i := 0; i < E; i++ {
		for j := 0; j < L; j++ {
			T[i][j] = d[i+j]
		}
	}
	return T
}

func phaseWeight(A [][]float64, E float64) [][]float64 {
	m := len(A)
	n := len(A[0])
	scale := 1.0 / E
	newA := makeM(m, n)
	for j := 0; j < n; j++ {
		colSum := sumCol(A, j) * scale
		for i := 0; i < m; i++ {
			newA[i][j] = A[i][j] - colSum
		}
	}
	return newA
}

func Kernel(A [][]float64, nu float64) [][]float64 {
	m := len(A[0])
	K := makeM(m, m)
	gamma := 2. * nu * nu
	for i := 0; i < m; i++ {
		//      fmt.Printf("M[%d,%d] ", i, i)
		K[i][i] = 1.0
		for j := i + 1; j < m; j++ {
			//          fmt.Printf("M[%d,%d] ", i, j)
			dprod := prodV(A, i, j)
			dprod /= -gamma
			val := math.Exp(dprod)
			K[i][j] = val
			K[j][i] = val
		}
		println()
	}
	return K
}

func KernelSums(A [][]float64) []float64 {
	Ksums := make([]float64, len(A))
	for i := 0; i < len(A); i++ {
		row := A[i]
		sum := float64(0.0)
		for j := 0; j < len(row); j++ {
			sum = sum + row[j]
		}
		Ksums[i] = sum
	}
	return Ksums
}

type SVMOpts struct {
	WriteLP string
	Debug   bool
}

func SolveSVM(K [][]float64, Ksum []float64, opts SVMOpts) (*glpk.Prob, error) {
	m := len(Ksum)
	if opts.Debug {
		log.Printf("Solve solving for %d alphas\n", m)
	}
	lp := glpk.New()
	lp.SetProbName("svm")
	lp.SetObjName("F")
	lp.SetObjDir(glpk.MIN)

	// prepare m bound
	// ALPHA(j) * K(X(i),X(j)) + b >= 0
	lp.AddRows(m + 1)
	if opts.Debug {
		log.Printf("Solve Adding %d rows\n", m+1)
	}
	for i := 0; i < m; i++ {
		lp.SetRowName(i+1, fmt.Sprintf("BV%d", i+1))
		lp.SetRowBnds(i+1, glpk.LO, 0.0, 0.0)
	}
	// unity SUM(alpha(i)) = 1
	lp.SetRowName(m+1, "unity")
	lp.SetRowBnds(m+1, glpk.FX, 1.0, 1.0)

	// set up bound cols
	// alpha(i) >= 0.0
	// coefficients
	// there are m alphas and one for the basis b
	lp.AddCols(m + 1)
	if opts.Debug {
		log.Printf("Solve adding %d cols\n", m+1)
	}
	for i := 0; i < m; i++ {
		n := fmt.Sprintf("alpha%d", i+1)
		if opts.Debug {
			fmt.Printf("setting col %s\n", n)
		}
		lp.SetColName(i+1, n)
		lp.SetColBnds(i+1, glpk.LO, 0.0, 0.0)
		// coefficient for objective function
		lp.SetObjCoef(i+1, Ksum[i])
	}
	// add our coeff for b
	lp.SetColName(m+1, "b")
	lp.SetObjCoef(m+1, float64(m))
	lp.SetColBnds(m+1, glpk.DB, -100.0, 100.0)

	// set up m constraints
	// set up indices
	ind := make([]int32, m+2)
	for i := 0; i < m+2; i++ {
		ind[i] = int32(i)
	}
	// distinct ??
	tmpRow := make([]float64, m+2)
	for i := 0; i < m; i++ {
		tmpRow[0] = 0.0
		for j := 0; j < m; j++ {
			tmpRow[j+1] = K[i][j]
		}
		tmpRow[m+1] = 1.0
		if opts.Debug {
			fmt.Printf("SetMatRow(%d %+v %+v\n", i+1, ind, tmpRow)
		}
		lp.SetMatRow(i+1, ind, tmpRow)
	}
	tmpRow[0] = 0.0
	tmpRow[m+1] = 0.0
	for i := 1; i < m+1; i++ {
		tmpRow[i] = 1.0
	}
	if opts.Debug {
		fmt.Printf("SetMatRow(%d %+v %+v\n", m+1, ind, tmpRow)
	}
	lp.SetMatRow(m+1, ind, tmpRow)

	// set our matrix rows

	err := lp.Simplex(nil)
	if err != nil {
		return nil, err
	}

	if opts.WriteLP != "" {
		lp.WriteLP(nil, opts.WriteLP)
	}
	return lp, nil
}

func ExtractSolution(lp *glpk.Prob, m int) ([]float64, float64) {
	alpha := make([]float64, m)
	var b float64
	for i := 0; i < m; i++ {
		alpha[i] = lp.ColPrim(i + 1)
	}
	b = lp.ColPrim(m + 1)
	return alpha, b
}

const (
	EPSERROR = 1.0
)

// add m - eps(i) errors variables
func SolveSVMWithErrors(K [][]float64, Ksum []float64, lambda float64, opts SVMOpts) (*glpk.Prob, error) {
	m := len(Ksum)
	if opts.Debug {
		log.Printf("Solve solving for %d alphas\n", m)
	}
	lp := glpk.New()
	lp.SetProbName("svm")
	lp.SetObjName("F")
	lp.SetObjDir(glpk.MIN)

	// prepare m bound
	// ALPHA(j) * K(X(i),X(j)) + b + eps(i) >= 0
	lp.AddRows(m + 1)
	if opts.Debug {
		log.Printf("Solve Adding %d rows\n", m+1)
	}
	for i := 0; i < m; i++ {
		lp.SetRowName(i+1, fmt.Sprintf("BV%d", i+1))
		lp.SetRowBnds(i+1, glpk.LO, 0.0, 0.0)
	}
	// unity SUM(alpha(i)) = 1
	lp.SetRowName(m+1, "unity")
	lp.SetRowBnds(m+1, glpk.FX, 1.0, 1.0)

	// set up bound cols
	// alpha(i) >= 0.0
	// coefficients
	// there are m alphas and one for the basis b
	//	lp.AddCols(m + 1)

	// now we have m alphas and m eps
	lp.AddCols(2*m + 1)
	if opts.Debug {
		log.Printf("Solve adding %d cols\n", m+1)
	}

	// set up column coefficients
	// start with alphas
	for i := 0; i < m; i++ {
		alpha := fmt.Sprintf("alpha%d", i+1)
		if opts.Debug {
			fmt.Printf("setting col %s\n", alpha)
		}
		lp.SetColName(i+1, alpha)
		lp.SetColBnds(i+1, glpk.LO, 0.0, 0.0)
		// coefficient for objective function
		lp.SetObjCoef(i+1, Ksum[i])
	}
	// add our coeff for b
	lp.SetColName(m+1, "b")
	lp.SetObjCoef(m+1, float64(m))
	lp.SetColBnds(m+1, glpk.DB, -100.0, 100.0)

	// add coeff for our errors variables eps(i)
	for i := 0; i < m; i++ {
		eps := fmt.Sprintf("eps%d", i+1)
		lp.SetColName(m+i+2, eps)
		// no col bounds on eps(i)
		lp.SetColBnds(m+i+2, glpk.DB, -EPSERROR, EPSERROR)
		lp.SetObjCoef(m+i+2, lambda)
	}

	// set up m constraints
	// set up indices

	// FEI FIGURE THIS OUT FOR OUR eps(i)
	//	ind := make([]int32, m+2)
	ind := make([]int32, 2*m+2)
	for i := 0; i < 2*m+2; i++ {
		ind[i] = int32(i)
	}
	// distinct ??
	tmpRow := make([]float64, 2*m+2)
	for i := 0; i < m; i++ {
		tmpRow[0] = 0.0
		for j := 0; j < m; j++ {
			tmpRow[j+1] = K[i][j]
		}
		tmpRow[m+1] = 1.0
		tmpRow[m+i+2] = lambda
		if opts.Debug {
			fmt.Printf("SetMatRow(%d %+v %+v\n", i+1, ind, tmpRow)
		}
		lp.SetMatRow(i+1, ind, tmpRow)
		// unset eps value for this row
		tmpRow[m+i+2] = 0.0
	}
	tmpRow[0] = 0.0
	tmpRow[m+1] = 0.0
	for i := 1; i < m+1; i++ {
		tmpRow[i] = 1.0
	}
	if opts.Debug {
		fmt.Printf("SetMatRow(%d %+v %+v\n", m+1, ind, tmpRow)
	}
	lp.SetMatRow(m+1, ind, tmpRow)

	// set our matrix rows

	if opts.WriteLP != "" {
		lp.WriteLP(nil, opts.WriteLP)
	}

	return lp, nil

	err := lp.Simplex(nil)
	if err != nil {
		return nil, err
	}

	return lp, nil
}

type F func() ([][]float64, []float64)

var nilF F = func() ([][]float64, []float64) { return nil, nil }

func Solve(data [][]float64, nu float64, wantLP string) (*glpk.Prob, F, error) {
	K := Kernel(data, nu)
	Ksums := KernelSums(K)
	m := len(Ksums)
	opts := SVMOpts{}
	if wantLP != "" {
		opts.WriteLP = wantLP
	}
	lp, err := SolveSVM(K, Ksums, opts)
	DisplaySolution(lp, m)
	f := func() ([][]float64, []float64) {
		return K, Ksums
	}
	return lp, f, err
}
