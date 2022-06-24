package appconfig

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	configFileName    = "tikiserver"
	defaultConfigPath = []string{"./testconfig", ".", "./config"}
	defaultConfig     = map[string]interface{}{
		"IMO_NET_PORT":                                 8671,
		"TIKI_DEPLOYMENT":                              "local",
		"DEVELOPER_EMAIL":                              "",
		"LISTENER_PORT":                                9090,
		"LISTENER_HOST":                                "0.0.0.0",
		"LOG_LEVEL":                                    "INFO",
		"LOGGER_OUTPUT_FILE":                           "/tmp/tikiserver.log",
		"SESSION_MAX_LENGTH":                           604800, // 1 week, in seconds
		"SESSION_MAX_SIMULTANEOUS_USERS":               3,
		"SESSION_KEEP_EXPIRED_SESSIONS_FOR":            30, // In days
		"PEER_COMMUNICATION.DISCOVERY_METHOD":          "",
		"PEER_COMMUNICATION.PORT":                      8671,
		"PEER_COMMUNICATION.PROTOCOL":                  "udp",
		"PEER_COMMUNICATION.NAMESPACE":                 "",
		"PEER_COMMUNICATION.PEERS":                     []string{},
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
		"AUTHENTICATION_PROVIDER.SCOPES":               []string{"https://www.googleapis.com/auth/userinfo.email"},
		"AUTHENTICATION_PROVIDER.GT_ISS":               []string{"https://accounts.google.com", "accounts.google.com"},
		"AUTHENTICATION_PROVIDER.GT_HD":                "",
	}
)

type AppConfig struct {
	ListenerPort                     string                  `mapstructure:"LISTENER_PORT"`
	ListenerHost                     string                  `mapstructure:"LISTENER_HOST"`
	LoggerOutputFile                 string                  `mapstructure:"LOGGER_OUTPUT_FILE"`
	PeerCommunication                PeerCommunicationConfig `mapstructure:"PEER_COMMUNICATION"`
	TikiDBConfig                     DBConfig                `mapstructure:"DB_CONFIG"`
	TikiSTSConfig                    STSConfig               `mapstructure:"STS_CONFIG"`
	TikiOpenSearchConfig             OpenSearchConfig        `mapstructure:"OPEN_SEARCH_CONFIG"`
	TikiAuthenticationProviderConfig AuthenticationProvider  `mapstructure:"AUTHENTICATION_PROVIDER"`
	LogLevel                         string                  `mapstructure:"LOG_LEVEL"`
	Deployment                       string                  `mapstructure:"TIKI_DEPLOYMENT"`
	DeveloperEmail                   string                  `mapstructure:"DEVELOPER_EMAIL"`
	IMONetPort                       int                     `mapstructure:"IMO_NET_PORT"`
	SessionMaxLength                 int64                   `mapstructure:"SESSION_MAX_LENGTH"`
	SessionMaxSimultaneousUsers      int                     `mapstructure:"SESSION_MAX_SIMULTANEOUS_USERS"`
	SessionKeepExpiredSessionsFor    int                     `mapstructure:"SESSION_KEEP_EXPIRED_SESSIONS_FOR"`
}

type PeerCommunicationConfig struct {
	Port            int      `mapstructure:"PORT"`
	Protocol        string   `mapstructure:"PROTOCOL"`
	DiscoveryMethod string   `mapstructure:"DISCOVERY_METHOD"`
	Namespace       string   `mapstructure:"NAMESPACE"`
	Peers           []string `mapstructure:"PEERS"`
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

// type InMemoryStoreConfig struct {
// 	StoreType     string `mapstructure:"STORE_TYPE"`
// 	ClusterName   string `mapstructure:"CLUSTER_NAME"`
// 	HazelcastAddr string `mapstructure:"HAZELCAST_ADDR"`
// }

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
		os.Exit(1)
	}

	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}

	if viper.GetString("DB_CONFIG.DB_TYPE") == "" {
		fmt.Printf("DB_CONFIG.DB_TYPE is required but has not been set\n")
		os.Exit(2)
	}

	if err := viper.Unmarshal(&Appconf); err != nil {
		fmt.Printf("Failed to process config file: %s\n", err)

		fmt.Println(Appconf)
		os.Exit(3)
	}

	if viper.Get("DB_CONFIG") == nil {
		fmt.Printf("DB_CONFIG is required but has not been set\n")
		os.Exit(4)
	}

	if viper.Get("STS_CONFIG") == nil {
		fmt.Printf("STS_CONFIG is required but has not been set\n")
		os.Exit(5)
	}

	if viper.Get("PEER_COMMUNICATION.PROTOCOL") == "tcp" {
		fmt.Printf("PEER_COMMUNICATION.PROTOCOL is tcp, but only udp is supported\n")
		os.Exit(6)
	}

	return Appconf
}
