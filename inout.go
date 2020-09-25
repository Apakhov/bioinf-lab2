package main

import (
	"fmt"
	"io"
	"lab2/sequence"
	"log"
	"math"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func readSeqsFromFile(filename string) ([]*AminoSequence, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "opening file "+filename)
	}
	seqs := make([]*AminoSequence, 0)
	p := NewFastaParser(f)
	for {
		seq, err := p.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			fatal("processing error: %s", err)
		}
		seqs = append(seqs, seq)
	}
	return seqs, nil
}

func readSeqsFromFiles(files []string) (*AminoSequence, *AminoSequence) {
	var seq1, seq2 *AminoSequence
	if len(files) == 0 || len(files) > 2 {
		fatal("bad amount of files - %d", len(files))
	}
	if len(files) == 1 {
		seqs, err := readSeqsFromFile(files[0])
		if err != nil {
			fatal(err.Error())
		}
		if len(seqs) != 2 {
			fatal("bad amount of sequeces %d", len(seqs))
		}
		seq1, seq2 = seqs[0], seqs[1]
	}
	if len(files) == 2 {
		seqs1, err := readSeqsFromFile(files[0])
		if err != nil {
			fatal(err.Error())
		}
		seqs2, err := readSeqsFromFile(files[1])
		if err != nil {
			fatal(err.Error())
		}
		if len(seqs1) != 1 || len(seqs2) != 1 {
			fatal("bad amount of sequeces %d", len(seqs1)+len(seqs2))
		}
		seq1, seq2 = seqs1[0], seqs2[1]
	}
	return seq1, seq2
}

func formatRes(alg sequence.Alligner, res1, res2 string, v int, withColor bool) string {
	bld1 := strings.Builder{}
	bldMid := strings.Builder{}
	bld2 := strings.Builder{}
	bld1.WriteString("seq1: ")
	bldMid.WriteString("      ")
	bld2.WriteString("seq2: ")

	matchColor := color.New(color.FgBlue)
	mismatchColor := color.New(color.FgGreen)
	gapColor := color.New(color.FgRed)
	if !withColor || noColor {
		matchColor.DisableColor()
		mismatchColor.DisableColor()
		gapColor.DisableColor()
	}
	i := 0
	l := int(math.Min(float64(len(res1)), float64(len(res2))))
	for ; i < l; i++ {
		col := mismatchColor
		conn := "."
		if res1[i] == alg.Gap() || res2[i] == alg.Gap() {
			col = gapColor
			conn = " "
		} else if res1[i] == res2[i] {
			col = matchColor
			conn = "|"
		}
		if outAlignment != 0 && (i)%int(outAlignment) == 0 {
			bld1.WriteByte('\n')
			bld2.WriteByte('\n')
		}
		col.Fprintf(&bld1, "%c", res1[i])
		col.Fprintf(&bldMid, conn)
		col.Fprintf(&bld2, "%c", res2[i])
	}
	if !noConnectios && outAlignment == 0 {
		bld1.WriteByte('\n')
		bld1.WriteString(bldMid.String())
	}
	bld1.WriteByte('\n')
	bld1.WriteString(bld2.String())
	bld1.WriteByte('\n')
	bld1.WriteString(fmt.Sprintf("score: %d\n", v))

	return bld1.String()
}

func printRes(alg sequence.Alligner, res1, res2 string, v int) {
	if outFile != "" {
		f, err := os.Create(outFile)
		defer f.Close()
		if err != nil {
			log.Fatal(errors.Wrap(err, "opening file "+outFile).Error())
		}
		fmt.Fprintf(f, formatRes(alg, res1, res2, v, false))
		return
	}
	fmt.Print(formatRes(alg, res1, res2, v, true))
}
