package helpers

import (
	"errors"
	"regexp"
	"strings"

	"github.com/khoa5773/go-server/src/constant"
)

type RoleScopes struct {
	PersonalScopes []string
	ProjectScopes  []string
}

func GetRoleScopes(userScopes []string) (*RoleScopes, error) {
	var personalScopes []string
	var projectScopes []string

	for _, v := range append(constant.Scopes["BASIC"], userScopes...) {
		match, _ := regexp.MatchString("proj:", v)
		switch {
		case match:
			projectScopes = append(projectScopes, v)
		default:
			personalScopes = append(personalScopes, v)
		}
	}
	return &RoleScopes{
		PersonalScopes: personalScopes,
		ProjectScopes:  projectScopes,
	}, nil
}

func GetMimeTypeFromFilename(fileName string) (string, error) {
	if fileName == "" || !strings.Contains(fileName, ".") {
		return ",", errors.New("invalid filename")
	}

	s := strings.Split(fileName, ".")
	if len(s) != 2 {
		return ",", errors.New("invalid filename")
	}

	return strings.ToLower(s[1]), nil
}
