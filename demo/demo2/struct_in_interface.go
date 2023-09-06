package demo2

import (
	"fmt"
)

// AnimalSkill 动物能力
type (
	Debug interface {
		Print(content string)
	}

	AnimalSkill interface {
		Debug
		Say()  // 吃饭
		Eat()  // 吃饭
		Walk() // 行走
	}

	defaultAnimal struct {
		Name string
	}
)

func (a *defaultAnimal) Print(content string) {
	fmt.Println(fmt.Sprintf("%s, %s", a.Name, content))
}

func (a *defaultAnimal) Say() {
	a.Print("say...")
}

func (a *defaultAnimal) Eat() {
	a.Print("eat...")
}

func (a *defaultAnimal) Walk() {
	a.Print("walk...")
}

func NewAnimal(name string) AnimalSkill {
	return &defaultAnimal{
		Name: name,
	}
}

type (
	// PersonSkill 人类能力
	PersonSkill interface {
		AnimalSkill //匿名字段
		Think()     // 思考
		Study()     // 学习
	}

	defaultPerson1 struct {
		AnimalSkill //匿名字段
	}

	defaultPerson2 struct {
		Base AnimalSkill //非匿名字段
	}
)

func (p *defaultPerson1) Think() {
	p.Print("think...")
}

func (p *defaultPerson2) Think() {
	p.Base.Print("think...")
}

func (p *defaultPerson1) Study() {
	p.Print("study...")
}

func (p *defaultPerson2) Study() {
	p.Base.Print("study...")
}

func NewPerson1(name string) PersonSkill {
	return &defaultPerson1{
		NewAnimal(name),
	}
}

// NewPerson2 注意这里的NewPerson2和NewPerson1的区别
func NewPerson2(name string) *defaultPerson2 {
	return &defaultPerson2{
		Base: NewAnimal(name),
	}
}

func StructInInterTest1() {
	fmt.Println("\nStructInInterTest1 >>> ")
	p := NewPerson1("张三")
	p.Say()
	p.Eat()
	p.Walk()
	p.Think()
	p.Study()
}

func StructInInterTest2() {
	fmt.Println("\nStructInInterTest2 >>> ")
	p := NewPerson2("张三")
	p.Base.Say()
	p.Base.Eat()
	p.Base.Walk()

	p.Think()
	p.Study()
}
