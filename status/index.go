package main

import (
	//"fmt"
	"net/http"
	//"os"
	"strconv"
	//"strings"
	//"github.com/mmaedel/code2go/pb"

	//f "github.com/fauna/faunadb-go/faunadb"
)

type Access struct {
	//Reference *f.RefV `fauna:"ref"`
	Timestamp int    `fauna:"ts"`
	Secret    string `fauna:"secret"`
	Role      string `fauna:"role"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":

		//strings.TrimSuffix(r.Host, ".code2go.dev")

		/* fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

		x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database(sl[0]), "role": "server"}))

		if err != nil {

			fmt.Fprint(w, err)

			return

		}

		var access *Access

		if err := x.Get(&access); err != nil {

			fmt.Fprint(w, err)

			return

		} */

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
	}

}
