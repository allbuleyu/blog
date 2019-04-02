package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", hh)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)

	}

}

func hh(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	//for key, param := range r.Form {
	//	fmt.Println("param:", param)
	//	//fmt.Fprintf(w, "%s:%s", key, param)
	//}

	//temp, err := template.New("index").Parse("hello, my blog!!, {{.}}")
	temp, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println("err:", err)
	}

	err = temp.Execute(w, r.Form)
	if err != nil {
		fmt.Println("err1:", err)
	}
}