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
		resA, resB, score, err := Allign(allg, tc.a, tc.b, 1)
		require.NoError(t, err)
		require.Equal(t, tc.resA, resA, "failed seq A test %d", i)
		require.Equal(t, tc.resB, resB, "failed seq B test %d", i)
		require.Equal(t, tc.score, score, "failed SCORE test %d", i)
	}
	for i, tc := range tcs {
		resA, resB, score, err := Allign(allg, tc.a, tc.b, 8)
		require.NoError(t, err)
		require.Equal(t, tc.resA, resA, "failed seq A test %d", i)
		require.Equal(t, tc.resB, resB, "failed seq B test %d", i)
		require.Equal(t, tc.score, score, "failed SCORE test %d", i)
	}
	for i, tc := range tcs {
		_, _, score, err := AllignMemoryOpt(allg, tc.a, tc.b, 1)
		require.NoError(t, err)
		require.Equal(t, tc.score, score, "failed SCORE test %d", i)
	}
	for i, tc := range tcs {
		_, _, score, err := AllignMemoryOpt(allg, tc.a, tc.b, 8)
		require.NoError(t, err)
		require.Equal(t, tc.score, score, "failed SCORE test %d", i)
	}
}

func TestAllg(t *testing.T) {
	allg := testAlligner()
	tcs := []testCaseAllg{
		{
			a:     "",
			b:     "",
			resA:  "",
			resB:  "",
			score: 0,
		},
		{
			a:     "A",
			b:     "A",
			resA:  "A",
			resB:  "A",
			score: 5,
		},
		{
			a:     "AA",
			b:     "AA",
			resA:  "AA",
			resB:  "AA",
			score: 10,
		},
		{
			a:     "AAA",
			b:     "AAA",
			resA:  "AAA",
			resB:  "AAA",
			score: 15,
		},
		{
			a:     "A",
			b:     "C",
			resA:  "A",
			resB:  "C",
			score: -4,
		},
		{
			a:     "AAA",
			b:     "AAB",
			resA:  "AAA",
			resB:  "AAB",
			score: 6,
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
			a:     "AAC",
			b:     "BCC",
			resA:  "AAC",
			resB:  "BCC",
			score: -3,
		},
		{
			a:     "AA",
			b:     "BC",
			resA:  "AA",
			resB:  "BC",
			score: -8,
		},
		{
			a:     "AA",
			b:     "B",
			resA:  "AA",
			resB:  "-B",
			score: -9,
		},
		{
			a:     "AACD",
			b:     "BCD",
			resA:  "AACD",
			resB:  "-BCD",
			score: 1,
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
		{ // verified on https://www.ebi.ac.uk/Tools/psa/emboss_needle/
			a:     "GCGCGTGCGCGGAAGGAGCCAAGGTGAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGGCACCATGCTGTCCCCCGAGGCGGAGCGGGTGCTGCGGTACCTGGTCGAAGTAGAGGAGTTG",
			b:     "GACTTGTGGAACCTACTTCCTGAAAATAACCTTCTGTCCTCCGAGCTCTCCGCACCCGTGGATGACCTGCTCCCGTACACAGATGTTGCCACCTGGCTGGATGAATGTCCGAATGAAGCG",
			resA:  "GCGCGTGCGCGGAAGGAGCCAAGGTGAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGGCA-CCAT-GCTGTCCCCCGAGGCGGA-GCGGGTGCTG-C-GGTACCTGGTCGAA-GT-AG-AGGAGTTG",
			resB:  "G-AC-T-TGTGGAA-CCTACTTCCTGAA--AATAACCTTCTGTCCTCCGAGCT-CTCCGCACCCGTGGATGACCTGC-TCCCGTACACAGATGTTGCCACCTGGCTGGATGAATGTCCGAATGAAGCG",
			score: -41,
		},
	}
	runTestCases(t, allg, tcs)
}
