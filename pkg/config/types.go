package config

// ServiceConfig service配置
type ServiceConfig struct {
	Port int    `mapstructure:"port" yaml:"port"`
	Host string `mapstructure:"host" yaml:"host"`
}

// Postgres 数据库配置
type Postgres struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	User     string `mapstructure:"user" yaml:"user"`
	Password string `mapstructure:"password" yaml:"password"`
	Database string `mapstructure:"database" yaml:"database"`
	SSLMode  string `mapstructure:"ssl_mode" yaml:"ssl_mode"`
}

// Redis 缓存配置
type Redis struct {
	Host            string `mapstructure:"host" yaml:"host"`
	Port            int    `mapstructure:"port" yaml:"port"`
	Password        string `mapstructure:"password" yaml:"password"`
	DB              int    `mapstructure:"db" yaml:"db"`
	PoolSize        int    `mapstructure:"pool_size" yaml:"pool_size"`                   // 连接池的大小
	PoolTimeout     int    `mapstructure:"pool_timeout" yaml:"pool_timeout"`             // 连接池内获取可用连接超时 单位秒
	MinIdleConns    int    `mapstructure:"min_idle_conns" yaml:"min_idle_conns"`         // 最小空闲连接数
	MaxIdleConns    int    `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`         // 最大空闲连接数
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time" yaml:"conn_max_idle_time"` // 连接的最大空闲时间 单位秒
}

type ServerConfig struct {
	WalletService ServiceConfig `mapstructure:"wallet_service" yaml:"wallet_service"`
	Postgres      Postgres      `mapstructure:"postgres" yaml:"postgres"`
	Redis         Redis         `mapstructure:"redis" yaml:"redis"`
}
