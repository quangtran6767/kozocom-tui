package messages

type AuthCheckSuccessMsg struct {
	UserID string
	Email  string
}

type AuthCheckFailMsg struct{}

type NoTokenMsg struct{}

type LoginSuccessMsg struct {
	Token  string
	UserID string
	Email  string
}

type LoginFailMsg struct {
	Error string
}
