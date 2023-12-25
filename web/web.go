package web

type Config struct {
	PrevResource int
	NextResource int
	Url          string
}

func makeFullUrl(endPointUrl, resource string) string {
	if endPointUrl[0] != '/' {
		endPointUrl = "/" + endPointUrl
	}
	if endPointUrl[len(endPointUrl)-1] != '/' {
		endPointUrl = endPointUrl + "/"
	}
	return "https://pokeapi.co" + endPointUrl + resource + "/"
}
