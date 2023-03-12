package util

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func InitializeFirebaseApp(ctx context.Context, config *Config) (*firebase.App, error) {
	opt := option.WithCredentialsFile("./sweet-dreams-15907-firebase-adminsdk-v3f4i-de8e605c3c.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		err = fmt.Errorf("could not initialize firebase:%v", err)
		return nil, err
	}

	return app, nil
}
