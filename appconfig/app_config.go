package appconfig

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	configFileName    = "tikiserver"
	defaultConfigPath = []string{".", "./config"}
	defaultConfig     = map[string]interface{}{
		"TIKI_DEPLOYMENT":                              "local",
		"LISTENER_PORT":                                9090,
		"LISTENER_HOST":                                "0.0.0.0",
		"LOG_LEVEL":                                    "INFO",
		"LOGGER_OUTPUT_FILE":                           "/tmp/tikiserver.log",
		"DB_CONFIG.DB_TYPE":                            "DynamoDB",
		"DB_CONFIG.DB_PROFILE_REGION":                  "us-west-1",
		"DB_CONFIG.DB_PROFILE_ID":                      "",
		"DB_CONFIG.DB_PROFILE_SECRET":                  "",
		"DB_CONFIG.LOCAL_SUFFIX":                       "dev",
		"STS_CONFIG.STS_PROFILE_ID":                    "",
		"STS_CONFIG.STS_PROFILE_SECRET":                "",
		"STS_CONFIG.STS_REGION":                        "us-west-1",
		"AUTHENTICATION_PROVIDER.NAME":                 "Google",
		"AUTHENTICATION_PROVIDER.GOOGLE_CLIENT_ID":     "",
		"AUTHENTICATION_PROVIDER.GOOGLE_CLIENT_SECRET": "",
		"AUTHENTICATION_PROVIDER.REDIRECT_URI":         "",
		"AUTHENTICATION_PROVIDER.SCOPES":               "",
		"AUTHENTICATION_PROVIDER.GT_ISS":               "",
		"AUTHENTICATION_PROVIDER.GT_HD":                "",
	}
)

type AppConfig struct {
	ListenerPort                     string                 `mapstructure:"LISTENER_PORT"`
	ListenerHost                     string                 `mapstructure:"LISTENER_HOST"`
	LoggerOutputFile                 string                 `mapstructure:"LOGGER_OUTPUT_FILE"`
	TikiDBConfig                     DBConfig               `mapstructure:"DB_CONFIG"`
	TikiSTSConfig                    STSConfig              `mapstructure:"STS_CONFIG"`
	TikiInMemoryStoreConfig          InMemoryStoreConfig    `mapstructure:"IN_MEMORY_STORE_CONFIG"`
	TikiOpenSearchConfig             OpenSearchConfig       `mapstructure:"OPEN_SEARCH_CONFIG"`
	TikiAuthenticationProviderConfig AuthenticationProvider `mapstructure:"AUTHENTICATION_PROVIDER"`
	LogLevel                         string                 `mapstructure:"LOG_LEVEL"`
	Deployment                       string                 `mapstructure:"TIKI_DEPLOYMENT"`
}

type DBConfig struct {
	DbType          string `mapstructure:"DB_TYPE"`
	DbProfileId     string `mapstructure:"DB_PROFILE_ID"`
	DbProfileSecret string `mapstructure:"DB_PROFILE_SECRET"`
	DbRegion        string `mapstructure:"DB_REGION"`
	LocalSuffix     string `mapstructure:"LOCAL_SUFFIX"`
}

type STSConfig struct {
	STSProfileId     string `mapstructure:"STS_PROFILE_ID"`
	STSProfileSecret string `mapstructure:"STS_PROFILE_SECRET"`
	STSRegion        string `mapstructure:"STS_REGION"`
}

type InMemoryStoreConfig struct {
	StoreType     string `mapstructure:"STORE_TYPE"`
	ClusterName   string `mapstructure:"CLUSTER_NAME"`
	HazelcastAddr string `mapstructure:"HAZELCAST_ADDR"`
}

type OpenSearchConfig struct {
	Username string `mapstructure:"USERNAME"`
	Url      string `mapstructure:"URL"`
}

type AuthenticationProvider struct {
	Name         string   `mapstructure:"NAME"`
	ClientId     string   `mapstructure:"GOOGLE_CLIENT_ID"`
	ClientSecret string   `mapstructure:"GOOGLE_CLIENT_SECRET"`
	RedirectUri  string   `mapstructure:"REDIRECT_URI"`
	Scopes       []string `mapstructure:"SCOPES"`
	GtIss        []string `mapstructure:"GT_ISS"`
	GtHd         string   `mapstructure:"GT_HD"`
}

func GetAppConfig() *AppConfig {
	var Appconf *AppConfig

	viper.SetConfigName(configFileName)
	viper.SetConfigType("yml")
	for _, v := range defaultConfigPath {
		viper.AddConfigPath(v)
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Failed to read config file: %s\n", err)
		fmt.Printf("Please specify the config file in the working directory\n")
		fmt.Printf("or define environment variable to specify it\n")
		os.Exit(-1)
	}

	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}

	if viper.Get("DB_CONFIG") == nil {
		fmt.Printf("DB_CONFIG is required but has not been set\n")
		os.Exit(0)
	}

	if viper.Get("STS_CONFIG") == nil {
		fmt.Printf("STS_CONFIG is required but has not been set\n")
		os.Exit(0)
	}

	if viper.GetString("DB_CONFIG.DB_TYPE") == "" {
		fmt.Printf("DB_CONFIG.DB_TYPE is required but has not been set\n")
		os.Exit(0)
	}

	if err := viper.Unmarshal(&Appconf); err != nil {
		fmt.Printf("Failed to process config file: %s\n", err)

		fmt.Println(Appconf)
		os.Exit(-1)
	}

	return Appconf
}
