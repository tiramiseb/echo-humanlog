Echo Humanlog
=============

Go package humanlog provides a logger output that makes
[Echo](https://echo.labstack.com/) logs human-readable.

Instead of reimplementing log capabilities, this package receives output from
Labstack's gommon logger and writes them in a human-readable form, thus
allowing maximum compatibility.

This package provides two things:

* a logreader, where the Echo logger must write
* a config for the Logger middleware

How to use it
-------------

```go
e.Logger.SetOutput(humanlog.New(e.Logger.Output()))
e.Use(middleware.LoggerWithConfig(humanlog.LoggerConfig))
```
