package sequence

import (
	"math"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

type Alligner interface {
	Compare(a, b byte) float64
	GapVal() float64
	InAlphabet(b byte) bool
	Gap() byte
}

type AllignerImpl struct {
	alg Alligner
}

func checkSeq(alg Alligner, a string) bool {
	for _, v := range a {
		if !alg.InAlphabet(byte(v)) || byte(v) == alg.Gap() {
			return false
		}
	}
	return true
}

type allgAction int

const (
	actionUp = allgAction(iota)
	actionLeft
	actionUpLeft
)

type allgDinTable struct {
	acts [][]allgAction
	vals [][]float64
}

func initDinTable(alg Alligner, a, b string) allgDinTable {
	vals := make([][]float64, len(a)+1)
	acts := make([][]allgAction, len(a)+1)
	for i := 0; i <= len(a); i++ {
		vals[i] = make([]float64, len(b)+1)
		acts[i] = make([]allgAction, len(b)+1)
	}
	vals[0][0] = 0
	for i := 1; i <= len(a); i++ {
		vals[i][0] = vals[i-1][0] + alg.GapVal()
		acts[i][0] = actionUp
	}
	for i := 1; i <= len(b); i++ {
		vals[0][i] = vals[0][i-1] + alg.GapVal()
		acts[0][i] = actionLeft
	}

	return allgDinTable{
		acts: acts,
		vals: vals,
	}
}

func (dt allgDinTable) max(up, left, upLeft float64) (float64, allgAction) {
	m := math.Max(up, math.Max(left, upLeft))
	act := actionUp
	if m == left {
		act = actionLeft
	}
	if m == upLeft {
		act = actionUpLeft
	}
	return m, act
}

func (dt allgDinTable) calc(alg Alligner, i, j int, a, b byte) {
	dt.vals[i][j], dt.acts[i][j] =
		dt.max(
			dt.vals[i-1][j]+alg.GapVal(),
			dt.vals[i][j-1]+alg.GapVal(),
			dt.vals[i-1][j-1]+alg.Compare(a, b))
}

func reverse(s string) string {
	res := strings.Builder{}
	for i := len(s) - 1; i >= 0; i-- {
		res.WriteByte(s[i])
	}
	return res.String()
}

func (dt allgDinTable) allign(alg Alligner, a, b string) (string, string, float64) {
	resA := strings.Builder{}
	resB := strings.Builder{}
	i, j := len(a), len(b)
	for i != 0 || j != 0 {
		switch dt.acts[i][j] {
		case actionUp:
			i--
			resA.WriteByte(a[i])
			resB.WriteByte(alg.Gap())
		case actionLeft:
			j--
			resA.WriteByte(alg.Gap())
			resB.WriteByte(b[j])
		case actionUpLeft:
			i--
			j--
			resA.WriteByte(a[i])
			resB.WriteByte(b[j])
		}
	}
	return reverse(resA.String()), reverse(resB.String()), dt.vals[len(a)][len(b)]
}

func (dt allgDinTable) calcRow(
	alg Alligner,
	beg, end int,
	lBound, rBound *int32,
	a, b string,
) {
	if beg == end {
		return
	}
	for i := 1; i <= len(a); i++ {
		for {
			v := atomic.LoadInt32(lBound)
			if int32(i) < v+1 {
				break
			}
		}
		for j := beg; j < end; j++ {
			dt.calc(alg, i, j, a[i-1], b[j-1])
		}
		atomic.AddInt32(rBound, 1)
	}
}

func (dt allgDinTable) calcTable(alg Alligner, a, b string, amThreads int) {
	rowsPoints := make([]int, amThreads+1)
	rowsPoints[0] = 1
	for i := 1; i < amThreads+1; i++ {
		rowsPoints[i] = len(b) / amThreads
		if len(b)%amThreads > i-1 {
			rowsPoints[i]++
		}
		rowsPoints[i] += rowsPoints[i-1]
	}
	bounds := make([]int32, amThreads)
	bounds[0] = int32(len(b))
	wg := sync.WaitGroup{}
	wg.Add(amThreads)
	for i := 0; i < amThreads; i++ {
		beg, end := rowsPoints[i], rowsPoints[i+1]
		if beg == end {
			wg.Done()
			continue
		}
		lBound, rBound := &bounds[i], &bounds[(i+1)%amThreads]
		go func(beg, end int, lBound, rBound *int32) {
			dt.calcRow(alg, beg, end, lBound, rBound, a, b)
			wg.Done()
		}(beg, end, lBound, rBound)
	}
	wg.Wait()
}

func Allign(alg Alligner, a, b string, amThreads int) (string, string, float64, error) {
	if !checkSeq(alg, a) || !checkSeq(alg, b) {
		return "", "", 0, errors.New("bad seq")
	}
	dt := initDinTable(alg, a, b)
	dt.calcTable(alg, a, b, amThreads)
	resA, resB, v := dt.allign(alg, a, b)
	return resA, resB, v, nil
}
