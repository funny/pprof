介绍
====

这个包可以帮你监控你的Go应用的总体运行情况。

安装
====

使用`go get github.com/funny/pprof`命令把本项目安装到本地.

然后在你的代码中引用`github.com/funny/pprof`。

GC综合情况
=========

GC综合状况用来监控程序的GC总体情况，保护GC暂停时间，内存申请量等。

获取GC综合情况：

```go
summary := pprof.GCSummary()
```

输出信息：

```go

// 方便阅读的格式
println(summmary.String())


// CSV格式
println(summary.CSV())
```

输出CSV之前应该会用到CSV字段列表：

```go
println(pprof.GCSummaryColumns)
println(summary.CSV())
```

有时候会需要写入到文件：

```go

// 方便阅读的格式
summary.Write(file)

// CSV格式
summary.WriteCSV(file)
```

监控执行时间
==========

`TimeRecorder`可以帮助你监控API或者函数的执行时间。

首先你需要实例化`TimeRecorder`。

```go
recorder := pprof.NewTimeRecorder()
```

然后在任意地方记录执行时间。

```go
t1 := time.Now()

your_application.do_some_thing()

recorder.Record("do_some_thing", time.Since(t1))
```

保存结果到CSV文件中。

```go
recorder.SaveCSV("time.csv")
```

保存下来的CSV文件有以下六个字段：

```
name - 条目名称，等于Record()方法的第一个参数，可以是API名称或函数名称等等。

times - 当前条目的记录次数。

avg - 当前条目的平均执行时间。

min - 当前条目的最短执行时间。

max - 当前条目的最长执行时间。

total - 当前条目的总执行时间。
```

输出的表格默认按`avg` + `times`排序.

如果表格显示某个条目有较长的执行时间并且调用次数不低，大概就意味着你需要想办法优化了。
