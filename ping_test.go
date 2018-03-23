package thresh

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type fake struct {
	Value int `json:"value"`
}

type fakeInvalid struct {
	ID int `json:"id"`
}

type S struct{}

var _ = check.Suite(&S{})

func (s *S) TestAddrsSplit(c *check.C) {
	os.Setenv("ADDRS", "www.globo.com,www.facebook.com,www.leagueoflegends.com")
	defer os.Unsetenv("ADDRS")
	addrs, err := Addrs()
	c.Assert(err, check.IsNil)
	c.Assert(len(addrs), check.Equals, 3)
}

func (s *S) TestAddrsNil(c *check.C) {
	os.Setenv("ADDRS", "")
	defer os.Unsetenv("ADDRS")
	addrs, err := Addrs()
	c.Assert(err, check.IsNil)
	c.Assert(len(addrs), check.Equals, 0)
}

func (s *S) TestAddrsInvalidURL(c *check.C) {
	os.Setenv("ADDRS", "invalid%%%%")
	defer os.Unsetenv("ADDRS")
	_, err := Addrs()
	c.Assert(err, check.NotNil)
}

func (s *S) TestStructFromBody(c *check.C) {
	f := new(fake)
	body := ioutil.NopCloser(bytes.NewBufferString("{\"value\":10}"))
	err := structFromBody(body, f)
	c.Assert(err, check.IsNil)
	c.Assert(((fake{}) == *f), check.Equals, false)
}

func (s *S) TestStructFromEmptyBody(c *check.C) {
	f := new(fake)
	body := ioutil.NopCloser(bytes.NewBufferString("{}"))
	err := structFromBody(body, f)
	c.Assert(err, check.IsNil)
	c.Assert(((fake{}) == *f), check.Equals, true)
}

func (s *S) TestStructFromBodyIsInvalid(c *check.C) {
	f := new(fakeInvalid)
	body := ioutil.NopCloser(bytes.NewBufferString("{\"value\":10}"))
	err := structFromBody(body, f)
	c.Assert(err, check.IsNil)
	c.Assert(((fakeInvalid{}) == *f), check.Equals, true)
}

func (s *S) TestPingAddrStatusOK(c *check.C) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()
	addr, err := url.Parse(testServer.URL)
	c.Assert(err, check.IsNil)
	_, err = pingAddr(*addr)
	c.Assert(err, check.IsNil)
}

func (s *S) TestPingAddrEmpty(c *check.C) {
	addr := new(url.URL)
	_, err := pingAddr(*addr)
	c.Assert(err, check.NotNil)
}

func (s *S) TestPingAddrStatusNotFound(c *check.C) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer testServer.Close()
	addr, err := url.Parse(testServer.URL)
	c.Assert(err, check.IsNil)
	_, err = pingAddr(*addr)
	c.Assert(err, check.NotNil)
}

func (s *S) TestStartAddrStatusOK(c *check.C) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"value\":10}"))
		w.Header().Set("Content-Type", "application/json")
	}))
	defer testServer.Close()
	addr, err := url.Parse(testServer.URL)
	c.Assert(err, check.IsNil)
	p := new(Ping)
	p.Addr = *addr
	f := new(fakeInvalid)
	err = p.StartAddr(f)
	c.Assert(err, check.IsNil)
}

func (s *S) TestStartAddrStatusNotFound(c *check.C) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer testServer.Close()
	addr, err := url.Parse(testServer.URL)
	c.Assert(err, check.IsNil)
	p := new(Ping)
	p.Addr = *addr
	f := new(fakeInvalid)
	err = p.StartAddr(f)
	c.Assert(err, check.NotNil)
}
