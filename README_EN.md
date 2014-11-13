Introduction
============

This package helps you to monitor overall situation of your Go application.

Install
=======

Use `go get github.com/funny/overall` command to install it into your project.

And import `github.com/funny/overall` in your code.

Get GC summary
==============

Just invoke `overall.GCSummary(writer)`, you will get the GC summary of your application.

Monitor Execution Time
======================

The `TimeRecoder` helps you to monitor execution time of APIs or functions.

First you need to a `TimeRecoder` instance.

```go
recoder := overall.NewTimeRecoder()
```

Then record execution time at any where you want.


```go
t1 := time.Now()

your_application.do_some_thing()

recoder.Record("do_some_thing", time.Since(t1))
```

Save records into a CSV file.

```go
recoder.SaveFile("time.csv")
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
