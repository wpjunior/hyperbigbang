package fetch

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"

	"github.com/jtacoma/uritemplates"
)

var defaultFetcher = New(nil)

type M map[string]interface{}

type Fetcher struct {
	*http.Client
}

func (f *Fetcher) URL(url string, header http.Header, includes []string) (M, error) {
	resource := M{}
	resultHeaders, err := f.simpleFetch(url, header, &resource)
	if err != nil {
		return nil, err
	}
	schemaURL, err := extractSchemaURL(resultHeaders.Get("content-type"))
	if err != nil {
		return nil, err
	}
	schema := &Schema{}
	_, err = f.simpleFetch(schemaURL, header, &schema)
	if err != nil {
		return nil, err
	}

	if len(includes) > 0 {
		resource["_rels"] = f.fetchRels(header, resource, schema, includes)
	}
	return resource, nil
}

func (f *Fetcher) fetchRels(header http.Header, resource map[string]interface{}, schema *Schema, includes []string) M {
	rels := M{}

	for _, include := range includes {
		method, relHref := schema.Links.FindRel(include)
		if method != "GET" {
			log.Printf("Not allowed %s method", method)
			continue
		}
		if relHref == "" {
			log.Printf("Link: %s is not found", include)
			continue
		}

		result, err := f.fetchRel(relHref, header, resource)
		if err != nil {
			log.Printf("Failed to get %s: %s", relHref, err.Error())
			continue
		}

		rels[include] = result
	}
	return rels
}

func (f *Fetcher) fetchRel(relHref string, header http.Header, resource map[string]interface{}) (M, error) {
	result := M{}
	template, err := uritemplates.Parse(relHref)

	if err != nil {
		return nil, err
	}

	url, err := template.Expand(resource)
	if err != nil {
		return nil, err
	}

	_, err = f.simpleFetch(url, header, &result)
	return result, err
}

func (f *Fetcher) simpleFetch(url string, header http.Header, result interface{}) (http.Header, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header = header

	if err != nil {
		return nil, err
	}

	resp, err := f.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(result)
	return resp.Header, err
}

func New(client *http.Client) *Fetcher {
	if client == nil {
		client = http.DefaultClient
	}
	return &Fetcher{Client: client}
}

func URL(url string, header http.Header, includes []string) (M, error) {
	return defaultFetcher.URL(url, header, includes)
}

func extractSchemaURL(contentType string) (string, error) {
	_, params, err := mime.ParseMediaType(contentType)
	return params["profile"], err
}
