package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	tea "charm.land/bubbletea/v2"
	"github.com/quangtran6767/kozocom-tui/config"
	"github.com/quangtran6767/kozocom-tui/messages"
)

type DayOffRequestQuery struct {
	DateStart string
	DateEnd   string
}

// FetchLeaveBalance calls GET /api/v1/user/day-off-request/leave-balance
func FetchLeaveBalance(token string) tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest("GET", config.BaseURL+"/user/day-off-request/leave-balance", nil)
		if err != nil {
			return messages.LeaveBalanceFailMsg{Error: err.Error()}
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := httpClient.Do(req)
		if err != nil {
			return messages.LeaveBalanceFailMsg{Error: "Cannot connect to server"}
		}
		defer resp.Body.Close()

		var result struct {
			Data messages.IDayLeaveBalance `json:"data"`
			Meta struct {
				Code         int    `json:"code"`
				ErrorMessage string `json:"error_message"`
			} `json:"meta"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return messages.LeaveBalanceFailMsg{Error: "Failed to decode response"}
		}
		if result.Meta.Code != http.StatusOK {
			return messages.LeaveBalanceFailMsg{Error: result.Meta.ErrorMessage}
		}

		return messages.LeaveBalanceMsg{Data: result.Data}
	}
}

// FetchApprovers calls GET /api/v1/user/approver-members
func FetchApprovers(token string) tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest("GET", config.BaseURL+"/user/approver-members", nil)
		if err != nil {
			return messages.ApproversFailMsg{Error: err.Error()}
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := httpClient.Do(req)
		if err != nil {
			return messages.ApproversFailMsg{Error: "Cannot connect to server"}
		}
		defer resp.Body.Close()

		var result struct {
			Data []messages.EmployeeItem `json:"data"`
			Meta struct {
				Code         int    `json:"code"`
				ErrorMessage string `json:"error_message"`
			} `json:"meta"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return messages.ApproversFailMsg{Error: "Failed to decode response"}
		}
		if result.Meta.Code != http.StatusOK {
			return messages.ApproversFailMsg{Error: result.Meta.ErrorMessage}
		}

		return messages.ApproversMsg{Data: result.Data}
	}
}

// FetchAllEmployees calls GET /api/v1/user/all-employees
func FetchAllEmployees(token string) tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest("GET", config.BaseURL+"/user/all-employees", nil)
		if err != nil {
			return messages.AllEmployeesFailMsg{Error: err.Error()}
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := httpClient.Do(req)
		if err != nil {
			return messages.AllEmployeesFailMsg{Error: "Cannot connect to server"}
		}
		defer resp.Body.Close()

		var result struct {
			Data []messages.EmployeeItem `json:"data"`
			Meta struct {
				Code         int    `json:"code"`
				ErrorMessage string `json:"error_message"`
			} `json:"meta"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return messages.AllEmployeesFailMsg{Error: "Failed to decode response"}
		}
		if result.Meta.Code != http.StatusOK {
			return messages.AllEmployeesFailMsg{Error: result.Meta.ErrorMessage}
		}

		return messages.AllEmployeesMsg{Data: result.Data}
	}
}

// FetchDayOffRequests calls GET /api/v1/user/day-off-requests
func FetchDayOffRequests(token string, query DayOffRequestQuery) tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest("GET", buildDayOffRequestsURL(query), nil)
		if err != nil {
			return messages.DayOffRequestsFailMsg{Error: err.Error()}
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := httpClient.Do(req)
		if err != nil {
			return messages.DayOffRequestsFailMsg{Error: "Cannot connect to server"}
		}
		defer resp.Body.Close()

		items, errMsg := decodeDayOffRequestsResponse(resp.Body)
		if errMsg != "" {
			return messages.DayOffRequestsFailMsg{Error: errMsg}
		}

		return messages.DayOffRequestsMsg{Data: items}
	}
}

func decodeDayOffRequestsResponse(body io.Reader) ([]messages.DayOffRequestItem, string) {
	var result struct {
		Data struct {
			Items []messages.DayOffRequestItem `json:"data"`
		} `json:"data"`
		Meta struct {
			Code         int    `json:"code"`
			ErrorMessage string `json:"error_message"`
		} `json:"meta"`
	}

	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, "Failed to decode response"
	}

	if result.Meta.Code != http.StatusOK {
		return nil, result.Meta.ErrorMessage
	}

	return result.Data.Items, ""
}

func buildDayOffRequestsURL(query DayOffRequestQuery) string {
	values := url.Values{}
	if query.DateStart != "" {
		values.Set("date_start", query.DateStart)
	}
	if query.DateEnd != "" {
		values.Set("date_end", query.DateEnd)
	}
	if encoded := values.Encode(); encoded != "" {
		return config.BaseURL + "/user/day-off-requests?" + encoded
	}
	return config.BaseURL + "/user/day-off-requests"
}

// CalculateLeaveDays calls GET /api/v1/user/day-off-request/leave-days-calculate
func CalculateLeaveDays(token, startTime, endTime string) tea.Cmd {
	return func() tea.Msg {
		urlStr := fmt.Sprintf("%s/user/day-off-request/leave-days-calculate?start_time=%s&end_time=%s", config.BaseURL, url.QueryEscape(startTime), url.QueryEscape(endTime))
		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			return messages.LeaveDaysCalcFailMsg{Error: err.Error()}
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := httpClient.Do(req)
		if err != nil {
			return messages.LeaveDaysCalcFailMsg{Error: "Cannot connect to server"}
		}
		defer resp.Body.Close()

		var result struct {
			Data struct {
				Result float64 `json:"result"`
			} `json:"data"`
			Meta struct {
				Code         int    `json:"code"`
				ErrorMessage string `json:"error_message"`
			} `json:"meta"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return messages.LeaveDaysCalcFailMsg{Error: "Failed to decode response"}
		}
		if result.Meta.Code != http.StatusOK {
			return messages.LeaveDaysCalcFailMsg{Error: result.Meta.ErrorMessage}
		}

		return messages.LeaveDaysCalcMsg{Result: result.Data.Result}
	}
}

// CreateDayOffRequest calls POST /api/v1/user/day-off-requests
func CreateDayOffRequest(token string, payload map[string]interface{}) tea.Cmd {
	return func() tea.Msg {
		body, err := json.Marshal(payload)
		if err != nil {
			return messages.CreateDayOffFailMsg{Error: "Failed to marshal payload"}
		}

		req, err := http.NewRequest("POST", config.BaseURL+"/user/day-off-requests", bytes.NewBuffer(body))
		if err != nil {
			return messages.CreateDayOffFailMsg{Error: err.Error()}
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			return messages.CreateDayOffFailMsg{Error: "Cannot connect to server"}
		}
		defer resp.Body.Close()

		var result struct {
			Meta struct {
				Code         int    `json:"code"`
				ErrorMessage string `json:"error_message"`
			} `json:"meta"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return messages.CreateDayOffFailMsg{Error: "Failed to decode response"}
		}
		if result.Meta.Code != http.StatusOK {
			return messages.CreateDayOffFailMsg{Error: result.Meta.ErrorMessage}
		}

		return messages.CreateDayOffSuccessMsg{}
	}
}
