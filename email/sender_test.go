package email

import (
	"testing"

	"github.com/chensheep/simple-bank-backend/util"
	"github.com/stretchr/testify/require"
)

func TestXxx(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	err = sender.SendEmail(
		[]string{"chensheep1005@gmail.com"},
		[]string{},
		[]string{},
		"A test email",
		`
		<h1>Hello</h1>
		<p>This is a test email from <a href="https://github.com/chensheep/simple-bank-backend">fancigo simple bank</a>.</p>
		`,
		[]string{"../README.md"},
	)
	require.NoError(t, err)
}
