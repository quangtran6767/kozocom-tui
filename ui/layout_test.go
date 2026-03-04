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
				SidebarHeight:         40, // 50 - 10
				SidebarUserInfoHeight: 10, // 50 * 0.20
				ContentWidth:          75, // 100 - 25
				TopHeight:             32, // 50 * 0.65 = 32.5 → int = 32
				BottomHeight:          18, // 50 - 32
			},
		},
		{
			name:        "zero size terminal",
			totalWidth:  0,
			totalHeight: 0,
			want: LayoutDimemsions{
				SidebarWidth:          0,
				SidebarHeight:         0,
				SidebarUserInfoHeight: 0,
				ContentWidth:          0,
				TopHeight:             0,
				BottomHeight:          0,
			},
		},
		{
			name:        "wide terminal 200x40",
			totalWidth:  200,
			totalHeight: 40,
			want: LayoutDimemsions{
				SidebarWidth:          50,  // 200 * 0.25
				SidebarHeight:         32,  // 40 - 8
				SidebarUserInfoHeight: 8,   // 40 * 0.20
				ContentWidth:          150, // 200 - 50
				TopHeight:             26,  // 40 * 0.65
				BottomHeight:          14,  // 40 - 26
			},
		},
		{
			name:        "odd-number terminal 79x33",
			totalWidth:  79,
			totalHeight: 33,
			want: LayoutDimemsions{
				SidebarWidth:          19, // int(79 * 0.25) = int(19.75) = 19
				SidebarHeight:         27, // 33 - 6
				SidebarUserInfoHeight: 6,  // int(33 * 0.20) = int(6.6) = 6
				ContentWidth:          60, // 79 - 19
				TopHeight:             21, // int(33 * 0.65) = int(21.45) = 21
				BottomHeight:          12, // 33 - 21
			},
		},
		{
			name:        "very small terminal 4x3",
			totalWidth:  4,
			totalHeight: 3,
			want: LayoutDimemsions{
				SidebarWidth:          1, // int(4 * 0.25) = 1
				SidebarHeight:         3, // 3 - 0
				SidebarUserInfoHeight: 0, // int(3 * 0.20) = int(0.6) = 0
				ContentWidth:          3, // 4 - 1
				TopHeight:             1, // int(3 * 0.65) = int(1.95) = 1
				BottomHeight:          2, // 3 - 1
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

			// Horizontal: sidebar + content must equal total width
			if dims.SidebarWidth+dims.ContentWidth != s.width {
				t.Errorf(
					"%dx%d: horizontal mismatch: SidebarWidth(%d) + ContentWidth(%d) = %d, want %d",
					s.width, s.height,
					dims.SidebarWidth, dims.ContentWidth,
					dims.SidebarWidth+dims.ContentWidth, s.width,
				)
			}

			// Vertical right: top + bottom must equal total height
			if dims.TopHeight+dims.BottomHeight != s.height {
				t.Errorf(
					"%dx%d: vertical mismatch: TopHeight(%d) + BottomHeight(%d) = %d, want %d",
					s.width, s.height,
					dims.TopHeight, dims.BottomHeight,
					dims.TopHeight+dims.BottomHeight, s.height,
				)
			}

			// Vertical left: sidebar menu + user info must equal total height
			if dims.SidebarHeight+dims.SidebarUserInfoHeight != s.height {
				t.Errorf(
					"%dx%d: sidebar vertical mismatch: SidebarHeight(%d) + SidebarUserInfoHeight(%d) = %d, want %d",
					s.width, s.height,
					dims.SidebarHeight, dims.SidebarUserInfoHeight,
					dims.SidebarHeight+dims.SidebarUserInfoHeight, s.height,
				)
			}
		})
	}
}
