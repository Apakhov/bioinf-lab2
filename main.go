package main

import (
	"flag"
	"lab2/sequence"
)

func main() {
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		fatal("bad amount of files - 0")
	}

	seq1, seq2 := readSeqsFromFiles(files)

	var allg sequence.Alligner
	switch tableType {
	case useBlosum:
		allg = sequence.NewAlligerBLOSUM62(gap)
	case useDefault:
		allg = sequence.NewDefault(gap)
	case useDNA:
		allg = sequence.NewAlligerDNA(gap)
	default:
		fatal("bad table type %s", tableType)
	}
	res1, res2, v, err := sequence.Allign(allg, seq1.Value, seq2.Value)
	if err != nil {
		fatal("alligning %s", err.Error())
	}
	printRes(allg, res1, res2, int(v))
}