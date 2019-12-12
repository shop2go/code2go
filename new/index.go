package main

import (
	"fmt"

	//"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	f "github.com/fauna/faunadb-go/faunadb"
	/* 	"github.com/mschneider82/problem"
	   	"github.com/aerogo/packet"
	   	"github.com/mmaedel/code2go/pb" */)

type Post struct {
	ID       string
	Date     string
	Password string
	Title    string
	Content  interface{}
}

type Access struct {
	//Reference *f.RefV `fauna:"ref"`
	Timestamp int    `fauna:"ts"`
	Secret    string `fauna:"secret"`
	Role      string `fauna:"role"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	u := url.URL.String()

	var access *Access

	fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

	x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database(time.Now().Format("2006")), "role": "server"}))

	if err != nil {

		fmt.Fprint(w, err)

		return

	}

	if err := x.Get(&access); err != nil {

		fmt.Fprint(w, err)

		return

	}

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
	`

		+u +

			`
	<form class="form-inline" role="form" method="POST">
	<input readonly="true" class="form-control-plaintext" id="Schedule" aria-label="Schedule" name ="Schedule" value="">
	<input readonly="true" class="form-control-plaintext" id="Password" aria-label="Password" name ="Password" value="">
	<input class="form-control mr-sm-2" type="text" placeholder="Title" aria-label="Title" id ="Title" name ="Title" required>
	<!--input class="form-control mr-sm-2" type="text" placeholder="entry" aria-label="Entry" id ="Entry" name ="Entry" required-->
	<input class="form-control mr-sm-2" type="text" placeholder="Tags" aria-label="Tags" id ="Tags" name ="Tags">
	<textarea class="form-control  mr-sm-2" id="Content" rows="2" placeholder="Content"></textarea>
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

	}

}
