Introduction
============

This package helps you to monitor pprof situation of your Go application.

Install
=======

Use `go get github.com/funny/pprof` command to install it into your project.

And import `github.com/funny/pprof` in your code.

GC summary
==========

GC summary used to monitor GC status like GC pause time and allocation rate etc.

Get GC summary:

```go
summary := pprof.GCSummary()
```

Display GC summary：

```go

// Humman readable format
println(summmary.String())

// CSV format
println(summary.CSV())
```

Some time you need to CSV column names：

```go
println(pprof.GCSummaryColumns)
println(summary.CSV())
```

Some time you need to save into file：

```go

// Humman readable format
summary.Write(file)

// CSV format
summary.WriteCSV(file)
```

Monitor Execution Time
======================

The `TimeRecorder` helps you to monitor execution time of APIs or functions.

First you need to a `TimeRecorder` instance.

```go
recorder := pprof.NewTimeRecorder()
```

Then record execution time at any where you want.


```go
t1 := time.Now()

your_application.do_some_thing()

recorder.Record("do_some_thing", time.Since(t1))
```

Save records into a CSV file.

```go
recorder.SaveCSV("time.csv")
```

There have 6 fields in the CSV file.

```
name - Item name, equals the Record() method's first parameter, like request name、function name、operation name etc.

times - This field shows how many times the item recorded.

avg - The average execution time of the item.

min - The minmal execution time of the item.

max - The maximum execution time of the item.

total - The total execution time of the item.
```

The output table sort by `avg` + `times` in default.

If the table shows an item have long execution time and execute many times. It means maybe you need to check the execution point or make some optimization.
