package main

import (
	
	"net/http"
	"strconv"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	/* url := strings.TrimPrefix(r.URL.Path, "/")

	f := strings.Replace(url, "form", "", -1)

	fmt.Fprint(w, f) */

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
						<body style="background-color: #bcbcbc;">



						<div class="container" id="data" style="color:white; font-size:30px;">
						<form class="form-inline" role="form" method="POST">
		<input class="form-control mr-sm-2" type="text" placeholder="topic" aria-label="Topic" id ="Topic" name ="Topic" required><br>
		<input class="form-control mr-sm-2" type="text" placeholder="event" aria-label="Event" id ="Event" name ="Event" required><br>
		<input class="form-control mr-sm-2" type="text" placeholder="tag" aria-label="Tag" id ="Tag" name ="Tag"><br>
		
		<button class="btn btn-outline-light my-2 my-sm-1" type="submit">set</button><br>
	  </div>


	  <script src="https://assets.medienwerk.now.sh/material.min.js"></script>
		</body>
		</html>
		`

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(str)))
		w.Write([]byte(str))

	case "POST":

		url := strings.TrimPrefix(r.URL.Path, "/")


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
						<body style="background-color: #bcbcbc;">



						<div class="container" id="data" style="color:white; font-size:30px;">

						<p>` + url + `</p>

						</div>

						<script src="https://assets.medienwerk.now.sh/material.min.js"></script>
		</body>
		</html>
		`

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(str)))
		w.Write([]byte(str))

	}

}
