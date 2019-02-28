package context

var ContextKey = struct{}{}

type Context struct {
	vals map[string]interface{}
}

func NewContext(initial map[string]interface{}) *Context {
	if initial == nil {
		initial = make(map[string]interface{})
	}
	return &Context{
		vals: initial,
	}
}

func (c *Context) Set(key string, val interface{}) {
	c.vals[key] = val
}

func (c *Context) Map() map[string]interface{} {
	return c.vals
}
