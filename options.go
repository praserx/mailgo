package mailgo

import "github.com/emersion/go-sasl"

type MailerOptions struct {
	Host  string
	Port  string
	Name  string
	From  string
	Creds sasl.Client
}

type MailerOption func(*MailerOptions)

func WithHost(host string) MailerOption {
	return func(opts *MailerOptions) {
		opts.Host = host
	}
}

func WithPort(port string) MailerOption {
	return func(opts *MailerOptions) {
		opts.Port = port
	}
}

func WithName(name string) MailerOption {
	return func(opts *MailerOptions) {
		opts.Name = name
	}
}

func WithFrom(from string) MailerOption {
	return func(opts *MailerOptions) {
		opts.From = from
	}
}

func WithCredentials(username, password string) MailerOption {
	return func(opts *MailerOptions) {
		opts.Creds = sasl.NewPlainClient("", username, password)
	}
}
