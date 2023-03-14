package mail

import (
	"testing"

	"github.com/gost-codes/sweet_dreams/util"
	"github.com/stretchr/testify/require"
)

func TestSendMailWithGmail(t *testing.T) {
	config, err := util.LoadConfig("..")

	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	subject := "Test email from me"
	to := []string{"hopedorkenoo@gmail.com"}
	content := `
	<h1>Hello Email</h1>
	<p> This is a test email from sweet dreams golang backend<br> <a href="https://google.com">Check this out</a></p>
	`
	attachmentFiles := []string{"../ReadMe.md"}
	err = sender.SendEmail(subject, content, to, nil, nil, attachmentFiles)

	require.NoError(t, err)

}
