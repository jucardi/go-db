package dbx

import "time"

// M is an alias for map[string]interface{}
type M map[string]interface{}

type DbConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
	Options  string `json:"options,omitempty" yaml:"options,omitempty"`
	// DialMaxRetries defines the maximum amount of retries to attempt when dialing to a
	// connection to a mongodb instance
	DialMaxRetries int `json:"dial_max_retries" yaml:"dial_max_retries"`

	// DialRetryTimeout defines the timeout in milliseconds between retries when dialing
	// for a connection to a mongodb instance.
	DialRetryTimeout time.Duration `json:"dial_retry_timeout" yaml:"dial_max_retries"`
}
