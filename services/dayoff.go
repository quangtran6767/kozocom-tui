package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/quangtran6767/kozocom-tui/config"
	"github.com/quangtran6767/kozocom-tui/messages"
)

type responseMeta struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"error_message"`
}

type floatValue float64

func (f *floatValue) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		*f = 0
		return nil
	}

	var number float64
	if err := json.Unmarshal(data, &number); err == nil {
		*f = floatValue(number)
		return nil
	}

	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return err
	}
	if text == "" {
		*f = 0
		return nil
	}

	number, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return err
	}
	*f = floatValue(number)
	return nil
}

type leaveBalanceResponse struct {
	Data leaveBalancePayload `json:"data"`
	Meta responseMeta        `json:"meta"`
}

type leaveBalancePayload struct {
	DaysOffTaken       floatValue         `json:"days_off_taken"`
	DaysOffTakenText   string             `json:"days_off_taken_text"`
	LeaveBalanceTotal  floatValue         `json:"leave_balance_total"`
	DaysOffTakenUnpaid floatValue         `json:"days_off_taken_unpaid"`
	LeaveBalance       []leaveBalanceItem `json:"leave_balance"`
	Approver           string             `json:"approver"`
	InvolvingPersons   []string           `json:"involving_persons"`
}

type leaveBalanceItem struct {
	ID              int        `json:"id"`
	Year            int        `json:"year"`
	AnnualLeaveDays floatValue `json:"annual_leave_days"`
	ExpiryDate      string     `json:"expiry_date"`
	LeaveTypeID     int        `json:"leave_type_id"`
	LeaveType       string     `json:"leave_type"`
	RemainingDays   floatValue `json:"remaining_days"`
	Text            string     `json:"text"`
}

type employeeListResponse struct {
	Data []employeePayload `json:"data"`
	Meta responseMeta      `json:"meta"`
}

type employeePayload struct {
	AccountID   int    `json:"account_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	AccountName string `json:"account_name"`
	FullName    string `json:"fullName"`
	FullNameAlt string `json:"full_name"`
	Nickname    string `json:"nickname"`
}

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

		data, errMsg := decodeLeaveBalanceResponse(resp.Body)
		if errMsg != "" {
			return messages.LeaveBalanceFailMsg{Error: errMsg}
		}

		return messages.LeaveBalanceMsg{Data: data}
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

		data, errMsg := decodeEmployeesResponse(resp.Body)
		if errMsg != "" {
			return messages.ApproversFailMsg{Error: errMsg}
		}

		return messages.ApproversMsg{Data: data}
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

		data, errMsg := decodeEmployeesResponse(resp.Body)
		if errMsg != "" {
			return messages.AllEmployeesFailMsg{Error: errMsg}
		}

		return messages.AllEmployeesMsg{Data: data}
	}
}

func decodeLeaveBalanceResponse(body io.Reader) (messages.IDayLeaveBalance, string) {
	var result leaveBalanceResponse
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return messages.IDayLeaveBalance{}, "Failed to decode response"
	}
	if result.Meta.Code != http.StatusOK {
		return messages.IDayLeaveBalance{}, result.Meta.ErrorMessage
	}

	items := make([]messages.ILeaveBalance, 0, len(result.Data.LeaveBalance))
	for _, item := range result.Data.LeaveBalance {
		items = append(items, item.toMessage())
	}

	return messages.IDayLeaveBalance{
		DaysOffTaken:       float64(result.Data.DaysOffTaken),
		DaysOffTakenText:   result.Data.DaysOffTakenText,
		LeaveBalanceTotal:  float64(result.Data.LeaveBalanceTotal),
		DaysOffTakenUnpaid: float64(result.Data.DaysOffTakenUnpaid),
		LeaveBalance:       items,
		Approver:           result.Data.Approver,
		InvolvingPersons:   result.Data.InvolvingPersons,
	}, ""
}

func decodeEmployeesResponse(body io.Reader) ([]messages.EmployeeItem, string) {
	var result employeeListResponse
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, "Failed to decode response"
	}
	if result.Meta.Code != http.StatusOK {
		return nil, result.Meta.ErrorMessage
	}

	items := make([]messages.EmployeeItem, 0, len(result.Data))
	for _, item := range result.Data {
		items = append(items, item.toMessage())
	}

	return items, ""
}

func (item leaveBalanceItem) toMessage() messages.ILeaveBalance {
	return messages.ILeaveBalance{
		ID:              item.ID,
		Year:            item.Year,
		AnnualLeaveDays: float64(item.AnnualLeaveDays),
		ExpiryDate:      item.ExpiryDate,
		LeaveTypeID:     item.LeaveTypeID,
		LeaveType:       item.LeaveType,
		RemainingDays:   float64(item.RemainingDays),
		Text:            item.Text,
	}
}

func (item employeePayload) toMessage() messages.EmployeeItem {
	return messages.EmployeeItem{
		AccountID:   item.AccountID,
		FirstName:   item.FirstName,
		LastName:    item.LastName,
		Email:       item.Email,
		AccountName: item.AccountName,
		FullName:    firstNonEmpty(item.FullName, item.FullNameAlt),
		Nickname:    item.Nickname,
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
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
