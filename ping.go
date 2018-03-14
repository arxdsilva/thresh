package thresh

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Addrs() (addrs []url.URL, err error) {
	envAddrs := os.Getenv("ADDRS")
	addrsStr := strings.Split(envAddrs, ",")
	for _, a := range addrsStr {
		ad, errP := url.Parse(a)
		if errP != nil {
			return nil, errP
		}
		if (url.URL{}) == *ad {
			continue
		}
		addrs = append(addrs, *ad)
	}
	return
}

func pingAddr(addr url.URL) (body io.ReadCloser, err error) {
	resp, err := http.Get(addr.String())
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status not ok: %v", resp.StatusCode)
	}
	return resp.Body, err
}

func structFromBody(body io.ReadCloser, st interface{}) (err error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	return json.Unmarshal(b, &st)
}

type Ping struct {
	Addr   url.URL
	Status bool
}

// StartAddr receives a structure and a address and unmarshals a body into it's structure
func (p *Ping) StartAddr(s interface{}) (err error) {
	b, err := pingAddr(p.Addr)
	if err != nil {
		return
	}
	defer b.Close()
	err = structFromBody(b, &s)
	if err != nil {
		return
	}
	return
}

func (p *Ping) CheckStatus() {
	if !p.Status {
		// send slack msg
	}
}
