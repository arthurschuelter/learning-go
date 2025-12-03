package models

type TypeChart struct {
	Table [18][18]float32
}

func (t *TypeChart) Evaluate(t1 int, t2 int) float32 {
	return t.Table[t1][t2]
}
