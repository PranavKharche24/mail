package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pranavKharche24/mail/mailer"
)

// Terminal colors
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

// CLI handles the command line interface
type CLI struct {
	reader *bufio.Reader
	mailer *mailer.Mailer
}

// New creates a new CLI instance
func New() *CLI {
	return &CLI{
		reader: bufio.NewReader(os.Stdin),
		mailer: mailer.New(),
	}
}

// SetMailer sets an existing mailer instance
func (c *CLI) SetMailer(m *mailer.Mailer) {
	c.mailer = m
}

// Run starts the interactive CLI
func (c *CLI) Run() {
	for {
		c.printMenu()
		choice := c.prompt("Select option")

		switch choice {
		case "1":
			c.sendPlainEmail()
		case "2":
			c.sendHTMLEmail()
		case "3":
			c.configureCredentials()
		case "4":
			c.showCredentials()
		case "5":
			c.listTemplates()
		case "6", "q", "quit", "exit":
			fmt.Println("\n  Goodbye.\n")
			return
		default:
			c.showError("Invalid option")
		}
	}
}

func (c *CLI) printMenu() {
	fmt.Println()
	fmt.Printf("  %s%sMAIN MENU%s\n", Bold, Cyan, Reset)
	fmt.Println("  " + strings.Repeat("-", 40))
	fmt.Println()
	fmt.Printf("  %s[1]%s  Send Plain Text Email\n", Green, Reset)
	fmt.Printf("  %s[2]%s  Send HTML Email\n", Green, Reset)
	fmt.Printf("  %s[3]%s  Configure Credentials\n", Yellow, Reset)
	fmt.Printf("  %s[4]%s  Show Current Credentials\n", Yellow, Reset)
	fmt.Printf("  %s[5]%s  List Available Templates\n", Blue, Reset)
	fmt.Printf("  %s[6]%s  Exit\n", Dim, Reset)
	fmt.Println()
}

func (c *CLI) prompt(label string) string {
	fmt.Printf("  %s> %s:%s ", Bold, label, Reset)
	input, _ := c.reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func (c *CLI) promptMultiline(label string) string {
	fmt.Printf("  %s> %s (end with empty line):%s\n", Bold, label, Reset)
	var lines []string
	for {
		fmt.Print("    ")
		line, _ := c.reader.ReadString('\n')
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (c *CLI) showSuccess(msg string) {
	fmt.Printf("\n  %s[OK]%s %s\n", Green, Reset, msg)
}

func (c *CLI) showError(msg string) {
	fmt.Printf("\n  %s[ERROR]%s %s\n", Red, Reset, msg)
}

func (c *CLI) showInfo(msg string) {
	fmt.Printf("  %s[INFO]%s %s\n", Cyan, Reset, msg)
}

func (c *CLI) sendPlainEmail() {
	fmt.Println()
	fmt.Printf("  %s%sSEND PLAIN TEXT EMAIL%s\n", Bold, Green, Reset)
	fmt.Println("  " + strings.Repeat("-", 40))
	fmt.Println()

	if !c.mailer.IsConfigured() {
		c.showError("Credentials not configured. Use option [3] first.")
		return
	}

	to := c.prompt("To (comma-separated)")
	if to == "" {
		c.showError("Recipient is required")
		return
	}

	cc := c.prompt("CC (optional)")
	bcc := c.prompt("BCC (optional)")
	subject := c.prompt("Subject")
	message := c.promptMultiline("Message")
	attachments := c.prompt("Attachments (paths, optional)")

	var attachmentList []string
	if attachments != "" {
		for _, a := range strings.Split(attachments, ",") {
			attachmentList = append(attachmentList, strings.TrimSpace(a))
		}
	}

	c.showInfo("Sending...")

	err := c.mailer.SendPlain(
		splitEmails(to),
		subject,
		message,
		splitEmails(cc),
		splitEmails(bcc),
		attachmentList,
	)

	if err != nil {
		c.showError(fmt.Sprintf("Send failed: %v", err))
		return
	}

	c.showSuccess("Email sent successfully")
}

func (c *CLI) sendHTMLEmail() {
	fmt.Println()
	fmt.Printf("  %s%sSEND HTML EMAIL%s\n", Bold, Green, Reset)
	fmt.Println("  " + strings.Repeat("-", 40))
	fmt.Println()

	if !c.mailer.IsConfigured() {
		c.showError("Credentials not configured. Use option [3] first.")
		return
	}

	c.listTemplates()

	to := c.prompt("To (comma-separated)")
	if to == "" {
		c.showError("Recipient is required")
		return
	}

	cc := c.prompt("CC (optional)")
	bcc := c.prompt("BCC (optional)")
	subject := c.prompt("Subject")
	htmlFile := c.prompt("HTML file path")

	if _, err := os.Stat(htmlFile); os.IsNotExist(err) {
		c.showError(fmt.Sprintf("File not found: %s", htmlFile))
		return
	}

	attachments := c.prompt("Attachments (paths, optional)")

	var attachmentList []string
	if attachments != "" {
		for _, a := range strings.Split(attachments, ",") {
			attachmentList = append(attachmentList, strings.TrimSpace(a))
		}
	}

	c.showInfo("Sending...")

	err := c.mailer.SendHTML(
		splitEmails(to),
		subject,
		htmlFile,
		splitEmails(cc),
		splitEmails(bcc),
		attachmentList,
	)

	if err != nil {
		c.showError(fmt.Sprintf("Send failed: %v", err))
		return
	}

	c.showSuccess("HTML email sent successfully")
}

func (c *CLI) configureCredentials() {
	fmt.Println()
	fmt.Printf("  %s%sCONFIGURE CREDENTIALS%s\n", Bold, Yellow, Reset)
	fmt.Println("  " + strings.Repeat("-", 40))
	fmt.Println()

	email := c.prompt("Gmail address")
	password := c.prompt("App password")

	if email == "" || password == "" {
		c.showError("Both email and password are required")
		return
	}

	c.mailer.SetCredentials(email, password)

	envContent := fmt.Sprintf("# Gomail Configuration\nEMAIL_FROM=%s\nEMAIL_PASSWORD=%s\nPORT=8080\n", email, password)

	err := os.WriteFile(".env", []byte(envContent), 0600)
	if err != nil {
		c.showError(fmt.Sprintf("Failed to save: %v", err))
		return
	}

	c.showSuccess("Credentials saved to .env")
}

func (c *CLI) showCredentials() {
	fmt.Println()
	fmt.Printf("  %s%sCURRENT CREDENTIALS%s\n", Bold, Yellow, Reset)
	fmt.Println("  " + strings.Repeat("-", 40))
	fmt.Println()

	email, _ := c.mailer.GetCredentials()

	if email == "" {
		c.showError("No credentials configured")
		return
	}

	fmt.Printf("  Email:    %s\n", email)
	fmt.Printf("  Password: ********\n")
	fmt.Printf("  Status:   %sConfigured%s\n", Green, Reset)
}

func (c *CLI) listTemplates() {
	fmt.Println()
	fmt.Printf("  %s%sAVAILABLE TEMPLATES%s\n", Bold, Blue, Reset)
	fmt.Println("  " + strings.Repeat("-", 40))

	found := false

	templates, _ := filepath.Glob("templates/*.html")
	if len(templates) > 0 {
		fmt.Println("\n  templates/")
		for _, t := range templates {
			fmt.Printf("    - %s\n", filepath.Base(t))
		}
		found = true
	}

	uploads, _ := filepath.Glob("uploads/*.html")
	if len(uploads) > 0 {
		fmt.Println("\n  uploads/")
		for _, u := range uploads {
			fmt.Printf("    - %s\n", filepath.Base(u))
		}
		found = true
	}

	current, _ := filepath.Glob("*.html")
	if len(current) > 0 {
		fmt.Println("\n  ./")
		for _, c := range current {
			fmt.Printf("    - %s\n", c)
		}
		found = true
	}

	if !found {
		c.showInfo("No HTML templates found")
	}
	fmt.Println()
}

func splitEmails(s string) []string {
	if s == "" {
		return []string{}
	}
	var result []string
	for _, p := range strings.Split(s, ",") {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
