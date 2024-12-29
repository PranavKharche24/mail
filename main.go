package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

// Global admin-configurable values
var fromEmail = "youremail@gmail.com"
var fromPass = "yourapppassword"

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/send", handleSend)
	http.HandleFunc("/admin", handleAdmin)
	http.HandleFunc("/admin/save", handleAdminSave)

	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/admin.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	data := struct {
		FromEmail string
		FromPass  string
	}{
		FromEmail: fromEmail,
		FromPass:  fromPass,
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func handleAdminSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	fromEmail = r.FormValue("fromEmail")
	fromPass = r.FormValue("fromPass")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	mailType := r.FormValue("mailType")
	to := splitEmails(r.FormValue("to"))
	cc := splitEmails(r.FormValue("cc"))
	bcc := splitEmails(r.FormValue("bcc"))
	subject := r.FormValue("subject")
	message := r.FormValue("message")

	htmlFilePath, err := saveUploadedFile(r, "htmlFile")
	if err != nil {
		http.Error(w, "Error saving HTML file", http.StatusInternalServerError)
		return
	}

	attachments, err := saveUploadedFiles(r, "attachments")
	if err != nil {
		http.Error(w, "Error saving attachments", http.StatusInternalServerError)
		return
	}

	if mailType == "html" && htmlFilePath != "" {
		if err := sendMailSimpleHTML(to, subject, htmlFilePath, cc, bcc, attachments); err != nil {
			http.Error(w, "Error sending email", http.StatusInternalServerError)
			return
		}
	} else {
		if err := sendMailPlain(to, subject, message, cc, bcc, attachments); err != nil {
			http.Error(w, "Error sending email", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func sendMailSimpleHTML(to []string, subject, htmlFile string, cc, bcc []string, attachments []string) error {
	var body bytes.Buffer
	t, err := template.ParseFiles(htmlFile)
	if err != nil {
		return fmt.Errorf("error parsing HTML file: %v", err)
	}
	if err := t.Execute(&body, struct{ Name string }{Name: "Pranav"}); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}
	return sendEmailWithAttachments(to, subject, body.String(), cc, bcc, attachments, true)
}

func sendMailPlain(to []string, subject, msg string, cc, bcc []string, attachments []string) error {
	return sendEmailWithAttachments(to, subject, msg, cc, bcc, attachments, false)
}

func sendEmailWithAttachments(to []string, subject, content string, cc, bcc []string, attachments []string, isHTML bool) error {
	auth := smtp.PlainAuth("", fromEmail, fromPass, "smtp.gmail.com")

	recipients := append(append([]string{}, to...), cc...)
	recipients = append(recipients, bcc...)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	fmt.Fprintf(&buf, "From: %s\r\n", fromEmail)
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
		if err := attachFile(writer, filePath); err != nil {
			return fmt.Errorf("error attaching file: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("error closing writer: %v", err)
	}

	if err := smtp.SendMail("smtp.gmail.com:587", auth, fromEmail, recipients, buf.Bytes()); err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	return nil
}

func attachFile(w *multipart.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	filename := filepath.Base(filePath)
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	partHeader.Set("Content-Type", "application/octet-stream")
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

func saveUploadedFile(r *http.Request, formFileName string) (string, error) {
	file, header, err := r.FormFile(formFileName)
	if err != nil {
		return "", nil // No file uploaded
	}
	defer file.Close()

	os.MkdirAll("uploads", 0755)
	filePath := filepath.Join("uploads", header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("error saving file: %v", err)
	}
	return filePath, nil
}

func saveUploadedFiles(r *http.Request, formFileName string) ([]string, error) {
	var filePaths []string
	files := r.MultipartForm.File[formFileName]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("error opening file: %v", err)
		}
		defer file.Close()

		os.MkdirAll("uploads", 0755)
		filePath := filepath.Join("uploads", fileHeader.Filename)
		out, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("error creating file: %v", err)
		}
		defer out.Close()
		if _, err := io.Copy(out, file); err != nil {
			return nil, fmt.Errorf("error saving file: %v", err)
		}
		filePaths = append(filePaths, filePath)
	}
	return filePaths, nil
}
