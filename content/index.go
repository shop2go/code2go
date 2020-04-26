package main

import (
	"context"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

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

func Handler(w http.ResponseWriter, r *http.Request) {

	var content, pbid string

	id := r.Host

	id = strings.TrimSuffix(id, "code2go.dev")

	id = strings.TrimSuffix(id, ".")

	switch id {

	//various stages
	default:

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

		pbid = string(q.FindAssetByID.PbID)

		if q.FindAssetByID.Policy == "signed" {

			client := muxgo.NewAPIClient(
				muxgo.NewConfiguration(
					muxgo.WithBasicAuth(os.Getenv("MUX_ID"), os.Getenv("MUX_SECRET")),
				))

			k, err := client.URLSigningKeysApi.CreateUrlSigningKey()

			if err != nil {

				fmt.Fprintf(w, "%v", err)

			}

			pk, _ := base64.StdEncoding.DecodeString(k.Data.PrivateKey)

			/* block, _ := pem.Decode(pk)
			if block.Type != "RSA PRIVATE KEY" {
				fmt.Fprintf(w, "%s %s", block.Type, "error!")
			} */

			type Claim struct {
				//Kid string `json:"kid"`
				jwt.StandardClaims
			}

			claims := Claim{
				//k.Data.Id,
				jwt.StandardClaims{
					Subject:   pbid,
					Audience:  "v",
					ExpiresAt: 15000,
					Issuer:    r.Host,
					Id:        string(q.FindAssetByID.AssetID),
				},
			}

			t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			//privKey, _ := base64.StdEncoding.DecodeString(k.Data.PrivateKey)

			/* 			block, _ := pem.Decode(privKey)
			   			if block.Type != "RSA PRIVATE KEY" {
			   				fmt.Fprintf(w, "%s %s", block.Type, "err!")
			   			}
			*/
			t.Header = map[string]interface{}{
				"kid": k.Data.Id,
			}

			token, err := t.SignedString(pk)

			if err != nil {

				fmt.Fprintf(w, "%v", err)

			}

			/* 		client := muxgo.NewAPIClient(
			   			muxgo.NewConfiguration(
			   				muxgo.WithBasicAuth(os.Getenv("MUX_ID"), os.Getenv("MUX_SECRET")),
			   			))

			   		asset, err := client.AssetsApi.GetAsset(string(q.FindAssetByID.AssetID))

			   		if err != nil {

			   			fmt.Fprintf(w, "%v", err)

			   		} */

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
						
			<h1>content owner:</h1>	
			
			<form role="form" method="POST">

			<input type="email" class="form-control" placeholder="name@example.com" aria-label="Email" id ="Email" name ="Email" required><br>
			
			<br>
			
			<button type="submit" class="btn btn-light">submit</button>
			
			</form>
			
			</div>

			<video id="myVideo" controls></video>

<!-- Use HLS.js to support the HLS format in browsers. -->
<script src="https://cdn.jsdelivr.net/npm/hls.js@0.8.2"></script>
<script>
  (function(){
    // Replace with your asset's playback ID
    var playbackId = "` + pbid + `";
    var url = "https://stream.mux.com/"+playbackId+".m3u8?token={` + token + `}";

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

		switch r.Method {

		case "GET":

			switch id {

			//start
			default:

				w.Header().Set("Content-Type", "text/html")
				w.Header().Set("Content-Length", strconv.Itoa(len(content)))
				w.Write([]byte(content))

			}

			/* 	case "POST":

				switch id {

				case "public":

					r.ParseForm()

					id = r.Form.Get("ID")

					//fmt.Fprintf(w, "id: %v\n", i)
					http.Redirect(w, r, "https://"+id+".code2go.dev/video", http.StatusSeeOther)

				case "signed":

					r.ParseForm()

					id = r.Form.Get("ID")

					//fmt.Fprintf(w, "id: %v\n", i)
					http.Redirect(w, r, "https://"+id+".code2go.dev/video", http.StatusSeeOther)

				case "":

					r.ParseForm()

					c := r.Form.Get("Access")

					if c == "" {

						c = "public"

					}

					http.Redirect(w, r, "https://"+c+".code2go.dev/video", http.StatusSeeOther)

				default:

					r.ParseForm()

					first := r.Form.Get("First")
					last := r.Form.Get("Last")
					email := strings.ToLower(r.Form.Get("Email"))
					title := strings.ToLower(r.Form.Get("Title"))
					category := strings.ToLower(r.Form.Get("Category"))
					content := strings.ToLower(r.Form.Get("Content"))

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

					httpClient := oauth2.NewClient(context.Background(), src)

					caller := graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

				PREV:

					if _, err = strconv.Atoi(id); err == nil {

						var q struct {
							FindAssetByID struct {
								AssetEntry
							} `graphql:"findAssetByID(id: $ID)"`
						}

						v := map[string]interface{}{
							"ID": graphql.ID(id),
						}

						if err := caller.Query(context.Background(), &q, v); err != nil {
							fmt.Fprintf(w, "error with asset source: %v\n", err)
						}

						//basic auth assign
						switch string(q.FindAssetByID.Email) {

						default:

							http.Redirect(w, r, "https://"+id+".code2go.dev/video", http.StatusSeeOther)

						case "":

							t := time.Unix(int64(access.Timestamp)/1e6, 0)

							s := t.Format("20060102150405")

							var m struct {
								UpdateAsset struct {
									AssetEntry
								} `graphql:"updateAsset(id: $ID, data:{title: $Title, category: $Category, pbID: $PbID, email: $Email, first: $First, last: $Last, content: $Content, checked: $Checked})"`
							}

							v = map[string]interface{}{
								"ID":       graphql.ID(id),
								"Email":    graphql.String(email),
								"First":    graphql.String(first),
								"Last":     graphql.String(last),
								"Category": graphql.String(category),
								"Title":    graphql.String(title + "_" + s),
								"Content":  graphql.String(content),
								"PbID":     graphql.String(pbid),
								"Checked":  graphql.Boolean(false),
							}

							if err := caller.Mutate(context.Background(), &m, v); err != nil {
								fmt.Fprintf(w, "error with asset update: %v\n", err)
							} else {

								http.Redirect(w, r, "https://"+title+"_"+s+".code2go.dev/video", http.StatusSeeOther)

							}

						//todo auth
						case email:

							//todo with token
							var m struct {
								UpdateAsset struct {
									AssetEntry
								} `graphql:"updateAsset(id: $ID, data:{checked: $Checked})"`
							}

							v = map[string]interface{}{
								"ID":      graphql.ID(id),
								"Checked": graphql.Boolean(true),
							}

							if err := caller.Mutate(context.Background(), &m, v); err != nil {
								fmt.Fprintf(w, "error with asset confirmation: %v\n", err)
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

							<a href="https://` + id + `.code2go.dev/content"><img src="https://image.mux.com/` + string(m.UpdateAsset.PbID) + `/thumbnail.jpg?width=214&height=121&fit_mode=pad"></a>
							<br>

							<p>` + string(m.UpdateAsset.First) + `<br>` + string(m.UpdateAsset.Title) + ` is "` + string(m.UpdateAsset.Policy) + `" content:<br>` + string(m.UpdateAsset.Content) + `</p>
							</div>

							<script src="https://assets.medienwerk.now.sh/material.min.js"></script>
							</body>
							</html>

							`

							w.Header().Set("Content-Type", "text/html")
							w.Header().Set("Content-Length", strconv.Itoa(len(content)))
							w.Write([]byte(content))

						}

					} else {

						var q struct {
							AssetsByEmail struct {
								Data []struct {
									AssetEntry
								}
							} `graphql:"assetsByEmail(email: $Email, checked: $Checked)"`
						}

						v := map[string]interface{}{
							"Email":   graphql.String(email),
							"Checked": graphql.Boolean(false),
						}

						if err := caller.Query(context.Background(), &q, v); err != nil {
							fmt.Fprintf(w, "error with asset source: %v\n", err)
						}

						for _, s := range q.AssetsByEmail.Data {

							if string(s.Title) == id {

								id = fmt.Sprintf("%s", s.ID)

								break

							}

						}

						goto PREV

					}

				}

			} */

		}
	}

}
