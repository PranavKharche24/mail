package web

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pranavKharche24/mail/mailer"
)

// Server handles the web interface
type Server struct {
	mailer *mailer.Mailer
	port   string
}

// New creates a new web server
func New(port string, m *mailer.Mailer) *Server {
	return &Server{
		mailer: m,
		port:   port,
	}
}

// Start starts the web server
func (s *Server) Start() error {
	http.HandleFunc("/", s.handleHome)
	http.HandleFunc("/send", s.handleSend)
	http.HandleFunc("/admin", s.handleAdmin)
	http.HandleFunc("/admin/save", s.handleAdminSave)
	http.HandleFunc("/api/status", s.handleAPIStatus)

	addr := ":" + s.port
	log.Printf("Web server listening on http://localhost%s", addr)

	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	email, _ := s.mailer.GetCredentials()
	data := struct {
		IsConfigured bool
		FromEmail    string
	}{
		IsConfigured: s.mailer.IsConfigured(),
		FromEmail:    email,
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}
	tmpl.Execute(w, data)
}

func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	email, password := s.mailer.GetCredentials()
	data := struct {
		FromEmail    string
		FromPass     string
		IsConfigured bool
	}{
		FromEmail:    email,
		FromPass:     password,
		IsConfigured: s.mailer.IsConfigured(),
	}

	tmpl, err := template.ParseFiles("templates/admin.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func (s *Server) handleAdminSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form error", http.StatusBadRequest)
		return
	}

	email := r.FormValue("fromEmail")
	password := r.FormValue("fromPass")

	s.mailer.SetCredentials(email, password)

	// Save to .env
	envContent := fmt.Sprintf("# Gomail Configuration\nEMAIL_FROM=%s\nEMAIL_PASSWORD=%s\nPORT=%s\n", email, password, s.port)
	os.WriteFile(".env", []byte(envContent), 0600)

	http.Redirect(w, r, "/?saved=true", http.StatusSeeOther)
}

func (s *Server) handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if !s.mailer.IsConfigured() {
		http.Redirect(w, r, "/admin?error=credentials", http.StatusSeeOther)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Form error", http.StatusBadRequest)
		return
	}

	mailType := r.FormValue("mailType")
	to := splitEmails(r.FormValue("to"))
	cc := splitEmails(r.FormValue("cc"))
	bcc := splitEmails(r.FormValue("bcc"))
	subject := r.FormValue("subject")
	message := r.FormValue("message")

	htmlFilePath, _ := s.saveUploadedFile(r, "htmlFile")
	attachments, _ := s.saveUploadedFiles(r, "attachments")

	var sendErr error
	if mailType == "html" && htmlFilePath != "" {
		sendErr = s.mailer.SendHTML(to, subject, htmlFilePath, cc, bcc, attachments)
	} else {
		sendErr = s.mailer.SendPlain(to, subject, message, cc, bcc, attachments)
	}

	if sendErr != nil {
		log.Printf("Send error: %v", sendErr)
		http.Redirect(w, r, "/?error=send", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/?success=true", http.StatusSeeOther)
}

func (s *Server) handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if s.mailer.IsConfigured() {
		w.Write([]byte(`{"status":"configured","ready":true}`))
	} else {
		w.Write([]byte(`{"status":"not_configured","ready":false}`))
	}
}

func (s *Server) saveUploadedFile(r *http.Request, name string) (string, error) {
	file, header, err := r.FormFile(name)
	if err != nil {
		return "", nil
	}
	defer file.Close()

	os.MkdirAll("uploads", 0755)
	path := filepath.Join("uploads", header.Filename)
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()
	io.Copy(out, file)
	return path, nil
}

func (s *Server) saveUploadedFiles(r *http.Request, name string) ([]string, error) {
	var paths []string
	files := r.MultipartForm.File[name]
	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		os.MkdirAll("uploads", 0755)
		path := filepath.Join("uploads", fh.Filename)
		out, err := os.Create(path)
		if err != nil {
			continue
		}
		defer out.Close()
		io.Copy(out, file)
		paths = append(paths, path)
	}
	return paths, nil
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
