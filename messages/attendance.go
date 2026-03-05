package messages

type CheckinStatusMsg struct {
	IsCheckedIn bool
}

type CheckinStatusFailMsg struct {
	Error string
}

type CheckinSuccessMsg struct{}

type CheckinFailMsg struct {
	Error string
}

type AttendanceLogsMsg struct {
	Data interface{}
}

type AttendanceLogsFailMsg struct {
	Error string
}
