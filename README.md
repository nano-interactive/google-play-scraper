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

### Get app details

Retrieves the full detail of an application.

```go
package main

import (
    "github.com/davecgh/go-spew/spew"
    "github.com/nano-interactive/google-play-scraper/pkg/app"
)

func main() {
    a := app.New("com.google.android.googlequicksearchbox", app.Options{
        Country:  "us",
        Language: "us",
    })
    err := a.LoadDetails()
    if err != nil {
        panic(err)
    }
    err = a.LoadPermissions()
    if err != nil {
        panic(err)
    }
    spew.Dump(a)
}
```

> [!WARNING]
> 
> Methods other than LoadDetails are not maintained and can return wrong response or error