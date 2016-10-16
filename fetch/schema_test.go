package fetch

import (
	"testing"

	"gopkg.in/check.v1"
)

var _ = check.Suite(&LinksSuite{})

type LinksSuite struct{}

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s *LinksSuite) TestFindRel(c *check.C) {
	links := Links{
		Link{Rel: "parent", Href: "http://parent"},
		Link{Rel: "other", Href: "http://other", Method: "GET"},
	}
	method, url := links.FindRel("other")
	c.Assert(method, check.Equals, "GET")
	c.Assert(url, check.Equals, "http://other")

	method, url = links.FindRel("parent")
	c.Assert(method, check.Equals, "GET")
	c.Assert(url, check.Equals, "http://parent")

	method, url = links.FindRel("not-found")
	c.Assert(method, check.Equals, "")
	c.Assert(url, check.Equals, "")
}
