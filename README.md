iso-disk-builder
[![Go Reference](https://pkg.go.dev/badge/github.com/gentlemanautomaton/iso-disk-builder.svg)](https://pkg.go.dev/github.com/gentlemanautomaton/iso-disk-builder)
[![Go Report Card](https://goreportcard.com/badge/github.com/gentlemanautomaton/iso-disk-builder)](https://goreportcard.com/report/github.com/gentlemanautomaton/iso-disk-builder)
====

A basic command line utility to create ISO 9660 disk images from a source
directory. Written in Go.

Uses the [go-diskfs](https://pkg.go.dev/github.com/diskfs/go-diskfs) package.

Many thanks to `go-diskfs` for providing the `create-iso-from-folder` example,
which was the starting point for this project.

The `go-diskfs` project supports the Rock Ridge extension for ISO 9660, but
not the Juliet extension. As a result, file names in Windows will be
encumbered by the limitations of the stock ISO 9660 file name format.
