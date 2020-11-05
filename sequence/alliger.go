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
	actionUp     allgAction = 0b001
	actionLeft   allgAction = 0b010
	actionUpLeft allgAction = 0b100

	stShift  allgAction = 0
	insShift allgAction = 3
	delShift allgAction = 6

	actionUpSt     = actionUp << stShift
	actionLeftSt   = actionLeft << stShift
	actionUpLeftSt = actionUpLeft << stShift

	actionUpIns     = actionUp << insShift
	actionLeftIns   = actionLeft << insShift
	actionUpLeftIns = actionUpLeft << insShift

	actionUpDel     = actionUp << delShift
	actionLeftDel   = actionLeft << delShift
	actionUpLeftDel = actionUpLeft << delShift

	upLeftMask = (actionUpLeftSt | actionUpLeftIns | actionUpLeftDel)
	leftMask   = (actionLeftSt | actionLeftIns | actionLeftDel)
	upMask     = (actionUpSt | actionUpIns | actionUpDel)
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
		dt.acts[i][0] = actionUpDel
	}
	for i := 1; i <= len(b); i++ {
		dt.vals[0][i] = inf
		dt.inss[0][i] = alg.GapOpen() + float64(i-1)*alg.GapExtend()
		dt.dels[0][i] = inf
		dt.acts[0][i] = actionLeftIns
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
		dt.vals[i-1][j-1]+cmp, actionUpLeftSt,
		dt.inss[i-1][j-1]+cmp, actionUpLeftIns,
		dt.dels[i-1][j-1]+cmp, actionUpLeftDel,
	)
	dt.inss[i][j], actIns = maxFloat3DirAlt(
		dt.vals[i][j-1]+open, actionLeftSt,
		dt.inss[i][j-1]+ext, actionLeftIns,
		dt.dels[i][j-1]+open, actionLeftDel,
	)
	dt.dels[i][j], actDel = maxFloat3DirAlt(
		dt.vals[i-1][j]+open, actionUpSt,
		dt.inss[i-1][j]+open, actionUpIns,
		dt.dels[i-1][j]+ext, actionUpDel,
	)

	dt.acts[i][j] = actSt | actIns | actDel
}

func choosePrevAction(a allgAction, v allgAction) (hor allgAction, vShift allgAction) {
	switch v {
	case stShift:
		switch a & upLeftMask {
		case actionUpLeftSt:
			return actionUpLeft, stShift
		case actionUpLeftIns:
			return actionUpLeft, insShift
		case actionUpLeftDel:
			return actionUpLeft, delShift
		}
	case insShift:
		switch a & leftMask {
		case actionLeftSt:
			return actionLeft, stShift
		case actionLeftIns:
			return actionLeft, insShift
		case actionLeftDel:
			return actionLeft, delShift
		}
	case delShift:
		switch a & upMask {
		case actionUpSt:
			return actionUp, stShift
		case actionUpIns:
			return actionUp, insShift
		case actionUpDel:
			return actionUp, delShift
		}
	}
	panic(SwitchErr)
}

func (dt allgDinTable) allignExtend(alg Alligner, a, b string) (string, string, float64) {
	if len(a) == 0 && len(b) == 0 {
		return "", "", 0
	}

	resA := strings.Builder{}
	resB := strings.Builder{}

	i, j := len(a), len(b)
	m := dt.vals[i][j]
	vShift := stShift
	if dt.inss[i][j] > m {
		m = dt.inss[i][j]
		vShift = insShift
	}
	if dt.dels[i][j] > m {
		m = dt.dels[i][j]
		vShift = delShift
	}
	var act allgAction
	for i != 0 || j != 0 {
		switch act, vShift = choosePrevAction(dt.acts[i][j], vShift); act {
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
