package main

type AllParticules struct {
	Particules []Particule `split:"\n"`
}

// Day 20

type Particule struct {
	Position     Triplet `prefix:"p=<" suffix:">"`
	Velocity     Triplet `prefix:"v=<" suffix:">"`
	Acceleration Triplet `prefix:"a=<" suffix:">"`
}

type Triplet struct {
	X int `split:"," index:"0"`
	Y int `split:"," index:"1"`
	Z int `split:"," index:"2"`
}

// Day 18

type Program struct {
	Instructions []Instruction `split:"\n"`
}

type Instruction struct {
	Name string `split:" " index:"0"`
	Arg1 string `split:" " index:"1"`
	Arg2 string `split:" " index:"2"`
}

// Day 16

type Dance struct {
	Moves []Move `split:","`
}

type Move struct {
	Code string `end:"1"`
	Arg1 string `start:"1" split:"/" index:"0"`
	Arg2 string `split:"/" index:"1"`
}

func main() {

	input := `p=< 3,0,0>, v=< 2,0,0>, a=<-1,0,0>    -4 -3 -2 -1  0  1  2  3  4
	p=< 4,0,0>, v=< 0,0,0>, a=<-2,0,0>                         (0)(1)

	p=< 4,0,0>, v=< 1,0,0>, a=<-1,0,0>    -4 -3 -2 -1  0  1  2  3  4
	p=< 2,0,0>, v=<-2,0,0>, a=<-2,0,0>                      (1)   (0)

	p=< 4,0,0>, v=< 0,0,0>, a=<-1,0,0>    -4 -3 -2 -1  0  1  2  3  4
	p=<-2,0,0>, v=<-4,0,0>, a=<-2,0,0>`

	_ = input
}
