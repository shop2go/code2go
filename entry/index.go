package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	//"strings"
	"time"

	"golang.org/x/oauth2"

	f "github.com/fauna/faunadb-go/faunadb"
	"github.com/shurcooL/graphql"
)

type Cal struct {
	Year  int
	Month int
	Days  map[int]string
}

type Access struct {
	//Reference *f.RefV `fauna:"ref"`
	Timestamp int    `fauna:"ts"`
	Secret    string `fauna:"secret"`
	Role      string `fauna:"role"`
}

type Cache struct {
	ID    graphql.String   `graphql:"_id"`
	Month graphql.String   `graphql:"month"`
	Posts []graphql.String `graphql:"posts"`
}

type Post struct {
	ID         graphql.String   `graphql:"_id"`
	Date       graphql.String   `graphql:"date"`
	Iscommited graphql.Boolean  `graphql:"iscommited`
	Salt       graphql.String   `graphql:"salt`
	Tags       []graphql.String `graphql:"tags`
	Topics     []graphql.String `graphql:"topics`
	Content    graphql.String   `graphql:"content`
	Isparent   []graphql.String `graphql:"isparent`
}

/* func getCache(a *Access) ([]Cache, error) {

	var result []Cache = make([]Cache, 0)

	dir := "allCaches"

	s := `{"query":"query{` + dir + `{data{month posts{_id date password topics tags content}}}}"}`
	body := strings.NewReader(s)
	req, _ := http.NewRequest("POST", "https://graphql.fauna.com/graphql", body)

	req.Header.Set("Authorization", "Bearer "+a.Secret)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Schema-Preview", "partial-update-mutation")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {

		return nil, err

	}

	defer resp.Body.Close()

	bdy, _ := ioutil.ReadAll(resp.Body)

	var i interface{}

	json.Unmarshal(bdy, &i)

	if i == nil {

		return nil, errOnData

	}

	x := i.(map[string]interface{})

	b := x["data"]

	if b == nil {

		return nil, errOnData

	}

	c := b.(map[string]interface{})

	d := c[dir]

	if d == nil {

		return nil, errOnData

	}

	e := d.(map[string]interface{})

	f := e["data"]

	if f == nil {

		return nil, errOnData

	}

	g := f.([]interface{})

	if g == nil {

		return nil, errOnData

	} else {

		l := len(g)

		h := make([]map[string]interface{}, l)

		cache := make([]Cache, l)

		for j := 0; j < l; j++ {

			//h == Caches
			h[j] = g[j].(map[string]interface{})

			if h[j] != nil {

				cache[j].Month = h[j]["month"].(string)

			}

		}

		if h != nil {

			posts := make([]interface{}, l)

			for j := 0; j < l; j++ {

				posts[j] = h[j]["posts"].([]interface{})

				o := posts[j].([]interface{})

				for k := 0; k < len(o); k++ {

					//q :=  o[k].(interface{})
					if o[k] == nil {

						return nil, errOnData


					}

					p := o[k].(map[string]interface{})

					resultP := make([]Post, len(p))

					resultP[k].ID = p["_id"].(string)

					resultP[k].Date = p["date"].(string)

					resultP[k].Password = p["password"].(string)

					resultP[k].Topics = p["topics"]

					resultP[k].Tags = p["tags"]

					resultP[k].Content = p["content"]

					cache[j].Posts = append(cache[j].Posts, resultP[k])

				}

				result = append(result, cache[j])

			}

		}

	}

	return result, nil

} */

var hits map[string][]graphql.String = make(map[string][]graphql.String, 21)

func Handler(w http.ResponseWriter, r *http.Request) {

	var c Cal

	now := time.Now()

	c.Year = now.Year()
	c.Month = int(now.Month())
	day := map[int]string{now.Day(): now.Weekday().String()}

	c.Days = day

	var years []string = []string{now.Format("2006")}

	i := 0

	for i < 21 {

		i++

		loc, _ := time.LoadLocation("")

		t := time.Date(c.Year, now.AddDate(0, i, 0).Month(), 1, 0, 0, 0, 0, loc).Format("2006")

		//t := now.AddDate(0, i, 0).Format("2006")

		if t != time.Date(c.Year, now.AddDate(0, i-1, 0).Month(), 1, 0, 0, 0, 0, loc).Format("2006") {

			years = append(years, t)

		}

	}

	z := 1

	for z < 32 {

		d := now.AddDate(0, 0, z)

		m := int(d.Month())

		if m != c.Month {

			break

		}

		e := d.Day()

		c.Days[e] = d.Weekday().String()

		z++

	}

	j := 1

	for j < 32 {

		d := now.AddDate(0, 0, -j)

		m := int(d.Month())

		if m != c.Month {

			break

		}

		e := d.Day()

		c.Days[e] = d.Weekday().String()

		j++

	}

	var q int

	l := len(c.Days)

	p := now.Day()

	for i := l; i >= p; i-- {

		q = i

	}

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

	<ul class="list-group">
	`

	fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

	l = len(years)

	fx := make(map[string]f.Value, l)

	//sort.Slice(years, func(i, j int) bool { return years[i] < years[j] })

	for i := 0; i < l; i++ {

		x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database(years[i]), "role": "server-readonly"}))

		if err != nil {

			fmt.Fprint(w, err)

			return

		}

		fx[years[i]] = x

		var access *Access

		if err = x.Get(&access); err != nil {

			fmt.Fprint(w, err)

		}

		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: access.Secret},
		)

		httpClient := oauth2.NewClient(context.Background(), src)

		call := graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

		var query struct {
			AllCaches struct {
				Data []Cache
			}
		}

		if err = call.Query(context.Background(), &query, nil); err != nil {
			fmt.Fprintf(w, "get cache error: %v\n", err)
		}

		result := query.AllCaches.Data

		if result != nil {

			for _, v := range result {

				hits[string(v.Month)] = v.Posts

			}

		}

	}

	var result Post

	//expose the anchor of specified date++; list apropriate entries for that date whithin the actual month from persitence layer

	for k := q; k <= l; k++ {

		m := fmt.Sprintf("%02d", c.Month)

		n := fmt.Sprintf("%02d", k)

		schedule := strconv.Itoa(c.Year) + `-` + m + `-` + n

		str = str + `
			<br>
			<button type="button" class="btn btn-link" onclick="window.location.href='` + c.Days[k] + `'">
			<span class="badge badge-pill badge-dark">
			` + c.Days[k] + `
			</span>
			</button>
			<button type="button" class="btn btn-light">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" id="` + schedule + `" value="` + schedule + `" placeholder="` + schedule + `">
			</span><button><br>

			<div class="container" id="threads">
			<form class="form-inline" role="form">
			<input readonly class="form-control-plaintext list-group-item-action" id="thread` + schedule + `" value="new thread" placeholder="new thread" onclick="window.location.href='https://` + schedule + `.code2go.dev/new'">
						
			`

		if v, ok := hits[strconv.Itoa(c.Year) + `-` + m]; ok {

			//sort.Slice(v, func(i, j int) bool { return v[i] > v[j] })

			x := fx[strconv.Itoa(c.Year)]

			var access *Access

			if err := x.Get(&access); err != nil {

				fmt.Fprint(w, err)
	
			}

			src := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: access.Secret},
			)

			httpClient := oauth2.NewClient(context.Background(), src)

			call := graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

			for _, postID := range v {

				var q struct {
					FindPostByID struct {
						Data Post
					} `graphql:"postsByDate(id: $id)"`
				}

				v1 := map[string]interface{}{
					"id": graphql.String(postID),
				}

				if err := call.Query(context.Background(), &q, v1); err != nil {
					fmt.Fprintf(w, "get post error: %v\n", err)
				}

				result = q.FindPostByID.Data

				if string(result.Salt) == "" {

					var s string

					for _, v := range result.Topics {

						s = s + string(v)
					}

					str = str + `
							<input readonly="true" class="form-control-plaintext list-group-item-action" id="` + string(result.ID) + `" value="` + s + `" onclick="window.location.href='https://` + string(result.ID) + `.code2go.dev/status'">
							`

				}

			}

		}

	}

	o := 0

LOOP:

	for o < 21 {

		o++

		now = time.Now().AddDate(0, o, 0)

		c.Year = now.Year()

		c.Month = int(now.Month())

		for k := 1; k <= 31; k++ {

			m := fmt.Sprintf("%02d", c.Month)

			n := fmt.Sprintf("%02d", k)

			loc, _ := time.LoadLocation("")

			date := time.Date(c.Year, time.Now().AddDate(0, o, 0).Month(), k, 0, 0, 0, 0, loc).Month()

			if c.Month != int(date) {

				goto LOOP

			} else {

				schedule := strconv.Itoa(c.Year) + `-` + m + `-` + n

				str = str + `
				<br>
				<button type="button" class="btn btn-link" onclick="window.location.href='` + c.Days[k] + `'">
				<span class="badge badge-pill badge-dark">
				` + c.Days[k] + `
				</span>
				</button>
				<button type="button" class="btn btn-light">
				<span class="badge badge-pill badge-light">
				<input readonly class="form-control-plaintext list-group-item-action" id="` + schedule + `" value="` + schedule + `" placeholder="` + schedule + `">
				</span><button><br>
	
				<div class="container" id="threads">
				<form class="form-inline" role="form">
				<input readonly class="form-control-plaintext list-group-item-action" id="thread` + schedule + `" value="new thread" placeholder="new thread" onclick="window.location.href='https://` + schedule + `.code2go.dev/new'">
							
				`

				if v, ok := hits[strconv.Itoa(c.Year) + `-` + m]; ok {

					sort.Slice(v, func(i, j int) bool { return v[i] > v[j] })

					x := fx[strconv.Itoa(c.Year)]

					var access *Access

					if err := x.Get(&access); err != nil {

						fmt.Fprint(w, err)
			
					}

					src := oauth2.StaticTokenSource(
						&oauth2.Token{AccessToken: access.Secret},
					)

					httpClient := oauth2.NewClient(context.Background(), src)

					call := graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

					for _, postID := range v {

						var q struct {
							FindPostByID struct {
								Data Post
							} `graphql:"postsByDate(id: $id)"`
						}

						v1 := map[string]interface{}{
							"id": graphql.String(postID),
						}

						if err := call.Query(context.Background(), &q, v1); err != nil {
							fmt.Fprintf(w, "get post error: %v\n", err)
						}

						result = q.FindPostByID.Data

						if string(result.Salt) == "" {

							var s string

							for _, v := range result.Topics {

								s = s + string(v)
							}

							str = str + `
								<input readonly="true" class="form-control-plaintext list-group-item-action" id="` + string(result.ID) + `" value="` + s + `" onclick="window.location.href='https://` + string(result.ID) + `.code2go.dev/status'">
								`

						}

					}

				}

			}

		}

	}

	str = str + `
	
		</form>
		</div>
	
 		</ul>
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
