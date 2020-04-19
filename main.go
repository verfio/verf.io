package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

// Message is message
type Message struct {
	From    string
	To      string
	Cc      string
	Subject string
	Body    string
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func sendMail(w http.ResponseWriter, r *http.Request) {

	var body string

	if err := r.ParseMultipartForm(0); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	if r.FormValue("body") != "" {
		body = r.FormValue("body") + "\r\n" + " my name: " + r.FormValue("name") + " my email: " + r.FormValue("email")
	} else {
		body = "hello"
	}

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	var message gmail.Message
	msg := Message{
		"oleg@verf.io",
		"hi@verf.io",
		"oleg.chorny@gmail.com",
		"get in touch",
		body,
	}

	temp := []byte("From:" + msg.From + "\r\n" +
		"reply-to:" + msg.From + "\r\n" +
		"To:  " + msg.To + "\r\n" +
		"Cc:  " + msg.Cc + "\r\n" +
		"Subject: " + msg.Subject + "\r\n" +
		"\r\n" + msg.Body)

	message.Raw = base64.StdEncoding.EncodeToString(temp)
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)

	// imgFile, err := os.Open("image.png") // a QR code image

	// if err != nil {
	// 	log.Fatalf("Error in opening file")
	// }
	// defer imgFile.Close()

	//mediaOptions := googleapi.ContentType("message/rfc822")
	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Fatalf("Unable to send. %v", err)
	}

	body2 := "Thanks for getting in touch with us!\r\nWe will respond to your request soon.\r\n\r\nRegards,\r\nOleg."

	msg = Message{
		"oleg@verf.io",
		r.FormValue("email"),
		"hi@verf.io",
		"VERF.IO: Get in Touch",
		body2,
	}

	temp = []byte("From:" + msg.From + "\r\n" +
		"reply-to:" + msg.From + "\r\n" +
		"To:  " + msg.To + "\r\n" +
		"Cc:  " + msg.Cc + "\r\n" +
		"Subject: " + msg.Subject + "\r\n" +
		"\r\n" + msg.Body)

	message.Raw = base64.StdEncoding.EncodeToString(temp)
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)

	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Fatalf("Unable to send. %v", err)
	}

	// user := "me"
	// r, err := srv.Users.Labels.List(user).Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve labels: %v", err)
	// }
	// if len(r.Labels) == 0 {
	// 	fmt.Println("No labels found.")
	// 	return
	// }
	// fmt.Println("Labels:")
	// for _, l := range r.Labels {
	// 	fmt.Printf("- %s\n", l.Name)
	// }
}

func sendGrid(w http.ResponseWriter, r *http.Request) {
	var body string

	if err := r.ParseMultipartForm(0); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	if r.FormValue("body") != "" {
		body = r.FormValue("body") + "\r\n" + " my name: " + r.FormValue("name") + " my email: " + r.FormValue("email")
	} else {
		body = "hello"
	}

	from := mail.NewEmail("oleg", "oleg@verf.io")
	subject := "get in touch"
	to := mail.NewEmail("hi", "hi@verf.io")
	//plainTextContent := "and easy to do anywhere, even with Go"
	plainTextContent := body
	htmlContent := body
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

func mailJet(w http.ResponseWriter, r *http.Request) {
	var body string

	if err := r.ParseMultipartForm(0); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	if r.FormValue("body") != "" {
		body = r.FormValue("body") + "\r\n" + " my name: " + r.FormValue("name") + " my email: " + r.FormValue("email")
	} else {
		body = "hello"
	}

	m := mailjet.NewMailjetClient(os.Getenv("MAILJET_PUBLIC_KEY"), os.Getenv("MAILJET_PRIVATE_KEY"))
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "oleg@verf.io",
				Name:  "Oleg",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: "hi@verf.io",
					Name:  "hi",
				},
			},
			Subject:  "get in touch",
			TextPart: body,
			HTMLPart: body,
			CustomID: "AppGettingStartedTest",
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := m.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data: %+v\n", res)
}

func main() {

	//fmt.Println(t)
	//r := mux.NewRouter()
	srv := http.NewServeMux()
	srv.Handle("/", http.FileServer(http.Dir(".")))
	//srv.HandleFunc("/webhooks/send", sendMail)

	//srv.HandleFunc("/webhooks/send", sendGrid)
	srv.HandleFunc("/webhooks/send", mailJet)

	log.Printf("server started")

	log.Fatal(http.ListenAndServe(":8080", srv))

}
