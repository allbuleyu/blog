package main

import (
	"fmt"
	"github.com/allbuleyu/blog/framework"
	"html/template"
	"log"
	"net/http"
)

type MainController struct {
	framework.Controller
}

func (c *MainController) Get() {
	err := c.Ctx.Request.ParseForm()
	if err != nil {
		panic(err)
	}


	c.Tpl,err = template.ParseFiles("public/index.html")
	if err != nil {
		panic(err)
	}
	c.TplNames = c.Tpl.Name()

	c.Data["Name"] = "hyl"
	c.Data["Email"] = "hyl.gmail.com"
	c.Data["User"] = c.Ctx.Params

	mgr := framework.NewSessionMgr("qqqqq", 10)
	mgr.StartSession(c.Ctx.ResponseWriter, c.Ctx.Request)

}

func main() {
	routes := framework.RegistorController{}
	routes.Add("/", &MainController{})
	routes.Add("/users/:id([0-9]+)/:xxx(\\w+)", &MainController{})



	//http.HandleFunc("/", hh)
	err := http.ListenAndServe(":8080", &routes)

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