package mailgun

import (
	"os"

	"github.com/sirsean/go-mailgun/mailgun"
)

type Client struct {
	*mailgun.Client
}

func NewClient() *Client {
	apiKey := os.Getenv("MAILGUN_APIKEY")
	domain := os.Getenv("MAILGUN_DOMAIN")
	client := mailgun.NewClient(apiKey, domain)
	return &Client{client}
}

func (c *Client) Send(name, from, to, subject, body string) (string, error) {
	return c.Client.Send(mailgun.Message{
		FromName:    name,
		FromAddress: from,
		ToAddress:   to,
		Subject:     subject,
		Body:        body,
	})
}
