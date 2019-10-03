package src

import (
	//"errors"
	"fmt"
	//"log"
	//"log"
	"net/http"
	"os"
	"strconv"
	"time"

	f "github.com/fauna/faunadb-go/faunadb"
)

type Booking struct {
	Device string
	ID     uint
	Room   string
	Log    int64
}

type Access struct {
	Reference *f.RefV `fauna:"ref"`
	Timestamp int     `fauna:"ts"`
	Secret    string  `fauna:"secret"`
	Role      string  `fauna:"role"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	//var id f.RefV
	os.Setenv()

	c := f.NewFaunaClient(os.Getenv("FAUNA"))

	s, err := c.Query(f.CreateKey(f.Obj{"database": f.Database("code2go"), "role": "server-readonly"}))

	if err != nil {

		fmt.Fprintf(w, err.Error())

	}

	var access *Access

	s.Get(&access)

	t := time.Unix(int64(access.Timestamp)/1e6, 0)

	//log.Printf("%v %v\n", access.Reference.ID, access.Role)

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
		</form>
		</div>
		<br>
		<div class="container" id="nav" style="color:white;">
		` + t.Format("Mon Jan 2 15:04:05 -0700 MST 2006") + `
		<br>
		` + access.Reference.ID + `
		<br>`

	switch r.Method {

	case "GET":

		str = str + `
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
