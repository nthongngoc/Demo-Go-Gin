package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	oidc "github.com/coreos/go-oidc"
	"github.com/khoa5773/go-server/src/configs"
	"github.com/khoa5773/go-server/src/domains/users"
	"github.com/khoa5773/go-server/src/shared"
)

func handleLoginCallback(c context.Context, code string) (string, error) {
	authenticator, err := shared.NewAuthenticator("login")
	if err != nil {
		return "", err
	}

	token, err := authenticator.Config.Exchange(c, code)
	if err != nil {
		return "", err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", errors.New("token invalid")
	}

	oidcConfig := &oidc.Config{
		ClientID: configs.ConfigsService.Auth0ID,
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(c, rawIDToken)

	if err != nil {
		return "", err
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		return "", err
	}

	credentials := shared.Credentials{
		Id:                profile["sub"].(string),
		HasPersonalScopes: true,
	}

	findOneUserDto := &users.FindOneUserDto{ID: fmt.Sprintf("%v", profile["sub"])}

	_, err = users.FindOneUser(findOneUserDto, credentials)

	if err != nil {
		return "", err
	}
	lastAccess, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", profile["updated_at"]))

	if err != nil {
		return "", err
	}

	updateUserDto := users.UpdateUserDto{
		LastAccess: lastAccess,
	}

	isSuccess, err := users.UpdateUser(fmt.Sprintf("%v", profile["sub"]), &updateUserDto, credentials)

	if err != nil || !isSuccess {
		return "", err
	}

	return rawIDToken, nil
}

func handleSignupCallback(c context.Context, code string) (string, error) {
	authenticator, err := shared.NewAuthenticator("login")
	if err != nil {
		return "", err
	}

	token, err := authenticator.Config.Exchange(c, code)
	if err != nil {
		return "", err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", errors.New("token invalid")
	}

	oidcConfig := &oidc.Config{
		ClientID: configs.ConfigsService.Auth0ID,
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(c, rawIDToken)

	if err != nil {
		return "", err
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		return "", err
	}

	credentials := shared.Credentials{
		Id:                profile["sub"].(string),
		HasPersonalScopes: true,
	}

	findOneUserDto := &users.FindOneUserDto{ID: fmt.Sprintf("%v", profile["sub"])}

	_, err = users.FindOneUser(findOneUserDto, credentials)

	if err == nil {
		return rawIDToken, nil
	}

	lastAccess, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", profile["updated_at"]))

	if err != nil {
		return "", err
	}

	createUserDto := users.CreateUserDto{
		ID:         fmt.Sprintf("%v", profile["sub"]),
		LastAccess: lastAccess,
		Name:       fmt.Sprintf("%v", profile["name"]),
		Email:      fmt.Sprintf("%v", profile["email"]),
		Picture:    fmt.Sprintf("%v", profile["picture"]),
	}
	InitUserSetting(&createUserDto)

	isSuccess, err := users.CreateUser(&createUserDto)

	if err != nil || !isSuccess {
		return "", err
	}

	return rawIDToken, nil
}

func InitUserSetting(userInfo *users.CreateUserDto) {
	userInfo.Access = []map[string]interface{}{}
	userInfo.SettingSelection = map[string]interface{}{
		"RBF_BLUF": map[string]interface{}{
			"default": "auto",
			"gamma":   10,
		},
		"ELASTIC_NET": map[string]interface{}{
			"default": "auto",
			"alpha":   0.2,
		},
		"LASSO": map[string]interface{}{
			"default": "time",
			"time":    1.5,
		},
		"G_BLUP": map[string]interface{}{
			"default": "",
		},
	}
	userInfo.SettingMethod = map[string]interface{}{
		"N_FOLD": map[string]interface{}{
			"n":          0,
			"repetition": 10,
		},
	}
}
