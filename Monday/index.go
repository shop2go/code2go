package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	url := strings.TrimPrefix(r.URL.Path, "/")

	now := time.Now()

	var day time.Time

	for i := 0; i < 7; i++ {

		j := now.AddDate(0, 0, i)

		if j.Weekday().String() == url {

			day = j

			break

		}

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
		<link href="https://fonts.googleapis.com/css?family=Roboto:300,300i,400,400i,500,500i,700,700i|Roboto+Mono:300,400,700|Roboto+Slab:300,400,700" rel="stylesheet">
		<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
   		<link href="https://assets.medienwerk.now.sh/material.min.css" rel="stylesheet">
		</head>
		<body style="background-color: #bcbcbc;">
   		<div class="container" id="search" style="color:white; font-size:30px;">
		<form class="form-inline" role="form">
	   	<input class="form-control mr-sm-2" type="text" placeholder="Search" aria-label="Search" id ="find" name ="find">
	   	<button class="btn btn-outline-light my-2 my-sm-1" type="submit">Search</button><br>
		</div>
		<br>
		<div class="container" id="data" style="color:white; font-size:30px;">
		<form class="form-inline" role="form" method="POST">
		<ul class="list-group">
		`

		for i := 1; i < 53; i++ {

			str = str + `
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='entry#` + strconv.Itoa(day.AddDate(0, 0, i*7).Year()) + `-` + strconv.Itoa(day.AddDate(0, 0, i*7).Month()) + `-` + strconv.Itoa(day.AddDate(0, 0, i*7).Day()) + `'">` + strconv.Itoa(day.AddDate(0, 0, i*7).Year()) + `-` + strconv.Itoa(day.AddDate(0, 0, i*7).Month()) + `-` + strconv.Itoa(day.AddDate(0, 0, i*7).Day()) + `
			</button>
			`
		}

		str = str + `
		</ul>
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
	}

}
