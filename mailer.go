package mailgo

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/smtp"
	"strings"
	"time"

	gosasl "github.com/emersion/go-sasl"
	gosmtp "github.com/emersion/go-smtp"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var mailer *Mailer

// ErrNilMailer specifies nil mailer instance error
var ErrNilMailer = errors.New("no mailer instance")
var ErrNoFrom = errors.New("from option missing")

// Mailer struct definition.
type Mailer struct {
	Host       string
	Port       string
	Name       string
	From       string
	Creds      gosasl.Client
	Localizer  *i18n.Localizer
	ReturnPath string
}

func init() {
	mailer, _ = NewMailer()
}

// SetupMailer set ups package mailer with specified parameters.
func SetupMailer(opts ...MailerOption) (err error) {
	mailer, err = NewMailer(opts...)
	return err
}

// NewMailer creates new mailer instance. If we can't create new mailer than
// error is returned.
func NewMailer(opts ...MailerOption) (*Mailer, error) {
	var options = &MailerOptions{
		Name:  "MailGo",
		Port:  "25",
		Creds: nil,
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.From == "" {
		return nil, ErrNoFrom
	}

	if options.Host == "" {
		options.Host = options.From[strings.Index(options.From, "@"):]
	}

	mlr := &Mailer{
		Host:       options.Host,
		Port:       options.Port,
		Name:       options.Name,
		From:       options.From,
		Creds:      options.Creds,
		ReturnPath: options.RetPath,
	}

	return mlr, nil
}

// SendMail sends mail to given recipients with specified subject and plaintext
// and HTML message (that's why function requires both). If plaintext of html
// message is empty, then that format is not used.
func SendMail(recipients []string, subject, plain, html string) error {
	return mailer.SendMail(recipients, subject, plain, html)
}

// SendMail sends mail to given recipients with specified subject and plaintext
// and HTML message (that's why function requires both). If plaintext of html
// message is empty, then that format is not used.
func (m *Mailer) SendMail(recipients []string, subject, plain, html string) (err error) {
	var boundary string
	if boundary, err = generateRandomString(16); err != nil {
		return err
	}

	var body string
	body += m.getGeneralHeader(subject, boundary, recipients)
	body += "\r\n--" + boundary + "\r\n"
	if len(plain) != 0 {
		body += getPlainTextHeader()
		body += "\r\n" + lineSplit(base64.StdEncoding.EncodeToString([]byte(plain))) + "\r\n"
		body += "\r\n--" + boundary + "\r\n"
	}
	if len(html) != 0 {
		body += getHTMLHeader()
		body += "\r\n" + lineSplit(base64.StdEncoding.EncodeToString([]byte(html))) + "\r\n"
		body += "\r\n--" + boundary + "--\r\n"
	}

	if m.Creds == nil {
		for _, recipient := range recipients {
			m.sendMailWithoutAuth(recipient, body)
		}
	}
	return m.sendMail(recipients, body)
}

// sendMail sends e-mail via smtp-go with authentication.
func (m *Mailer) sendMail(recipients []string, body string) error {
	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	return gosmtp.SendMail(addr, m.Creds, m.From, recipients, strings.NewReader(body))
}

// sendMail sends e-mail via standard go smtp library without authentication.
func (m *Mailer) sendMailWithoutAuth(recipient string, body string) error {
	var err error
	var conn *smtp.Client
	var wc io.WriteCloser

	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)

	// Connect to the remote SMTP server.
	if conn, err = smtp.Dial(addr); err != nil {
		return err
	}

	// Set the sender
	if err = conn.Mail(m.From); err != nil {
		return err
	}

	// Set the recipient
	if err = conn.Rcpt(recipient); err != nil {
		return err
	}

	// Send the email body
	if wc, err = conn.Data(); err != nil {
		return err
	}

	if _, err = fmt.Fprint(wc, body); err != nil {
		return err
	}

	if err = wc.Close(); err != nil {
		return err
	}

	// Send the QUIT command and close the connection
	if err = conn.Quit(); err != nil {
		return err
	}

	return nil
}

// getGeneralHeader returns e-mail general header which contains from, to and
// subject part. It also contains information about e-mail multipart formatting
// (it means that e-mail should contains both html and plaintext version).
func (m *Mailer) getGeneralHeader(subject, boundary string, to []string) string {
	var recipients string

	const rfc2822 = "Mon, 02 Jan 2006 15:04:05 -0700"

	domain := m.From[strings.Index(m.From, "@"):]
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%d-%s", time.Now().UnixNano(), subject)))

	for i, rcpt := range to {
		if i == 0 {
			recipients += rcpt
		} else {
			recipients += ", " + rcpt
		}
	}

	content := "From: " + mime.QEncoding.Encode("utf-8", m.Name) + " <" + m.From + ">\r\n"
	content += "To: " + recipients + "\r\n"
	content += "Reply-To: " + mime.QEncoding.Encode("utf-8", m.Name) + " <no-reply" + domain + ">\r\n"
	content += "Subject: " + mime.QEncoding.Encode("utf-8", subject) + "\r\n"
	content += "MIME-Version: 1.0\r\n"
	content += "Message-ID: <" + fmt.Sprintf("%x", hasher.Sum(nil)) + domain + ">\r\n"
	content += "Date: " + time.Now().Format(rfc2822) + "\r\n"
	content += "Content-Type: multipart/alternative; boundary=\"" + boundary + "\"\r\n"
	content += getReturnPath(m.ReturnPath)
	return content
}

// getReturnPath return valid return if some is set. Otherwise empty string is
// returned.
func getReturnPath(returnPath string) string {
	if returnPath != "" {
		return "Return-Path: " + mime.QEncoding.Encode("utf-8", returnPath) + "\r\n"
	}
	return ""
}

// getPlainTextHeader returns header for plaintext part of message.
func getPlainTextHeader() string {
	content := "Content-Type: text/plain; charset=\"utf-8\"\r\n"
	content += "Content-Transfer-Encoding: BASE64\r\n"
	content += "Content-Disposition: inline\r\n"
	return content
}

// getHTMLHeader returns header for HTML part of message.
func getHTMLHeader() string {
	content := "Content-Type: text/html; charset=\"utf-8\"\r\n"
	content += "Content-Transfer-Encoding: BASE64\r\n"
	content += "Content-Disposition: inline\r\n"
	return content
}

// generateRandomString returns random string.
func generateRandomString(length int) (string, error) {
	var err error
	var randomBytes []byte

	randomString := ""
	stringSeed := getLetters(48, 57) + getLetters(65, 90) + getLetters(97, 122)

	if randomBytes, err = generateRandomBytes(length); err == nil {
		for _, b := range randomBytes {
			randomString += string(stringSeed[b%byte(len(stringSeed))])
		}
	}

	return randomString, err
}

// generateRandomBytes returns random byte stream.
func generateRandomBytes(length int) ([]byte, error) {
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err
	}

	return randomBytes, nil
}

// Returns string of letters of interval [from, to].
func getLetters(from, to int) string {
	letters := ""
	for i := from; i <= to; i++ {
		letters += string(rune(i))
	}
	return letters
}

// Returns splitted string by lines of length 70 chars.
func lineSplit(text string) string {
	var splitted string

	for i, c := range text {
		if i%70 == 0 {
			splitted += "\r\n"
		}

		splitted += string(c)
	}

	return splitted
}
