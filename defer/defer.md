考察对defer的理解，defer函数属延迟执行，延迟到调用者函数执行 return 命令前被执行。多个defer之间按LIFO先进后出顺序执行。

故考题中，在Panic触发时结束函数运行，在return前先依次打印:打印后、打印中、打印前 。最后由runtime运行时抛出打印panic异常信息。

需要注意的是，函数的return value 不是原子操作.而是在编译器中分解为两部分：返回值赋值 和 return 。而defer刚好被插入到末尾的return前执行。故可以在derfer函数中修改返回值。如下示例：

```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println(doubleScore(0))    //0
	fmt.Println(doubleScore(20.0)) //40
	fmt.Println(doubleScore(50.0)) //50
}
func doubleScore(source float32) (score float32) {
	defer func() {
		if score < 1 || score >= 100 {
			//将影响返回值
			score = source
		}
	}()
	score = source * 2
	return
	//或者
	//return source * 2
}

```