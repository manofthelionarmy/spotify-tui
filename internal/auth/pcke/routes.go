package pcke

import "net/http"

func (p pcke) routes() *http.ServeMux {
	routes := http.NewServeMux()
	routes.HandleFunc("/callback", p.completeAuth)
	return routes
}
