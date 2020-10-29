package sequence

import (
	"math"
	"strings"

	"github.com/pkg/errors"
)

type cell struct {
	i int
	j int
}

type allgDinTableMem struct {
	alg        Alligner
	a          string
	b          string
	upBuf      []float64
	downBuf    []float64
	reserveBuf []float64
	threads    int
}

func initDinTableMem(alg Alligner, a, b string, amThreads int) allgDinTableMem {
	if amThreads <= 0 {
		amThreads = 1
	}
	return allgDinTableMem{
		alg:        alg,
		a:          a,
		b:          b,
		upBuf:      make([]float64, len(a)+1),
		downBuf:    make([]float64, len(a)+1),
		reserveBuf: make([]float64, len(a)+1),
		threads:    amThreads,
	}
}

func (dt *allgDinTableMem) max(up, left, upLeft float64) (float64, allgAction) {
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

func (dt *allgDinTableMem) calcFromUp(from, to cell) {

	cur := float64(0)
	for i := from.i; i <= to.i; i++ {
		dt.upBuf[i] = cur
		cur += dt.alg.GapVal()
	}

	for j := from.j; j < to.j; j++ {
		dt.reserveBuf[from.i] = dt.upBuf[from.i] + dt.alg.GapVal()
		for i := from.i + 1; i <= to.i; i++ {
			dt.reserveBuf[i], _ = dt.max(
				dt.upBuf[i]+dt.alg.GapVal(),
				dt.reserveBuf[i-1]+dt.alg.GapVal(),
				dt.upBuf[i-1]+dt.alg.Compare(dt.a[i-1], dt.b[j]),
			)
		}
		for i := from.i; i <= to.i; i++ {
			dt.upBuf[i] = dt.reserveBuf[i]
		}

	}
}

func (dt *allgDinTableMem) calcFromDown(from, to cell) {

	cur := float64(0)
	for i := to.i; i >= from.i; i-- {
		dt.downBuf[i] = cur
		cur += dt.alg.GapVal()
	}

	for j := to.j; j > from.j; j-- {
		dt.reserveBuf[to.i] = dt.downBuf[to.i] + dt.alg.GapVal()
		for i := to.i - 1; i >= from.i; i-- {
			dt.reserveBuf[i], _ = dt.max(
				dt.downBuf[i]+dt.alg.GapVal(),
				dt.reserveBuf[i+1]+dt.alg.GapVal(),
				dt.downBuf[i+1]+dt.alg.Compare(dt.a[i], dt.b[j-1]),
			)
		}
		for i := from.i; i <= to.i; i++ {
			dt.downBuf[i] = dt.reserveBuf[i]
		}

	}
}

func (dt *allgDinTableMem) calcPart(from, to cell) []allgAction {

	if from.j == to.j {
		res := make([]allgAction, to.i-from.i)
		for i := 0; i < len(res); i++ {
			res[i] = actionLeft
		}
		return res
	}

	size := to.j - from.j
	sizeFromUp := size / 2
	sizeFromDown := (size - (size+1)%2) / 2

	dt.calcFromUp(
		from,
		cell{
			i: to.i,
			j: from.j + sizeFromUp,
		},
	)
	dt.calcFromDown(
		cell{
			i: from.i,
			j: to.j - sizeFromDown,
		},
		to,
	)

	// default
	upI := from.i
	j := from.j + sizeFromUp
	action := actionUp
	val := dt.upBuf[upI] + dt.downBuf[upI] + dt.alg.GapVal()
	// actionUp
	for i := from.i; i <= to.i; i++ {
		curVal := dt.upBuf[i] + dt.downBuf[i] + dt.alg.GapVal()

		if curVal > val {
			upI = i
			val = curVal
		}
	}
	// actionUpLeft
	for i := from.i; i < to.i; i++ {
		curVal := dt.upBuf[i] + dt.downBuf[i+1] + dt.alg.Compare(dt.a[i], dt.b[j])

		if curVal > val {
			upI = i
			val = curVal
			action = actionUpLeft
		}
	}

	nextTo := cell{
		i: upI,
		j: j,
	}
	nextFrom := cell{
		i: upI,
		j: j + 1,
	}
	if action == actionUpLeft {
		nextFrom.i++
	}

	return append(append(dt.calcPart(from, nextTo), action), dt.calcPart(nextFrom, to)...)
}

func (dt allgDinTableMem) allign(path []allgAction) (string, string, float64) {

	resA := strings.Builder{}
	resB := strings.Builder{}
	val := float64(0)
	i, j := 0, 0

	for c := 0; i < len(dt.a) || j < len(dt.b); c++ {
		switch path[c] {
		case actionUp:
			resA.WriteByte(dt.alg.Gap())
			resB.WriteByte(dt.b[j])
			j++
			val += dt.alg.GapVal()
		case actionLeft:
			resA.WriteByte(dt.a[i])
			resB.WriteByte(dt.alg.Gap())
			i++
			val += dt.alg.GapVal()
		case actionUpLeft:
			resA.WriteByte(dt.a[i])
			resB.WriteByte(dt.b[j])
			val += dt.alg.Compare(dt.a[i], dt.b[j])
			i++
			j++
		}
	}
	return resA.String(), resB.String(), val
}

func AllignMemoryOpt(alg Alligner, a, b string, amThreads int) (string, string, float64, error) {
	if !checkSeq(alg, a) || !checkSeq(alg, b) {
		return "", "", 0, errors.New("bad seq")
	}

	dt := initDinTableMem(alg, a, b, amThreads)
	path := dt.calcPart(
		cell{
			i: 0,
			j: 0,
		},
		cell{
			i: len(a),
			j: len(b),
		},
	)
	resA, resB, v := dt.allign(path)
	return resA, resB, v, nil
}
