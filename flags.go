package main

import "flag"

func init() {
	flag.StringVar(&tableType, "type", useDefault, "table type, one of Blosum64, DNA, Default")
	flag.StringVar(&tableType, "t", useDefault, "table type, one of Blosum64, DNA, Default")
	flag.StringVar(&outFile, "out", "", "output file")
	flag.StringVar(&outFile, "o", "", "output file")
	flag.Float64Var(&gap, "gap", -2, "gap value")
	flag.Float64Var(&gap, "g", -2, "gap value")
	flag.BoolVar(&noColor, "no-color", false, "disables colored output in cosole")
	flag.BoolVar(&noConnectios, "no-connections", false, "disables connections in output")
}
