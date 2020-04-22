package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	//"github.com/mmaedel/code2go/pb"
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
		<div class="container" id="nav" style="color:white;">
		<br>`

		switch day.Weekday() {

		case 0:

			str = str + `
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Monday'">Monday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Tuesday'">Tuesday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Wednesday'">Wednesday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Thursday'">Thursday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Friday'">Friday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Saturday'">Saturday</button>
			<button type="button" class="btn btn-light" onclick="window.location.href='Sunday'">Sunday</button>
			`

		case 1:

			str = str + `
			<button type="button" class="btn btn-light" onclick="window.location.href='Monday'">Monday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Tuesday'">Tuesday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Wednesday'">Wednesday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Thursday'">Thursday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Friday'">Friday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Saturday'">Saturday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Sunday'">Sunday</button>
			`

		case 2:

			str = str + `
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Monday'">Monday</button>
			<button type="button" class="btn btn-light" onclick="window.location.href='Tuesday'">Tuesday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Wednesday'">Wednesday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Thursday'">Thursday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Friday'">Friday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Saturday'">Saturday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Sunday'">Sunday</button>
			`

		case 3:

			str = str + `
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Monday'">Monday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Tuesday'">Tuesday</button>
			<button type="button" class="btn btn-light" onclick="window.location.href='Wednesday'">Wednesday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Thursday'">Thursday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Friday'">Friday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Saturday'">Saturday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Sunday'">Sunday</button>
			`

		case 4:

			str = str + `
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Monday'">Monday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Tuesday'">Tuesday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Wednesday'">Wednesday</button>
			<button type="button" class="btn btn-light" onclick="window.location.href='Thursday'">Thursday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Friday'">Friday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Saturday'">Saturday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Sunday'">Sunday</button>
			`

		case 5:

			str = str + `
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Monday'">Monday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Tuesday'">Tuesday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Wednesday'">Wednesday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Thursday'">Thursday</button>
			<button type="button" class="btn btn-light" onclick="window.location.href='Friday'">Friday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Saturday'">Saturday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Sunday'">Sunday</button>
			`

		case 6:

			str = str + `
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Monday'">Monday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Tuesday'">Tuesday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Wednesday'">Wednesday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Thursday'">Thursday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Friday'">Friday</button>
			<button type="button" class="btn btn-light" onclick="window.location.href='Saturday'">Saturday</button>
			<button type="button" class="btn btn-outline-dark" onclick="window.location.href='Sunday'">Sunday</button>
			`

		}

		str = str + `
		</div>
		<br>
		<div class="container" id="data" style="color:white; font-size:30px;">
		<form class="form-inline" role="form" method="POST">
		<ul class="list-group">
		<button type="button" class="btn btn-link btn-outline-dark" onclick="window.location.href='entry#` + day.Format("2006") + `-` + day.Format("1") + `-` + day.Format("2") + `'">
		<span class="badge badge-pill badge-light">
		` + day.Format("2006") + `-` + day.Format("1") + `-` + day.Format("2") + `
		</span>
		</button>		
		`

		for i := 1; i < 53; i++ {

			str = str + `
			<button type="button" class="btn btn-link btn-outline-dark" onclick="window.location.href='entry#` + day.AddDate(0, 0, i*7).Format("2006") + `-` + day.AddDate(0, 0, i*7).Format("1") + `-` + day.AddDate(0, 0, i*7).Format("2") + `'">
			<span class="badge badge-pill badge-light">
			` + day.AddDate(0, 0, i*7).Format("2006") + `-` + day.AddDate(0, 0, i*7).Format("1") + `-` + day.AddDate(0, 0, i*7).Format("2") + `
			</span>
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
