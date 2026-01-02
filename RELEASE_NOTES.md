# Gomail v1.0.0

**Professional Email Utility for Linux, macOS, and Windows**

---

## Release Title
**Gomail v1.0.0 - Initial Release**

---

## What's New

### Features

- **Dual Interface** - Run both CLI and Web interface simultaneously
- **Gmail SMTP Integration** - Send emails securely via Gmail with App Password authentication
- **HTML Email Support** - Send rich HTML formatted emails
- **File Attachments** - Attach files up to 25MB
- **CC/BCC Support** - Send to multiple recipients
- **Secure Configuration** - Credentials stored in local .env file (never committed to Git)
- **Zero Dependencies** - Single static binary, no external libraries required

### Interfaces

**Command Line Interface (CLI)**
- Interactive menu-driven interface
- Colored output for better readability
- Quick send and credential management

**Web Interface**
- Clean, professional design
- Real-time form validation
- Settings panel for credential management
- Responsive layout

---

## Installation

### Debian/Ubuntu (.deb package)
```bash
sudo dpkg -i gomail_1.0.0_amd64.deb
```

### Linux (binary)
```bash
tar -xzf gomail-1.0.0-linux-amd64.tar.gz
cd gomail-1.0.0-linux-amd64
./gomail
```

### macOS
```bash
tar -xzf gomail-1.0.0-darwin-amd64.tar.gz
cd gomail-1.0.0-darwin-amd64
./gomail
```

### Windows
Extract the zip file and run `gomail.exe`

---

## Configuration

Create a `.env` file or edit `~/.config/gomail/.env`:

```env
EMAIL_FROM=your-email@gmail.com
EMAIL_PASSWORD=your-app-password
PORT=8080
```

**Get Gmail App Password:** https://myaccount.google.com/apppasswords

---

## Usage

```bash
gomail          # Start both CLI and Web interface
gomail cli      # CLI mode only
gomail web      # Web mode only
gomail help     # Show help
gomail version  # Show version
```

---

## Downloads

| Platform | Architecture | File |
|----------|-------------|------|
| Debian/Ubuntu | x64 | `gomail_1.0.0_amd64.deb` |
| Linux | x64 | `gomail-1.0.0-linux-amd64.tar.gz` |
| Linux | ARM64 | `gomail-1.0.0-linux-arm64.tar.gz` |
| macOS | x64 (Intel) | `gomail-1.0.0-darwin-amd64.tar.gz` |
| macOS | ARM64 (Apple Silicon) | `gomail-1.0.0-darwin-arm64.tar.gz` |
| Windows | x64 | `gomail-1.0.0-windows-amd64.zip` |

---

## Checksums

Verify your download with SHA256:
```bash
sha256sum -c checksums.txt
```

---

## Requirements

- Gmail account with App Password enabled
- No external dependencies required

---

## License

MIT License

---

## Links

- **Repository:** https://github.com/pranavKharche24/mail
- **Issues:** https://github.com/pranavKharche24/mail/issues
