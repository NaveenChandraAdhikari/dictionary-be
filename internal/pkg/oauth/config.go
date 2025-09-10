package oauth

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig = &oauth2.Config{
	RedirectURL:  "http://127.0.0.1:3000/auth/google/callback",
	ClientID:     "",
	ClientSecret: "",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

var GitHubOAuthConfig = &oauth2.Config{
	RedirectURL:  "http://127.0.0.1:3000/auth/github/callback",
	ClientID:     "",
	ClientSecret: "",
	Scopes:       []string{"user:email"},
	Endpoint:     github.Endpoint,
}

func InitOAuthConfigs() {
	GoogleOAuthConfig.ClientID = os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	GoogleOAuthConfig.ClientSecret = os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")

	GitHubOAuthConfig.ClientID = os.Getenv("GITHUB_OAUTH_CLIENT_ID")
	GitHubOAuthConfig.ClientSecret = os.Getenv("GITHUB_OAUTH_CLIENT_SECRET")
}

const (
	GoogleUserInfoAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	GitHubUserInfoAPI = "https://api.github.com/user"
)
