package models

import (
	"fmt"
	"math/rand"
)

type Gene struct {
	Gene         []int
	GeneSize     int
	Fitness      float32
	MutationRate float32
}

func (g *Gene) Print() {
	fmt.Printf("[")
	for _, v := range g.Gene {
		fmt.Printf("%2d, ", v)
	}
	fmt.Printf("]")
}

func (g *Gene) Mutate() {
	rate := g.MutationRate

	for i := range g.Gene {
		chance := rand.Float32()
		if chance < rate {
			// fmt.Print(i, " Mutei ", g.gene[i], " -> ")
			g.Gene[i] = rand.Intn(18)
			// fmt.Println(g.gene[i])
		}
	}
}
