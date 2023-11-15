package main

import (
	"flag"
	"fmt"

	"io"
	"os"
)

var wantLP string

var debug *string = flag.String("d", "", "output LP problem")
var orig *bool = flag.Bool("o", false, "eval F for original data")
var tffile *string = flag.String("tff", "", "output for training set function file")

var ofile *string = flag.String("of", "", "output file")
var ifile *string = flag.String("if", "data.txt", "input data file")

var lbx *float64 = flag.Float64("lbx", -.5, "lower bound for X")
var ubx *float64 = flag.Float64("ubx", .5, "upper bound for X")
var num *int = flag.Int("n", 50, "X step increment")
var lby *float64 = flag.Float64("lby", -.5, "lower bound for X")
var uby *float64 = flag.Float64("uby", .5, "upper bound for X")

var flagNu *float64 = flag.Float64("nu", 0.05, "numerical parameter nu")

func main() {
	flag.Parse()
	if *debug != "" {
		wantLP = *debug
	}
	// nu := 0.05 // rbf fitness param
	nu := *flagNu

	data, err := ReadData(*ifile, -1)
	if err != nil {
		fmt.Printf("unable to read from %s error %v\n", *ifile, err)
		return
	}
	n := len(data[0])
	fmt.Printf("data read %d points\n", n)
	fmt.Printf("data read from %s\n", *ifile)
	lp, _, err := Solve(data, nu, wantLP)
	if err != nil {
		fmt.Printf("lp error : %+v\n", err)
		return
	}
	alpha, b := ExtractSolution(lp, n)

	if *orig {
		fd, err := os.Create("orig.txt")
		if err != nil {
			fd = os.Stdout
		}
		defer fd.Close()
		var datum [2]float64
		X := datum[:]
		for i := 0; i < n; i++ {
			X[0] = data[0][i]
			X[1] = data[1][i]
			f := EvaluateF(data, X, alpha, b, nu)
			fmt.Fprintf(fd, "%f %f %f\n", X[0], X[1], f)
		}

	}

	if *ofile != ".orig" {

		var ofd io.WriteCloser
		if *ofile == "" {
			ofd = os.Stdout
		} else {
			ofd, err = os.Create(*ofile)
			if err != nil {
				fmt.Printf("unable to open %s for writing %+v\n", *ofile)
				ofd = os.Stdout
			}
			defer ofd.Close()
		}

		XMIN := float64(*lbx)
		XMAX := float64(*ubx)
		nx := *num
		YMIN := float64(*lby)
		YMAX := float64(*uby)
		ny := *num

		x := XMIN
		dx := (XMAX - XMIN) / float64(nx-1)
		y := YMIN
		dy := (YMAX - YMIN) / float64(ny-1)

		fmt.Printf("Grid : (%f,%f), (%f,%f) (%f,%f)\n", XMIN, XMAX, YMIN, YMAX, dx, dy)

		var datum [2]float64
		X := datum[:]
		for j := 0; j < ny; j++ {
			x = XMIN
			for i := 0; i < nx; i++ {
				X[0] = x
				X[1] = y
				f := EvaluateF(data, X, alpha, b, nu)
				f = Fsmall(f)
				if f >= 0.0 {
					fmt.Printf("x[%d]: %f y[%d]: %f -> %g\n", i, x, j, y, f)
				}
				if *ofile != "" {
					fmt.Fprintf(ofd, "%f %f %f\n", x, y, f)
				}
				x += dx
			}
			y += dy
		}
	}
}
