package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	Username           string `mapstructure:"USERNAME"`
	Password           string `mapstructure:"PASSWORD"`
	Hostname           string `mapstructure:"HOSTNAME"`
	Port               int32  `mapstructure:"PORT"`
	AdminPort          int32  `mapstructure:"ADMIN_PORT"`
	ChangeUserPassword bool   `mapstructure:"CHANGE_USER_PASSWORD"`
	GenerateToken      bool   `mapstructure:"GENERATE_TOKEN"`
	MaxRetry           int16  `mapstructure:"READINESS_MAX_RETRY"`
	KubeConfigPath     string `mapstructure:"KUBE_CONFIG_PATH"`
	Namespace          string `mapstructure:"NAMESPACE"`
}

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userId struct {
	Id string `json:"id"`
}

type jwtToken struct {
	AccessToken    string `json:"accessToken"`
	ExpiryDuration int    `json:"expiryDuration"`
	RefreshToken   string `json:"refreshToken"`
	TokenType      string `json:"tokenType"`
}

type botJwtToken struct {
	JWTToken          string `json:"JWTToken"`
	JWTTokenExpiresAt int    `json:"JWTTokenExpiresAt"`
	JWTTokenExpiry    string `json:"JWTTokenExpiry"`
}

func httpClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	client := &http.Client{
		Transport: t,
		Timeout:   10 * time.Second,
	}

	return client
}

func main() {
	log.Println("Openmetadata Initilizer")
	log.Println("Loading configuration")
	username := os.Getenv("ADMIN_PORT")
	println(username)
	config, err := loadConfig()
	if err != nil {
		log.Fatalln("unable to load configuration: ", err)
	}
	log.Println("Username: ", config.Username)
	log.Println("Password: ", config.Password)
	log.Printf("Url: %s:%d", config.Hostname, config.Port)
	readyChan := make(chan bool)
	c := httpClient()
	go readinessProbe(&config, readyChan, c)
	log.Printf("Openmetadata ready: %t", <-readyChan)
	var token jwtToken
	getTokenWithUserPass(&config, &token, c)
	log.Println(token.AccessToken)
	id := getUserIdByName(&config, &token, c)
	log.Println(id)
	botToken := getBotTokenById(&config, &token, c, id)
	log.Println(botToken)
	kubeConfig := initK8sClient(&config)
	secretClient := initK8sSecretClient(kubeConfig, config.Namespace)
	createSecret(secretClient, botToken)
}
