package mail

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net"
	"net/mail"
	"os"
	"strings"

	"github.com/mhale/smtpd"
	"github.com/thebaer/burner/validate"
)

var (
	mailInfo  *log.Logger
	mailError *log.Logger
)

type MailConfig struct {
	Port int
	Host string
}

var mailCfg *MailConfig

func Serve(host string, port int) error {
	mailCfg = &MailConfig{
		Port: port,
		Host: host,
	}

	// Set up variables
	mailInfo = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	mailError = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Start mail server
	mailInfo.Printf("Starting %s mail server on :%d", mailCfg.Host, mailCfg.Port)
	err := smtpd.ListenAndServe(fmt.Sprintf(":%d", mailCfg.Port), mailHandler, mailCfg.Host, mailCfg.Host)
	if err != nil {
		mailError.Printf("Couldn't start mail server: %v", err)
		return err
	}

	return nil
}

func mailHandler(origin net.Addr, from string, to []string, data []byte) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		mailError.Printf("Couldn't read email message: %v", err)
		return
	}

	var content []byte
	subject := msg.Header.Get("Subject")
	contentType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		mailError.Printf("Couldn't parse Content-Type: %v", err)
	}
	if strings.HasPrefix(contentType, "multipart/") {
		mr := multipart.NewReader(msg.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				mailError.Printf("Error getting NextPart(): %v", err)
				continue
			}
			if !strings.HasPrefix(p.Header.Get("Content-Type"), "text/plain") {
				continue
			}

			// We correctly read the text/plain section of email
			content, err = ioutil.ReadAll(p)
			if err != nil {
				mailError.Printf("Error in ReadAll on part: %v", err)
			}
			break
		}
	} else {
		var err error
		content, err = ioutil.ReadAll(msg.Body)
		if err != nil {
			mailError.Printf("Couldn't read email body: %v", err)
			return
		}
	}

	mailInfo.Printf("Received mail from %s for %s with subject %s", from, to[0], subject)

	createPostFromEmail(to[0], subject, from, content)
}

func createPostFromEmail(to, subject, from string, content []byte) {
	if err := validate.Email(to); err != nil {
		return
	}
}
