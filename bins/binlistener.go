package bins

import (
	"fmt"
	"net/http"
)

func BinListener(IP string, PORT string) {

	mux := http.NewServeMux()
	mux.HandleFunc("/", BinHandler)

	fmt.Println("bins listening on", PORT)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", IP, PORT),
		Handler: mux,
	}

	server.ListenAndServe()
}
