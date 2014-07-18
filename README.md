Introduction
============

This is an execution time profiling/monitor tool for Go projects.

Hot to use
==========

Use `go install github.com/funny/tprof` command to install it into your project.

And import `github.com/funny/tprof` in your code.

Create a profiler for your application.

```go
prof := tprof.New()
```

And record execution time at any point, like request processing.


```go
t1 := time.Now()

your_server.process(your_request)

prof.Record(your_request.Name, time.Since(t1))
```

Save profiling result as a CSV file.

```go
prof.SaveFile("tprof.csv")
```

The output
==========

There have 6 fields in the profile result.

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
