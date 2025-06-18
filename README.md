# mailgo

mailgo is a simple and easy-to-use SMTP wrapper library for Go applications. It provides a convenient interface for sending emails, including support for attachments, multiple recipients, and authentication.

## Features
- Send emails via SMTP with minimal setup
- Support for multiple recipients
- Send emails with or without authentication
- Add attachments to emails (since v1.4.0)
- Customizable email headers
- Error handling for each recipient

## Installation

```
go get github.com/mailgo
```

## Usage

```go
package main

import (
    "github.com/mailgo"
)

func main() {
    m := mailgo.NewMailer(mailgo.Options{
        Host:     "smtp.example.com",
        Port:     587,
        Username: "your_username",
        Password: "your_password",
        From:     "your@email.com",
    })

    err := m.SendMail(mailgo.Message{
        To:      []string{"recipient1@example.com", "recipient2@example.com"},
        Subject: "Test Email",
        Body:    "This is a test email sent using mailgo.",
        Attachments: []string{"/path/to/file.pdf"}, // Optional: add attachments
    })
    if err != nil {
        // handle error
    }
}
```

## Options

- `Host`: SMTP server address
- `Port`: SMTP server port
- `Username`: SMTP username (optional for unauthenticated sending)
- `Password`: SMTP password (optional for unauthenticated sending)
- `From`: Sender email address

## Message

- `To`: List of recipient email addresses
- `Subject`: Email subject
- `Body`: Email body (plain text)
- `Attachments`: List of file paths to attach (optional)

## Changelog
See [CHANGELOG.md](CHANGELOG.md) for release notes.

## License
This project is licensed under the MIT License.
