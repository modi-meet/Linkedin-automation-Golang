# LinkedIn Automation (Go + Rod)

Automates a basic LinkedIn workflow using Go and the Rod browser automation library:

1. Launches a Chromium instance (non-headless by default)
2. Restores a logged-in session using cookies (if available)
3. If cookies are missing/expired, logs in using credentials from `.env`
4. Searches for people by keyword
5. Visits each profile and attempts to send a connection request with a note

> **Important**: Automated interaction with LinkedIn may violate LinkedIn’s Terms of Service and can lead to account restrictions. Use responsibly, on your own account, and at your own risk.

Explanation video - https://drive.google.com/file/d/1xocEAjwKOh64OasL7DALhuCSl_eFtA7u/view?usp=sharing

---

## Prerequisites

- Go 1.20+ (recommended)
- A LinkedIn account
- Internet access

Rod will download/launch a browser (Chromium) via the launcher. If your environment restricts browser downloads/launch, you may need to configure your system accordingly.

---

## Project structure

```
.
├── actions/
│   └── connect.go         # Sends connection requests
├── api/
│   └── server.go          # HTTP API server
├── auth/
│   └── auth.go            # Login + session handling
├── pkg/
│   ├── logger/            # Logging with SSE broadcast
│   └── workflow/          # Main automation workflow
├── search/
│   └── search.go          # Search for profile URLs
├── utils/
│   └── mouse.go           # Stealth techniques (8 methods)
├── main.go                # Entry point
└── go.mod
```

<img width="4225" height="2230" alt="Go-linkedIn_Autmation-workflow" src="https://github.com/user-attachments/assets/5e7e1cf4-5983-4183-8cf8-b85e050cb54c" />


---

## Setup

### 1) Install dependencies

```powershell
go mod tidy
```

### 2) Create a `.env`

Create a file named `.env` in the project root:

```env
LINKEDIN_EMAIL=your_email@example.com
LINKEDIN_PASSWORD=your_password
```

**Notes**
- Use an account you control.
- Avoid committing `.env` to source control.

### 3) Cookies file

The app persists session cookies to:

- `linkedin_cookies.json`

On the next run, it will attempt to restore the session from this file.

---

## Running

### Start the server

```powershell
go run .
```

### Trigger automation via API

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/start" `
  -Method POST -ContentType "application/json" `
  -Body '{"keyword":"Go Developer","limit":3}'
```

### API Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `keyword` | string | Search keyword for profiles |
| `limit` | int | Max profiles to process |
| `email` | string | (Optional) Override .env email |
| `password` | string | (Optional) Override .env password |
| `connectMessage` | string | Custom connection note |
| `headless` | bool | Run browser headless |

---

## Configuration

### Search keyword

In `main.go` the keyword is currently set in code:

```go
keyword := "Golang Developer"
profiles := search.Run(page, keyword, 3)
```

- The third argument is the **limit** (used for testing).
- Increase the limit once you’re confident it works.

### Connection message

In `main.go`, the connection note is also set in code:

```go
msg := "Hi, I am a Go developer expanding my network. Would love to connect!"
```

---

## How it works (high level)

### Authentication (cookies + login)

- `auth.LoadCookies(browser, cookieFile)` attempts to load cookies from disk and set them on the Rod browser.
- If cookie load succeeds, `main.go` navigates to the feed and checks whether LinkedIn redirects to `/login`.
- If redirected (or cookie load fails), it uses `auth.Login(page)`.
- After a successful login it persists cookies using `auth.SaveCookies(browser, cookieFile)`.

### Search

- `search.Run(page, keyword, limit)` navigates to LinkedIn’s people search results and extracts profile URLs.

### Send connection request

- `actions.SendConnectionRequest(page, profileURL, message)`:
  - Opens the profile
  - Tries to find a visible **Connect** button
  - If not found, opens **More actions** and tries to click **Connect** from the dropdown
  - If the **Add a note** dialog is available, inputs the provided message and sends

### Stealth Techniques (Anti-Detection)

The automation implements 8 stealth techniques to avoid detection:

1. **Bézier Mouse Movement** - Natural curves with overshoot and micro-corrections
2. **Randomized Timing** - Variable delays with "thinking" pauses
3. **Browser Fingerprint Masking** - Hides `navigator.webdriver`, spoofs plugins/WebGL
4. **Random Scrolling** - Variable speeds with occasional scroll-back
5. **Realistic Typing** - 3% typo rate with corrections, variable keystroke timing
6. **Mouse Hovering** - Idle cursor wandering and element hover
7. **Persistent Browser Profile** - Session persists at `~/.linkedin-automation-profile`
8. **Stealth Chrome Flags** - Disables automation indicators

---

## Troubleshooting

### Build errors

Run:

```powershell
go build ./...
```

If you see type errors around cookies, they usually come from Rod `proto` type changes. This repository currently targets Rod `v0.116.2`.

### “Cookies loaded” but still redirected to login

Common causes:

- Cookies expired / invalidated
- LinkedIn requires additional verification (CAPTCHA/2FA)
- Cookies were saved from a different browser context

The code already detects redirect to `/login` and re-runs login.

### Connect button not found

LinkedIn UI varies by account, region, A/B tests, and relationship state. Reasons include:

- The profile only shows **Follow**
- You’ve already sent a request
- You’re out of connection requests / rate limited

### It’s too fast / gets flagged

Increase sleeps in:

- `utils.RandomSleep(...)`

And reduce the number of profiles per run.

---

## Safety / responsible use

- Keep limits low and ramp up gradually.
- Avoid running this continuously or aggressively.
- Consider adding strong backoff, daily caps, and manual review steps if you extend this project.

---



