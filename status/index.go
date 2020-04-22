package main

import (
	//"context"
	//"fmt"
	"net/http"
	//"os"
	"strconv"
	//"strings"

/* 	"golang.org/x/oauth2"

	f "github.com/fauna/faunadb-go/faunadb"
	"github.com/shurcooL/graphql" */
)

type Access struct {
	//Reference *f.RefV `fauna:"ref"`
	Timestamp int    `fauna:"ts"`
	Secret    string `fauna:"secret"`
	Role      string `fauna:"role"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
/* 
	fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

	x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database("users"), "role": "server"}))

	if err != nil {

		fmt.Fprint(w, err)

		return

	}

	var access *Access

	x.Get(&access)

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: access.Secret},
	)

	httpClient := oauth2.NewClient(context.Background(), src)

	call := graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

	u := r.Host

	u = strings.TrimSuffix(u, "code2go.dev")

	if u != "" {

		u = strings.TrimSuffix(u, ".")

	}
 */
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
		<div class="container" id="data" style="color:white;">
		<br>
		<form class="form-inline" role="form" method="POST">
		<input type="email" class="form-control" placeholder="name@example.com" aria-label="Email" id="Email" name ="Email">
		<br>
		<input type="text" class="form-control" aria-label="User" id="User" name ="User">
		<input class="form-control mr-sm-2" type="password" placeholder="Secret" aria-label="Secret" id="Secret" name ="Secret" value="">
		<br>
		<button type="submit" class="btn btn-light">login</button>
		</form>
		</div>
		`

		str = str + `
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
