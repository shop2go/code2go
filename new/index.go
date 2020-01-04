package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/oauth2"

	f "github.com/fauna/faunadb-go/faunadb"
	"github.com/shurcooL/graphql"
)

type Access struct {
	//Reference *f.RefV `fauna:"ref"`
	Timestamp int    `fauna:"ts"`
	Secret    string `fauna:"secret"`
	Role      string `fauna:"role"`
}

/* type Post struct {
	ID         graphql.String   `graphql:"_id"`
	Date       graphql.String   `graphql:"date"`
	Iscommited graphql.Boolean  `graphql:"iscommited`
	Salt       graphql.String   `graphql:"salt`
	Tags       []graphql.String `graphql:"tags`
	Topics     []graphql.String `graphql:"topics`
	Content    graphql.String   `graphql:"content`
	Isparent   []graphql.String `graphql:"isparent`
	Ischild    graphql.String   `graphql:"ischild`
} */

func Handler(w http.ResponseWriter, r *http.Request) {

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
		<!-- CSS -->
		<!-- Add Material font (Roboto) and Material icon as needed -->
		<link href="https://fonts.googleapis.com/css?family=Roboto:300,300i,400,400i,500,500i,700,700i|Roboto+Mono:300,400,700|Roboto+Slab:300,400,700" rel="stylesheet">
		<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">

		<!-- Add Material CSS, replace Bootstrap CSS -->
		<link href="https://assets.medienwerk.now.sh/material.min.css" rel="stylesheet">
		</head>

		<body style="background-color:#adebad">

		<div class="container" id="data" style="color:white;">
		<form class="form-inline" role="form" method="POST">
		<input type="email" class="form-control" placeholder="name@example.com" aria-label="Email" id ="Email" name ="Email">
		<br>
		<input readonly="true" class="form-control-plaintext" id="Schedule" aria-label="Schedule" name ="Schedule" value="` + strings.TrimSuffix(r.Host, ".code2go.dev") + `">
		<input class="form-control mr-sm-2" type="text" placeholder="Secret" aria-label="Secret" id ="Secret" name ="Secret" value="">
		<input class="form-control mr-sm-2" type="text" placeholder="Title" aria-label="Title" id ="Title" name ="Title" required>
		<!--input class="form-control mr-sm-2" type="text" placeholder="Entry" aria-label="Entry" id ="Entry" name ="Entry" required-->
		<input class="form-control mr-sm-2" type="text" placeholder="Tags" aria-label="Tags" id ="Tags" name ="Tags">
		<input class="form-control mr-sm-2" tyoe="text" aria-label="Content" id ="Content" name ="Content" placeholder="Content"></textarea>
		<br>
		<button type="submit" class="btn btn-light">submit</button>
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

	case "POST":

		r.ParseForm()

		//email := r.FormValue("Email")
		pw := r.FormValue("Secret")
		date := r.FormValue("Schedule")
		topics := r.FormValue("Title")
		tags := r.FormValue("Tags")
		content := r.FormValue("Entry")

		sl := strings.SplitN(date, "-", -1)

		fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

		x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database(sl[0]), "role": "server"}))

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

		var m1 struct {
			CreatePost struct {
				ID         graphql.String   `graphql:"_id"`
				Date       graphql.String   `graphql:"date"`
				Iscommited graphql.Boolean  `graphql:"iscommited`
				Salt       graphql.String   `graphql:"salt`
				Tags       []graphql.String `graphql:"tags`
				Topics     []graphql.String `graphql:"topics`
				Content    graphql.String   `graphql:"content`
				Ischild    graphql.String   `graphql:"ischild`
			} `graphql:"createPost(data:{date: $date, iscommited: false, salt: $salt, tags: $tags, topics: $topics, content: $content, ischild: $ischild})"`
		}

		qsl := make([]graphql.String, 0)

		v1 := map[string]interface{}{
			"date":    graphql.String(date),
			"salt":    graphql.String(pw),
			"content": graphql.String(content),
			"ischild": graphql.String(""),
		}

		topic := strings.SplitN(topics, " ", -1)

		for _, v := range topic {

			qsl = append(qsl, graphql.String(v))
		}

		v1["topics"] = qsl

		qsl = nil

		tag := strings.SplitN(tags, " ", -1)

		for _, v := range tag {

			qsl = append(qsl, graphql.String(strings.ToLower(v)))
		}

		v1["tags"] = qsl

		qsl = nil

		if err = call.Mutate(context.Background(), &m1, v1); err != nil {
			fmt.Fprintf(w, "create Post error: %v\n", err)
		}

		id := m1.CreatePost.ID /*  + graphql.String(":") + graphql.String(date) */

		x, err = fc.Query(f.CreateKey(f.Obj{"database": f.Database("users"), "role": "server"}))

		if err != nil {

			fmt.Fprint(w, err)

			return

		}

		x.Get(&access)

		src = oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: access.Secret},
		)

		httpClient = oauth2.NewClient(context.Background(), src)

		call = graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

		var m2 struct {
			CreatePost struct {
				ID   graphql.String `graphql:"_id"`
				Post graphql.String `graphql:"post"`
			} `graphql:"createPost(data:{post: $post})"`
		}

		v2 := map[string]interface{}{
			"post": id,
		}

		if err = call.Mutate(context.Background(), &m2, v2); err != nil {
			fmt.Fprintf(w, "create Post error: %v\n", err)
		}

		http.Redirect(w, r, "https://"+string(id)+"_notAproved.code2go.dev/status", http.StatusFound)

		/*
			//to := strings.ReplaceAll(topics, " ", "\", \"")

			//to = "\"" + to + "\""

			tags = strings.ToLower(tags)

			//ta := strings.ReplaceAll(tags, " ", "\", \"")

			//ta = "\"" + ta + "\""

			dir := "createPost"

			s := `{"query":"mutation{` + dir + `(data:{iscommited: false password: \"` + pw + `\" date: \"` + date + `\" topics: \"` + topics + `\" content: \"` + content + `\" tags: \"` + tags + `\"}) {_id}}"}`

			body := strings.NewReader(s)
			req, _ := http.NewRequest("POST", "https://graphql.fauna.com/graphql", body)

			req.Header.Set("Authorization", "Bearer "+access.Secret)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
			req.Header.Set("X-Schema-Preview", "partial-update-mutation")

			resp, err := http.DefaultClient.Do(req)

			if err != nil {

				fmt.Fprint(w, err)

				return

			}

			defer resp.Body.Close()

			bdy, _ := ioutil.ReadAll(resp.Body)

			var i interface{}

			json.Unmarshal(bdy, &i)

			if i != nil {

				a := i.(map[string]interface{})

				b := a["data"]

				if b != nil {

					c := b.(map[string]interface{})

					d := c[dir]

					if d != nil {

						e := d.(map[string]interface{})

						f := e["_id"]

						id := f.(string)

						if id != "" {

							fmt.Fprint(w, sl[0]+"_"+id)

							http.Redirect(w, r, "https://"+id+".code2go.dev/status", 301)

						} else {

							fmt.Fprint(w, errOnData)

						}

					}

				}

			} else {

				fmt.Fprint(w, errOnData)

			}
		*/
	}

}
