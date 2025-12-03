package main

import (
	"fmt"
	"gene-algo/internal/helpers"
	"gene-algo/internal/models"
	"math/rand"
	"time"
)

// type Gene struct {
// 	gene         []int
// 	geneSize     int
// 	fitness      float32
// 	mutationRate float32
// }

// type Population struct {
// 	population     []Gene
// 	populationSize int
// 	survivorSize   int
// }

// type TypeChart struct {
// 	table [18][18]float32
// }

const (
	NORMAL = iota
	FIRE
	WATER
	ELECTRIC
	GRASS
	ICE
	FIGHTING
	POISON
	GROUND
	FLYING
	PSYCHIC
	BUG
	ROCK
	GHOST
	DRAGON
	DARK
	STEEL
	FAIRY
)

func main() {
	fmt.Println("Genetic Algorithm")
	var geneSize int = 6
	var populationSize int = 500
	var generations int = 1000
	var survivorSize int = 20
	var mutationRate float32 = 0.08

	best_solution := models.Gene{}

	start := time.Now()

	typeChart := MakeTypeChart()
	p := MakePopulation(populationSize, geneSize, survivorSize, mutationRate)

	for j := range generations {
		// fmt.Printf(" === Generation %d === \n", j)

		for i := range p.Population {
			gene := &p.Population[i]
			gene.Fitness = Fitness(*gene, typeChart)
		}

		p.Population = helpers.SortList(p.Population)
		UpdateBestSolution(p.Population[0], &best_solution)

		if j+1 == generations {
			break
		}

		p.Population = MakeNewGeneration(p, typeChart)
	}

	elapsed := time.Since(start)
	p.Population = helpers.SortList(p.Population)
	p.Print()
	ShowFinalSolution(best_solution, typeChart)
	Fitness(best_solution, typeChart)
	fmt.Println("Execution time:", elapsed)
}

func UpdateBestSolution(current models.Gene, best *models.Gene) {
	if current.Fitness > best.Fitness {
		*best = Clone(current)
	}
}

func ShowFinalSolution(g models.Gene, t models.TypeChart) {
	fmt.Println("Final Solution:")
	for _, value := range g.Gene {
		typeName, err := helpers.GetType(value)
		if err != nil {
			panic(err)
		}
		fmt.Print(typeName + " ")
	}
	fmt.Println()
	fmt.Printf("Fitness: %.2f\n", g.Fitness)
}

func Fitness(g models.Gene, t models.TypeChart) float32 {
	var sum float32
	atk := make([]float32, 18)
	def := make([]float32, 18)

	unique := VerifyUnique(g.Gene)

	var atkSum float32
	var defSum float32

	for i, v := range g.Gene {
		atk[i] = t.Evaluate(v, i)

		switch atk[i] {
		case 0.5:
			atk[i] = -1
		case 0:
			atk[i] = -2
		}

		atkSum += atk[i]

		def[i] = t.Evaluate(i, v)

		switch def[i] {
		case 0.5:
			def[i] = +1
		case 0:
			def[i] = +2
		}

		defSum += def[i]
		sum += float32(v)
	}

	result := (atkSum + defSum) * unique

	return result
}

func MakeNewGeneration(p models.Population, t models.TypeChart) []models.Gene {
	// p.population = SortList(p.population)
	survivors := Selection(p)
	newPop := make([]models.Gene, 0)

	newSurvivorNum := p.PopulationSize - len(survivors)

	for range newSurvivorNum {
		n1 := rand.Intn(len(survivors))
		n2 := rand.Intn(len(survivors))

		c1 := Crossover(survivors[n1], survivors[n2])
		c1.Fitness = Fitness(c1, t)
		newPop = append(newPop, c1)
	}

	newPop = append(newPop, survivors...)

	for i := range newPop {
		// newPop[i].Print(i)
		// fmt.Println()
		newPop[i].Mutate()
		// newPop[i].Print(i)
		// fmt.Println()
	}

	return newPop
}

func Selection(p models.Population) []models.Gene {
	// p.population = SortList(p.population)
	return p.Population[0:p.SurvivorSize]
}

func Crossover(p1 models.Gene, p2 models.Gene) models.Gene {
	c1 := Clone(p1)
	c2 := Clone(p2)

	idx := rand.Intn(6)

	for i := range p1.Gene {
		if i >= idx {
			c1.Gene[i] = p2.Gene[i]
			c2.Gene[i] = p1.Gene[i]
		}
	}
	return c1
}

func VerifyUnique(gene []int) float32 {
	collision := make(map[int]bool)
	var multiplier float32 = 1.0
	for _, v := range gene {
		_, ok := collision[v]
		if !ok {
			collision[v] = true
		} else {
			multiplier *= 0.5
		}
	}

	if multiplier == 1.0 {
		return 1.5
	}

	return multiplier
}

func MakeGene(size int, mutationRate float32) models.Gene {
	gene := make([]int, size)

	for i := range gene {
		gene[i] = rand.Intn(18)
	}

	return models.Gene{
		Gene:         gene,
		GeneSize:     size,
		Fitness:      0.00,
		MutationRate: mutationRate,
	}
}

func Clone(g1 models.Gene) models.Gene {
	newGene := make([]int, g1.GeneSize)
	copy(newGene, g1.Gene)

	return models.Gene{
		Gene:         newGene,
		GeneSize:     g1.GeneSize,
		Fitness:      g1.Fitness,
		MutationRate: g1.MutationRate,
	}
}

func MakePopulation(populationSize int, geneSize int, survivorSize int, mutationRate float32) models.Population {
	population := make([]models.Gene, populationSize)
	for i := range population {
		population[i] = MakeGene(geneSize, mutationRate)
	}
	return models.Population{
		Population:     population,
		PopulationSize: populationSize,
		SurvivorSize:   survivorSize,
	}
}

func MakeTypeChart() models.TypeChart {
	matchupChart := [18][18]float32{
		NORMAL:   {1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 0.5, 0.0, 1.0, 1.0, 0.5, 1.0},
		FIRE:     {1.0, 0.5, 0.5, 1.0, 2.0, 2.0, 1.0, 1.0, 1.0, 1.0, 1.0, 2.0, 0.5, 1.0, 0.5, 1.0, 2.0, 1.0},
		WATER:    {1.0, 2.0, 0.5, 1.0, 0.5, 1.0, 1.0, 1.0, 2.0, 1.0, 1.0, 1.0, 2.0, 1.0, 0.5, 1.0, 1.0, 1.0},
		ELECTRIC: {1.0, 1.0, 2.0, 0.5, 0.5, 1.0, 1.0, 1.0, 0.0, 2.0, 1.0, 1.0, 1.0, 1.0, 0.5, 1.0, 1.0, 1.0},
		GRASS:    {1.0, 0.5, 2.0, 1.0, 0.5, 1.0, 1.0, 0.5, 2.0, 0.5, 1.0, 0.5, 2.0, 1.0, 0.5, 1.0, 0.5, 1.0},
		ICE:      {1.0, 0.5, 0.5, 1.0, 2.0, 0.5, 1.0, 1.0, 2.0, 2.0, 1.0, 1.0, 1.0, 1.0, 2.0, 1.0, 0.5, 1.0},
		FIGHTING: {2.0, 1.0, 1.0, 1.0, 1.0, 2.0, 1.0, 0.5, 1.0, 0.5, 0.5, 0.5, 2.0, 0.0, 1.0, 2.0, 2.0, 0.5},
		POISON:   {1.0, 1.0, 1.0, 1.0, 2.0, 1.0, 1.0, 0.5, 0.5, 1.0, 1.0, 1.0, 0.5, 0.5, 1.0, 1.0, 0.0, 2.0},
		GROUND:   {1.0, 2.0, 1.0, 2.0, 0.5, 1.0, 1.0, 2.0, 1.0, 0.0, 1.0, 0.5, 2.0, 1.0, 1.0, 1.0, 2.0, 1.0},
		FLYING:   {1.0, 1.0, 1.0, 0.5, 2.0, 1.0, 2.0, 1.0, 1.0, 1.0, 1.0, 2.0, 0.5, 1.0, 1.0, 1.0, 0.5, 1.0},
		PSYCHIC:  {1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 2.0, 2.0, 1.0, 1.0, 0.5, 1.0, 1.0, 1.0, 1.0, 0.0, 0.5, 1.0},
		BUG:      {1.0, 0.5, 1.0, 1.0, 2.0, 1.0, 0.5, 0.5, 1.0, 0.5, 2.0, 1.0, 1.0, 0.5, 1.0, 2.0, 0.5, 0.5},
		ROCK:     {1.0, 2.0, 1.0, 1.0, 1.0, 2.0, 0.5, 1.0, 0.5, 2.0, 1.0, 2.0, 1.0, 1.0, 1.0, 1.0, 0.5, 1.0},
		GHOST:    {0.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 2.0, 1.0, 1.0, 2.0, 1.0, 0.5, 1.0, 1.0},
		DRAGON:   {1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 2.0, 1.0, 0.5, 0.0},
		DARK:     {1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 0.5, 1.0, 1.0, 1.0, 2.0, 1.0, 1.0, 2.0, 1.0, 0.5, 1.0, 0.5},
		STEEL:    {1.0, 0.5, 0.5, 0.5, 1.0, 2.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 2.0, 1.0, 1.0, 1.0, 0.5, 2.0},
		FAIRY:    {1.0, 0.5, 1.0, 1.0, 1.0, 1.0, 2.0, 0.5, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 2.0, 2.0, 0.5, 1.0},
	}
	return models.TypeChart{
		Table: matchupChart,
	}
}
