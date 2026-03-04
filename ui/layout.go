package ui

import "charm.land/lipgloss/v2"

const (
	SidebarRatio         = 0.25
	SidebarUserInfoRatio = 0.20
	TopContentRatio      = 0.65
)

type LayoutDimemsions struct {
	// Left
	SidebarWidth          int
	SidebarHeight         int
	SidebarMenuHeight     int
	SidebarUserInfoHeight int
	// Right
	ContentWidth int
	TopHeight    int
	BottomHeight int
}

// CalculateLayout calculates the layout dimensions based on terminal size.
// @param totalWidth - the total width of the terminal
// @param totalHeight - the total height of the terminal
// @return LayoutDimemsions - the layout dimensions
func CalculateLayout(totalWidth, totalHeight int) LayoutDimemsions {
	sidebarW := int(float64(totalWidth) * SidebarRatio)
	contentW := totalWidth - sidebarW

	topH := int(float64(totalHeight) * TopContentRatio)
	bottomH := totalHeight - topH

	sidebarUserInforH := int(float64(totalHeight) * SidebarUserInfoRatio)
	sidebarMenuH := totalHeight - sidebarUserInforH

	return LayoutDimemsions{
		SidebarWidth:          sidebarW,
		SidebarHeight:         sidebarMenuH,
		SidebarUserInfoHeight: sidebarUserInforH,
		ContentWidth:          contentW,
		TopHeight:             topH,
		BottomHeight:          bottomH,
	}
}

// RenderLayout reassemable all the panels into a single layout.
func RenderLayout(sidebarMenu, sidebarUserInfo, topContent, bottomContent string) string {
	leftSide := lipgloss.JoinVertical(lipgloss.Left, sidebarUserInfo, sidebarMenu)
	rightSide := lipgloss.JoinVertical(lipgloss.Left, topContent, bottomContent)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftSide, rightSide)
}
