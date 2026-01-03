package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig        `mapstructure:"app"`
	Database   DatabaseConfig   `mapstructure:"database"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Encryption EncryptionConfig `mapstructure:"encryption"`
	Docker     DockerConfig     `mapstructure:"docker"`
	Traefik    TraefikConfig    `mapstructure:"traefik"`
	Logger     LoggerConfig     `mapstructure:"logger"`
}

type AppConfig struct {
	Name                string `mapstructure:"name"`
	Env                 string `mapstructure:"env"`
	Port                int    `mapstructure:"port"`
	Debug               bool   `mapstructure:"debug"`
	RegistrationEnabled bool   `mapstructure:"registration_enabled"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type JWTConfig struct {
	Secret        string        `mapstructure:"secret"`
	AccessExpiry  time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
}

type EncryptionConfig struct {
	Key string `mapstructure:"key"`
}

type DockerConfig struct {
	Host string `mapstructure:"host"`
}

type TraefikConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	DashboardPort int    `mapstructure:"dashboard_port"`
	HTTPPort      int    `mapstructure:"http_port"`
	HTTPSPort     int    `mapstructure:"https_port"`
	ACMEEmail     string `mapstructure:"acme_email"`
	Network       string `mapstructure:"network"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("PODORU")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	bindEnvVariables()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	setDefaults(&cfg)

	return &cfg, nil
}

func bindEnvVariables() {
	viper.BindEnv("app.name", "APP_NAME")
	viper.BindEnv("app.env", "APP_ENV")
	viper.BindEnv("app.port", "APP_PORT")
	viper.BindEnv("app.debug", "APP_DEBUG")
	viper.BindEnv("app.registration_enabled", "REGISTRATION_ENABLED")

	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.ssl_mode", "DB_SSL_MODE")

	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.access_expiry", "JWT_ACCESS_EXPIRY")
	viper.BindEnv("jwt.refresh_expiry", "JWT_REFRESH_EXPIRY")

	viper.BindEnv("encryption.key", "ENCRYPTION_KEY")

	viper.BindEnv("docker.host", "DOCKER_HOST")

	viper.BindEnv("traefik.enabled", "TRAEFIK_ENABLED")
	viper.BindEnv("traefik.dashboard_port", "TRAEFIK_DASHBOARD_PORT")
	viper.BindEnv("traefik.http_port", "TRAEFIK_HTTP_PORT")
	viper.BindEnv("traefik.https_port", "TRAEFIK_HTTPS_PORT")
	viper.BindEnv("traefik.acme_email", "TRAEFIK_ACME_EMAIL")
	viper.BindEnv("traefik.network", "TRAEFIK_NETWORK")
}

func setDefaults(cfg *Config) {
	if cfg.App.Port == 0 {
		cfg.App.Port = 8080
	}
	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 25
	}
	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 5
	}
	if cfg.Database.ConnMaxLifetime == 0 {
		cfg.Database.ConnMaxLifetime = 5 * time.Minute
	}
	if cfg.JWT.AccessExpiry == 0 {
		cfg.JWT.AccessExpiry = 15 * time.Minute
	}
	if cfg.JWT.RefreshExpiry == 0 {
		cfg.JWT.RefreshExpiry = 7 * 24 * time.Hour
	}
	if cfg.Docker.Host == "" {
		cfg.Docker.Host = "unix:///var/run/docker.sock"
	}
	if cfg.Logger.Level == "" {
		cfg.Logger.Level = "info"
	}
	if cfg.Logger.Format == "" {
		cfg.Logger.Format = "json"
	}
	if cfg.Traefik.Network == "" {
		cfg.Traefik.Network = "podoru_traefik"
	}
}

func (c *AppConfig) IsDevelopment() bool {
	return c.Env == "development"
}

func (c *AppConfig) IsProduction() bool {
	return c.Env == "production"
}
