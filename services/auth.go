package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/quangtran6767/kozocom-tui/config"
	"github.com/quangtran6767/kozocom-tui/messages"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// CheckAuth call GET /me endpoint to check if the token is valid
// return tea.Cmd - run async in bubble tea runtime
func CheckAuth(token string) tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest("GET", config.BaseURL+"/user/infomation", nil)
		if err != nil {
			config.DebugLog.Println("CheckAuth: failed to create request", err)
			return messages.AuthCheckFailMsg{}
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := httpClient.Do(req)
		if err != nil {
			config.DebugLog.Println("CheckAuth: failed to send request", err)
			return messages.AuthCheckFailMsg{}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return messages.AuthCheckFailMsg{}
		}

		var result struct {
			Data struct {
				ID    string `json:"id"`
				Email string `json:"email"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			config.DebugLog.Println("CheckAuth: failed to decode response", err)
			return messages.AuthCheckFailMsg{}
		}

		return messages.AuthCheckSuccessMsg{UserID: result.Data.ID, Email: result.Data.Email}
	}
}

// Login call POST /login endpoint to login
// return tea.Cmd - run async in bubble tea runtime
func Login(email, password string) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": password,
		})

		resp, err := httpClient.Post(
			config.BaseURL+"/user/login",
			"application/json",
			bytes.NewBuffer(body),
		)

		if err != nil {
			return messages.LoginFailMsg{Error: "Can not connect to server"}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return messages.LoginFailMsg{Error: fmt.Sprintf("Login failed with status code: %d", resp.StatusCode)}
		}

		var result struct {
			Data struct {
				Token struct {
					AccessToken string `json:"access_token"`
				} `json:"token"`
				UserInfo struct {
					ID    int    `json:"id"`
					Email string `json:"id"`
				} `json:"user_info"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return messages.LoginFailMsg{Error: "Cannot decode response"}
		}

		return messages.LoginSuccessMsg{
			Token:  result.Data.Token.AccessToken,
			UserID: fmt.Sprintf("%d", result.Data.UserInfo.ID),
			Email:  result.Data.UserInfo.Email,
		}
	}
}
