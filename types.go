package dbx

// M is an alias for map[string]interface{}
type M map[string]interface{}

// DbConfig contains all database configuration fields.
type DbConfig struct {
	// Host is the hostname where the database is located
	Host string `json:"host" yaml:"host"`

	// Port is the port where the database is listening to connections
	Port int `json:"port" yaml:"port"`

	// Username is the username to authenticate to the database
	Username string `json:"username" yaml:"username"`

	// Password is the password to authenticate to the database
	Password string `json:"password" yaml:"password"`

	// Database indicates the database name to connect to
	Database string `json:"database" yaml:"database"`

	// Options is any additional options to be added to the connection string
	Options string `json:"options,omitempty" yaml:"options,omitempty"`

	// DialMaxRetries defines the maximum amount of retries to attempt when dialing to a db
	DialMaxRetries int `json:"dial_max_retries" yaml:"dial_max_retries"`

	// DialRetryTimeout defines the timeout in milliseconds between retries when dialing to a db
	DialRetryTimeout int64 `json:"dial_retry_timeout" yaml:"dial_retry_timeout"`
}