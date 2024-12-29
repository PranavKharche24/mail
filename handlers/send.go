package handlers

import (
	"net/http"

	"github.com/Pranavkharche24/mail/mail"
	"github.com/Pranavkharche24/mail/utils"
)

func HandleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	mailType := r.FormValue("mailType")
	to := utils.SplitEmails(r.FormValue("to"))
	cc := utils.SplitEmails(r.FormValue("cc"))
	bcc := utils.SplitEmails(r.FormValue("bcc"))
	subject := r.FormValue("subject")
	message := r.FormValue("message")

	htmlFilePath, err := utils.SaveUploadedFile(r, "htmlFile")
	if err != nil {
		http.Error(w, "Error saving HTML file", http.StatusInternalServerError)
		return
	}

	attachments, err := utils.SaveUploadedFiles(r, "attachments")
	if err != nil {
		http.Error(w, "Error saving attachments", http.StatusInternalServerError)
		return
	}

	if mailType == "html" && htmlFilePath != "" {
		if err := mail.SendMailSimpleHTML(to, subject, htmlFilePath, cc, bcc, attachments, fromEmail, fromPass); err != nil {
			http.Error(w, "Error sending email", http.StatusInternalServerError)
			return
		}
	} else {
		if err := mail.SendMailPlain(to, subject, message, cc, bcc, attachments, fromEmail, fromPass); err != nil {
			http.Error(w, "Error sending email", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
