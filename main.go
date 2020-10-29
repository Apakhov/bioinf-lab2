package main

import (
	"flag"
	"lab2/sequence"
	"log"
	"time"
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
	allign := sequence.Allign
	if memOpt {
		allign = sequence.AllignMemoryOpt
		log.Println("using mem-opt")
	}
	t := time.Now()
	res1, res2, v, err := allign(allg, seq1.Value, seq2.Value, amThreads)
	if logTime {
		log.Print("calculation time: ", time.Now().Sub(t))
	}
	if err != nil {
		fatal("alligning %s", err.Error())
	}
	printRes(allg, res1, res2, int(v))
}
