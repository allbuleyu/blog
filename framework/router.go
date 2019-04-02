// http 路由
// 1. 存储路由
// 2. 转发路由
package framework

import (
	"reflect"
	"regexp"
	"strings"
)

type Register interface {
	Add(pattern string, c ControllerInterface)
}

type Route struct {
	regexp *regexp.Regexp
	params map[int]string
	controllerType reflect.Type
}

type RegisterController struct {
	routers []*Route
}

func (r *RegisterController) Add(pattern string, c ControllerInterface) {
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

	if len(r.routers) == 0 {
		r.routers = make([]*Route, 0)
	}

	r.routers = append(r.routers, route)

}




