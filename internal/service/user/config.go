package user

// Config store the config of [user.Service]
type Config struct {
	secretKey secret
}

// SecretAuthKey returns a pointer to the [flag.Value] to set up the [Server]
func (c *Config) SecretAuthKey() *secret { // revive:disable-line:unexported-return provides the interface to the caller
	return &c.secretKey
}

type secret string

func (sec secret) String() string {
	return ""
}

// Set sets the value of a [*secret]
func (sec *secret) Set(s string) error {
	*sec = secret(s)
	return nil
}
