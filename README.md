# gago
Google Analytics for Go

[![Build Status](https://travis-ci.com/MarkEdmondson1234/gago.svg?branch=master)](https://travis-ci.com/MarkEdmondson1234/gago)
[![codecov](https://codecov.io/gh/MarkEdmondson1234/gago/branch/master/graph/badge.svg)](https://codecov.io/gh/MarkEdmondson1234/gago)
[![Go Report Card](https://goreportcard.com/badge/github.com/MarkEdmondson1234/gago)](https://goreportcard.com/report/github.com/MarkEdmondson1234/gago)

[gago documentation on godoc.com](https://godoc.org/github.com/MarkEdmondson1234/gago/gago)

## Mission

Create a CLI that will download GA multi-threaded, using anti-sampling, auto-paging etc. as developed with [`googleAnalyticsR::google_analytics()`](https://code.markedmondson.me/googleAnalyticsR/articles/v4.html#anti-sampling)

Intended use case is for creating executables that can run on any machine without installing another program first, such as R or Python.  This should give more options for running scheduled scripts etc. 

## Install

There is a Go library for use in your own Go projects, or a command line interface (CLI) for end users.

### CLI

Find the latest binary for your system (Windows/MacOS/linux) in the [release page](https://github.com/MarkEdmondson1234/gago/releases)

Download the binary of the appropriate type for your system and put it in your bin folder such as /usr/local/bin - rename it to `gagocli` and chmod to 755

e.g. on MacOS

```
curl -o gagocli https://github.com/MarkEdmondson1234/gago/releases/download/v0.1.0/gagocli-vlatest-darwin-amd64
# from folder where download is
mv gagocli /usr/local/bin/gagocli
sudo chmod 755 /usr/local/bin/gagocli

# should now be able to use via
gagocli
```

Read the [CLI Readme](https://github.com/MarkEdmondson1234/gago/blob/master/gagocli/README.md) for usage.

### Go

Add the gago library to your Go project via `go get github.com/MarkEdmondson1234/gago/gago`

Read the [gago library Readme](https://github.com/MarkEdmondson1234/gago/blob/master/gago/README.md) for usage.

