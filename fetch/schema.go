package fetch

const defaultLinkMethod = "GET"

type Schema struct {
	Links Links `json:"links"`
}

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type Links []Link

func (l Links) FindRel(rel string) (string, string) {
	for _, link := range l {
		if link.Rel == rel {
			if link.Method == "" {
				return defaultLinkMethod, link.Href
			}
			return link.Method, link.Href
		}
	}

	return "", ""
}
