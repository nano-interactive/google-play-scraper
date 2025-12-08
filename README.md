# google-play-scraper

[![GoDoc](https://godoc.org/github.com/nano-interactive/google-play-scraper/pkg?status.svg)](https://godoc.org/github.com/nano-interactive/google-play-scraper/pkg)
[![Go Report Card](https://goreportcard.com/badge/github.com/nano-interactive/google-play-scraper)](https://goreportcard.com/report/github.com/nano-interactive/google-play-scraper)
[![Coverage Status](https://coveralls.io/repos/github/n0madic/google-play-scraper/badge.svg?branch=master)](https://coveralls.io/github/n0madic/google-play-scraper?branch=master)

Golang scraper to get data from Google Play Store

This project is inspired by the [google-play-scraper](https://github.com/facundoolano/google-play-scraper) node.js project

## Installation

```shell
go get -u github.com/nano-interactive/google-play-scraper/
```

## Usage

> [!WARNING]
>
> Methods other than LoadDetails are not maintained and can return wrong response or error

### Get app details

Retrieves the full detail of an application.

```go
package main

import (
	"fmt"
	scraper "github.com/nano-interactive/google-play-scraper"
)

func main() {
	appDetails := scraper.New("com.google.android.googlequicksearchbox", scraper.Options{
		Country:  "us",
		Language: "us",
	})
	err := appDetails.LoadDetails()
	if err != nil {
		panic(err)
	}

	fmt.Println(appDetails)
}
```