package main

import (
	"context"
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

	//with public policy(stage 2)
	case "public":

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
			} `graphql:"createAsset(data:{sourceID: $SourceID, checked: $Checked})"`
		}

		v := map[string]interface{}{
			"SourceID": graphql.String(sourceID),
			"Checked":  graphql.Boolean(false),
		}

		if err = caller.Mutate(context.Background(), &m, v); err != nil {
			fmt.Fprintf(w, "error with asset source: %v\n", err)
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
			
			<div class="container" id="content" style="color:rgb(255, 255, 255); font-size:30px;">
			
			<br>
			<br>
						
			<h1>content upload:</h1>
			<b>
			<form>
				<input id="picker" type="file" accept="video/*" /><br>
				<p>please wait for content upload completion --> confirm!</p>
				</form>		
				
				<form role="form" method="POST">
							
				<input readonly="true" class="form-control-plaintext" id="ID" aria-label="ID" name ="ID" value="` + i + `" hidden>
				<br>
							
				<p>when file upload done, submit content:</p>
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
					chunkSize: 20480,
				  });
				  
				  upload.on('error', err => {
					console.error('something went wrong', err.detail);
				  });	
				  
				  upload.on('progress', progress => {
					console.log('So far we've uploaded ${progress.detail}');
				  });
				
				  upload.on('success', () => {
					alert('file upload completed.');
				  });
				};
				 
			</script>
			
			<script src="https://assets.medienwerk.now.sh/material.min.js"></script>
			</body>
			</html>
					
			`

	//with signed policy(stage 2)
	case "signed":

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

		car := muxgo.CreateAssetRequest{PlaybackPolicy: []muxgo.PlaybackPolicy{muxgo.SIGNED}}

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
			} `graphql:"createAsset(data:{sourceID: $SourceID, checked: $Checked})"`
		}

		v := map[string]interface{}{
			"SourceID": graphql.String(sourceID),
			"Checked":  graphql.Boolean(false),
		}

		if err = caller.Mutate(context.Background(), &m, v); err != nil {
			fmt.Fprintf(w, "error with asset source: %v\n", err)
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
				
				<div class="container" id="content" style="color:rgb(255, 255, 255); font-size:30px;">
				
				<br>
				<br>
							
				<h1>content upload:</h1>
				<b>
				<form>
				<input id="picker" type="file" accept="video/*" /><br>
				<p>please wait for upload completion --> confirm!</p>
				</form>		
				
				<form role="form" method="POST">
							
				<input readonly="true" class="form-control-plaintext" id="ID" aria-label="ID" name ="ID" value="` + i + `" hidden>
				<br>
							
				<p>when file upload done, submit content:</p>
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
						alert('file upload completed.');
					  });
					};
					 
				</script>
				
				<script src="https://assets.medienwerk.now.sh/material.min.js"></script>
				</body>
				</html>
						
				`

	//policy option (stage 1)
	case "":

		break

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

		client := muxgo.NewAPIClient(
			muxgo.NewConfiguration(
				muxgo.WithBasicAuth(os.Getenv("MUX_ID"), os.Getenv("MUX_SECRET")),
			))

		assets, err := client.AssetsApi.ListAssets()

		if err != nil {

			fmt.Fprintf(w, "%v", err)

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

			<script src="https://cdn.jsdelivr.net/npm/magic-sdk/dist/magic.js"></script>
			</head>
			<body style="background-color: #a1b116;">

			
			<div class="container" id="content" style="color:rgb(255, 255, 255); font-size:30px;">
			
			<br>
			<br>
						
			<h1>file owner:</h1>	
			
			<form role="form" method="POST">
			<label for="Email">Email Address</label>
			<input type="email" class="form-control" placeholder="name@example.com" aria-label="Email" id ="Email" name ="Email" required><br>

			`

		//assigning to db
		if _, err = strconv.Atoi(id); err == nil {

			for _, a := range assets.Data {

				inputInfo, _ := client.AssetsApi.GetAssetInputInfo(a.Id)

				for _, b := range inputInfo.Data {

					url := b.Settings.Url

					url = strings.TrimPrefix(url, "https://storage.googleapis.com/video-storage-us-east1-uploads/")

					sl := strings.SplitN(url, "?", -1)

					var q struct {
						AssetBySourceID struct {
							AssetEntry
						} `graphql:"assetBySourceID(sourceID: $SourceID, checked: $Checked)"`
					}

					v := map[string]interface{}{
						"SourceID": graphql.String(sl[0]),
						"Checked":  graphql.Boolean(false),
					}

					if err := caller.Query(context.Background(), &q, v); err != nil {
						fmt.Fprintf(w, "error with asset source: %v\n", err)
					}

					if q.AssetBySourceID.ID == graphql.ID(id) {

						var m struct {
							UpdateAsset struct {
								AssetEntry
							} `graphql:"updateAsset(id: $ID, data:{assetID: $AssetID, checked: $Checked, policy: $Policy})"`
						}

						v := map[string]interface{}{
							"ID":      q.AssetBySourceID.ID,
							"AssetID": graphql.String(a.Id),
							"Checked": graphql.Boolean(false),
							"Policy":  graphql.String(a.PlaybackIds[0].Policy),
						}

						if err := caller.Mutate(context.Background(), &m, v); err != nil {
							fmt.Fprintf(w, "error with asset ID: %v\n", err)
						}

						pbid = a.PlaybackIds[0].Id

						content = content + `

						<label for="Last">Last Name</label>
						<input class="form-control mr-sm-2" type="text" placeholder="Last" aria-label="Last" id ="Last" name ="Last">
						<label for="First">First Name</label>
						<input class="form-control mr-sm-2" type="text" placeholder="First" aria-label="First" id ="First" name ="First"><br>
						<br>
						<label for="Title">Video Title</label>
						<input class="form-control mr-sm-2" type="text" placeholder="Title" aria-label="Title" id ="Title" name ="Title" required>
						<label for="Category">Video Category</label>
						<input class="form-control mr-sm-2" type="text" placeholder="Category" aria-label="Category" id ="Category" name ="Category" required>
						<label for="Content">description of content</label>
						<textarea class="form-control" id="Content" name ="Content" rows="3"></textarea>
												
						`

						goto NEXT

					}

				}

			}

		} /*  else {

			content = content + `

			<script>

			import { Magic } from 'magic-sdk';

			const magic = new Magic(` + os.Getenv("MAGIC_KEY") + `);

			await magic.auth.loginWithMagicLink({ email: 'your.email@example.com' });

			</script>

			`

		} */

	NEXT:

		content = content + `
			
			<br>
			
			<button type="submit" class="btn btn-light">submit</button>
			
			</form>
			
			</div>
						
			<script src="https://assets.medienwerk.now.sh/material.min.js"></script>
			</body>
			</html>
					
			`

	}

	switch r.Method {

	case "GET":

		switch id {

		//start
		case "":

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
			
			<div class="container" id="content" style="color:rgb(255, 255, 255); font-size:30px;">
			
			<br>
			<br>	
			
			<form role="form" method="POST">

			<input type="checkbox" id="Access" name="Access" value="SIGNED">
 			<label for="Access"><h2>I do not want public content.</h2></label><br>
  			<button type="submit" class="btn btn-light">submit</button>
			
			</form>


			<script src="https://assets.medienwerk.now.sh/material.min.js"></script>
			</body>
			</html>
					
			`

			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("Content-Length", strconv.Itoa(len(content)))
			w.Write([]byte(content))

		//stages
		default:

			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("Content-Length", strconv.Itoa(len(content)))
			w.Write([]byte(content))

		}

	case "POST":

		switch id {

		case "public":

			r.ParseForm()

			id = r.Form.Get("ID")

			//fmt.Fprintf(w, "id: %v\n", i)
			http.Redirect(w, r, "https://"+id+".code2go.dev/content", http.StatusSeeOther)

		case "signed":

			r.ParseForm()

			id = r.Form.Get("ID")

			//fmt.Fprintf(w, "id: %v\n", i)
			http.Redirect(w, r, "https://"+id+".code2go.dev/content", http.StatusSeeOther)

		case "":

			r.ParseForm()

			c := r.Form.Get("Access")

			if c == "" {

				c = "public"

			}

			http.Redirect(w, r, "https://"+c+".code2go.dev/content", http.StatusSeeOther)

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

					http.Redirect(w, r, "https://"+id+".code2go.dev/content", http.StatusSeeOther)

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

						http.Redirect(w, r, "https://"+title+"_"+s+".code2go.dev/content", http.StatusSeeOther)

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

					<div class="container" id="content" style="color:rgb(255, 255, 255); font-size:30px;">

					<br>
					<br>

					`

					if string(m.UpdateAsset.Policy) == "public" {

						content = content +

							`

					<img src="https://image.mux.com/` + string(m.UpdateAsset.PbID) + `/thumbnail.jpg?width=214&height=121&fit_mode=pad">
					<br>

					`

					}

					content = content +

						`
					<p>` + string(m.UpdateAsset.First) + `<br><a href="https://` + id + `.code2go.dev/video">` + string(m.UpdateAsset.Title) + `</a> is "` + string(m.UpdateAsset.Policy) + `" content:<br>` + string(m.UpdateAsset.Content) + `</p>
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

	}

}
