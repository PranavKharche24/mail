package mailer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

// Mailer handles email sending operations
type Mailer struct {
	email    string
	password string
	smtpHost string
	smtpPort string
}

// New creates a new Mailer instance
func New() *Mailer {
	return &Mailer{
		smtpHost: "smtp.gmail.com",
		smtpPort: "587",
	}
}

// SetCredentials sets the email credentials
func (m *Mailer) SetCredentials(email, password string) {
	m.email = email
	m.password = password
}

// GetCredentials returns the current credentials
func (m *Mailer) GetCredentials() (string, string) {
	return m.email, m.password
}

// IsConfigured returns true if credentials are set
func (m *Mailer) IsConfigured() bool {
	return m.email != "" && m.password != ""
}

// SendPlain sends a plain text email
func (m *Mailer) SendPlain(to []string, subject, message string, cc, bcc, attachments []string) error {
	return m.sendEmail(to, subject, message, cc, bcc, attachments, false)
}

// SendHTML sends an HTML email from a file
func (m *Mailer) SendHTML(to []string, subject, htmlFile string, cc, bcc, attachments []string) error {
	var body bytes.Buffer
	t, err := template.ParseFiles(htmlFile)
	if err != nil {
		return fmt.Errorf("error parsing HTML file: %v", err)
	}
	if err := t.Execute(&body, struct{ Name string }{Name: "User"}); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}
	return m.sendEmail(to, subject, body.String(), cc, bcc, attachments, true)
}

// SendHTMLContent sends HTML content directly
func (m *Mailer) SendHTMLContent(to []string, subject, htmlContent string, cc, bcc, attachments []string) error {
	return m.sendEmail(to, subject, htmlContent, cc, bcc, attachments, true)
}

func (m *Mailer) sendEmail(to []string, subject, content string, cc, bcc, attachments []string, isHTML bool) error {
	if !m.IsConfigured() {
		return fmt.Errorf("email credentials not configured")
	}

	auth := smtp.PlainAuth("", m.email, m.password, m.smtpHost)

	recipients := append(append([]string{}, to...), cc...)
	recipients = append(recipients, bcc...)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	fmt.Fprintf(&buf, "From: %s\r\n", m.email)
	fmt.Fprintf(&buf, "To: %s\r\n", strings.Join(to, ","))
	if len(cc) > 0 {
		fmt.Fprintf(&buf, "Cc: %s\r\n", strings.Join(cc, ","))
	}
	fmt.Fprintf(&buf, "Subject: %s\r\n", subject)
	fmt.Fprintf(&buf, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&buf, "Content-Type: multipart/mixed; boundary=%s\r\n\r\n", writer.Boundary())

	bodyHeader := make(textproto.MIMEHeader)
	if isHTML {
		bodyHeader.Set("Content-Type", `text/html; charset="UTF-8"`)
	} else {
		bodyHeader.Set("Content-Type", `text/plain; charset="UTF-8"`)
	}
	bodyPart, err := writer.CreatePart(bodyHeader)
	if err != nil {
		return fmt.Errorf("error creating body part: %v", err)
	}
	if _, err := bodyPart.Write([]byte(content)); err != nil {
		return fmt.Errorf("error writing body part: %v", err)
	}

	for _, filePath := range attachments {
		if filePath == "" {
			continue
		}
		if err := m.attachFile(writer, filePath); err != nil {
			return fmt.Errorf("error attaching file: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("error closing writer: %v", err)
	}

	addr := fmt.Sprintf("%s:%s", m.smtpHost, m.smtpPort)
	if err := smtp.SendMail(addr, auth, m.email, recipients, buf.Bytes()); err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	return nil
}

func (m *Mailer) attachFile(w *multipart.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	filename := filepath.Base(filePath)
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	partHeader.Set("Content-Type", "application/octet-stream")
	partHeader.Set("Content-Transfer-Encoding", "base64")
	part, err := w.CreatePart(partHeader)
	if err != nil {
		return fmt.Errorf("error creating part: %v", err)
	}
	contentBytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(contentBytes)))
	base64.StdEncoding.Encode(encoded, contentBytes)
	if _, err := part.Write(encoded); err != nil {
		return fmt.Errorf("error writing part: %v", err)
	}
	return nil
}
