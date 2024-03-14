package senders

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kayprogrammer/socialnet-v6/config"
	"github.com/kayprogrammer/socialnet-v6/models"
	"gopkg.in/gomail.v2"
)

func sortEmail(user *models.User, emailType string, code *uint32) map[string]interface{} {
	templateFile := "templates/welcome.html"
	subject := "Account verified"
	data := make(map[string]interface{})
	data["template_file"] = templateFile
	data["subject"] = subject

	// Sort different templates and subject for respective email types
	if emailType == "activate" {
		templateFile = "templates/email-activation.html"
		subject = "Activate your account"
		data["template_file"] = templateFile
		data["subject"] = subject
		data["otp"] = code

	} else if emailType == "reset" {
		templateFile = "templates/password-reset.html"
		subject = "Reset your password"
		data["template_file"] = templateFile
		data["subject"] = subject
		data["otp"] = code

	} else if emailType == "reset-success" {
		templateFile = "templates/password-reset-success.html"
		subject = "Password reset successfully"
		data["template_file"] = templateFile
		data["subject"] = subject
	}
	return data
}

type EmailContext struct {
	Name string
	Otp  *uint32
}

func SendEmail(user *models.User, emailType string, code *uint32) {
	if os.Getenv("ENVIRONMENT") == "TESTING" {
		return
	}
	cfg := config.GetConfig()

	emailData := sortEmail(user, emailType, code)
	templateFile := emailData["template_file"]
	subject := emailData["subject"]

	// Create a context with dynamic data
	data := EmailContext{
		Name: user.FirstName,
	}
	if otp, ok := emailData["otp"]; ok {
		code := otp.(*uint32)
		data.Otp = code
	}

	// Read the HTML file content
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		log.Println("Unable to identify current directory (needed to load templates)", os.Stderr)
		os.Exit(1)
	}
	basepath := filepath.Dir(file)
	tempfile := fmt.Sprintf("../%s", templateFile.(string))
	htmlContent, err := os.ReadFile(filepath.Join(basepath, tempfile))
	if err != nil {
		log.Fatal("Error reading HTML file:", err)
	}

	// Create a new template from the HTML file content
	tmpl, err := template.New("email_template").Parse(string(htmlContent))
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	// Execute the template with the context and set it as the body of the email
	var bodyContent bytes.Buffer
	if err := tmpl.Execute(&bodyContent, data); err != nil {
		log.Fatal("Error executing template:", err)
	}

	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.MailSenderEmail)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", subject.(string))
	m.SetBody("text/html", bodyContent.String())

	// Create a new SMTP client
	d := gomail.NewDialer(cfg.MailSenderHost, cfg.MailSenderPort, cfg.MailSenderEmail, cfg.MailSenderPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Error sending email:", err)
	}
}
