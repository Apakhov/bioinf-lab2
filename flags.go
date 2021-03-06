package main

import (
	"flag"
	"fmt"
	"os"
)

func init() {
	flag.StringVar(&tableType, "type", useDefault, "table type, one of Blosum64, DNA, Default")
	flag.StringVar(&tableType, "t", useDefault, "table type, one of Blosum64, DNA, Default")
	flag.StringVar(&outFile, "out", "", "output file")
	flag.StringVar(&outFile, "o", "", "output file")
	flag.Float64Var(&gap, "gap", -2, "gap value")
	flag.Float64Var(&gap, "g", -2, "gap value")
	flag.Float64Var(&gap, "gap-open", -2, "gap value")
	flag.Float64Var(&gapExt, "gap-extend", 0, "gap extand value")
	flag.Float64Var(&gapExt, "ge", 0, "gap extand value")
	flag.BoolVar(&noColor, "no-color", false, "disables colored output in cosole")
	flag.BoolVar(&noConnectios, "no-connections", false, "disables connections in output")
	flag.UintVar(&outAlignment, "outalignment", 0, "alignment of result sequences, if 0 no alignment used")
	flag.UintVar(&outAlignment, "oa", 0, "alignment of result sequences, if 0 no alignment used")
	flag.BoolVar(&logTime, "log-time", false, "print time of processing in log")
	flag.IntVar(&amThreads, "threads", 8, "amount of threads for computing, for optimal speed use available amount of cpu")
	flag.BoolVar(&memOpt, "mem-opt", false, "run with memory usage optimized algorithm. it is slower but uses far less memory")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %[1]s:\n%[1]s {-flag [val]} file [file2]\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func isGapExtPassed() bool {
	return isFlagPassed("gap-extend") || isFlagPassed("ge")
}
