package cluster

type Config struct {
	name    string
	conf    string
	pool    string
	user    string
	keyring string
}

func (c *Config) SetName(name string) {
	c.name = name
}

func (c *Config) GetName() string {
	return c.name
}

func (c *Config) SetConfPath(path string) {
	c.conf = path
}

func (c *Config) GetConfPath() string {
	return c.conf
}

func (c *Config) SetPool(pool string) {
	c.pool = pool
}

func (c *Config) GetPool() string {
	return c.pool
}

func (c *Config) SetUser(user string) {
	c.user = user
}

func (c *Config) GetUser() string {
	return c.user
}

func (c *Config) SetKeyringPath(path string) {
	c.keyring = path
}

func (c *Config) GetKeyringPath() string {
	return c.keyring
}
