package sequence

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testAlligner() Alligner {
	return &tableAlliger{
		gapVal: -5,
		table: [][]float64{
			{5, -4, -4, -4},
			{-4, 5, -4, -4},
			{-4, -4, 5, -4},
			{-4, -4, -4, 5},
		},
		byteToIdx: map[byte]int{
			'A': 0,
			'B': 1,
			'C': 2,
			'D': 3,
		},
	}
}

type testCaseAllg struct {
	a     string
	b     string
	resA  string
	resB  string
	score float64
}

func runTestCases(t *testing.T, allg Alligner, tcs []testCaseAllg) {
	for i, tc := range tcs {
		resA, resB, score, err := Allign(allg, tc.a, tc.b)
		require.NoError(t, err)
		require.Equal(t, tc.resA, resA, "failed seq a test %d", i)
		require.Equal(t, tc.resB, resB, "failed seq b test %d", i)
		require.Equal(t, tc.score, score, "failed score test %d", i)
	}
}

func TestAllg(t *testing.T) {
	allg := testAlligner()
	tcs := []testCaseAllg{
		{
			a:     "AAA",
			b:     "AAA",
			resA:  "AAA",
			resB:  "AAA",
			score: 15,
		},
		{
			a:     "AAAA",
			b:     "AAAB",
			resA:  "AAAA",
			resB:  "AAAB",
			score: 11,
		},
		{
			a:     "AAAA",
			b:     "AAA",
			resA:  "AAAA",
			resB:  "-AAA",
			score: 10,
		},
		{
			a:     "ABBBCDDD",
			b:     "ABCDDD",
			resA:  "ABBBCDDD",
			resB:  "A--BCDDD",
			score: 20,
		},
		{
			a:     "ABCDDD",
			b:     "ABBBCDDD",
			resA:  "A--BCDDD",
			resB:  "ABBBCDDD",
			score: 20,
		},
		{
			a:     "AACD",
			b:     "BCD",
			resA:  "AACD",
			resB:  "-BCD",
			score: 1,
		},
	}
	runTestCases(t, allg, tcs)
}

func TestAllgBlosum(t *testing.T) {
	allg := NewAlligerBLOSUM62(-10)
	tcs := []testCaseAllg{
		{ // verified on https://www.ebi.ac.uk/Tools/psa/emboss_needle/
			a:     "SPETVIHSGWVIWRELFSHWPDQCKLLFGDWFAWIHWTYLVYYSAGPPCQGQSDIVVMMQKKLRTNFCQCYKYWYQ",
			b:     "SPSDQFFTVIHSCLYWVIWRDLMSHLFMNGAAIDIHWTWDSIAIGPPLVYPIEEVFAGPSTIVVMMQKMLRTNFCQCYKPWYQ",
			resA:  "SP--E--TVIHS--GWVIWRELFSH-WPDQCKL-LFGDWFAWIHWTYLVYYSAGPPCQGQSDIVVMMQKKLRTNFCQCYKYWYQ",
			resB:  "SPSDQFFTVIHSCLYWVIWRDLMSHLFMNGAAIDIHWTWDSIAIGPPLV-YPIEEVFAGPSTIVVMMQKMLRTNFCQCYKPWYQ",
			score: 116,
		},
	}
	runTestCases(t, allg, tcs)
}

func TestAllgDNA(t *testing.T) {
	allg := NewAlligerDNA(-10)
	tcs := []testCaseAllg{
		{
			a:     "GCGCGTGCGCGGAAGGAGCCAAGGTGAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGGCACCATGCTGTCCCCCGAGGCGGAGCGGGTGCTGCGGTACCTGGTCGAAGTAGAGGAGTTG",
			b:     "GACTTGTGGAACCTACTTCCTGAAAATAACCTTCTGTCCTCCGAGCTCTCCGCACCCGTGGATGACCTGCTCCCGTACACAGATGTTGCCACCTGGCTGGATGAATGTCCGAATGAAGCG",
			resA:  "GCGCGTGCGCGGAAGGAGCCAAGGTGAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGGCA-CCAT-GCTGTCCCCCGAGGCGGA-GCGGGTGCTG-C-GGTACCTGGTCGAA-GT-AG-AGGAGTTG",
			resB:  "G-AC-T-TGTGGAA-CCTACTTCCTGAA--AATAACCTTCTGTCCTCCGAGCT-CTCCGCACCCGTGGATGACCTGC-TCCCGTACACAGATGTTGCCACCTGGCTGGATGAATGTCCGAATGAAGCG",
			score: -41,
		},
	}
	runTestCases(t, allg, tcs)
}
