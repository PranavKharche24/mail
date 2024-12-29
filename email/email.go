package email

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

var fromEmail = "youremail@gmail.com"
var fromPass = "yourapppassword"

func SendMailSimpleHTML(to []string, subject, htmlFile string, cc, bcc []string, attachments []string) error {
    var body bytes.Buffer
    t, err := template.ParseFiles(htmlFile)
    if (err != nil) {
        return fmt.Errorf("error parsing HTML file: %v", err)
    }
    if err := t.Execute(&body, struct{ Name string }{Name: "Pranav"}); err != nil {
        return fmt.Errorf("error executing template: %v", err)
    }
    return sendEmailWithAttachments(to, subject, body.String(), cc, bcc, attachments, true)
}

func SendMailPlain(to []string, subject, msg string, cc, bcc []string, attachments []string) error {
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