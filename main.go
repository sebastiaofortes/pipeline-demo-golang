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
		_, err := w.Write([]byte("Hello World"))
		if err != nil{
			fmt.Println(err)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		fmt.Println(err)
	}
}
