package webscoket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jcamiloguz/goftp/internal/model"
)

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	wsConn *websocket.Conn
)

type spaHandler struct {
	staticPath string
	indexPath  string
}
type wsHandler struct {
	listener chan *model.Payload
	caller   chan any
}

func Start(listener chan *model.Payload, caller chan any) {

	router := mux.NewRouter()
	ws := wsHandler{listener, caller}
	router.PathPrefix("/socket").Handler(ws)

	spa := spaHandler{staticPath: "static", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path, err := filepath.Abs(r.URL.Path)
	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.staticPath, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {

		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func (ws wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsUpgrader.CheckOrigin = func(req *http.Request) bool {
		env := os.Getenv("APP_ENV")
		if env == "development" {
			return true
		}
		if req.Header.Get("Origin") != "http://"+req.Host {
			fmt.Printf("Origin %s is not allowed\n %s \n", req.Header.Get("Origin"), req.Host)
			http.Error(w, "Origin not allowed", http.StatusForbidden)
			return false
		}
		return true
	}
	var err error
	wsConn, err = wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("could not upgrade: %s\n", err.Error())
		return
	}

	defer wsConn.Close()

	// Refactor call to payload
	ws.caller <- wsConn

	for {
		payload := <-ws.listener
		payloadJson, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("could not marshal: %s\n", err.Error())
			return
		}
		err = wsConn.WriteMessage(websocket.TextMessage, payloadJson)
		if err != nil {
			fmt.Printf("could not write: %s\n", err.Error())
			return
		}

	}
}
