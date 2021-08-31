### pointer
> go 语言的pointer并不像其他语言那么灵活，可以直接作运算。这导致的问题很简单，使用上不方便，但安全性更高

- go 语言的pointer不能作运算
- 不同类型的指针不能相互转换
- 不同类型的指针不能使用 == 或者 != 比较
- 不同类型的指针变量不能相互赋值

### unsafe
> 因为普通指针并不能作运算，golang提供了unsafe包中可能实现指针的操作

```go
// pointer可指向任意类型
type ArbitraryType int
type Pointer *ArbitraryType

//  返回所占据的字节数，但不包含所指向的内容的大小
func Sizeof(x ArbitraryType) uintptr

// 构体成员在内存中的位置离结构体起始处的字节数
func Offsetof(x ArbitraryType) uintptr

// 内存对齐使用，能够返回分配到的内存地址能整除 m
func Alignof(x ArbitraryType) uintptr
```
### uintptr
> pointer 不能直接进行数学运算，但可以把它转换成 uintptr，对 uintptr 类型进行数学运算，再转回 pointer 类型。
uintptr 并没有指针的语义，意思就是 uintptr 所指向的对象会被 gc 无情地回收( 不要让uintptr变量出现临时变量，不然有被GC的风险)。而 unsafe.Pointer 有指针语义，可以保护它所指向的对象在“有用”的时候不会被垃圾回收