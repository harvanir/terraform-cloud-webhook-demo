package config

// AppConfig root model
type AppConfig struct {
	Server    ServerConfig    `yaml:"server"`
	Terraform TerraformConfig `yaml:"terraform"`
}

// ServerConfig server related configuration
type ServerConfig struct {
	Port int `yaml:"port"`
}

// TerraformConfig terraform related configuration
type TerraformConfig struct {
	Token TokenConfig `yaml:"token"`
}

// TokenConfig terraform cloud token authorization
type TokenConfig struct {
	Value string `yaml:"value"`
}
