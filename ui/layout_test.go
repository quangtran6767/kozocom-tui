package ui

import "testing"

func TestCalculateLayout(t *testing.T) {
	tests := []struct {
		name        string
		totalWidth  int
		totalHeight int
		want        LayoutDimemsions
	}{
		{
			name:        "standard terminal 100x50",
			totalWidth:  100,
			totalHeight: 50,
			want: LayoutDimemsions{
				SidebarWidth:          25, // 100 * 0.25
				SidebarHeight:         40, // usable(49) - 9
				SidebarUserInfoHeight: 9,  // int(49 * 0.20) = int(9.8) = 9
				ContentWidth:          75, // 100 - 25
				TopHeight:             31, // int(49 * 0.65) = int(31.85) = 31
				BottomHeight:          18, // 49 - 31
			},
		},
		{
			name:        "zero size terminal",
			totalWidth:  0,
			totalHeight: 0,
			want: LayoutDimemsions{
				SidebarWidth:          0,
				SidebarHeight:         -1, // usable(-1) - 0
				SidebarUserInfoHeight: 0,  // int(-1 * 0.20) = 0
				ContentWidth:          0,
				TopHeight:             0,  // int(-1 * 0.65) = 0
				BottomHeight:          -1, // -1 - 0
			},
		},
		{
			name:        "wide terminal 200x40",
			totalWidth:  200,
			totalHeight: 40,
			want: LayoutDimemsions{
				SidebarWidth:          50,  // 200 * 0.25
				SidebarHeight:         32,  // usable(39) - 7
				SidebarUserInfoHeight: 7,   // int(39 * 0.20) = int(7.8) = 7
				ContentWidth:          150, // 200 - 50
				TopHeight:             25,  // int(39 * 0.65) = int(25.35) = 25
				BottomHeight:          14,  // 39 - 25
			},
		},
		{
			name:        "odd-number terminal 79x33",
			totalWidth:  79,
			totalHeight: 33,
			want: LayoutDimemsions{
				SidebarWidth:          19, // int(79 * 0.25) = int(19.75) = 19
				SidebarHeight:         26, // usable(32) - 6
				SidebarUserInfoHeight: 6,  // int(32 * 0.20) = int(6.4) = 6
				ContentWidth:          60, // 79 - 19
				TopHeight:             20, // int(32 * 0.65) = int(20.8) = 20
				BottomHeight:          12, // 32 - 20
			},
		},
		{
			name:        "very small terminal 4x3",
			totalWidth:  4,
			totalHeight: 3,
			want: LayoutDimemsions{
				SidebarWidth:          1, // int(4 * 0.25) = 1
				SidebarHeight:         2, // usable(2) - 0
				SidebarUserInfoHeight: 0, // int(2 * 0.20) = int(0.4) = 0
				ContentWidth:          3, // 4 - 1
				TopHeight:             1, // int(2 * 0.65) = int(1.3) = 1
				BottomHeight:          1, // 2 - 1
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateLayout(tc.totalWidth, tc.totalHeight)

			if got != tc.want {
				t.Errorf("\ngot: %+v\nwant: %+v", got, tc.want)
			}
		})
	}
}

func TestCalculateLayout_Invariants(t *testing.T) {
	sizes := []struct {
		width, height int
	}{
		{100, 50},
		{79, 33},
		{4, 3},
		{0, 0},
		{220, 55},
	}

	for _, s := range sizes {
		t.Run("invariants", func(t *testing.T) {
			dims := CalculateLayout(s.width, s.height)
			usableHeight := s.height - HelpBarHeight

			// Horizontal: sidebar + content must equal total width
			if dims.SidebarWidth+dims.ContentWidth != s.width {
				t.Errorf(
					"%dx%d: horizontal mismatch: SidebarWidth(%d) + ContentWidth(%d) = %d, want %d",
					s.width, s.height,
					dims.SidebarWidth, dims.ContentWidth,
					dims.SidebarWidth+dims.ContentWidth, s.width,
				)
			}

			// Vertical right: top + bottom must equal usable height (total - help bar)
			if dims.TopHeight+dims.BottomHeight != usableHeight {
				t.Errorf(
					"%dx%d: vertical mismatch: TopHeight(%d) + BottomHeight(%d) = %d, want %d (usable)",
					s.width, s.height,
					dims.TopHeight, dims.BottomHeight,
					dims.TopHeight+dims.BottomHeight, usableHeight,
				)
			}

			// Vertical left: sidebar menu + user info must equal usable height
			if dims.SidebarHeight+dims.SidebarUserInfoHeight != usableHeight {
				t.Errorf(
					"%dx%d: sidebar vertical mismatch: SidebarHeight(%d) + SidebarUserInfoHeight(%d) = %d, want %d (usable)",
					s.width, s.height,
					dims.SidebarHeight, dims.SidebarUserInfoHeight,
					dims.SidebarHeight+dims.SidebarUserInfoHeight, usableHeight,
				)
			}
		})
	}
}
