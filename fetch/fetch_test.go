package fetch

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jarcoal/httpmock"

	"gopkg.in/check.v1"
)

var _ = check.Suite(&FetcherSuite{})

const myAPIRoot = "http://my-api/resource/"
const schemaAPIRoot = "http://my-api/_schema"

type FetcherSuite struct {
	*Fetcher

	transport *httpmock.MockTransport
}

func (f *FetcherSuite) SetUpSuite(c *check.C) {
	t := httpmock.NewMockTransport()
	f.transport = t
	f.Fetcher = New(&http.Client{Transport: t})

	responder, _ := httpmock.NewJsonResponder(http.StatusOK, map[string]interface{}{
		"links": []Link{
			Link{Method: "GET", Rel: "father", Href: myAPIRoot + "{fatherId}"},
			Link{Rel: "mother", Href: myAPIRoot + "{motherId}"},
		},
	})

	t.RegisterResponder("GET", schemaAPIRoot, responder)
	f.registerResource(1, "Son", 2, 3)
	f.registerResource(2, "Father", 0, 0)
	f.registerResource(3, "Mother", 0, 0)
}

func (f *FetcherSuite) registerResource(resourceID int, name string, fatherId, motherId int) {
	url := myAPIRoot + strconv.Itoa(resourceID)
	f.transport.RegisterResponder("GET", url, func(*http.Request) (*http.Response, error) {
		response, err := httpmock.NewJsonResponse(http.StatusOK, map[string]interface{}{
			"id":       resourceID,
			"name":     name,
			"fatherId": fatherId,
			"motherId": motherId,
		})
		if err == nil {
			response.Header.Set("content-type", fmt.Sprintf("application/json; profile=\"%s\"", schemaAPIRoot))
		}
		return response, err
	})
}

func (f *FetcherSuite) TestFetchURL(c *check.C) {
	result, err := f.Fetcher.URL(myAPIRoot+"1", http.Header{}, []string{"father", "mother"})
	resource := result["resource"].(M)

	c.Assert(err, check.IsNil)
	c.Assert(resource["id"], check.Equals, float64(1))
	c.Assert(resource["name"], check.Equals, "Son")

	father := result["rels"].(M)["father"].(M)
	c.Assert(father["id"], check.Equals, float64(2))
	c.Assert(father["name"], check.Equals, "Father")

	mother := result["rels"].(M)["mother"].(M)
	c.Assert(mother["id"], check.Equals, float64(3))
	c.Assert(mother["name"], check.Equals, "Mother")
}
