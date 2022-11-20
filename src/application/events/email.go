package events

type Attachments struct {
	Base64Content string
	Filename      string
	Size          int64
	ContentType   string
}

type Email struct {
	From        string
	To          string
	PlainText   string
	HTML        string
	Attachments []Attachments
	Subject     string
}

type ReceivedEmail = Email

type SendEmail = Email
