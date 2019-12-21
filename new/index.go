package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	//"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	f "github.com/fauna/faunadb-go/faunadb"
	/* 	"github.com/mschneider82/problem"
	   	"github.com/aerogo/packet"
	   	"github.com/mmaedel/code2go/pb" */)

type Access struct {
	//Reference *f.RefV `fauna:"ref"`
	Timestamp int    `fauna:"ts"`
	Secret    string `fauna:"secret"`
	Role      string `fauna:"role"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":

		str := `

		<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>CODE2GO</title>
		<!-- CSS -->
		<!-- Add Material font (Roboto) and Material icon as needed -->
		<link href="https://fonts.googleapis.com/css?family=Roboto:300,300i,400,400i,500,500i,700,700i|Roboto+Mono:300,400,700|Roboto+Slab:300,400,700" rel="stylesheet">
		<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">

		<!-- Add Material CSS, replace Bootstrap CSS -->
		<link href="https://assets.medienwerk.now.sh/material.min.css" rel="stylesheet">
		</head>

		<body style="background-color:#adebad">

		<div class="container" id="data" style="color:white;">
		<form class="form-inline" role="form" method="POST">
		<input readonly="true" class="form-control-plaintext" id="Schedule" aria-label="Schedule" name ="Schedule" value="` + strings.TrimSuffix(r.Host, ".code2go.dev") + `">
		<input class="form-control mr-sm-2" type="password" placeholder="Password" aria-label="Password" id ="Password" name ="Password" value="">
		<input class="form-control mr-sm-2" type="text" placeholder="Title" aria-label="Title" id ="Title" name ="Title" required>
		<!--input class="form-control mr-sm-2" type="text" placeholder="Entry" aria-label="Entry" id ="Entry" name ="Entry" required-->
		<input class="form-control mr-sm-2" type="text" placeholder="Tags" aria-label="Tags" id ="Tags" name ="Tags">
		<input class="form-control mr-sm-2" tyoe="text" aria-label="Content" id ="Content" name ="Content" placeholder="Content"></textarea>
		<br>
		<button type="submit" class="btn btn-light">submit</button>
		</form>
		</div>
		
		<script src="https://assets.medienwerk.now.sh/material.min.js">
		</script>
		</body>
		</html>
		`

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(str)))
		w.Write([]byte(str))

	case "POST":

		r.ParseForm()

/* 		pw := r.Form.Get("Password")
		date := r.Form.Get("Schedule")
		topics := r.Form.Get("Topic")
		content := r.Form.Get("Content")
		tags := r.Form.Get("Tags") */

		pw := r.FormValue("Password")
		date := r.FormValue("Schedule")
		topics := r.FormValue("Title")
		tags := r.FormValue("Tags")
		content := r.FormValue("Entry")


		sl := strings.SplitN(date, "-", -1)

		fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

		x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database(sl[0]), "role": "server"}))

		if err != nil {

			fmt.Fprint(w, err)

			return

		}

		var access *Access

		if err := x.Get(&access); err != nil {

			fmt.Fprint(w, err)

			return

		}

		to := strings.ReplaceAll(topics, " ", "\", \"")

		to = "\"" + to + "\""

		tags = strings.ToLower(tags)

		ta := strings.ReplaceAll(tags, " ", "\", \"")

		ta = "\"" + ta + "\""

		dir := "createPost"

		s := `{"query":"mutation{` + dir + `(data:{iscommited: false password: \"` + pw + `\" date: \"` + date + `\" topics: ` + to + ` content: \"` + content + `\" tags: ` + ta + `}) {_id}}"}`

		body := strings.NewReader(s)
		req, _ := http.NewRequest("POST", "https://graphql.fauna.com/graphql", body)

		req.Header.Set("Authorization", "Bearer "+access.Secret)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("X-Schema-Preview", "partial-update-mutation")

		resp, err := http.DefaultClient.Do(req)

		if err != nil {

			fmt.Fprint(w, err)

			return

		}

		defer resp.Body.Close()

		bdy, _ := ioutil.ReadAll(resp.Body)

		var i interface{}

		json.Unmarshal(bdy, &i)

		if i != nil {

			a := i.(map[string]interface{})

			b := a["data"]

			if b != nil {

				c := b.(map[string]interface{})

				d := c[dir]

				if d != nil {

					e := d.(map[string]interface{})

					f := e["_id"]

					id := f.(string)

					if pw == "" {

						http.Redirect(w, r, "https://"+id+".code2go.dev/public", 301)

					} else {

						http.Redirect(w, r, "https://"+id+".code2go.dev/password", 301)

					}

					fmt.Fprint(w, "checking post id: "+id)

					return

				}

				fmt.Fprint(w, "error dir: ", i)

				return

			}

			fmt.Fprint(w, "error data: ", i)

			return

		}

		fmt.Fprint(w, i)

	}

}
