package calendar

type Checkin struct {
	CheckinTime  string `json:"checkin_time"`
	CheckoutTime string `json:"checkout_time"`
}

type DayOff struct {
	StartTime             string `json:"start_time"`
	EndTime               string `json:"end_time"`
	TypeName              string `json:"type_name"`
	Type                  string `json:"type"`
	Reason                string `json:"reason"`
	SalaryEntitlementRate int    `json:"salary_entitlement_rate"`
}

type Overtime struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Type      string `json:"type"`
	Content   string `json:"content"`
}

type Holiday struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type DayLog struct {
	Date            string     `json:"-"`
	Checkins        []Checkin  `json:"checkin,omitempty"`
	HoursAttendance float64    `json:"hours_attendance,omitempty"`
	DayOffs         []DayOff   `json:"day_off,omitempty"`
	HoursDayOff     float64    `json:"hours_day_off,omitempty"`
	Overtimes       []Overtime `json:"overtime,omitempty"`
	HoursOvertime   float64    `json:"hours_overtime,omitempty"`
	Holidays        []Holiday  `json:"holiday,omitempty"`
}

type AttendanceData struct {
	ID             string            `json:"id"`
	EmployeeID     string            `json:"employee_id"`
	FirstName      string            `json:"first_name"`
	LastName       string            `json:"last_name"`
	StartDate      string            `json:"start_date"`
	ActualWork     float64           `json:"actual_work"`
	UnpaidLeave    float64           `json:"unpaid_leave"`
	DaysOffTaken   float64           `json:"days_off_taken"`
	HoursOverTime  float64           `json:"hours_over_time"`
	AttendanceLogs map[string]DayLog `json:"attendance_logs"`
}
