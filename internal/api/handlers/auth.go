package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"restapi/internal/models"
	"restapi/internal/pkg/oauth"
	"restapi/internal/pkg/utils"
	"restapi/internal/repository/sqlconnect"
)

func OAuthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := oauth.GenerateOAuthState(w)
	authURL := oauth.GoogleOAuthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func OAuthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	stateParam := r.FormValue("state")
	if !oauth.ValidateOAuthState(r, stateParam) {
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}

	token, err := oauth.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.ErrorHandler(err, "OAuth token exchange failed")
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	userInfo, err := getGoogleUserInfo(token.AccessToken)
	if err != nil {
		utils.ErrorHandler(err, "Failed to get user info from Google")
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	username := userInfo.Email
	user, err := sqlconnect.FindOrCreateOAuthUser(
		"google",
		userInfo.ID,
		userInfo.Email,
		userInfo.GivenName,
		userInfo.FamilyName,
		username,
	)
	if err != nil {
		utils.ErrorHandler(err, "Failed to create/find OAuth user")
		http.Error(w, "Failed to process user", http.StatusInternalServerError)
		return
	}

	jwtToken, err := utils.SignToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.ErrorHandler(err, "Failed to create JWT token")
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	})

	oauth.ClearOAuthState(w)
	http.Redirect(w, r, "http://dictionary-bucker.s3-website-us-east-1.amazonaws.com/index.html", http.StatusTemporaryRedirect)
}

func getGoogleUserInfo(accessToken string) (*models.GoogleUserInfo, error) {
	resp, err := http.Get(oauth.GoogleUserInfoAPI + accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var userInfo models.GoogleUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}

	return &userInfo, nil
}

func OAuthGitHubLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := oauth.GenerateOAuthState(w)
	authURL := oauth.GitHubOAuthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func OAuthGitHubCallback(w http.ResponseWriter, r *http.Request) {
	stateParam := r.FormValue("state")
	if !oauth.ValidateOAuthState(r, stateParam) {
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}

	token, err := oauth.GitHubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.ErrorHandler(err, "GitHub OAuth token exchange failed")
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	userInfo, err := getGitHubUserInfo(token.AccessToken)
	if err != nil {
		utils.ErrorHandler(err, "Failed to get user info from GitHub")
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	username := userInfo.Login
	user, err := sqlconnect.FindOrCreateOAuthUser(
		"github",
		strconv.Itoa(userInfo.ID),
		userInfo.Email,
		userInfo.Name,
		"",
		username,
	)
	if err != nil {
		utils.ErrorHandler(err, "Failed to create/find GitHub OAuth user")
		http.Error(w, "Failed to process user", http.StatusInternalServerError)
		return
	}

	jwtToken, err := utils.SignToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.ErrorHandler(err, "Failed to create JWT token")
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	})

	oauth.ClearOAuthState(w)
	http.Redirect(w, r, "http://dictionary-bucker.s3-website-us-east-1.amazonaws.com/index.html", http.StatusTemporaryRedirect)
}

func getGitHubUserInfo(accessToken string) (*models.GitHubUserInfo, error) {
	req, err := http.NewRequest("GET", oauth.GitHubUserInfoAPI, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var userInfo models.GitHubUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}

	return &userInfo, nil
}
