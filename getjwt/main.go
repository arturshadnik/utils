package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/go-yaml/yaml"
	"google.golang.org/api/option"
)

type VerifyTokenResponse struct {
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
}

type Config struct {
	FirebaseWebApiKey string `yaml:"apiKey"`
	FirebaseUserId    string `yaml:"userId"`
	PathToGoogleJson  string `yaml:"googleCreds"`
}

func main() {

	configFile := flag.String("f", "", "Path to YAML config file")

	flag.Parse()

	var config Config
	if *configFile != "" {
		config = *parseYaml(configFile, &config)
	} else {
		if len(os.Args) != 4 {
			fmt.Println("Usage: jwt.exe -f <PATH_TO_CONFIG_YAML")
			fmt.Println("Alternate Usage: jwt.exe <FIREBASE_WEB_API_KEY> <FIREBASE_USER_ID> <PATH_TO_GOOGLE_JSON>")
			fmt.Println("Prerequisites:\n - Firebase Auth enabled\n - IAM API enabled\n - Sufficient permission for service account")
			os.Exit(1)
		}

		config.FirebaseWebApiKey = os.Args[1]
		config.FirebaseUserId = os.Args[2]
		config.PathToGoogleJson = os.Args[3]
	}

	// add a check for valid args: api_keys starts with A, user id length, json endswith .json

	c := context.Background()

	app := initializeApp(c, config.PathToGoogleJson)

	token := getToken(app, c, config.FirebaseUserId)

	token = verifyToken(config.FirebaseWebApiKey, token)

	fmt.Print(token)
}

func initializeApp(c context.Context, pathToGoogleJson string) *firebase.App {
	opts := option.WithCredentialsFile(pathToGoogleJson)

	app, err := firebase.NewApp(c, nil, opts)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase App: %v\n", err)
	}

	return app
}

func getToken(app *firebase.App, c context.Context, uid string) string {
	client, err := app.Auth(c)

	if err != nil {
		log.Fatalf("Failed to start auth client: %v", err)
	}

	token, err := client.CustomToken(c, uid)

	if err != nil {
		log.Fatalf("Error minting custom token: %v", err)
	}

	return token
}

func verifyToken(apiKey, token string) string {
	url := "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=" + apiKey
	data := map[string]any{"token": token, "returnSecureToken": true}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Something went wrong: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("Failed to validate token: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to validate token: %v", err)
	}

	var respMap VerifyTokenResponse
	err = json.Unmarshal(respBody, &respMap)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	token = respMap.IdToken

	return token
}

func parseYaml(configFile *string, config *Config) *Config {
	data, err := os.ReadFile(*configFile)

	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	return config
}
