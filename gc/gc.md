### 相关文章
https://draveness.me/golang/docs/part3-runtime/ch07-memory/golang-garbage-collector/
https://www.cnblogs.com/luozhiyun/p/14564903.html
https://www.kancloud.cn/aceld/golang/1958308#Go_V13mark_and_sweep_21
https://www.bilibili.com/video/BV1wz4y1y7Kd?p=14&spm_id_from=pageDriver
https://cloud.tencent.com/developer/article/1756163
https://blog.haohtml.com/archives/26358

https://mp.weixin.qq.com/s/5xjH-LJ53XiNm2sMNQZiGQ

GCMark	标记准备阶段，为并发标记做准备工作，启动写屏障	STW
GCMark	扫描标记阶段，与赋值器并发执行，写屏障开启	并发
GCMarkTermination	标记终止阶段，保证一个周期内标记任务完成，停止写屏障	STW
GCoff	内存清扫阶段，将需要回收的内存归还到堆中，写屏障关闭	并发
GCoff	内存归还阶段，将过多的内存归还给操作系统，写屏障关闭	并发

Go 语言中对 GC 的触发时机存在两种形式：

主动触发，通过调用 runtime.GC 来触发 GC，此调用阻塞式地等待当前 GC 运行完毕。

被动触发，分为两种方式：

使用系统监控，当超过两分钟没有产生任何 GC 时，强制触发 GC。

使用步调（Pacing）算法，其核心思想是控制内存增长的比例。