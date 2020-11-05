package sequence

import "testing"

func testAllignerExt() Alligner {
	return &tableAlliger{
		gapOpen:   -10,
		gapExtend: -10,
		extended:  true,
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

func TestAllgExt(t *testing.T) {
	allg := testAllignerExt()
	tcs := []testCaseAllg{
		{
			a:     "AB",
			b:     "C",
			resA:  "AB",
			resB:  "-C",
			score: -14,
		},
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
			resB:  "AA-A",
			score: 5,
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
			score: -14,
		},
		{
			a:     "AACD",
			b:     "BCD",
			resA:  "AACD",
			resB:  "B-CD",
			score: -4,
		},
		{
			a:     "ABBBCDDD",
			b:     "ABCDDD",
			resA:  "ABBBCDDD",
			resB:  "AB--CDDD",
			score: 10,
		},
		{
			a:     "ABCDDD",
			b:     "ABBBCDDD",
			resA:  "AB--CDDD",
			resB:  "ABBBCDDD",
			score: 10,
		},
	}
	runTestCases(t, allg, tcs)
}

func TestAllgBlosumExt(t *testing.T) {
	allg := NewAlligerBLOSUM62(-10, -1)
	tcs := []testCaseAllg{
		{ // verified on https://www.ebi.ac.uk/Tools/psa/emboss_needle/
			a:     "PSPETVIHSGWVIWRELFSHWPDQCKLLFGDWFAWIHWTYLVYYSAGPPCQGQSDIVVMMQKKLRTNFCQCYKYWYQ",
			b:     "PSPSDQFFTVIHSCLYWVIWRDLMSHLFMNGAAIDIHWTWDSIAIGPPLVYPIEEVFAGPSTIVVMMQKMLRTNFCQCYKPWYQ",
			resA:  "PSPE----TVIHSG--WVIWRELFSHWPDQCKLLFGDWFAW-IHWTYLVYYSAGPPC--------QGQSDIVVMMQKKLRTNFCQCYKYWYQ",
			resB:  "PSPSDQFFTVIHSCLYWVIWRDLMSH-------LFMNGAAIDIHWTW-DSIAIGPPLVYPIEEVFAGPSTIVVMMQKMLRTNFCQCYKPWYQ",
			score: 183,
		},
	}
	runTestCases(t, allg, tcs)
}

func TestAllgDNAExt(t *testing.T) {
	allg := NewAlligerDNA(-10, -1)
	tcs := []testCaseAllg{
		{ // verified on https://www.ebi.ac.uk/Tools/psa/emboss_needle/
			a:     "GCGCGTGCGCGGAAGGAGCCAAGGTGAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGGCACCATGCTGTCCCCCGAGGCGGAGCGGGTGCTGCGGTACCTGGTCGAAGTAGAGGAGTTG",
			b:     "GACTTGTGGAACCTACTTCCTGAAAATAACCTTCTGTCCTCCGAGCTCTCCGCACCCGTGGATGACCTGCTCCCGTACACAGATGTTGCCACCTGGCTGGATGAATGTCCGAATGAAGCG",
			resA:  "G-CGCGTGCGCGGAAGGAGCCAAGGT---GAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGGCACCATGCTGTCC---CCCGAGGCGGAGCGGGTGCTGCGGTAC--------------CTGG-TCGAAGTA-----GA--GGAGTTG",
			resB:  "GACTTGTG----GAA----CCTACTTCCTGAAAAT--AACCTTCTGTC---------------CTCCGAGCTCTCCGCACCCGTGGATGACC---TGCTCCCGTACACAGATGTTGCCACCTGGCTGGATGAATGTCCGAATGAAGC-G",
			score: 46,
		},
	}
	runTestCases(t, allg, tcs)
}
