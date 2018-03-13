package handlers

import (
	"net/http"

	"github.com/arxdsilva/thresh"
	"github.com/labstack/echo"
)

type SomeStruct struct {
	Value int
}

func HealthCheck(c echo.Context) (err error) {
	addrs, err := thresh.Addrs()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
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
		// create a checker function to your structure
		CheckFields(str, p)
		p.CheckStatus()
	}
	return c.String(http.StatusOK, "OK")
}

func CheckFields(s *SomeStruct, p *thresh.Ping) {
	if s.Value > 100 {
		p.Status = true
	}
}
