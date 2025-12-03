package models

import "fmt"

type Population struct {
	Population     []Gene
	PopulationSize int
	SurvivorSize   int
}

func (p *Population) Print() {
	fmt.Println("Population size:", p.PopulationSize)
	for i, g := range p.Population {
		fmt.Print(i, " ")
		g.Print()
		fmt.Printf(" Fit: %.2f\n", g.Fitness)
	}
}
