package framework

import (
	"html/template"
)

type ControllerInterface interface {
	Init(ct *Context, cn string)    //初始化上下文和子类名称
	Prepare()                       //开始执行之前的一些处理
	Get()                           //method=GET的处理
	Post()                          //method=POST的处理
	Delete()                        //method=DELETE的处理
	Put()                           //method=PUT的处理
	Head()                          //method=HEAD的处理
	Patch()                         //method=PATCH的处理
	Options()                       //method=OPTIONS的处理
	Finish()                        //执行完成之后的处理
	Render() error                  //执行完method对应的方法之后渲染页面
}

type Controller struct {
	Ctx        *Context
	Tpl       *template.Template
	Data      map[interface{}]interface{}
	ChildName string
	TplNames  string
	Layout    []string
	TplExt    string
}

func (c *Controller) Init(ctx *Context, cn string) {
	c.Data = make(map[interface{}]interface{})
	c.Layout = make([]string, 0)
	c.TplNames = ""
	c.ChildName = cn
	c.Ctx = ctx
	c.TplExt = "html"
}

func (c *Controller) Prepare() {

}

func (c *Controller) Get() {
	c.Tpl.Execute(c.Ctx.ResponseWriter, c.Data)
}

func (*Controller) Post() {
	panic("implement me")
}

func (*Controller) Delete() {
	panic("implement me")
}

func (*Controller) Put() {
	panic("implement me")
}

func (*Controller) Head() {
	panic("implement me")
}

func (*Controller) Patch() {
	panic("implement me")
}

func (*Controller) Options() {
	panic("implement me")
}

func (*Controller) Finish() {
	panic("implement me")
}

func (c *Controller) Render() error {
	c.Tpl.Execute(c.Ctx.ResponseWriter, c.Data)

	return nil
}


