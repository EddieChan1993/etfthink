package main

import (
	"etfthink/core"
	"fmt"
	"github.com/EddieChan1993/gcore/utils/cast"
	"log"
	"net/http"
	"os"
)

const port = "1234"

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println("etfthink path isUp(bool)")
		return
	}
	path := cast.ToString(args[1])
	isUp := cast.ToBool(args[2])
	fmt.Println(path, isUp)
	core.Run(path, isUp)

	fs := http.FileServer(http.Dir("html"))
	log.Println("running server at http://localhost:" + port + "/line.html")
	log.Fatal(http.ListenAndServe("localhost:"+port, logRequest(fs)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
