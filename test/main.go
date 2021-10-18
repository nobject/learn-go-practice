package main

import "fmt"

// 定义接口
type Car interface {
	GetName() string
	Run()
}

// 定义结构体
type Tesla struct {
	Name string
}

// 实现接口的GetName()方法
func (t Tesla) GetName() string {
	t.Name = "test"
	return t.Name
}

// 实现接口的Run()方法
func (t *Tesla) Run() {
	fmt.Printf("%s is running\n", t.Name)
}

func main() {
	var c Car
	t := Tesla{"Tesla Model S"}
	c = &t  // 上面是用指针*Tesla实现了接口的方法，这里要传地址
	fmt.Println(c.GetName())
	c.Run()
}