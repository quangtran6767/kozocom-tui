package messages

type AuthCheckSuccessMsg struct {
	UserID int
}

type AuthCheckFailMsg struct{}

type NoTokenMsg struct{}

type LoginSuccessMsg struct {
	Token  string
	UserID int
}

type LoginFailMsg struct {
	Error string
}
