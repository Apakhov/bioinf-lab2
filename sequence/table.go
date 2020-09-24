package sequence

type tableAlliger struct {
	gapVal    float64
	table     [][]float64
	byteToIdx map[byte]int
}

func (t *tableAlliger) Compare(a, b byte) float64 {
	return t.table[t.byteToIdx[a]][t.byteToIdx[b]]
}

func (t *tableAlliger) GapVal() float64 {
	return t.gapVal
}
func (t *tableAlliger) InAlphabet(b byte) bool {
	_, ok := t.byteToIdx[b]
	return ok
}

func (t *tableAlliger) Gap() byte {
	return '-'
}

func NewTableAlliger(
	gapVal float64,
	table [][]float64,
	byteToIdx map[byte]int,
) Alligner {
	ta := &tableAlliger{
		gapVal:    gapVal,
		table:     table,
		byteToIdx: byteToIdx,
	}
	return ta
}

func NewAlligerBLOSUM62(gapVal float64) Alligner {
	return NewTableAlliger(
		gapVal,
		[][]float64{
			{4, -1, -2, -2, 0, -1, -1, 0, -2, -1, -1, -1, -1, -2, -1, 1, 0, -3, -2, 0},
			{-1, 5, 0, -2, -3, 1, 0, -2, 0, -3, -2, 2, -1, -3, -2, -1, -1, -3, -2, -3},
			{-2, 0, 6, 1, -3, 0, 0, 0, 1, -3, -3, 0, -2, -3, -2, 1, 0, -4, -2, -3},
			{-2, -2, 1, 6, -3, 0, 2, -1, -1, -3, -4, -1, -3, -3, -1, 0, -1, -4, -3, -3},
			{0, -3, -3, -3, 9, -3, -4, -3, -3, -1, -1, -3, -1, -2, -3, -1, -1, -2, -2, -1},
			{-1, 1, 0, 0, -3, 5, 2, -2, 0, -3, -2, 1, 0, -3, -1, 0, -1, -2, -1, -2},
			{-1, 0, 0, 2, -4, 2, 5, -2, 0, -3, -3, 1, -2, -3, -1, 0, -1, -3, -2, -2},
			{0, -2, 0, -1, -3, -2, -2, 6, -2, -4, -4, -2, -3, -3, -2, 0, -2, -2, -3, -3},
			{-2, 0, 1, -1, -3, 0, 0, -2, 8, -3, -3, -1, -2, -1, -2, -1, -2, -2, 2, -3},
			{-1, -3, -3, -3, -1, -3, -3, -4, -3, 4, 2, -3, 1, 0, -3, -2, -1, -3, -1, 3},
			{-1, -2, -3, -4, -1, -2, -3, -4, -3, 2, 4, -2, 2, 0, -3, -2, -1, -2, -1, 1},
			{-1, 2, 0, -1, -3, 1, 1, -2, -1, -3, -2, 5, -1, -3, -1, 0, -1, -3, -2, -2},
			{-1, -1, -2, -3, -1, 0, -2, -3, -2, 1, 2, -1, 5, 0, -2, -1, -1, -1, -1, 1},
			{-2, -3, -3, -3, -2, -3, -3, -3, -1, 0, 0, -3, 0, 6, -4, -2, -2, 1, 3, -1},
			{-1, -2, -2, -1, -3, -1, -1, -2, -2, -3, -3, -1, -2, -4, 7, -1, -1, -4, -3, -2},
			{1, -1, 1, 0, -1, 0, 0, 0, -1, -2, -2, 0, -1, -2, -1, 4, 1, -3, -2, -2},
			{0, -1, 0, -1, -1, -1, -1, -2, -2, -1, -1, -1, -1, -2, -1, 1, 5, -2, -2, 0},
			{-3, -3, -4, -4, -2, -2, -3, -2, -2, -3, -2, -3, -1, 1, -4, -3, -2, 11, 2, -3},
			{-2, -2, -2, -3, -2, -1, -2, -3, 2, -1, -1, -2, -1, 3, -3, -2, -2, 2, 7, -1},
			{0, -3, -3, -3, -1, -2, -2, -3, -3, 3, 1, -2, 1, -1, -2, -2, 0, -3, -1, 4},
		},
		map[byte]int{
			'A': 0,
			'R': 1,
			'N': 2,
			'D': 3,
			'C': 4,
			'Q': 5,
			'E': 6,
			'G': 7,
			'H': 8,
			'I': 9,
			'L': 10,
			'K': 11,
			'M': 12,
			'F': 13,
			'P': 14,
			'S': 15,
			'T': 16,
			'W': 17,
			'Y': 18,
			'V': 19,
		},
	)
}

func NewAlligerDNA(gapVal float64) Alligner {
	return NewTableAlliger(
		gapVal,
		[][]float64{
			{5, -4, -4, -4},
			{-4, 5, -4, -4},
			{-4, -4, 5, -4},
			{-4, -4, -4, 5},
		},
		map[byte]int{
			'A': 0,
			'T': 1,
			'G': 2,
			'C': 3,
		},
	)
}

type defAlligner struct {
	gapVal float64
}

func NewDefault(gapVal float64) Alligner {
	return &defAlligner{
		gapVal: gapVal,
	}
}

func (t *defAlligner) Compare(a, b byte) float64 {
	if a == b {
		return 1
	}
	return -1
}

func (t *defAlligner) GapVal() float64 {
	return t.gapVal
}
func (t *defAlligner) InAlphabet(b byte) bool {
	return b != t.Gap()
}

func (t *defAlligner) Gap() byte {
	return '-'
}
