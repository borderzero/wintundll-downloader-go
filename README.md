# wintundll-downloader-go

[![Go Report Card](https://goreportcard.com/badge/github.com/borderzero/wintundll-downloader-go)](https://goreportcard.com/report/github.com/borderzero/wintundll-downloader-go)
[![Documentation](https://godoc.org/github.com/borderzero/wintundll-downloader-go?status.svg)](https://godoc.org/github.com/borderzero/wintundll-downloader-go)
[![GitHub issues](https://img.shields.io/github/issues/borderzero/wintundll-downloader-go.svg)](https://github.com/borderzero/wintundll-downloader-go/issues)
[![license](https://img.shields.io/github/license/borderzero/wintundll-downloader-go.svg)](https://github.com/borderzero/wintundll-downloader-go/blob/master/LICENSE)

A small package to ensure the presence of the wintun dll from a golang program.


### Usage:

```
package main

import (
	"log"
	"time"

	"github.com/borderzero/wintundll-downloader-go/wintundll"
)

func main() {
	err := wintundll.Ensure(
		wintundll.WithDownloadURL("https://www.wintun.net/builds/wintun-0.14.1.zip"),
		wintundll.WithDownloadTimeout(time.Second*10),
	)
	if err != nil {
		log.Fatal("failed to ensure the presence of the wintun.ddl file")
	}
}
```