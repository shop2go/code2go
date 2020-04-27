package main

import (
	"context"
	"encoding/base64"
	//"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	f "github.com/fauna/faunadb-go/faunadb"
	"github.com/muxinc/mux-go"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"

	"github.com/dgrijalva/jwt-go"
)

type Access struct {
	//Reference *f.RefV `fauna:"ref"`
	Timestamp int    `fauna:"ts"`
	Secret    string `fauna:"secret"`
	Role      string `fauna:"role"`
}

type AssetEntry struct {
	ID       graphql.ID      `graphql:"_id"`
	Key      graphql.String  `graphql:"key"`
	SourceID graphql.String  `graphql:"sourceID"`
	AssetID  graphql.String  `graphql:"assetID"`
	PbID     graphql.String  `graphql:"pbID"`
	First    graphql.String  `graphql:"first"`
	Last     graphql.String  `graphql:"last"`
	Email    graphql.String  `graphql:"email"`
	Title    graphql.String  `graphql:"title"`
	Category graphql.String  `graphql:"category"`
	Content  graphql.String  `graphql:"content"`
	Policy   graphql.String  `graphql:"policy"`
	Checked  graphql.Boolean `graphql:"checked"`
}

type KeyEntry struct {
	Key   graphql.String `graphql:"key"`
	Token graphql.String `graphql:"token"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	var content, email, key, title, pbid, policy string

	id := r.Host

	id = strings.TrimSuffix(id, "code2go.dev")

	id = strings.TrimSuffix(id, ".")

	switch id {

	case "":

		http.Redirect(w, r, "https://code2go.dev/videos", http.StatusSeeOther)

	//various stages
	default:

		if _, err := strconv.Atoi(id); err == nil {

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

			var q struct {
				FindAssetByID struct {
					AssetEntry
				} `graphql:"findAssetByID(id: $ID)"`
			}

			v := map[string]interface{}{
				"ID": graphql.ID(id),
			}

			if err := caller.Query(context.Background(), &q, v); err != nil {
				fmt.Fprintf(w, "error with asset query: %v\n", err)
			}

			key = string(q.FindAssetByID.Key)

			policy = string(q.FindAssetByID.Policy)

			pbid = string(q.FindAssetByID.PbID)

			email = string(q.FindAssetByID.Email)

			title = string(q.FindAssetByID.Title)

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

			<br>
			<br>

			<div class="container" id="data" style="color:rgb(255, 255, 255); font-size:30px;">
						
			<h1>content owner:</h1>	
			
			<form role="form" method="POST">

			<br>
			<label for="Email">Email Address</label>
			<input type="email" class="form-control" placeholder="name@example.com" aria-label="Email" id ="Email" name ="Email" required><br>
			
			<br>
			
			<button type="submit" class="btn btn-light">submit</button>
			
			</form>
			
			</div>
			
					
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

			var q struct {
				AssetByTitle struct {
					AssetEntry
				} `graphql:"assetByTitle(title: $Title, checked: $Checked)"`
			}

			v := map[string]interface{}{
				"Title":   graphql.String(id),
				"Checked": graphql.Boolean(true),
			}

			if err := caller.Query(context.Background(), &q, v); err != nil {
				fmt.Fprintf(w, "error with asset query: %v\n", err)
			}

			policy = string(q.AssetByTitle.Policy)

			pbid = string(q.AssetByTitle.PbID)

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

			`

			if policy == "signed" {

				var tokenString string

				switch key {

				case "":

					client := muxgo.NewAPIClient(
						muxgo.NewConfiguration(
							muxgo.WithBasicAuth(os.Getenv("MUX_ID"), os.Getenv("MUX_SECRET")),
						))

					k, err := client.URLSigningKeysApi.CreateUrlSigningKey()

					if err != nil {

						fmt.Fprintf(w, "%v", err)

					}

					decodedKey, _ := base64.StdEncoding.DecodeString(k.Data.PrivateKey)

					signKey, err := jwt.ParseRSAPrivateKeyFromPEM(decodedKey)
					if err != nil {
						fmt.Fprintf(w, "%v", err)
					}

					token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
						"sub": pbid,
						"aud": "v",
						"exp": time.Now().Add(time.Minute * 15).Unix(),
						"kid": k.Data.Id,
					})

					tokenString, err = token.SignedString(signKey)
					if err != nil {
						fmt.Fprintf(w, "%v", err)
					}

					var m struct {
						UpdateAsset struct {
							AssetEntry
						} `graphql:"updateAsset(id: $ID, data:{key: $Key, checked: $Checked})"`
					}

					v := map[string]interface{}{
						"ID":      q.AssetByTitle.ID,
						"Key":     graphql.String(k.Data.Id),
						"Checked": q.AssetByTitle.Checked,
					}

					if err := caller.Mutate(context.Background(), &m, v); err != nil {
						fmt.Fprintf(w, "error with asset update: %v\n", err)
					}

					var n struct {
						CreateKey struct {
							KeyEntry
						} `graphql:"createKey(data:{key: $Key, token: $Token})"`
					}

					v = map[string]interface{}{
						"Key":   graphql.String(k.Data.Id),
						"Token": graphql.String(tokenString),
					}

					if err := caller.Mutate(context.Background(), &n, v); err != nil {
						fmt.Fprintf(w, "error with key create: %v\n", err)
					}

				default:

					var q struct {
						KeyByKey struct {
							KeyEntry
						} `graphql:"keyByKey(key: $Key)"`
					}

					v := map[string]interface{}{
						"Key": graphql.String(key),
					}

					if err := caller.Query(context.Background(), &q, v); err != nil {
						fmt.Fprintf(w, "error with key query: %v\n", err)
					}

					tokenString = string(q.KeyByKey.Token)

				}

				content = `

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
	
				<br>
				<br>

				<div class="container" id="data">

	
				<div class="media">
				<video class="align-self-start mr-3" id="Video" width="214" controls></video>		
			<div class="media-body"><br><br>
					
			<h2 class="mt-0">`+id+`</h2>
					
			<h3>`+string(q.AssetByTitle.Category)+`: </h3><p>`+string(q.AssetByTitle.Content)+`</p>
			</div>
			</div>

			</div>


			

<!-- Use HLS.js to support the HLS format in browsers. -->
<script src="https://cdn.jsdelivr.net/npm/hls.js@0.8.2"></script>
<script>
  (function(){
    // Replace with your asset's playback ID
    
    var url = "https://stream.mux.com/` + pbid + `.m3u8?token=` + tokenString + `";

    // HLS.js-specific setup code
    if (Hls.isSupported()) {
      var video = document.getElementById("Video");
      var hls = new Hls();
      hls.loadSource(url);
      hls.attachMedia(video);
    }
  })();
</script>

						
			<script src="https://assets.medienwerk.now.sh/material.min.js"></script>
			</body>
			</html>
					
			`

			} else if policy == "public" {

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
			
			<br>
			<br>

			<div class="container" id="data">

			<div class="media">
				<video class="align-self-start mr-3" id="Video" width="214" controls></video>		
			<div class="media-body"><br><br>
					
			<h2 class="mt-0">`+id+`</h2>
					
			<h3>`+string(q.AssetByTitle.Category)+`: </h3><p>`+string(q.AssetByTitle.Content)+`</p>
			</div>
			</div>
			</div>

<!-- Use HLS.js to support the HLS format in browsers. -->
<script src="https://cdn.jsdelivr.net/npm/hls.js@0.8.2"></script>
<script>
  (function(){
    // Replace with your asset's playback ID
    var url = "https://stream.mux.com/` + pbid + `.m3u8";

    // HLS.js-specific setup code
    if (Hls.isSupported()) {
      var video = document.getElementById("myVideo");
      var hls = new Hls();
      hls.loadSource(url);
      hls.attachMedia(video);
    }
  })();
</script>

						
			<script src="https://assets.medienwerk.now.sh/material.min.js"></script>
			</body>
			</html>
					
			`

			}

		}

	}

	switch r.Method {

	case "GET":

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(content)))
		w.Write([]byte(content))

	case "POST":

		switch id {

		//various stages
		default:

			if _, err := strconv.Atoi(id); err == nil {

				r.ParseForm()

				id = r.Form.Get("Email")

				if id == email {

					http.Redirect(w, r, "https://"+title+".code2go.dev/video", http.StatusSeeOther)

				} else {

					http.Redirect(w, r, "https://code2go.dev/videos", http.StatusSeeOther)

				}

			}

		}

	}

}
