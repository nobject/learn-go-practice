使用for range 特别需要注意的是指针操作，如果v的值是指针，很容易被最后一个v给覆盖

for range + goroutine的方式容易出问题