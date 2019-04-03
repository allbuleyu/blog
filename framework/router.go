// http 路由
// 1. 存储路由
// 2. 转发路由
package framework

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

type Registor interface {
	Add(pattern string, c ControllerInterface)
}

type Route struct {
	regexp *regexp.Regexp			// register router's regexp
	params map[int]string			// params value
	controllerType reflect.Type
}

type RegistorController struct {
	routers []*Route
}

func (rc *RegistorController) Add(pattern string, c ControllerInterface) {
	if len(pattern) == 0 {
		panic("register pattern can not null")
	}

	parts := strings.Split(pattern, "/")
	params := make(map[int]string)
	j := 0
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "[^/]+"

			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[1:index]
			}

			parts[i] = expr

			params[j] = part
			j++
		}

	}

	pattern = strings.Join(parts, "/")
	regex, regexpErr := regexp.Compile(pattern)
	if regexpErr != nil {
		panic(regexpErr)
	}

	t := reflect.Indirect(reflect.ValueOf(c)).Type()

	route := &Route{
		regexp:regex,
		params:params,
		controllerType:t,
	}

	if len(rc.routers) == 0 {
		rc.routers = make([]*Route, 0)
	}

	rc.routers = append(rc.routers, route)
}

func (rc *RegistorController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestPath := r.URL.Path

	var isFindRouter bool

	for _, route := range rc.routers {
		if !route.regexp.MatchString(requestPath) {
			continue
		}

		matchs := route.regexp.FindStringSubmatch(requestPath)
		if len(matchs[0]) != len(requestPath) {
			continue
		}

		params := make(map[string]string)
		if len(route.params) > 0 {
			values := r.URL.Query()

			// 路由参数
			for i, match := range matchs {
				values.Add(route.params[i], match)
				params[route.params[i]] = match
			}

			// URL的整体参数是路由参数与普通参数一起
			r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
		}

		vc := reflect.New(route.controllerType)

		// find method with bind router
		init := vc.MethodByName("Init")
		controllerCtx := &Context{ResponseWriter:w, Request:r, Params:params}

		in := make([]reflect.Value, 2)
		in[0] = reflect.ValueOf(controllerCtx)
		in[1] = reflect.ValueOf(route.controllerType.Name())
		init.Call(in)

		in = make([]reflect.Value, 0)
		method := vc.MethodByName("Prepare")
		method.Call(in)

		if r.Method == "GET" {
			method = vc.MethodByName("Get")
			method.Call(in)
		}

		// other case

		method = vc.MethodByName("Render")
		renderErr := method.Call(in)
		if renderErr != nil {
			fmt.Println(renderErr)

		}

		// finish

		isFindRouter = true
		break
	}

	if isFindRouter == false {
		http.NotFound(w, r)
	}
}

