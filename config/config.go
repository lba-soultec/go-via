package config

type Config struct {
	Debug       bool
	Port        int `default:"8443"`
	File        string
	Network     Network
	DisableDhcp bool `default:"false"`
}

type Network struct {
	Interfaces []string
}
