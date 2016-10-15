package handler

import (
	"testing"

	"gopkg.in/check.v1"
)

var _ = check.Suite(&URLSuite{})

type URLSuite struct{}

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s *URLSuite) TestConvertPathToURL(c *check.C) {
	url, err := convertPathToURL("/http/localhost:5000/api/test")
	c.Assert(err, check.IsNil)
	c.Assert(url, check.Equals, "http://localhost:5000/api/test")

	url, err = convertPathToURL("/https/localhost/api/test")
	c.Assert(err, check.IsNil)
	c.Assert(url, check.Equals, "https://localhost/api/test")

	url, err = convertPathToURL("/invalid")
	c.Assert(err, check.Equals, errInvalidURL)
}
