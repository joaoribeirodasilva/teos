package configuration

type ConfigCookie struct {
	Name     string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
	Secret   string
}
