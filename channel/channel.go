package channel

import (
	"fmt"
	"time"
)

// 带缓冲的channel，如果已放置入channel中，则通道关闭，还可以继续读取，但不能写入channel 在close
// channel之后如继续写入，编译器会抛panic
// panic: send on closed channel [recovered]
//        panic: send on closed channel
// 关闭后继续读数据，得到的是零值(对于int，就是0)。
func RunBufferChannel() {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	close(ch)
	// 遍历读取（只会读取出1，2）
	for v := range ch {
		fmt.Println(v)
	}
	// 继续读取，为0值
	fmt.Println(<-ch)
	// 关闭后继续写入，抛出panic
	ch <- 3
	for v := range ch {
		fmt.Println(v)
	}
}

// 对于不带缓冲的ch，和带缓冲的一样，channel close掉之后并不影响读，只影响写入
func RunChannel() {
	ch := make(chan int)
	go func() {
		ch <- 1
		close(ch)
	}()
	for i := 0; i < 2; i++ {
		v, ok := <-ch
		fmt.Println(v)
		fmt.Println(ok)
	}
}

func SolveZeroChannel() {
	ch := make(chan int)
	go func() {
		ch <- 1
		close(ch)
	}()
	for i := 0; i < 2; i++ {
		// solution 1: 使用断言判断，ok为false代表channel已关闭后读取的
		v, ok := <-ch
		fmt.Println(v)
		fmt.Println(ok)
	}
	// solution 2: 关闭时设置channel为nil值，判断时通过channel是否为nil判断
	ch1 := make(chan int)
	go func() {
		ch1 <- 1
		close(ch1)
		ch1 = nil
	}()
	// 置为nil后的channel只读取了写入的channel的值
	for i := 0; i < 2; i++ {
		v, ok := <-ch1
		fmt.Println(v)
		fmt.Println(ok)
	}
}

func NilChannel() {
	ch := make(chan int, 5)
	ch = nil
	go func() {
		fmt.Println("before")
		ch <- 5
		fmt.Println("after")
	}()
	time.Sleep(1000 * time.Second)
}