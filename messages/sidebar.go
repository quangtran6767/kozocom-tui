package messages

type MenuItemID int

const (
	MenuAttendanceLog MenuItemID = iota
)

type SidebarItemSelectedMsg struct {
	Item MenuItemID
}
