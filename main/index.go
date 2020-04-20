package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Items struct {
	Repos []Repo `json:"items"`
}

type Repo struct {
	Name        string `json:"name"`
	HtmlUrl     string `json:"html_url"`
	Description string `json:"description"`
	Owner       Owner  `json:"owner"`
}

type Owner struct {
	AvatarUrl string `json:"avatar_url"`
}



func Handler(w http.ResponseWriter, r *http.Request) {

	res, err := http.Get("https://api.github.com/search/repositories?q=language:go&sort=stars")

	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	if res.StatusCode != 200 {
		fmt.Fprint(w, res.StatusCode)
	}

	var data Items

	json.Unmarshal(body, &data)

	j := len(data.Repos)

	str := `

	<!DOCTYPE html>
	<html lang="en">
		 <head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<meta http-equiv="X-UA-Compatible" content="ie=edge">
				<meta name="google-site-verification" content="oclzAqcU7qq0dwMfTw3rlaNqZD-7VxcLUa5EA_YG8Io" />
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
					   <ul class="list-group">
	`

	for i := 0; i < j; i++ {

		s := strings.SplitAfterN(data.Repos[i].Owner.AvatarUrl, "?", 2)

		data.Repos[i].Owner.AvatarUrl = s[0] + "s=50&" + s[1]

		str = str + `
		<li class="list-group-item">
		<div class="media">
  <img class="mr-3" src="` + data.Repos[i].Owner.AvatarUrl + `" alt="` + data.Repos[i].Owner.AvatarUrl + `">
  <div class="media-body">
	<h5 class="mt-0"><a href="` + data.Repos[i].HtmlUrl + `">` + data.Repos[i].Name + `</a></h5>` + data.Repos[i].Description + `
	</div>
</div>

		</li><br>
		`

	}

	str = str + `
	</ul>
							  </div>
							  <!-- Then Material JavaScript on top of Bootstrap JavaScript -->
<script src="https://assets.medienwerk.now.sh/material.min.js"></script>
							  </body>
							  </html>
	`

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", strconv.Itoa(len(str)))
	w.Write([]byte(str))

}
