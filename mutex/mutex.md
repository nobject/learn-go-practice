mutex流程加锁流程：

如果锁未被持有，直接获取锁

正常模式：
  当前的goroutine会与被唤醒的goroutine进行抢锁，如果锁未抢到，则会进入自旋状态，自旋多次后，还未竞争到锁，如果是第1次未获取到锁，则加入到等待队列的尾部，如果超过阈值1毫
秒，那么，将这个Mutex设置为饥饿模式。

饥饿模式：饥饿模式下，mutex将锁直接交给等待队列的最前面的goroutine,新来的goroutine不会尝试获取锁，即使锁没有被持有，也不会去抢，也不会spin，会加入到等待队列的尾部.
如果当前等待的goroutine是最后一个waiter，没有其他等待的goroutine 或者 此goroutine等待的时间小于1ms，退出饥饿模式。


RWMutex: 读写锁（比较好理解）
// 结构：
```go
type RWMutex struct {
    w Mutex // 互斥锁解决多个writer的竞争
    writerSem uint32  // writer信号量
    readerSem uint32  // reader信号量
    readerCount int32 // reader的数量
    readerWait int32  // writer等待完成的reader的数量
}
```

const rwmutexMaxReaders = 1 << 30
写优先，如果已经有一个 writer 在等待请求锁的话，它会阻止新来的请求锁的 reader获取到锁，所以优先保障 writer。当然，如果有一些 reader 已经请求了锁的话，新请求的 writer 也会等待已经存在的 reader 都释放
锁之后才能获取。

基于这个思想，如果请求RLock()时，如果readerCount+1为负值(代表有写锁在)，则直接阻塞
请求Lock()时，对当前的readerCount取反，代表写操作中，然后将readerWait加上原来的readCount，表示在写锁之前等待的读锁数量，等这些读锁都Runlock之后就会唤醒阻塞写锁的这个goroutine


