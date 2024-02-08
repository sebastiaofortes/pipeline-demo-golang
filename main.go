package main

import (
	"fmt"
	"net/http"

	"github.com/sebastiaofortes/pipeline-demo-golang/hello"
)

func main()  {
	// comment
	fmt.Println(hello.HelloGitHub())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	http.ListenAndServe(":8080", nil)
}
