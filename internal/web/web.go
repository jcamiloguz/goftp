package web

import "net/http"

func Start() {
	// return index.html
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.ListenAndServe("0.0.0.0:8080", nil)
}
