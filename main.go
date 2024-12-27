package main

import (
	// "bufio"
	// "fmt"
	"bytes"
	"html/template"
	"net/smtp"
	// "os"
	// "strings"
)

// func sendMailSimple(recipient, subject, body string) error {
// 	auth := smtp.PlainAuth(
// 		"",
// 		"pranavkharche7@gmail.com",
// 		"zkhfypgajolgslhl",
// 		"smtp.gmail.com",
// 	)

// 	msg := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

//		return smtp.SendMail(
//			"smtp.gmail.com:587",
//			auth,
//			"pranavkharche7@gmail.com",
//			[]string{recipient},
//			[]byte(msg),
//		)
//	}
func sendMailSimpleHTML(recipient, subject, html string) error {
	var body bytes.Buffer
	t, err := template.ParseFiles(html)
	if err != nil {
		return err
	}
	t.Execute(&body, struct{ Name string }{Name: "Pranav"})
	auth := smtp.PlainAuth(
		"",
		"pranavkharche7@gmail.com",
		"zkhfypgajolgslhl",
		"smtp.gmail.com",
	)
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	msg := "Subject: " + subject + "\n" + headers + "\n\n" + body.String()

	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"pranavkharche7@gmail.com",
		[]string{recipient},
		[]byte(msg),
	)
}
func main() {
	// reader := bufio.NewReader(os.Stdin)

	// fmt.Print("Recipient: ")
	// recipient, _ := reader.ReadString('\n')
	// recipient = strings.TrimSpace(recipient)

	// fmt.Print("Subject: ")
	// subject, _ := reader.ReadString('\n')
	// subject = strings.TrimSpace(subject)

	// fmt.Print("Body: ")
	// body, _ := reader.ReadString('\n')
	// body = strings.TrimSpace(body)

	// if err := sendMailSimple(recipient, subject, body); err != nil {
	// 	fmt.Println("Failed to send:", err)
	// } else {
	// 	fmt.Println("Mail sent successfully")
	// }
	sendMailSimpleHTML(
		"2022bcs020@sggs.ac.in",
		"Another",
		"./test.html",
	)
}
