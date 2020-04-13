package main

import (
	"context"
	//"encoding/base64"
	"fmt"
	"net/http"
	"os"
	//"sort"
	"strconv"
	//"strings"

	f "github.com/fauna/faunadb-go/faunadb"
	"github.com/muxinc/mux-go"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
	//"github.com/plutov/paypal"
)

type Access struct {
	//Reference *f.RefV `fauna:"ref"`
	Timestamp int    `fauna:"ts"`
	Secret    string `fauna:"secret"`
	Role      string `fauna:"role"`
}

type InputEntry struct {
	ID  graphql.ID     `graphql:"_id"`
	Url graphql.String `graphql:"url"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	client := muxgo.NewAPIClient(
		muxgo.NewConfiguration(
			muxgo.WithBasicAuth(os.Getenv("MUX_ID"), os.Getenv("MUX_SECRET")),
		))

	car := muxgo.CreateAssetRequest{PlaybackPolicy: []muxgo.PlaybackPolicy{muxgo.SIGNED}}
	cur := muxgo.CreateUploadRequest{NewAssetSettings: car, Timeout: 3600, CorsOrigin: "code2go.dev"}
	res, err := client.DirectUploadsApi.CreateDirectUpload(cur)

	if err != nil {

		fmt.Fprintf(w, "%s %v", "something went wrong...\n", err)

	}

	i := res.Data.Status

	s := res.Data.Url

	//http.NewRequest("PUT", s, nil)

	switch r.Method {

	case "GET":

		str :=

			`
	<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title>vdo2go</title>
	<!-- CSS -->
	<!-- Add Material font (Roboto) and Material icon as needed -->
	<link href="https://fonts.googleapis.com/css?family=Roboto:300,300i,400,400i,500,500i,700,700i|Roboto+Mono:300,400,700|Roboto+Slab:300,400,700" rel="stylesheet">
	<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">

	<!-- Add Material CSS, replace Bootstrap CSS -->
	<link href="https://assets.medienwerk.now.sh/material.min.css" rel="stylesheet">

	</head>
	<body style="background-color: #a1b116;">

	<div class="container" id="shop" style="color:rgb(255, 255, 255); font-size:30px;">

	<br>
	<br>

	<h1>` + i + ` for video...</h1>

	<form role="form" method="POST">

	<button type="submit" class="btn btn-light">
	<input id="picker" type="file" />
	</button>

	</form>	

	</div>

	<script src="https://unpkg.com/@mux/upchunk@1.0.6/dist/upchunk.js"></script>

	<script>

	const picker = document.getElementById('picker');
	picker.onchange = () => {
	  const endpoint = '` + s + `';
	  const file = picker.files[0];
	
	  const upload = UpChunk.createUpload({
		endpoint,
		file,
	  });
	  upload.on('error', err => {
		console.error('something went wrong', err.detail);
	  });
	
	  upload.on('attempt', ({ detail }) => {
		alert('uploading...', detail);
	  });
	
	  upload.on('success', () => {
		alert('ready');

	
	  });
	};
	  </script>
				   
	<script src="https://assets.medienwerk.now.sh/material.min.js"></script>
	</body>
	</html>

	`
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(str)))
		w.Write([]byte(str))

	case "POST":

		fc := f.NewFaunaClient("FAUNA_ACCESS")

		x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database("assets"), "role": "server"}))

		if err != nil {

			fmt.Fprintf(w, "a connection error occured: %v\n", err)

		}

		var access *Access

		x.Get(&access)

		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: access.Secret},
		)

		httpClient := oauth2.NewClient(context.Background(), src)

		call := graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

		var m struct {
			CreateInput struct {
				InputEntry
			} `graphql:"createInput(data:{url: $Url})"`
		}

		v := map[string]interface{}{
			"Url": graphql.String(s),
		}

		if err = call.Mutate(context.Background(), &m, v); err != nil {
			fmt.Printf("error with input: %v\n", err)
		}

		s = fmt.Sprintf("%s", m.CreateInput.ID)

		http.Redirect(w, r, "https://"+s+"code2go.dev/video", http.StatusSeeOther)
		/* 	str :=

				`
			<!DOCTYPE html>
			<html lang="en">
			<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<title>shop2go</title>
			<!-- CSS -->
			<!-- Add Material font (Roboto) and Material icon as needed -->
			<link href="https://fonts.googleapis.com/css?family=Roboto:300,300i,400,400i,500,500i,700,700i|Roboto+Mono:300,400,700|Roboto+Slab:300,400,700" rel="stylesheet">
			<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">

			<!-- Add Material CSS, replace Bootstrap CSS -->
			<link href="https://assets.medienwerk.now.sh/material.min.css" rel="stylesheet">

			</head>
			<body style="background-color: #a1b116;">

			<div class="container" id="shop" style="color:rgb(255, 255, 255); font-size:30px;">

			<ul class="list-group">
			<br>
			<br>

			<h1>Video</h1>

			<form role="form">

			<input id="file-picker" type="file" />

			</form>

			<script type="module" src="https://unpkg.com/@mux/upchunk@1">

			import * as UpChunk from '@mux/upchunk';

			const filePicker = document.getElementById('file-picker');

			const url = '` + s + `'

			filePicker.onchange = () => {

			  const upload = UpChunk.createUpload({
				file: filePicker.files[0],
				endpoint: url,
			  });

			  upload.on('success', () => alert('file uploaded'));

			};

			  </script>
			`

			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("Content-Length", strconv.Itoa(len(str)))
			w.Write([]byte(str))

		}
		*/
	}
}
