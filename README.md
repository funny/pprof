Introduction
============

This is a execution time profiling tool for Go projects.

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

process(request)

prof.Record(request.Name(), time.Since(t1))
```

Save profiling result as a CSV file.

```go
prof.SaveFile("tprof.csv")
```
