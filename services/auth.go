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
		config.DebugLog.Println("CheckAuth: sending GET /me...")
		req, err := http.NewRequest("GET", config.BaseURL+"/me", nil)
		if err != nil {
			config.DebugLog.Printf("CheckAuth: failed to build request: %v", err)
			return messages.AuthCheckFailMsg{}
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := httpClient.Do(req)
		if err != nil {
			config.DebugLog.Printf("CheckAuth: HTTP error: %v", err)
			return messages.AuthCheckFailMsg{}
		}
		defer resp.Body.Close()

		config.DebugLog.Printf("CheckAuth: response status = %d", resp.StatusCode)

		if resp.StatusCode != http.StatusOK {
			return messages.AuthCheckFailMsg{}
		}

		var result struct {
			Success bool `json:"success"`
			User    int  `json:"user"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			config.DebugLog.Printf("CheckAuth: decode error: %v", err)
			return messages.AuthCheckFailMsg{}
		}

		config.DebugLog.Printf("CheckAuth: success=%v, user=%d", result.Success, result.User)

		if !result.Success {
			return messages.AuthCheckFailMsg{}
		}

		return messages.AuthCheckSuccessMsg{UserID: result.User}
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
			config.BaseURL+"/login",
			"application/json",
			bytes.NewBuffer(body),
		)

		if err != nil {
			return messages.LoginFailMsg{Error: "Can not connect to server"}
		}
		defer resp.Body.Close()

		var result struct {
			Success bool   `json:"success"`
			Token   string `json:"token"`
			User    int    `json:"user"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return messages.LoginFailMsg{Error: "Can not decode response"}
		}

		if !result.Success || resp.StatusCode != http.StatusOK {
			errMsg := result.Message
			if errMsg == "" {
				errMsg = fmt.Sprintf("Login failed with status code: %d", resp.StatusCode)
			}
			return messages.LoginFailMsg{Error: errMsg}
		}

		return messages.LoginSuccessMsg{
			Token:  result.Token,
			UserID: result.User,
		}
	}
}
