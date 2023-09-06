package demo2

import "fmt"

// Smellable 可以闻
type Smellable interface {
	smell()
}

// Eatable 可以吃
type Eatable interface {
	eat()
}

type Fruitable interface {
	Smellable
	Eatable
}

// Apple 苹果既可能闻又能吃
type Apple struct{}

func (a Apple) smell() {
	fmt.Println("apple can smell")
}

func (a Apple) eat() {
	fmt.Println("apple can eat")
}

// Flower 花只可以闻
type Flower struct{}

func (f Flower) smell() {
	fmt.Println("flower can smell")
}

func InterInInterTest() {
	var s1 Smellable
	var s2 Eatable
	var apple = Apple{}
	var flower = Flower{}
	s1 = apple
	s1.smell()

	s1 = flower
	s1.smell()

	s2 = apple
	s2.eat()

	fmt.Println("\n组合继承")
	var s3 Fruitable
	s3 = apple
	s3.smell()
	s3.eat()
}
