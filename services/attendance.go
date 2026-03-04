package services

import (
	"bytes"
	"encoding/json"
	"net/http"

	tea "charm.land/bubbletea/v2"
	"github.com/quangtran6767/kozocom-tui/config"
	"github.com/quangtran6767/kozocom-tui/messages"
)

// CheckTodayStatus calls GET /api/v1/user/attendance-logs/checkin-today-status
// to check if the user has already checked in today.
// Should be called right after transitioning to StateMain.
//
// @param token string - Bearer token
// @return tea.Cmd - returns CheckinStatusMsg or CheckinStatusFailMsg
func CheckinTodayStatus(token string) tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest("GET",
			config.BaseURL+"/user/attendance-logs/checkin-today-status", nil,
		)

		if err != nil {
			config.DebugLog.Println("CheckTodayStatus: failed to create request", err)
			return messages.CheckinStatusFailMsg{Error: "Failed to create request"}
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := httpClient.Do(req)
		if err != nil {
			config.DebugLog.Println("CheckTodayStatus: failed to send request", err)
			return messages.CheckinStatusFailMsg{Error: "Cannot connect to server"}
		}
		defer resp.Body.Close()

		var result struct {
			Data struct {
				CanCheckin bool `json:"can_checkin"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			config.DebugLog.Println("CheckTodayStatus: failed to decode response", err)
			return messages.CheckinStatusFailMsg{Error: "Cannot decode response"}
		}

		return messages.CheckinStatusMsg{IsCheckedIn: !result.Data.CanCheckin}
	}
}

// Checkin calls POST /api/v1/user/attendance-logs/checkin to perform check-in.
// Should only be called when isCheckedIn == false (guard ở phía TUI trước khi gọi).
//
// @param token string - Bearer token
// @return tea.Cmd - returns CheckinSuccessMsg or CheckinFailMsg
func Checkin(token string) tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest("POST",
			config.BaseURL+"/user/attendance-logs/checkin", bytes.NewBuffer(nil),
		)
		if err != nil {
			config.DebugLog.Println("Checkin: failed to create request", err)
			return messages.CheckinFailMsg{Error: "Failed to create request"}
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			config.DebugLog.Println("Checkin: failed to send request", err)
			return messages.CheckinFailMsg{Error: "Cannot connect to server"}
		}
		defer resp.Body.Close()

		var result struct {
			Meta struct {
				Code         int    `json:"code"`
				ErrorMessage string `json:"error_message"`
			} `json:"meta"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			config.DebugLog.Println("Checkin: failed to decode response", err)
			return messages.CheckinFailMsg{Error: "Cannot decode response"}
		}
		if result.Meta.Code != http.StatusOK {
			config.DebugLog.Println("Checkin: server returned error:", result.Meta.ErrorMessage)
			return messages.CheckinFailMsg{Error: result.Meta.ErrorMessage}
		}

		return messages.CheckinSuccessMsg{}
	}
}
