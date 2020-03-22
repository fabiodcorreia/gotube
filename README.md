# GoTube

The package allows to get information about youtube vidoes like the video title, duration, formats available and also download the video to any destination that implements the io.Writer interface.

## Examples

An example on how to use this package can be found on examples/main.go

## Benchmarks

The folder benchmarks contains the output benchmarks of each package internal or public

This is a benchmarks history of the example app during development

```bash
go run examples/main.go  1,23s user 1,30s system 12% cpu 20,883 total
go run examples/main.go  0,85s user 0,85s system 12% cpu 13,614 total
go run examples/main.go  0,92s user 0,90s system 13% cpu 13,715 total
go run examples/main.go  1,61s user 1,77s system 13% cpu 25,408 total
go run examples/main.go  0,89s user 0,87s system 13% cpu 13,080 total
```

## Dependencies

No external dependencies or modules