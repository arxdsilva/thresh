package thresh

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
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

func (s *S) TestStructFromBody(c *check.C) {
	type fake struct {
		Value int `json:"value"`
	}
	f := new(fake)
	body := ioutil.NopCloser(bytes.NewBufferString("{\"value\":10}"))
	err := structFromBody(body, f)
	c.Assert(err, check.IsNil)
	c.Assert(((fake{}) == *f), check.Equals, false)
}

func (s *S) TestStructFromEmptyBody(c *check.C) {
	type fake struct {
		Value int `json:"value"`
	}
	f := new(fake)
	body := ioutil.NopCloser(bytes.NewBufferString("{}"))
	err := structFromBody(body, f)
	c.Assert(err, check.IsNil)
	c.Assert(((fake{}) == *f), check.Equals, true)
}

func (s *S) TestStructFromBodyIsInvalid(c *check.C) {
	type fake struct {
		ID int `json:"id"`
	}
	f := new(fake)
	body := ioutil.NopCloser(bytes.NewBufferString("{\"value\":10}"))
	err := structFromBody(body, f)
	c.Assert(err, check.IsNil)
	c.Assert(((fake{}) == *f), check.Equals, true)
}
