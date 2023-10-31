package main

import (
	"flag"
	"log"
)

var ofile *string = flag.String("o", "data.txt", "output file name")
var ifile *string = flag.String("i", "", "input file name")
var meanx *float64 = flag.Float64("mx", 0.2, "mean x pos")
var meany *float64 = flag.Float64("my", 0.2, "mean y pos")
var num *int = flag.Int("n", 20, "num of 2d data points to generate")
var sigma = flag.Float64("s", 0.1, "sigma for distribution")

func main() {
	flag.Parse()

	m := 2
	//	num := 20
	//	meanx := 0.2
	//	meany := 0.2
	//	sigma := 0.1

	data := [][]float64{}
	rdata := [][]float64{}
	if *ifile != "" {
		// fail silently on error
		rdata, _ = readData(*ifile, -1)
	}

	wdata := makeData(m, *num, *sigma, *meanx, *meany)
	if len(rdata) == 0 {
		data = wdata
	} else {
		// append shortest
		if len(wdata) < len(rdata) {
			rdata, wdata = wdata, rdata
		}
		data = appendData(rdata, wdata)
	}

	err := writeData(*ofile, data)
	if err != nil {
		log.Printf("unable to write output file error : %+v\n", err)
	}

}
