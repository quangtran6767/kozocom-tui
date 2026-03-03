package messages

type AuthCheckSuccessMsg struct {
	UserID string
}

type AuthCheckFailMsg struct{}

type NoTokenMsg struct{}

type LoginSuccessMsg struct {
	Token  string
	UserID string
}

type LoginFailMsg struct {
	Error string
}
