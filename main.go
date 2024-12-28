package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

// Existing function; sends HTML email by parsing a template file.
func sendMailSimpleHTML(to []string, subject, htmlFile string, cc, bcc []string, attachments []string) error {
	var body bytes.Buffer
	t, err := template.ParseFiles(htmlFile)
	if err != nil {
		return err
	}
	t.Execute(&body, struct{ Name string }{Name: "Pranav"})

	// Build the message with MIME boundaries (supports attachments).
	return sendEmailWithAttachments(to, subject, body.String(), cc, bcc, attachments, true)
}

// This function sends plain text email instead of parsing a template file.
func sendMailPlain(to []string, subject, message string, cc, bcc []string, attachments []string) error {
	return sendEmailWithAttachments(to, subject, message, cc, bcc, attachments, false)
}

// Shared helper function for building & sending MIME messages for either HTML or plain text.
func sendEmailWithAttachments(to []string, subject, content string, cc, bcc []string, attachments []string, isHTML bool) error {
	from := "pranavkharche7@gmail.com" // Replace as needed
	pass := "zkhfypgajolgslhl"         // Gmail app password or environment variable
	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")

	// Combine all recipients
	recipients := append(append([]string{}, to...), cc...)
	recipients = append(recipients, bcc...)

	// Create a multi-part MIME email
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Basic headers
	fmt.Fprintf(&buf, "From: %s\r\n", from)
	fmt.Fprintf(&buf, "To: %s\r\n", strings.Join(to, ","))
	if len(cc) > 0 {
		fmt.Fprintf(&buf, "Cc: %s\r\n", strings.Join(cc, ","))
	}
	fmt.Fprintf(&buf, "Subject: %s\r\n", subject)
	fmt.Fprintf(&buf, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&buf, "Content-Type: multipart/mixed; boundary=%s\r\n\r\n", writer.Boundary())

	// Create the body part
	var contentType string
	if isHTML {
		contentType = "text/html; charset=\"UTF-8\""
	} else {
		contentType = "text/plain; charset=\"UTF-8\""
	}

	bodyHeader := make(textproto.MIMEHeader)
	bodyHeader.Set("Content-Type", contentType)
	bodyPart, _ := writer.CreatePart(bodyHeader)
	bodyPart.Write([]byte(content))

	// Attachments
	for _, filePath := range attachments {
		if filePath == "" {
			continue
		}
		attachFile(writer, filePath)
	}

	writer.Close()

	// Send
	return smtp.SendMail("smtp.gmail.com:587", auth, from, recipients, buf.Bytes())
}

// Attach file data into the MIME writer
func attachFile(w *multipart.Writer, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	filename := filepath.Base(filePath)
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	partHeader.Set("Content-Type", "application/octet-stream")
	part, err := w.CreatePart(partHeader)
	if err != nil {
		return
	}
	contentBytes, _ := io.ReadAll(file)
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(contentBytes)))
	base64.StdEncoding.Encode(encoded, contentBytes)
	part.Write(encoded)
}

// Helper to split comma-separated email lists
func splitEmails(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	var res []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseMultipartForm(10 << 20) // up to 10MB in memory

	mailType := r.FormValue("mailType")
	to := splitEmails(r.FormValue("to"))
	cc := splitEmails(r.FormValue("cc"))
	bcc := splitEmails(r.FormValue("bcc"))
	subject := r.FormValue("subject")
	message := r.FormValue("message")

	// Save any uploaded HTML file
	var htmlFilePath string
	file, header, err := r.FormFile("htmlFile")
	if err == nil && header != nil {
		defer file.Close()
		os.MkdirAll("uploads", 0755)
		htmlFilePath = filepath.Join("uploads", header.Filename)
		out, err := os.Create(htmlFilePath)
		if err == nil {
			defer out.Close()
			io.Copy(out, file)
		}
	}

	// Save attachments
	var attachments []string
	uploadedAttachments := r.MultipartForm.File["attachments"]
	for _, attachment := range uploadedAttachments {
		a, _ := attachment.Open()
		defer a.Close()
		os.MkdirAll("uploads", 0755)
		attachFilePath := filepath.Join("uploads", attachment.Filename)
		out, err := os.Create(attachFilePath)
		if err == nil {
			io.Copy(out, a)
			out.Close()
			attachments = append(attachments, attachFilePath)
		}
	}

	// Send email based on mail type
	if mailType == "html" && htmlFilePath != "" {
		sendMailSimpleHTML(to, subject, htmlFilePath, cc, bcc, attachments)
	} else {
		sendMailPlain(to, subject, message, cc, bcc, attachments)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/send", handleSend)
	fmt.Println("Server starting on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
