package auth

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/quangtran6767/kozocom-tui/messages"
	"github.com/quangtran6767/kozocom-tui/testutil"
)

// Helpers
func modelAtPhase(phase Phase) Model {
	m := New()
	m.phase = phase
	return m
}

// cmdIsNil check no command emitted.
func cmdIsNil(cmd tea.Cmd) bool {
	return cmd == nil
}

// -------------------------------------------------------------------
// New() — initial state
// -------------------------------------------------------------------
func TestNew_Phase_IsCheckingAuth(t *testing.T) {
	m := New()
	if m.phase != PhaseCheckingAuth {
		t.Errorf("expected PhaseCheckingAuth, got %v", m.phase)
	}
}

func TestNew_FocusIndex_IsEmail(t *testing.T) {
	m := New()
	if m.focusIndex != focusEmail {
		t.Errorf("expected focusEmail (%d), got %d", focusEmail, m.focusIndex)
	}
}

func TestNew_ErrorMsg_IsEmpty(t *testing.T) {
	m := New()
	if m.errMsg != "" {
		t.Errorf("expected empty errMsg, got %q", m.errMsg)
	}
}
func TestNew_Token_IsEmpty(t *testing.T) {
	m := New()
	if m.token != "" {
		t.Errorf("expected empty token, got %q", m.token)
	}
}

// -------------------------------------------------------------------
// Getters
// -------------------------------------------------------------------
func TestIsDone_ReturnsFalse_WhenNotPhaseDone(t *testing.T) {
	phases := []Phase{PhaseCheckingAuth, PhaseLoggingIn, PhaseLoginForm}
	for _, p := range phases {
		m := modelAtPhase(p)
		if m.IsDone() {
			t.Errorf("IsDone() should be false for phase %v", p)
		}
	}
}

func TestIsDone_ReturnsTrue_WhenPhaseDone(t *testing.T) {
	m := modelAtPhase(PhaseDone)
	if !m.IsDone() {
		t.Error("IsDone() should be true for PhaseDone")
	}
}

func TestToken_ReturnsStoredToken(t *testing.T) {
	m := New()
	m.token = "abc123"
	if m.Token() != "abc123" {
		t.Errorf("Token() = %q, want %q", m.Token(), "abc123")
	}
}
func TestUserID_ReturnsStoredUserID(t *testing.T) {
	m := New()
	m.userID = "user-42"
	if m.UserID() != "user-42" {
		t.Errorf("UserID() = %q, want %q", m.UserID(), "user-42")
	}
}
func TestEmail_ReturnsStoredEmail(t *testing.T) {
	m := New()
	m.email = "hello@world.com"
	if m.Email() != "hello@world.com" {
		t.Errorf("Email() = %q, want %q", m.Email(), "hello@world.com")
	}
}

// -------------------------------------------------------------------
// updateCheckingAuth — PhaseCheckingAuth
// -------------------------------------------------------------------
func TestUpdateCheckingAuth_Success_TransitionsToDone(t *testing.T) {
	m := New()

	msg := messages.AuthCheckSuccessMsg{UserID: "u123", Email: "a@b.com"}
	m2, cmd := m.Update(msg)

	if m2.phase != PhaseDone {
		t.Errorf("expected PhaseDone, got %v", m2.phase)
	}
	if m2.userID != "u123" {
		t.Errorf("userID = %q, want %q", m2.userID, "u123")
	}
	if m2.email != "a@b.com" {
		t.Errorf("email = %q, want %q", m2.email, "a@b.com")
	}
	if !cmdIsNil(cmd) {
		t.Error("expected no cmd after AuthCheckSuccessMsg")
	}
}

func TestUpdateCheckingAuth_Fail_TransitionsToLoginForm(t *testing.T) {
	m := New() // phase = PhaseCheckingAuth
	m2, cmd := m.Update(messages.AuthCheckFailMsg{})
	if m2.phase != PhaseLoginForm {
		t.Errorf("expected PhaseLoginForm, got %v", m2.phase)
	}
	if m2.token != "" {
		t.Errorf("token should be cleared, got %q", m2.token)
	}
	// Should emit Focus() cmd for email input
	if cmdIsNil(cmd) {
		t.Error("expected a Focus cmd, got nil")
	}
}

// -------------------------------------------------------------------
// updateLoginForm — PhaseLoginForm
// -------------------------------------------------------------------
func TestUpdateLoginForm_Tab_SwitchesFocusToPassword(t *testing.T) {
	m := modelAtPhase(PhaseLoginForm)
	tabMsg := tea.KeyPressMsg{Code: tea.KeyTab}
	m2, cmd := m.Update(tabMsg)
	if m2.focusIndex != focusPassword {
		t.Errorf("expected focusPassword (%d), got %d", focusPassword, m2.focusIndex)
	}
	if cmdIsNil(cmd) {
		t.Error("expected a Focus cmd for password input, got nil")
	}
}

func TestUpdateLoginForm_ShiftTab_SwitchesFocusBackToEmail(t *testing.T) {
	m := modelAtPhase(PhaseLoginForm)
	m.focusIndex = focusPassword // start at password
	shiftTabMsg := tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift}
	m2, cmd := m.Update(shiftTabMsg)
	if m2.focusIndex != focusEmail {
		t.Errorf("expected focusEmail (%d), got %d", focusEmail, m2.focusIndex)
	}
	if cmdIsNil(cmd) {
		t.Error("expected a Focus cmd for email input, got nil")
	}
}

func TestUpdateLoginForm_Enter_EmptyFields_ShowsError(t *testing.T) {
	m := modelAtPhase(PhaseLoginForm)
	// emailInput and passInput are empty by default
	enterMsg := tea.KeyPressMsg{Code: tea.KeyEnter}
	m2, cmd := m.Update(enterMsg)
	if m2.errMsg == "" {
		t.Error("expected errMsg to be set when fields are empty")
	}
	if m2.phase != PhaseLoginForm {
		t.Errorf("expected phase to stay PhaseLoginForm, got %v", m2.phase)
	}
	if !cmdIsNil(cmd) {
		t.Error("expected no cmd when validation fails")
	}
}

func TestUpdateLoginForm_Enter_ValidFields_TransitionsToLoggingIn(t *testing.T) {
	m := modelAtPhase(PhaseLoginForm)
	m.emailInput.SetValue("user@test.com")
	m.passInput.SetValue("secret")
	enterMsg := tea.KeyPressMsg{Code: tea.KeyEnter}
	m2, cmd := m.Update(enterMsg)
	if m2.phase != PhaseLoggingIn {
		t.Errorf("expected PhaseLoggingIn, got %v", m2.phase)
	}
	if m2.errMsg != "" {
		t.Errorf("expected no errMsg, got %q", m2.errMsg)
	}
	if cmdIsNil(cmd) {
		t.Error("expected a Batch cmd (spinner + login), got nil")
	}
}

// -------------------------------------------------------------------
// updateLoggingIn — PhaseLoggingIn
// -------------------------------------------------------------------
func TestUpdateLoggingIn_Success_TransitionsToDone(t *testing.T) {
	testutil.RedirectConfigDir(t) // prevent real file writes from config.SaveToken
	m := modelAtPhase(PhaseLoggingIn)
	msg := messages.LoginSuccessMsg{
		Token:  "tok-xyz",
		UserID: "u-99",
		Email:  "login@ok.com",
	}
	m2, _ := m.Update(msg)
	if m2.phase != PhaseDone {
		t.Errorf("expected PhaseDone, got %v", m2.phase)
	}
	if m2.token != "tok-xyz" {
		t.Errorf("token = %q, want %q", m2.token, "tok-xyz")
	}
	if m2.userID != "u-99" {
		t.Errorf("userID = %q, want %q", m2.userID, "u-99")
	}
	if m2.email != "login@ok.com" {
		t.Errorf("email = %q, want %q", m2.email, "login@ok.com")
	}
}
func TestUpdateLoggingIn_Fail_TransitionsToLoginForm(t *testing.T) {
	m := modelAtPhase(PhaseLoggingIn)
	msg := messages.LoginFailMsg{Error: "Invalid credentials"}
	m2, cmd := m.Update(msg)
	if m2.phase != PhaseLoginForm {
		t.Errorf("expected PhaseLoginForm, got %v", m2.phase)
	}
	if m2.errMsg != "Invalid credentials" {
		t.Errorf("errMsg = %q, want %q", m2.errMsg, "Invalid credentials")
	}
	if cmdIsNil(cmd) {
		t.Error("expected a Focus cmd for email input, got nil")
	}
}
func TestUpdateLoggingIn_Fail_EmptyError(t *testing.T) {
	m := modelAtPhase(PhaseLoggingIn)
	m2, _ := m.Update(messages.LoginFailMsg{Error: ""})
	if m2.phase != PhaseLoginForm {
		t.Errorf("expected PhaseLoginForm, got %v", m2.phase)
	}
}

// -------------------------------------------------------------------
// Update() dispatcher
// -------------------------------------------------------------------
func TestUpdate_WhenPhaseDone_ReturnsNoCmd(t *testing.T) {
	m := modelAtPhase(PhaseDone)
	m2, cmd := m.Update(messages.LoginSuccessMsg{Token: "x"})
	if m2.phase != PhaseDone {
		t.Errorf("phase should stay PhaseDone, got %v", m2.phase)
	}
	if !cmdIsNil(cmd) {
		t.Error("expected no cmd when phase is PhaseDone")
	}
}
