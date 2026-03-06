package messages

type ILeaveBalance struct {
	ID              int     `json:"id"`
	Year            int     `json:"year"`
	AnnualLeaveDays float64 `json:"annual_leave_days"`
	ExpiryDate      string  `json:"expiry_date"`
	LeaveTypeID     int     `json:"leave_type_id"`
	LeaveType       string  `json:"leave_type"`
	RemainingDays   float64 `json:"remaining_days"`
	Text            string  `json:"text"`
}

type IDayLeaveBalance struct {
	DaysOffTaken       float64         `json:"days_off_taken"`
	DaysOffTakenText   string          `json:"days_off_taken_text"`
	LeaveBalanceTotal  float64         `json:"leave_balance_total"`
	DaysOffTakenUnpaid float64         `json:"days_off_taken_unpaid"`
	LeaveBalance       []ILeaveBalance `json:"leave_balance"`
	Approver           string          `json:"approver"`
	InvolvingPersons   []string        `json:"involving_persons"`
}

type EmployeeItem struct {
	AccountID   int    `json:"account_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	AccountName string `json:"account_name"`
	FullName    string `json:"full_name"`
	Nickname    string `json:"nickname"`
}

type DayOffRequestItem struct {
	RequestCode      string   `json:"request_code"`
	StartTime        string   `json:"start_time"`
	EndTime          string   `json:"end_time"`
	Reason           string   `json:"reason"`
	Status           string   `json:"status"`
	LeaveType        string   `json:"leave_type"`
	ApproversName    []string `json:"approvers_name"`
	NumberLeavesDays float64  `json:"total_days"`
	CreatedAt        string   `json:"created_at"`
}

type LeaveBalanceMsg struct {
	Data IDayLeaveBalance
}
type LeaveBalanceFailMsg struct{ Error string }

type ApproversMsg struct {
	Data []EmployeeItem
}
type ApproversFailMsg struct{ Error string }

type AllEmployeesMsg struct {
	Data []EmployeeItem
}
type AllEmployeesFailMsg struct{ Error string }

type DayOffRequestsMsg struct {
	Data []DayOffRequestItem
}
type DayOffRequestsFailMsg struct{ Error string }

type LeaveDaysCalcMsg struct {
	Result float64
}
type LeaveDaysCalcFailMsg struct{ Error string }

type CreateDayOffSuccessMsg struct{}
type CreateDayOffFailMsg struct{ Error string }
