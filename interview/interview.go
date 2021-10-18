package interview

import (
	"fmt"
	"strings"
	"sync"
)

//使⽤两个goroutine 交替打印序列，⼀个 goroutine 打印数字， 另外⼀
//个goroutine 打印字⺟， 最终效果如下：
// 12AB34CD56EF78GH910IJ1112KL1314MN1516OP1718QR1920ST2122UV2324WX2526YZ2728
func AlternatePrint() {
	printNumChan := make(chan struct{})
	printLetterChan := make(chan struct{})
	wait := sync.WaitGroup{}
	wait.Add(1)
	go printNum(printNumChan, printLetterChan)
	go printLetter(printNumChan, printLetterChan, &wait)
	wait.Wait()
}

func printNum(printNumChan, printLetterChan chan struct{}) {
	i := 0
	for {
		if i > 0 && i%2 == 0 {
			printLetterChan <- struct{}{}
			<-printNumChan
		}
		i++
		fmt.Print(i)

	}
}

func printLetter(printNumChan, printLetterChan chan struct{}, wait *sync.WaitGroup) {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	index := 0
	for {
		if index == 0 {
			<-printLetterChan
		}
		if index > 0 && index%2 == 0 {
			printNumChan <- struct{}{}
			<-printLetterChan
		}
		if strings.Count(str, "")-1 == index {
			wait.Done()
			break
		}
		fmt.Printf("%c", ([]byte(str))[index])
		index++
	}
}