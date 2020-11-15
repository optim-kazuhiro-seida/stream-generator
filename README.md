# Use....

[https://github.com/optim-kazuhiro-seida/Go-Streamer](https://github.com/optim-kazuhiro-seida/Go-Streamer)

# Stream Generator

Golang generate command ndndnd

## Go get

```shell script
$ go get -u github.com/optim-kazuhiro-seida/stream-generator
```

## Source File ([sample](./sample/sample.go))

```shell script
//go:generate stream-generator -type=Sample
type Sample struct {
	Str string
	Int int
}
```


## Command

```shell script
$ go generate
```

## Result

You can look [source](./sample/sample_stream.go)
