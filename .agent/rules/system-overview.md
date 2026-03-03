# kozocom-tui — System Overview & Rules

## 1. Tổng quan

`kozocom-tui` là một ứng dụng **Terminal User Interface (TUI)** viết bằng **Go**, sử dụng bộ thư viện Charmbracelet:

| Thư viện                  | Vai trò                                               |
| ------------------------- | ----------------------------------------------------- |
| `charm.land/bubbletea/v2` | Framework TUI theo mô hình Elm (Init / Update / View) |
| `charm.land/bubbles/v2`   | UI components: textinput, spinner                     |
| `charm.land/lipgloss/v2`  | Styling, layout, căn chỉnh giao diện terminal         |

Module Go: `github.com/quangtran6767/kozocom-tui`

API Backend: `http://localhost:8000` (kz-portal-api)

---

## 2. Kiến trúc tổng thể

Ứng dụng theo mô hình **Elm Architecture** của Bubble Tea, với **App State Machine** để quản lý luồng Auth → Main:

```
main.go
  └── app.go (appModel)            ← Root model, điều phối toàn bộ app
        ├── state: StateAuth        ← Auth phase (login / kiểm tra session)
        │     └── components/auth  ← Auth component (login form + spinner)
        └── state: StateMain        ← Main phase (sau khi đăng nhập thành công)
              ├── components/sidebar   ← Panel bên trái (menu điều hướng)
              ├── components/content   ← Panel phía trên bên phải (nội dung chính)
              ├── components/footer    ← Panel phía dưới bên phải (thông tin phụ)
              └── ui/                  ← Layout engine & style definitions
                    ├── layout.go      ← Tính toán và render layout tổng thể
                    ├── panel.go       ← Render từng panel với border + legend title
                    └── styles.go      ← Màu sắc, border, style toàn cục
```

### Sơ đồ layout terminal (StateMain)

```
┌─────────────────────────────────────────────────────────────┐
│  [1] Sidebar  │  [2] Content                                │
│               │                                             │
│               │   (65% chiều cao)                           │
│               ├─────────────────────────────────────────────│
│               │  [3] Footer                                 │
│               │   (35% chiều cao còn lại)                   │
└───────────────┴─────────────────────────────────────────────┘
  ← 25% width →  ←────────── 75% width ──────────────────────→
```

---

## 3. Cấu trúc file

```
kozocom-tui/
├── main.go                        # Entry point — khởi tạo Bubble Tea Program
├── app.go                         # Root model (appModel): AppState machine, Init, Update, View
├── go.mod / go.sum                # Go module dependencies
├── components/
│   ├── auth/auth.go               # Auth component: login form, spinner, auth flow phases
│   ├── sidebar/sidebar.go         # Component sidebar (menu bên trái)
│   ├── content/content.go         # Component content (nội dung chính)
│   └── footer/footer.go           # Component footer (thông tin phụ)
├── services/
│   └── auth.go                    # HTTP service: CheckAuth (GET /me), Login (POST /login)
├── config/
│   └── token.go                   # Token persistence: LoadToken, SaveToken (OS config dir)
├── messages/
│   └── auth.go                    # Bubble Tea message types cho auth events
└── ui/
    ├── layout.go                  # CalculateLayout + RenderLayout
    ├── panel.go                   # RenderPanel (border + legend)
    └── styles.go                  # Màu sắc và style toàn cục
```

---

## 4. App State Machine

`appModel` có 2 state chính điều phối bởi `AppState`:

```go
type AppState int

const (
    StateAuth AppState = iota  // Đang xử lý auth (check token / login form)
    StateMain                  // Đã login thành công → hiển thị main layout
)
```

### Luồng chuyển state

```
App Start
    │
    ▼
StateAuth ──→ auth.Init()
    │            ├── Có token trong file → CheckAuth (GET /me) → spinner
    │            │       ├── Success → PhaseDone → StateMain
    │            │       └── Fail   → PhaseLoginForm (hiện form)
    │            └── Không có token → PhaseLoginForm (hiện form)
    │
    │  (auth.IsDone() == true)
    ▼
StateMain → Hiển thị 3-panel layout
```

---

## 5. Auth Component (`components/auth`)

Component tự quản lý luồng xác thực với Phase riêng:

```go
type Phase int

const (
    PhaseCheckingAuth Phase = iota  // Đang gọi GET /me (hiện spinner)
    PhaseLoginForm                  // Hiện form email/password
    PhaseLoggingIn                  // Đang gọi POST /login (hiện spinner)
    PhaseDone                       // Xác thực xong → root chuyển sang StateMain
)
```

**Model fields:**

- `phase Phase` — phase hiện tại của auth flow
- `spinner spinner.Model` — spinner (bubbles) khi đang chờ API
- `emailInput, passInput textinput.Model` — text inputs (bubbles)
- `focusIndex int` — input nào đang được focus
- `errMsg string` — thông báo lỗi
- `token string`, `userID int` — kết quả sau khi auth thành công

**Methods quan trọng:**

- `IsDone() bool` — root model dùng để biết khi nào chuyển sang StateMain
- `Token() string`, `UserID() int` — lấy kết quả auth

**Keyboard trong login form:**
| Phím | Hành động |
|---|---|
| `Tab` / `Shift+Tab` | Chuyển focus giữa Email và Password |
| `Enter` | Submit login |

---

## 6. Services (`services/`)

Các HTTP service function, đều trả về `tea.Cmd` để chạy async trong Bubble Tea runtime:

| Function                        | Method | Endpoint | Kết quả                                       |
| ------------------------------- | ------ | -------- | --------------------------------------------- |
| `CheckAuth(token string)`       | GET    | `/me`    | `AuthCheckSuccessMsg` hoặc `AuthCheckFailMsg` |
| `Login(email, password string)` | POST   | `/login` | `LoginSuccessMsg` hoặc `LoginFailMsg`         |

HTTP client có timeout 10 giây.

---

## 7. Messages (`messages/`)

Message types dùng để communicate giữa services và components qua Bubble Tea runtime:

```go
// Auth check
type AuthCheckSuccessMsg struct { UserID int }
type AuthCheckFailMsg    struct{}

// Login
type LoginSuccessMsg struct { Token string; UserID int }
type LoginFailMsg    struct { Error string }
```

---

## 8. Config (`config/`)

Token persistence — lưu/đọc JWT token ở OS config directory:

- **Path**: `~/.config/kozocom-tui/token` (Linux) / tương đương trên các OS khác
- `LoadToken() (string, error)` — đọc token, trả về `""` nếu file không tồn tại
- `SaveToken(token string) error` — ghi token, tự tạo thư mục nếu chưa có
- **Constant**: `BaseURL = "http://localhost:8000"` — URL của API backend

---

## 9. UI Package (`ui/`)

### `CalculateLayout(w, h int) LayoutDimensions`

Tính toán kích thước các vùng:

- Sidebar: 25% tổng chiều rộng (`SidebarRatio = 0.25`)
- Content (top): 65% tổng chiều cao (`TopContentRatio = 0.65`)
- Footer (bottom): 35% còn lại

### `RenderPanel(title, content, w, h, focused)`

Render panel với:

- Rounded border (`lipgloss.RoundedBorder()`)
- Border màu active (`#FF8899`) hoặc inactive (`#3D3D5C`)
- Legend title hiển thị đè lên border trên (dùng `lipgloss.NewCompositor`)

### `RenderLayout(sidebar, topContent, bottomContent)`

Ghép layout cuối cùng:

- `JoinVertical`: topContent + bottomContent → rightSide
- `JoinHorizontal`: sidebar + rightSide

---

## 10. Design Tokens (màu sắc)

| Token               | Hex       | Dùng cho                      |
| ------------------- | --------- | ----------------------------- |
| `BaseBg`            | `#0F0F1A` | Background cơ sở (dark)       |
| `BorderColor`       | `#3D3D5C` | Border panel không active     |
| `ActiveBorderColor` | `#FF8899` | Border panel đang được focus  |
| `LegendBg`          | `#FF8899` | Nền legend title (hiện tắt)   |
| `LegendFg`          | `#1A1A2E` | Chữ legend title (hiện tắt)   |
| `TitleForeground`   | `#FF8899` | Màu title trong login form    |
| `LabelForeground`   | `#888888` | Màu label (Email:, Password:) |
| `ErrorForeground`   | `#FF4444` | Màu thông báo lỗi             |
| `HintForeground`    | `#555555` | Màu hint text                 |
| `LoginForeground`   | `#3D3D5C` | Màu border login form box     |

---

## 11. Luồng hoạt động đầy đủ (Data Flow)

```
App Start
      │
      ▼
appModel.Init()
  ├── auth.Init() → LoadToken → CheckAuth / PhaseLoginForm
  └── RequestWindowSize → trigger layout

Terminal Input
      │
      ▼
tea.KeyPressMsg / tea.WindowSizeMsg
      │
      ▼
appModel.Update()
  ├── WindowSizeMsg → CalculateLayout → SetSize cho tất cả components
  ├── StateAuth → updateAuth() → delegate xuống auth.Update()
  │       └── auth.IsDone() == true → state = StateMain
  └── StateMain → updateMain()
        ├── Global keys: q / ctrl+c → Quit
        ├── Panel switch: 1, 2, 3 → switchPanel()
        └── Delegate msg xuống sidebar, content, footer
      │
      ▼
appModel.View()
  ├── StateAuth → auth.View() (spinner hoặc login form popup)
  └── StateMain
        ├── CalculateLayout() → lấy kích thước từng vùng
        ├── RenderPanel() × 3 → render từng panel với border và legend title
        └── RenderLayout() → ghép sidebar + contentPanel + footerPanel
      │
      ▼
Terminal Output (AltScreen)
```

---

## 12. Coding Rules cho project này

### Go Style

- **Indentation**: Tab (chuẩn Go `gofmt`)
- **Package naming**: `lowercase`, đặt theo tên folder (`auth`, `sidebar`, `content`, `footer`, `ui`, `services`, `config`, `messages`)
- **Exported vs unexported**: Chỉ export những gì cần thiết ra ngoài package

### Component Rules

- Mỗi component là một **package riêng** trong `components/`
- Component phải implement đủ 3 phương thức Bubble Tea: `Init`, `Update`, `View`
- Component phải có `SetSize`, `Focus`, `Blur` để root model điều phối
- **Không** để business logic ở `app.go` — delegate xuống component tương ứng

### Auth Rules

- Token được persist vào OS config directory qua `config.LoadToken` / `config.SaveToken`
- Auth flow có Phase riêng trong `components/auth` — **không** để auth logic ở `app.go`
- Khi auth xong, root model chỉ check `auth.IsDone()` để chuyển state
- HTTP calls luôn được wrap trong `tea.Cmd` (async) — **không** gọi HTTP synchronously

### Layout Rules

- Mọi tính toán kích thước tập trung trong `ui/layout.go` (`CalculateLayout`)
- Mọi style/màu sắc tập trung trong `ui/styles.go` — **không hardcode màu** ở nơi khác
- Khi thêm panel mới: cập nhật `LayoutDimensions`, `CalculateLayout`, và `RenderLayout`

### Message Rules

- Tất cả message types dùng để communicate giữa services và components đặt trong `messages/`
- Đặt theo domain: `messages/auth.go`, tương lai: `messages/jobs.go`, v.v.

### Bubble Tea Rules

- `Init()` của root model chỉ nên request window size để trigger layout lần đầu
- `Update()` xử lý message theo thứ tự: **WindowSize → state routing → delegate xuống components**
- `View()` phải kiểm tra `m.ready` trước khi render để tránh lỗi khi terminal chưa sẵn sàng
- Dùng `tea.Batch(cmds...)` để gộp commands từ nhiều component

### Thêm tính năng mới

1. **HTTP service mới**: Thêm function vào `services/` (trả về `tea.Cmd`)
2. **Message types mới**: Thêm struct vào `messages/` theo domain
3. **Component mới**: Tạo package trong `components/`, implement `Init/Update/View/SetSize/Focus/Blur`
4. **Panel mới**: Cập nhật `PanelID` enum và `switchPanel()` trong `app.go`, cập nhật `ui/layout.go`
5. **State mới**: Cập nhật `AppState` enum và routing trong `appModel.Update()`
