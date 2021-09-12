package configuration

type Configuration struct {
	Server      Server
	MiddleWares []string
}

type Server struct {
	ListenPort uint32
	ForwardUrl string
}

func (c *Configuration) HasMiddleware(name string) bool {
	for _, element := range c.MiddleWares {
		if element == name {
			return true
		}
	}
	return false
}
