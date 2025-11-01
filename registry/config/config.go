package config

type Container struct {
	Server Server
}

type Server struct {
	Port string
}

func NewConfig() *Container {
	return &Container{
		Server: Server{
			Port: ":" + getenv("PORT", "8080"),
		},
	}
}
