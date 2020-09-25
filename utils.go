package main

import (
	"flag"
	"log"
)

const (
	useBlosum  = "Blosum64"
	useDNA     = "DNA"
	useDefault = "Default"
)

var (
	tableType    string
	gap          float64
	outFile      string
	noColor      bool
	noConnectios bool
	logTime      bool
	amThreads    int
	outAlignment uint
)

func fatal(format string, v ...interface{}) {
	flag.Usage()
	log.Fatalf(format, v...)
}
