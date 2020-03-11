package shared

import (
	"context"
	"fmt"

	oidc "github.com/coreos/go-oidc"
	"github.com/khoa5773/go-server/src/configs"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func NewAuthenticator(purpose string) (*Authenticator, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, configs.ConfigsService.Auth0Provider)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     configs.ConfigsService.Auth0ID,
		ClientSecret: configs.ConfigsService.Auth0Secret,
		RedirectURL:  fmt.Sprintf("http://%s:%d/auth/%s/callback", configs.ConfigsService.Host, configs.ConfigsService.Port, purpose),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}
