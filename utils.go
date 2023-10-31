package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/lukpank/go-glpk/glpk"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func makeArray(m, n int) [][]float64 {
    data := make([][]float64, m)
    for i := 0; i < m; i++ {
        data[i] = make([]float64, n)
    }
    return data
}

func makeData(m, n int, sigma, meanx, meany float64) [][]float64 {
    data := makeArray(m, n)
	var means [2]float64
	means[0] = meanx
	means[1] = meany
	for j := 0; j < n; j++ {
		for i := 0; i < m; i++ {
            data[i][j] = r.NormFloat64()*sigma + means[i]
        }
    }
    return data
}

func appendData(x, y [][]float64) [][]float64 {
    for i:=0;i<len(x);i++ {
        x[i] = append(x[i], y[i]...)
    }
    return x
}

func writeData(name string, data [][]float64) error {
    fd, err := os.Create(name)
    if err != nil {
        return err
    }
    defer fd.Close()
    for j:=0;j<len(data[0]);j++ {
        for i:=0;i<len(data);i++ {
            fmt.Fprintf(fd, "%f ", data[i][j])
        }
        fmt.Fprintf(fd, "\n")
    }
    return nil
}

func readData(name string, n int) ([][]float64, error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	data := make([][]float64, 2)
	count := n
	for {
		var x, y float64
		_, err := fmt.Fscanf(fd, "%f %f\n", &x, &y)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		// inefficient
		data[0] = append(data[0], x)
		data[1] = append(data[1], y)
		count--
		if count == 0 {
			break
		}
	}
	return data, nil
}


func DisplaySolution(lp *glpk.Prob, m int) {
    fmt.Printf("%s = %g\n", lp.ObjName(), lp.ObjVal())
    for i := 0; i < m+1; i++ {
        if Fsmall(lp.ColPrim(i+1)) == 0.0 {
            continue
        }
        fmt.Printf("[%d] %s = %g\n", i+1, lp.ColName(i+1), lp.ColPrim(i+1))
    }
    println()
}
