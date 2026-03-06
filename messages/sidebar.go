package messages

type MenuItemID int

const (
	MenuAttendanceLog MenuItemID = iota
	MenuDayOffRequest
)

type SidebarItemSelectedMsg struct {
	Item MenuItemID
}
