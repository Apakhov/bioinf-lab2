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
	flag.BoolVar(&noColor, "no-color", false, "disables colored output in cosole")
	flag.BoolVar(&noConnectios, "no-connections", false, "disables connections in output")
	flag.BoolVar(&logTime, "log-time", false, "print time of processing in log")
	flag.IntVar(&amThreads, "threads", 8, "amount of threads for computing, for optimal speed use available amount of cpu")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %[1]s:\n%[1]s {-flag [val]} file [file2]\n", os.Args[0])
		flag.PrintDefaults()
	}
}
