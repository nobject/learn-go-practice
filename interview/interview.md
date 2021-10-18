### new与make的区别
https://mp.weixin.qq.com/s/tZg3zmESlLmefAWdTR96Tg
```text
make 函数：
能够分配并初始化类型所需的内存空间和结构，返回引用类型的本身。
具有使用范围的局限性，仅支持 channel、map、slice 三种类型。
具有独特的优势，make 函数会对三种类型的内部数据结构（长度、容量等）赋值。
new 函数：
能够分配类型所需的内存空间，返回指针引用（指向内存的指针）。
可被替代，能够通过字面值快速初始化。
```
### 单核 CPU，开两个 Goroutine，其中一个死循环，会怎么样？
https://mp.weixin.qq.com/s/h27GXmfGYVLHRG3Mu_8axw
```text
模拟场景：
设置 runtime.GOMAXPROCS 方法模拟了单核 CPU 下只有一个 P 的场景。
运行一个 Goroutine，内部跑一个 for 死循环，达到阻塞运行的目的。
运行一个 Goroutine，主函数（main）本身就是一个 Main Goroutine。

场景程序：
func main() {
    // 模拟单核 CPU
    runtime.GOMAXPROCS(1)
    
    // 模拟 Goroutine 死循环
    go func() {
        for {
        }
    }()

    time.Sleep(time.Millisecond)
    fmt.Println("脑子进煎鱼了")
}

结果：
在 Go1.14 前，不会输出任何结果。
在 Go1.14 及之后，能够正常输出结果。

场景分析：
显然，这段程序是有一个 Goroutine 是正在执行死循环，也就是说他肯定无法被抢占。
这段程序中更没有涉及主动放弃执行权的调用（runtime.Gosched），又或是其他调用（可能会导致执行权转移）的行为。因此这个 Goroutine 是没机会溜号的，只能一直打工...
那为什么主协程（Main Goroutine）会无法运行呢，其实原因是会优先调用休眠，但由于单核 CPU，其只有唯一的 P。唯一的 P 又一直在打工不愿意下班（执行 for 死循环，被迫无限加班）。
因此主协程永远没有机会呗调度，所以这个 Go 程序自然也就一直阻塞在了执行死循环的 Goroutine 中，永远无法下班（执行完毕，退出程序）。

在 Go1.14 实现了基于信号的抢占式调度，以此来解决上述一些仍然无法被抢占解决的场景。
主要原理是Go 程序在启动时，会在 runtime.sighandler 方法注册并且绑定 SIGURG 信号：
同时在调度的 runtime.sysmon 方法会调用 retake 方法处理一下两种场景：
抢占阻塞在系统调用上的 P。
抢占运行时间过长的 G。
该方法会检测符合场景的 P，当满足上述两个场景之一时，就会发送信号给 M。M 收到信号后将会休眠正在阻塞的 Goroutine，调用绑定的信号方法，并进行重新调度。以此来解决这个问题。
注：在 Go 语言中，sysmon 会用于检测抢占。sysmon 是 Go 的 Runtime 的系统检测器，sysmon 可进行 forcegc、netpoll、retake 等一系列骚操作（via @xiaorui）。
```


### defer的使用
https://mp.weixin.qq.com/s/lELMqKho003h0gfKkZxhHQ
```text
1. 多个defer的执行顺序为“后进先出”；
2. return其实应该包含前后两个步骤：
	4.1. 第一步是给返回值赋值（若为有名返回值则直接赋值，若为匿名返回值则先声明再赋值）；
	4.2. 第二步是调用RET返回指令并传入返回值，而RET则会检查defer是否存在，若存在就先逆序插播defer语句，最后RET携带返回值退出函数；
3. 同时使用闭包与defer，要注意闭包里的值是不是defer的局部变量
```
### struct的比较
https://mp.weixin.qq.com/s/HScH6nm3xf4POXVk774jUA
```text
1. 只有相同类型的结构体才可以比较，结构体是否相同不但与属性类型个数有关，还与属性顺序相关.
2. 结构体是相同的，但是结构体属性中有不可以比较的类型，如map,slice，function,则结构体不能用==比较。
3. 可以使用reflect.DeepEqual来深度比较里面的slice与map等字段
```
### Go 结构体和结构体指针调用有什么区别
https://mp.weixin.qq.com/s/g-D_eVh-8JaIoRne09bJ3Q
```go
type Person struct {
	Name string
}
func (p Person) GetName() string{
	return p.Name
}
func (p *Person) PGetName() string{
	return p.Name
}

func GetName(p Person, name string){
    return p.Name
}

func PGetName(p *Person,name string){
    return p.Name
}
```
似乎就是传指与传指针的区别。

在使用上的考虑：方法是否需要修改接收器？如果需要，接收器必须是一个指针。

在效率上的考虑：如果接收器很大，比如：一个大的结构，使用指针接收器会好很多。

在一致性上的考虑：如果类型的某些方法必须有指针接收器，那么其余的方法也应该有指针接收器，所以无论类型如何使用，方法集都是一致的。

但对接口实现上，两种还是有大区别：

(1)指针类型变量*T可以接收T和*T方法

(2)类型T只能接收T方法，不能接收*T实现的方法
### goroutine 泄漏的原因有哪些
https://mp.weixin.qq.com/s/ql01K1nOnEZpdbp--6EDYw
```text
Goroutine 内正在进行 channel/mutex 等读写操作，但由于逻辑问题，某些情况下会被一直阻塞。
Goroutine 内的业务逻辑进入死循环，资源一直无法释放。
Goroutine 内的业务逻辑进入长时间等待，有不断新增的 Goroutine 进入等待。
```
### map遍历为什么会是随机的
https://mp.weixin.qq.com/s/MzAktbjNyZD0xRVTPRKHpw
```text
for range map 在开始处理循环逻辑的时候，就做了随机播种, 从已选定的桶中开始进行遍历，寻找桶中的下一个元素进行处理
如果桶已经遍历完，则对溢出桶 overflow buckets 进行遍历处理
```
### go实现面向对象
```text
1.封装 利用结构体与结构的方法来实现
2.继承 利用结构体嵌套来实现
3.多态 利用接口来实现
```
