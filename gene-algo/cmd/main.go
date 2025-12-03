package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Gene struct {
	gene         []int
	geneSize     int
	fitness      float32
	mutationRate float32
}

type Population struct {
	population     []Gene
	populationSize int
	survivorSize   int
}

type TypeMatchup struct {
	table [18][18]float32
}

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
	var populationSize int = 300
	var generations int = 1000
	var survivorSize int = 20
	var mutationRate float32 = 0.08

	start := time.Now()

	typeChart := MakeTypeMatchup()
	p := MakePopulation(populationSize, geneSize, survivorSize, mutationRate)

	for j := range generations {
		fmt.Printf(" === Generation %d === \n", j)

		for i := range p.population {
			gene := &p.population[i]
			gene.fitness = Fitness(*gene, typeChart)
		}

		// p.Print()
		p.population = MakeNewGeneration(p, typeChart)
		// p.population = SortList(p.population)
	}

	elapsed := time.Since(start)
	p.population = SortList(p.population)
	p.Print()
	ShowFinalSolution(p.population[0], typeChart)
	fmt.Println("Execution time:", elapsed)

}

func ShowFinalSolution(g Gene, t TypeMatchup) {
	fmt.Println("Final Solution:")
	for _, value := range g.gene {
		typeName, err := GetType(value)
		if err != nil {
			panic(err)
		}
		fmt.Print(typeName + " ")
	}
	fmt.Println()
	fmt.Printf("Fitness: %.2f\n", g.fitness)
}

func Fitness(g Gene, t TypeMatchup) float32 {
	var sum float32
	atk := make([]float32, 18)
	def := make([]float32, 18)

	unique := VerifyUnique(g.gene)

	var atkSum float32
	var defSum float32

	for i, v := range g.gene {
		atk[i] = t.evaluate(v, i)

		switch atk[i] {
		case 0.5:
			atk[i] = -1
		case 0:
			atk[i] = -2
		}

		atkSum += atk[i]

		def[i] = t.evaluate(i, v)

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

func MakeNewGeneration(p Population, t TypeMatchup) []Gene {
	// p.population = SortList(p.population)
	survivors := Selection(p)
	newPop := make([]Gene, 0)

	newSurvivorNum := p.populationSize - len(survivors)

	for range newSurvivorNum {
		n1 := rand.Intn(len(survivors))
		n2 := rand.Intn(len(survivors))

		c1 := Crossover(survivors[n1], survivors[n2])
		c1.fitness = Fitness(c1, t)
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

func Selection(p Population) []Gene {
	p.population = SortList(p.population)
	return p.population[0:p.survivorSize]
}

func Crossover(p1 Gene, p2 Gene) Gene {
	c1 := Clone(p1)
	c2 := Clone(p2)

	idx := rand.Intn(6)

	for i := range p1.gene {
		if i >= idx {
			c1.gene[i] = p2.gene[i]
			c2.gene[i] = p1.gene[i]
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

func (p *Population) Print() {
	fmt.Println("Population size:", p.populationSize)
	for i, g := range p.population {
		g.Print(i)
		fmt.Printf("| Fit: %.2f\n", g.fitness)
	}
}

func (g *Gene) Print(i int) {
	fmt.Printf("%2d [", i)
	for _, v := range g.gene {
		// t, err := GetType(v)

		// if err != nil {
		// 	panic(err)
		// }

		fmt.Printf("%2d, ", v)
		// fmt.Printf("%s, \t", t)
	}
	fmt.Printf("]")
}

func (g *Gene) Mutate() {
	rate := g.mutationRate

	for i := range g.gene {
		chance := rand.Float32()
		if chance < rate {
			// fmt.Print(i, " Mutei ", g.gene[i], " -> ")
			g.gene[i] = rand.Intn(18)
			// fmt.Println(g.gene[i])
		}
	}
}

func MakeGene(size int, mutationRate float32) Gene {
	gene := make([]int, size)

	for i := range gene {
		gene[i] = rand.Intn(18)
	}

	return Gene{
		gene:         gene,
		geneSize:     size,
		fitness:      0.00,
		mutationRate: mutationRate,
	}
}

func Clone(g1 Gene) Gene {
	newGene := make([]int, g1.geneSize)
	copy(newGene, g1.gene)

	return Gene{
		gene:         newGene,
		geneSize:     g1.geneSize,
		fitness:      g1.fitness,
		mutationRate: g1.mutationRate,
	}
}

func MakePopulation(populationSize int, geneSize int, survivorSize int, mutationRate float32) Population {
	population := make([]Gene, populationSize)
	for i := range population {
		population[i] = MakeGene(geneSize, mutationRate)
	}
	return Population{
		population:     population,
		populationSize: populationSize,
		survivorSize:   survivorSize,
	}
}

func MakeTypeMatchup() TypeMatchup {
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
	return TypeMatchup{
		table: matchupChart,
	}
}

func GetType(Type int) (string, error) {
	switch Type {
	case 0:
		return "NORMAL", nil
	case 1:
		return "FIRE", nil
	case 2:
		return "WATER", nil
	case 3:
		return "ELECTRIC", nil
	case 4:
		return "GRASS", nil
	case 5:
		return "ICE", nil
	case 6:
		return "FIGHTING", nil
	case 7:
		return "POISON", nil
	case 8:
		return "GROUND", nil
	case 9:
		return "FLYING", nil
	case 10:
		return "PSYCHIC", nil
	case 11:
		return "BUG", nil
	case 12:
		return "ROCK", nil
	case 13:
		return "GHOST", nil
	case 14:
		return "DRAGON", nil
	case 15:
		return "DARK", nil
	case 16:
		return "STEEL", nil
	case 17:
		return "FAIRY", nil
	default:
		return "", fmt.Errorf("Invalid type")
	}

}

func (t *TypeMatchup) evaluate(t1 int, t2 int) float32 {
	return t.table[t1][t2]
}

func SortList(list []Gene) []Gene {
	sort.Slice(list, func(i, j int) bool {
		return list[i].fitness > list[j].fitness
	})
	return list
}
