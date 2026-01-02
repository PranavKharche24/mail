# Gomail

A professional email utility for Go with CLI and Web interfaces. Send plain text or HTML emails with attachments via Gmail SMTP.

## Features

- **Dual Interface** - CLI and Web server run simultaneously
- **Plain Text & HTML** - Send both types of emails
- **Attachments** - Multiple file attachments support
- **CC/BCC** - Full recipient management
- **Secure** - Credentials stored in `.env` file (gitignored)
- **Zero Dependencies** - Pure Go standard library

## Installation

```bash
git clone https://github.com/pranavKharche24/mail.git
cd mail
go build -o gomail
```

## Configuration

```bash
cp .env.example .env
```

Edit `.env`:
```
EMAIL_FROM=your-email@gmail.com
EMAIL_PASSWORD=your-app-password
PORT=8080
```

### Gmail App Password

1. Enable 2-Step Verification at [Google Account](https://myaccount.google.com/security)
2. Generate App Password at [App Passwords](https://myaccount.google.com/apppasswords)
3. Use the 16-character password in `.env`

## Usage

```bash
# Start both Web and CLI (default)
./gomail

# CLI only
./gomail cli

# Web only
./gomail web

# Help
./gomail help
```

### Web Interface

- Email Form: `http://localhost:8080`
- Admin Panel: `http://localhost:8080/admin`

### CLI Interface

```
MAIN MENU
------------------------------------------

[1]  Send Plain Text Email
[2]  Send HTML Email
[3]  Configure Credentials
[4]  Show Current Credentials
[5]  List Available Templates
[6]  Exit
```

## Project Structure

```
mail/
├── main.go           # Entry point
├── cli/
│   └── cli.go        # CLI interface
├── web/
│   └── server.go     # Web server
├── mailer/
│   └── mailer.go     # Email logic
├── config/
│   └── config.go     # Configuration
├── templates/
│   ├── index.html    # Email form
│   └── admin.html    # Settings page
├── uploads/          # Uploaded files
├── .env.example      # Config template
├── .gitignore
├── go.mod
└── README.md
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `EMAIL_FROM` | Gmail address | Yes |
| `EMAIL_PASSWORD` | App password | Yes |
| `PORT` | Web server port | No (default: 8080) |

## Security

- Credentials loaded from environment variables or `.env` file
- `.env` is excluded from version control via `.gitignore`
- File permissions set to 0600 for credential storage
- No secrets in source code

## License

MIT License
