package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	f "github.com/fauna/faunadb-go/faunadb"
	"github.com/muxinc/mux-go"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

type Access struct {
	//Reference *f.RefV `fauna:"ref"`
	Timestamp int    `fauna:"ts"`
	Secret    string `fauna:"secret"`
	Role      string `fauna:"role"`
}

type AssetEntry struct {
	ID       graphql.ID     `graphql:"_id"`
	SourceID graphql.String `graphql:"sourceID"`
	AssetID  graphql.String `graphql:"assetID"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	var content string

	id := r.Host

	id = strings.TrimSuffix(id, "code2go.dev")

	id = strings.TrimSuffix(id, ".")

	if id == "" {

		fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

		x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database("assets"), "role": "server"}))

		if err != nil {

			fmt.Fprintf(w, "a connection error occured: %v\n", err)

		}

		var access *Access

		x.Get(&access)

		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: access.Secret},
		)

		access = &Access{}

		httpClient := oauth2.NewClient(context.Background(), src)

		caller := graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

		client := muxgo.NewAPIClient(
			muxgo.NewConfiguration(
				muxgo.WithBasicAuth(os.Getenv("MUX_ID"), os.Getenv("MUX_SECRET")),
			))

		car := muxgo.CreateAssetRequest{PlaybackPolicy: []muxgo.PlaybackPolicy{muxgo.PUBLIC}}
		cur := muxgo.CreateUploadRequest{NewAssetSettings: car, Timeout: 3600, CorsOrigin: "code2go.dev"}

		res, err := client.DirectUploadsApi.CreateDirectUpload(cur)

		if err != nil {

			fmt.Fprintf(w, "%s %v", "something went wrong...\n", err)

		}

		sourceURL := res.Data.Url

		dul, _ := client.DirectUploadsApi.GetDirectUpload(res.Data.Id)

		sourceID := dul.Data.Id

		var m struct {
			CreateAsset struct {
				AssetEntry
			} `graphql:"createAsset(data:{sourceID: $SourceID})"`
		}

		v := map[string]interface{}{
			"SourceID": graphql.String(sourceID),
		}

		if err = caller.Mutate(context.Background(), &m, v); err != nil {
			fmt.Fprintf(w, "error with input: %v\n", err)
		}

		i := fmt.Sprintf("%s", m.CreateAsset.ID)

		content =

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
			
			<div class="container" id="video" style="color:rgb(255, 255, 255); font-size:30px;">
			
			<br>
			<br>
						
			<h1>video upload:</h1>
			
			<form role="form">
			<input id="picker" type="file" accept="video/*" />
			<p>when promted click ok</p>
			</form>		
			
			<form role="form" method="POST">
						
			<input readonly="true" class="form-control-plaintext" id="ID" aria-label="ID" name ="ID" value="` + i + `" hidden>
			<br>
						
			<p>after file upload completion:</p>
			<button type="submit" class="btn btn-light">submit</button>
			
			</form>
			
			</div>
			<script src="https://unpkg.com/@mux/upchunk@1.0.6/dist/upchunk.js"></script>
			
			<script>
	
				const picker = document.getElementById('picker');
				picker.onchange = () => {
				  const endpoint = '` + sourceURL + `';
				  const file = picker.files[0];
				
				  const upload = UpChunk.createUpload({
					endpoint,
					file,
					chunkSize: 5120,
				  });
				  upload.on('error', err => {
					console.error('something went wrong', err.detail);
				  });
				
				  upload.on('success', () => {
					alert('video file ready :)');
			
				
				  });
				};
				 
			</script>
			
			<script src="https://assets.medienwerk.now.sh/material.min.js"></script>
			</body>
			</html>
					
			`

	} else {

		fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

		x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database("assets"), "role": "server"}))

		if err != nil {

			fmt.Fprintf(w, "a connection error occured: %v\n", err)

		}

		var access *Access

		x.Get(&access)

		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: access.Secret},
		)

		access = &Access{}

		httpClient := oauth2.NewClient(context.Background(), src)

		caller := graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

		client := muxgo.NewAPIClient(
			muxgo.NewConfiguration(
				muxgo.WithBasicAuth(os.Getenv("MUX_ID"), os.Getenv("MUX_SECRET")),
			))

		assets, err := client.AssetsApi.ListAssets()

		if err != nil {

			fmt.Fprintf(w, "%v", err)

		}

		for _, a := range assets.Data {

			inputInfo, _ := client.AssetsApi.GetAssetInputInfo(a.Id)

			for _, b := range inputInfo.Data {

				url := b.Settings.Url

				url = strings.TrimPrefix(url, "https://storage.googleapis.com/video-storage-us-east1-uploads/")

				sl := strings.SplitN(url, "?", -1)

				var q struct {
					AssetBySourceID struct {
						AssetEntry
					} `graphql:"assetBySourceID(sourceID: $SourceID)"`
				}

				v := map[string]interface{}{
					"SourceID": graphql.String(sl[0]),
				}

				if err := caller.Query(context.Background(), &q, v); err != nil {
					fmt.Fprintf(w, "error with asset query: %v\n", err)
				}

				if q.AssetBySourceID.ID == graphql.ID(id) {

					var m struct {
						UpdateAsset struct {
							AssetEntry
						} `graphql:"updateAsset(id: $ID, data:{assetID: $AssetID})"`
					}

					v := map[string]interface{}{
						"ID":      q.AssetBySourceID.ID,
						"AssetID": graphql.String(a.Id),
					}

					if err := caller.Mutate(context.Background(), &m, v); err != nil {
						fmt.Fprintf(w, "error with asset mutation: %v\n", err)
					}

				}

			}

		}

		content =

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
			
			<div class="container" id="video" style="color:rgb(255, 255, 255); font-size:30px;">
			
			<br>
			<br>
						
			<h1>video upload for:</h1>	
			
			<form role="form" method="POST">
						
			<input type="email" class="form-control" placeholder="name@example.com" aria-label="Email" id ="Email" name ="Email" required>
			<input class="form-control mr-sm-2" type="text" placeholder="Last" aria-label="Last" id ="Last" name ="Last" required>
			<input class="form-control mr-sm-2" type="text" placeholder="First" aria-label="First" id ="First" name ="First" required>
			<input class="form-control mr-sm-2" tyoe="text" aria-label="Content" id ="Content" name ="Content" placeholder="Content" required></textarea>
			<br>
			
			<button type="submit" class="btn btn-light">submit</button>
			
			</form>
			
			</div>
						
			<script src="https://assets.medienwerk.now.sh/material.min.js"></script>
			</body>
			</html>
					
			`

	}

	//http.NewRequest("PUT", s, nil)

	switch r.Method {

	case "GET":

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(content)))
		w.Write([]byte(content))

	case "POST":

/* 		id := r.Host

		id = strings.TrimSuffix(id, "code2go.dev")

		id = strings.TrimSuffix(id, ".") */

		if id == "" {

			r.ParseForm()

			id = r.Form.Get("ID")

			//fmt.Fprintf(w, "id: %v\n", i)
			http.Redirect(w, r, "https://"+id+".code2go.dev/video", http.StatusSeeOther)

		} else {

			fmt.Fprint(w, id)

		}

	}

}
