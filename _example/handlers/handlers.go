package handlers

import (
	"net/http"
	"os"

	"github.com/arxdsilva/thresh"
	"github.com/labstack/echo"
)

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
			// Add a slack notifier function.
			if tok := os.Getenv("SLACK_TOKEN"); tok != "" {
				p.NotifierFunc = thresh.SlackNotifier(tok, "random")
			}
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
	}()
	return c.String(http.StatusOK, "OK")
}

func CheckFields(s *SomeStruct, p *thresh.Ping) {
	if s.Value > 100 {
		p.Status = true
	}
}
