package session


type Options struct {
	Domain string

	Path string

	// MaxAge=0 means no Max-Age attribute specified and the cookie will be
	// deleted after the browser session ends.
	// MaxAge<0 means delete cookie immediately.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge int

	Secure   bool

	HttpOnly bool
}

var DefaultOptions = &Options{
	Domain:"",
	Path:"/",
	MaxAge:0,
	Secure:false,
	HttpOnly:false,
}