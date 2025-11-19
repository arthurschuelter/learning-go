package main

type MacOs struct {
	Computer
}

func MakeMacOS() IComputer {
	return &MacOs{
		Computer: Computer{
			os:    "MacOs",
			price: 3000,
		},
	}
}

type Windows struct {
	Computer
}

func MakeWindows() IComputer {
	return &Windows{
		Computer: Computer{
			os:    "Windows",
			price: 2000,
		},
	}
}

type Linux struct {
	Computer
}

func MakeLinux() IComputer {
	return &Linux{
		Computer: Computer{
			os:    "Linux",
			price: 1000,
		},
	}
}
