package server

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

const pathSecretsFirebase = "secrets/serviceAccountKey.json"

// Must initializate services first
func InitFirebase() {
	opt := option.WithCredentialsFile(pathSecretsFirebase)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase: %v", err)
		panic("No secrets found in path: " + pathSecretsFirebase)
	}
	FirebaseApp = app
}
