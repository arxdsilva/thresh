[![GoDoc](https://godoc.org/github.com/arxdsilva/thresh?status.png)](https://godoc.org/github.com/arxdsilva/thresh)
[![Travis Build Status](https://api.travis-ci.org/arxdsilva/thresh.svg?branch=master)](https://travis-ci.org/arxdsilva/thresh)
[![Coverage Status](https://coveralls.io/repos/github/arxdsilva/thresh/badge.svg?branch=master)](https://coveralls.io/github/arxdsilva/thresh?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/arxdsilva/thresh)](https://goreportcard.com/report/github.com/arxdsilva/thresh)
[![LICENSE](https://img.shields.io/badge/license-MIT-orange.svg)](LICENSE)

# Thresh
<img src="thresh.jpg" alt="Drawing" height="200px"/>

Library that allows you to use a healthcheck route as a checker for other routes/APIs. This was made to keep checking the size of a channel that is returned in a route, but it's open enought to serve your own porpuse.

# [Example](https://github.com/arxdsilva/thresh/tree/master/example)

This is a simple example of a implementation in a handler, you can find the full source code [here](https://github.com/arxdsilva/thresh/tree/master/example).

```golang

type SomeStruct struct {
	Value int
}

func HealthCheck(c echo.Context) (err error) {
	go func() {
		addrs, err := thresh.Addrs()
		if err != nil {
			// log err somewhere
		}
		for _, ad := range addrs {
			p := new(thresh.Ping)
			p.Addr = ad
			// pass some structure
			str := new(SomeStruct)
			err = p.StartAddr(str)
			if err != nil {
				// log err somewhere
			}
			// create a checker function
			// to your structure
			CheckFields(str, p)
			// [WIP] this will warn 
			// your slack channel
			p.CheckStatus()
		}
	}()
	return c.String(http.StatusOK, "OK")
}

func CheckFields(s *SomeStruct, p *thresh.Ping) {
	if s.Value > 100 {
		p.Status = true
	}
}

```