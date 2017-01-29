package common

var Config struct {
	Auth0Secret string `short:"s" long:"auth-secret" description:"The secret from Auth0" required:"true"`
	ServerPort  int    `short:"p" long:"server-port" description:"The server port" default:"8000" required:"true"`
	Database    struct {
		Address  string `long:"db-address" description:"The database address" default:"localhost" required:"true"`
		Username string `long:"db-user" description:"The database username" required:"true"`
		Password string `long:"db-password" description:"The database password" required:"true"`
		Name     string `long:"db-name" description:"The database name" required:"true"`
	}
}
