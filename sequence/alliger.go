package sequence

import (
	"strings"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

type Alligner interface {
	IsExtended() bool
	Compare(a, b byte) float64
	GapOpen() float64
	GapExtend() float64
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
	actionUp = allgAction(iota + 1)
	actionLeft
	actionUpLeft

	dirMat = allgAction(iota - 2)
	dirIns
	dirDel

	dirMask  = 0b11
	shiftMat = 0
	shiftIns = 2
	shiftDel = 4
)

type allgDinTable struct {
	acts [][]allgAction
	vals [][]float64
	inss [][]float64
	dels [][]float64

	calcImpl func(alg Alligner, i, j int, a, b byte)
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
		vals[i][0] = vals[i-1][0] + alg.GapOpen()
		acts[i][0] = actionUp
	}
	for i := 1; i <= len(b); i++ {
		vals[0][i] = vals[0][i-1] + alg.GapOpen()
		acts[0][i] = actionLeft
	}

	dt := allgDinTable{
		acts: acts,
		vals: vals,
	}
	dt.calcImpl = dt.calc
	return dt
}

func (dt *allgDinTable) calc(alg Alligner, i, j int, a, b byte) {
	dt.vals[i][j], dt.acts[i][j] =
		maxFloat3Dir(
			dt.vals[i-1][j]+alg.GapOpen(),
			dt.vals[i][j-1]+alg.GapOpen(),
			dt.vals[i-1][j-1]+alg.Compare(a, b))
}

func reverse(s string) string {
	res := strings.Builder{}
	res.Grow(len(s))
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

func (dt *allgDinTable) initExtend(alg Alligner, a, b string) {
	inss := make([][]float64, len(a)+1)
	dels := make([][]float64, len(a)+1)
	for i := 0; i <= len(a); i++ {
		inss[i] = make([]float64, len(b)+1)
		dels[i] = make([]float64, len(b)+1)
	}
	dt.inss = inss
	dt.dels = dels

	inf := 2*alg.GapOpen() + float64(len(a)+len(b))*alg.GapExtend() - 10000
	dt.vals[0][0] = 0
	dt.inss[0][0] = inf
	dt.dels[0][0] = inf
	for i := 1; i <= len(a); i++ {
		dt.vals[i][0] = inf
		dt.inss[i][0] = inf
		dt.dels[i][0] = alg.GapOpen() + float64(i-1)*alg.GapExtend()
		dt.acts[i][0] = dirDel << shiftDel
	}
	for i := 1; i <= len(b); i++ {
		dt.vals[0][i] = inf
		dt.inss[0][i] = alg.GapOpen() + float64(i-1)*alg.GapExtend()
		dt.dels[0][i] = inf
		dt.acts[0][i] = dirIns << shiftIns
	}
	dt.calcImpl = dt.calcExtend
	return
}

func (dt *allgDinTable) calcExtend(alg Alligner, i, j int, a, b byte) {
	cmp := alg.Compare(a, b)
	open := alg.GapOpen()
	ext := alg.GapExtend()

	var actSt, actIns, actDel allgAction
	dt.vals[i][j], actSt = maxFloat3DirAlt(
		dt.vals[i-1][j-1]+cmp, dirMat,
		dt.inss[i-1][j-1]+cmp, dirIns,
		dt.dels[i-1][j-1]+cmp, dirDel,
	)
	dt.inss[i][j], actIns = maxFloat3DirAlt(
		dt.vals[i][j-1]+open, dirMat,
		dt.inss[i][j-1]+ext, dirIns,
		dt.dels[i][j-1]+open, dirDel,
	)
	dt.dels[i][j], actDel = maxFloat3DirAlt(
		dt.vals[i-1][j]+open, dirMat,
		dt.inss[i-1][j]+open, dirIns,
		dt.dels[i-1][j]+ext, dirDel,
	)
	dt.acts[i][j] = (actSt << shiftMat) | (actIns << shiftIns) | (actDel << shiftDel)
}

func (dt allgDinTable) allignExtend(alg Alligner, a, b string) (string, string, float64) {
	if len(a) == 0 && len(b) == 0 {
		return "", "", 0
	}

	resA := strings.Builder{}
	resB := strings.Builder{}

	i, j := len(a), len(b)
	m := dt.vals[i][j]
	dir := dirMat
	if dt.inss[i][j] > m {
		m = dt.inss[i][j]
		dir = dirIns
	}
	if dt.dels[i][j] > m {
		m = dt.dels[i][j]
		dir = dirDel
	}
	for i != 0 || j != 0 {
		n := dt.acts[i][j]
		switch dir {
		case dirDel:
			i--
			resA.WriteByte(a[i])
			resB.WriteByte(alg.Gap())
		case dirIns:
			j--
			resA.WriteByte(alg.Gap())
			resB.WriteByte(b[j])
		case dirMat:
			i--
			j--
			resA.WriteByte(a[i])
			resB.WriteByte(b[j])
		}
		switch dir {
		case dirMat:
			dir = (n >> shiftMat) & dirMask
		case dirDel:
			dir = (n >> shiftDel) & dirMask
		case dirIns:
			dir = (n >> shiftIns) & dirMask
		}
	}
	return reverse(resA.String()), reverse(resB.String()), m
}

func (dt *allgDinTable) calcRow(
	alg Alligner,
	beg, end int,
	lBound, rBound *int32,
	a, b string,
	jNum int,
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
			dt.calcImpl(alg, i, j, a[i-1], b[j-1])
		}
		atomic.AddInt32(rBound, 1)
	}
}

func (dt *allgDinTable) calcTable(alg Alligner, a, b string, amThreads int) {
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
	bounds[0] = int32(maxInt(len(b), len(a)))
	wg := sync.WaitGroup{}
	wg.Add(amThreads)
	for i := 0; i < amThreads; i++ {
		beg, end := rowsPoints[i], rowsPoints[i+1]
		if beg == end {
			wg.Done()
			continue
		}
		lBound, rBound := &bounds[i], &bounds[(i+1)%amThreads]
		go func(beg, end int, lBound, rBound *int32, i int) {
			dt.calcRow(alg, beg, end, lBound, rBound, a, b, i)
			wg.Done()
		}(beg, end, lBound, rBound, i)
	}
	wg.Wait()
}

func Allign(alg Alligner, a, b string, amThreads int) (resA, resB string, v float64, err error) {
	defer func() {
		if p, ok := recover().(int); ok {
			if p == SwitchErr {
				err = errors.Errorf("fatal error %d, contact developer", SwitchErr)
			}
		}
	}()

	if !checkSeq(alg, a) || !checkSeq(alg, b) {
		return "", "", 0, errors.New("bad seq")
	}
	dt := initDinTable(alg, a, b)
	allign := dt.allign
	if alg.IsExtended() {
		dt.initExtend(alg, a, b)
		allign = dt.allignExtend
	}
	dt.calcTable(alg, a, b, amThreads)
	resA, resB, v = allign(alg, a, b)
	return resA, resB, v, nil
}
