package main

// The imports
import (
	//"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	//"log"
	"net/http"
	//"net/smtp"
	//"net/url"
	"os"
	"strings"
	"time"

	f "github.com/fauna/faunadb-go/faunadb"
)

/* type Access struct {
	Reference *f.RefV `fauna:"ref"`
	Timestamp int     `fauna:"ts"`
	Secret    string  `fauna:"secret"`
	Role      string  `fauna:"role"`
} */

/*
// Constants
const (
	// The URL to validate reCAPTCHA
	recaptchaURL = "https://www.google.com/recaptcha/api/siteverify"
)

// Variables
var (
	// The reCAPTCHA Secret Token
	recaptchaSecret = os.Getenv("RECAPTCHA_SECRET")
	// The email address to send data to
	emailAddress = os.Getenv("EMAIL_ADDRESS")
	// The email password to use
	emailPassword = os.Getenv("EMAIL_PASSWORD")
	// The SMTP server
	smtpServer = os.Getenv("SMTP_SERVER")
	// The SMTP server port
	smtpPort = os.Getenv("SMTP_PORT")
)
*/
// Handler is the main entry point into tjhe function code as mandated by ZEIT
func Handler(w http.ResponseWriter, r *http.Request) {

	/* str := `

		<!DOCTYPE html>
			<html lang="en">
			<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<title>CODE2GO</title>
			<link href="https://fonts.googleapis.com/css?family=Roboto:300,300i,400,400i,500,500i,700,700i|Roboto+Mono:300,400,700|Roboto+Slab:300,400,700" rel="stylesheet">
			<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
	   		<link href="https://assets.medienwerk.now.sh/material.min.css" rel="stylesheet">
			</head>
			<body style="background-color: #bcbcbc;">
	   		<div class="container" id="search" style="color:white; font-size:30px;">
			<form class="form-inline" role="form" method="POST">
		   	<input class="form-control mr-sm-2" type="text" placeholder="Search" aria-label="Search" id ="search" name ="search">
		   	<button class="btn btn-outline-light my-2 my-sm-1" type="submit">Search</button><br>
			</form>
			</div>
			<br>
			<div class="container" id="nav" style="color:white; font-size:30px;">
			` + time.Now().Format("Monday, Jan 2 2006 15:04:05") + `
			<br>
			</div>

		<script src="https://assets.medienwerk.now.sh/material.min.js">
		</script>
		<br>
		</body>
		</html>
		`

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(str)))
		w.Write([]byte(str)) */

	f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

	// HTTPS will do a PreFlight CORS using the OPTIONS method.
	// To complete that a special response should be sent
	if r.Method == http.MethodOptions {
		response(w, true, "", r.Method)
		return
	}
	/*
		set, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database("code2go"), "role": "server"}))

		if err != nil {
			response(w, false, fmt.Sprintf(time.Now().Format("Monday, Jan 2 2006 15:04:05")+": %s", err.Error()), r.Method)
			return
		}

		var access *Access

		set.Get(&access)

		t := time.Unix(int64(access.Timestamp)/1e6, 0)
	*/
	/*
		// Parse the request body to a map
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		u, err := url.ParseQuery(buf.String())
		if err != nil {
			response(w, false, fmt.Sprintf("There was an error sending your form data: %s", err.Error()), r.Method)
			return
		}

		// Prepare the POST parameters
		urlData := url.Values{}
		urlData.Set("secret", recaptchaSecret)
		urlData.Set("response", u["g-recaptcha-response"][0])

		// Validate the reCAPTCHA
		resp, err := httpcall(recaptchaURL, "POST", "application/x-www-form-urlencoded", urlData.Encode(), nil)
		if err != nil {
			response(w, false, fmt.Sprintf("There was an error sending your form data: %s", err.Error()), r.Method)
			return
		}

		// Validate if the reCAPTCHA was successful
		if !resp.Body["success"].(bool) {
			response(w, false, fmt.Sprintf("There was an error sending your form data: %s", fmt.Sprintf("%v", resp.Body["error-codes"])), r.Method)
			return
		}

		// Set up email authentication information.
		auth := smtp.PlainAuth(
			"",
			emailAddress,
			emailPassword,
			smtpServer,
		)

		// Prepare the email
		mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
		subject := fmt.Sprintf("Subject: [BLOG] Message from %s %s!\n", u["name"][0], u["surname"][0])
		msg := []byte(fmt.Sprintf("%s%s\n%s\n\n%s", subject, mime, u["message"][0], u["email"][0]))

		// Connect to the server, authenticate, set the sender and recipient,
		// and send the email all in one step.
		err = smtp.SendMail(
			fmt.Sprintf("%s:%s", smtpServer, smtpPort),
			auth,
			emailAddress,
			[]string{emailAddress},
			msg,
		)
		if err != nil {
			fmt.Printf("[BLOG] Message from %s %s\n%s\n%s\nThe message was not sent: %s", u["name"][0], u["surname"][0], u["message"][0], u["email"][0], err.Error())
			response(w, false, "There was an error sending your email, but we've logged the data...", r.Method)
			return
		}
	*/
	// Return okay response
	response(w, true, time.Now().Format("Monday, Jan 2 2006 15:04:05"), r.Method)
	return

}

func response(w http.ResponseWriter, success bool, message string, method string) {
	// Create a map for the response body
	body := make(map[string]interface{})

	// Prepare the return data
	if success {
		body["type"] = "logged in"
	} else {
		body["type"] = "failure"
	}
	body["message"] = message
	bodyString, _ := json.Marshal(body)

	// Return the response
	if method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "*")
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyString)
}

// HTTPResponse is the response type for the HTTP requests
type HTTPResponse struct {
	Body       map[string]interface{}
	StatusCode int
	Headers    http.Header
}

// httpcall executes an HTTP request request to a URL and returns the response body as a JSON object
func httpcall(URL string, requestType string, encoding string, payload string, header http.Header) (HTTPResponse, error) {
	// Instantiate a response object
	httpresponse := HTTPResponse{}

	// Prepare placeholders for the request and the error object
	req := &http.Request{}
	var err error

	// Create a request
	if len(payload) > 0 {
		req, err = http.NewRequest(requestType, URL, strings.NewReader(payload))
		if err != nil {
			return httpresponse, fmt.Errorf("error while creating HTTP request: %s", err.Error())
		}
	} else {
		req, err = http.NewRequest(requestType, URL, nil)
		if err != nil {
			return httpresponse, fmt.Errorf("error while creating HTTP request: %s", err.Error())
		}
	}

	// Associate the headers with the request
	if header != nil {
		req.Header = header
	}

	// Set the encoding
	if len(encoding) > 0 {
		req.Header["Content-Type"] = []string{encoding}
	}

	// Execute the HTTP request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return httpresponse, fmt.Errorf("error while performing HTTP request: %s", err.Error())
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return httpresponse, err
	}

	httpresponse.Headers = res.Header
	httpresponse.StatusCode = res.StatusCode

	var data map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		return httpresponse, fmt.Errorf("error while unmarshaling HTTP response to JSON: %s", err.Error())
	}

	httpresponse.Body = data

	return httpresponse, nil
}
