package config

import (
	"io/ioutil"

	"github.com/spf13/viper"
)

func New() (*Config, error) {

	viper.SetConfigFile("config/config.toml")
	config := Config{}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	privateKey, err := readPublicKey("config/private_key.pem")
	if err != nil {
		return nil, err
	}

	publicKey, err := readPublicKey("config/public_key.pem")
	if err != nil {
		return nil, err
	}

	config.JWT.PrivateKey = privateKey
	config.JWT.PublicKey = publicKey

	return &config, nil

}

func readPublicKey(filepath string) (string, error) {
	keyBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return string(keyBytes), nil
}

type Config struct {
	Api        Api            `mapstructure:"api"`
	Postgres   Postgres       `mapstructure:"postgres"`
	Apm        Apm            `mapstructure:"apm"`
	Oauth      OauthProviders `mapstructure:"oauth"`
	Redis      Redis          `mapstructure:"redis"`
	Cloudinary Cloudinary     `mapstructure:"cloudinary"`
	Nats       Nats           `mapstructure:"nats"`
	Websocket  Websocket      `mapstructure:"websocket"`
	Scylla     Scylla         `mapstructure:"scylla"`
	JWT        JWT            `mapstructure:"jwt"`
	LiveKit    LiveKit        `mapstructure:"livekit"`
}

type Postgres struct {
	Host            string `mapstructure:"host"`
	Password        string `mapstructure:"password"`
	Port            int    `mapstructure:"port"`
	Database        string `mapstructure:"database"`
	Username        string `mapstructure:"username"`
	MaxOpenConn     int    `mapstructure:"max_open_conn"`
	MaxIdleConn     int    `mapstructure:"max_idle_con"`
	ConnMaxIdleTime int    `mapstructureL:"conn_max_idle_time"`
	ConnMaxLifeTime int    `mapstructure:"conn_max_life_time"`
}

type Api struct {
	Env          string `mapstructure:"env"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
}

type Apm struct {
	Enable      bool   `mapstructure:"enable"`
	Host        string `mapstructure:"host"`
	ServiceName string `mapstructure:"service_name"`
	LogLevel    string `mapstructure:"log_level"`
}

type OauthProviders struct {
	RedirectURI string `mapstructure:"redirect_uri"`
	Github      Oauth  `mapstructure:"github"`
	Google      Oauth  `mapstructure:"google"`
	Discord     Oauth  `mapstructure:"discord"`
}

type Oauth struct {
	Name             string `mapstructure:"name"`
	AuthURL          string `mapstructure:"auth_url"`
	ClientID         string `mapstructure:"client_id"`
	ClientSecret     string `mapstructure:"client_secret"`
	TokenExchangeURL string `mapstructure:"token_exchange_url"`
	Scope            string `mapstructure:"scope"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"passwors"`
	Database int    `mapstructure:"database"`
}

type Cloudinary struct {
	CloudName string `mapstructure:"cloud_name"`
	ApiKey    string `mapstructure:"api_key"`
	Secret    string `mapstructure:"secret"`
}

type Nats struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Websocket struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Scylla struct {
	Host     string `mapstructure:"host"`
	Keyspace string `mapstructure:"keyspace"`
}

type JWT struct {
	PublicKey            string `mapstructure:"public_key"`
	PrivateKey           string `mapstructure:"private_key"`
	AccessTokenDuration  int    `mapstructure:"access_token_duration"`
	RefreshTokenDuration int    `mapstructure:"refresh_token_duration"`
}

type LiveKit struct {
	ApiKey    string `mapstructure:"api_key"`
	ApiSecret string `mapstructure:"api_secret"`
}
