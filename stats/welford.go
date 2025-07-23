package stats

type WelfordMean struct {
	Mean        float64
	VarianceSum float64
	Count       uint
}

func (w *WelfordMean) In(a float64) {
	w.Count++
	if w.Count > 1 {
		delta := a - w.Mean
		w.Mean += delta / float64(w.Count)
		w.VarianceSum += delta * (a - w.Mean)
	} else {
		w.Mean = a
	}
}

func (w *WelfordMean) Reset() {
	w.Mean = 0
	w.VarianceSum = 0
	w.Count = 0
}
