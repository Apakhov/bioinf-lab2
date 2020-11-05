package sequence

import (
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type cell struct {
	i int
	j int
}

type allgDinTableMem struct {
	alg     Alligner
	a       string
	b       string
	upBuf   []float64
	downBuf []float64
	resBuf  []allgAction
	async   bool
	wg      sync.WaitGroup
}

func initDinTableMem(alg Alligner, a, b string, amThreads int) allgDinTableMem {
	if amThreads <= 0 {
		amThreads = 1
	}
	return allgDinTableMem{
		alg:     alg,
		a:       a,
		b:       b,
		upBuf:   make([]float64, len(a)+1),
		downBuf: make([]float64, len(a)+1),
		resBuf:  make([]allgAction, len(a)+len(b)),
		async:   amThreads > 1,
	}
}

func (dt *allgDinTableMem) calcFromUp(from, to cell, upBuf []float64) {
	cur := float64(0)
	for i := from.i; i <= to.i; i++ {
		upBuf[i] = cur
		cur += dt.alg.GapOpen()
	}

	var hold float64
	for j := from.j; j < to.j; j++ {
		hold, upBuf[from.i] = upBuf[from.i], upBuf[from.i]+dt.alg.GapOpen()
		for i := from.i + 1; i <= to.i; i++ {
			hold, upBuf[i] = upBuf[i], maxFloat3(
				upBuf[i]+dt.alg.GapOpen(),
				upBuf[i-1]+dt.alg.GapOpen(),
				hold+dt.alg.Compare(dt.a[i-1], dt.b[j]),
			)
		}
	}
}

func (dt *allgDinTableMem) calcFromDown(from, to cell, downBuf []float64) {
	cur := float64(0)
	for i := to.i; i >= from.i; i-- {
		downBuf[i] = cur
		cur += dt.alg.GapOpen()
	}

	var hold float64
	for j := to.j; j > from.j; j-- {
		hold, downBuf[to.i] = downBuf[to.i], downBuf[to.i]+dt.alg.GapOpen()
		for i := to.i - 1; i >= from.i; i-- {
			hold, downBuf[i] = downBuf[i], maxFloat3(
				downBuf[i]+dt.alg.GapOpen(),
				downBuf[i+1]+dt.alg.GapOpen(),
				hold+dt.alg.Compare(dt.a[i], dt.b[j-1]),
			)
		}
	}
}

func (dt *allgDinTableMem) calcFromUpDownAsync(upFrom, upTo, downFrom, downTo cell, upBuf, downBuf []float64) {
	dt.wg.Add(2)
	go func() {
		dt.calcFromDown(downFrom, downTo, downBuf)
		dt.wg.Done()
	}()
	go func() {
		dt.calcFromUp(upFrom, upTo, upBuf)
		dt.wg.Done()
	}()
	dt.wg.Wait()
}

func (dt *allgDinTableMem) calcPart(from, to cell) []allgAction {
	res := dt.resBuf[from.i+from.j : to.i+from.j]
	if from.j == to.j {
		l := to.i - from.i
		for i := 0; i < l; i++ {
			res[i] = actionLeft
		}
		return res[:l]
	}

	size := to.j - from.j
	sizeFromUp := size / 2
	sizeFromDown := (size - (size+1)%2) / 2

	if !dt.async {
		dt.calcFromUp(
			from,
			cell{
				i: to.i,
				j: from.j + sizeFromUp,
			},
			dt.upBuf,
		)
		dt.calcFromDown(
			cell{
				i: from.i,
				j: to.j - sizeFromDown,
			},
			to,
			dt.downBuf,
		)
	} else {
		dt.calcFromUpDownAsync(
			from,
			cell{
				i: to.i,
				j: from.j + sizeFromUp,
			},
			cell{
				i: from.i,
				j: to.j - sizeFromDown,
			},
			to,
			dt.upBuf,
			dt.downBuf,
		)
	}
	// default
	upI := from.i
	j := from.j + sizeFromUp
	action := actionUp
	val := dt.upBuf[upI] + dt.downBuf[upI] + dt.alg.GapOpen()
	// actionUp check
	for i := from.i; i <= to.i; i++ {
		curVal := dt.upBuf[i] + dt.downBuf[i] + dt.alg.GapOpen()

		if curVal > val {
			upI = i
			val = curVal
		}
	}
	// actionUpLeft check
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

	res = res[:0:cap(res)]
	res = append(res, dt.calcPart(from, nextTo)...)
	res = append(res, action)
	res = append(res, dt.calcPart(nextFrom, to)...)

	return res
}

func (dt allgDinTableMem) allign(path []allgAction) (string, string, float64) {

	resA := strings.Builder{}
	resB := strings.Builder{}
	resA.Grow(len(path))
	resB.Grow(len(path))
	val := float64(0)
	i, j := 0, 0
	for c := 0; i < len(dt.a) || j < len(dt.b); c++ {
		switch path[c] {
		case actionUp:
			resA.WriteByte(dt.alg.Gap())
			resB.WriteByte(dt.b[j])
			j++
			val += dt.alg.GapOpen()
		case actionLeft:
			resA.WriteByte(dt.a[i])
			resB.WriteByte(dt.alg.Gap())
			i++
			val += dt.alg.GapOpen()
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
