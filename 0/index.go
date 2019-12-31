package main

import (
	"context"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"net/http"
	"os"

	//"sort"
	"strconv"
	"strings"
	"time"

	//"github.com/mschneider82/problem"

	"golang.org/x/oauth2"

	f "github.com/fauna/faunadb-go/faunadb"
	"github.com/shurcooL/graphql"
	//ms "github.com/mitchellh/mapstructure"
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

type Post struct {
	ID graphql.String `graphql:"_id"`
}

var cache []graphql.String = make([]graphql.String, 0)

/* func response(w http.ResponseWriter, success bool, message string, method string) {
	// Create a map for the response body
	body := make(map[string]interface{})

	// Prepare the return data
	if success {
		body["type"] = "connected"
	} else {
		body["type"] = "failure"
	}
	body["message"] = message

	bodyString, _ := json.Marshal(body)

	// Return the response
	if method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "*")
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyString)

} */

func Handler(w http.ResponseWriter, r *http.Request) {

	//var id f.RefV
	url := strings.TrimPrefix(r.URL.Path, "/")

	n, err := strconv.Atoi(url)

	if err != nil {

		fmt.Fprint(w, "... an error occured ... please reload browser window ...")

		return

	}

	now := time.Now().AddDate(0, n, 0)

	var c Cal

	//c.Year = now.Year()

	c.Month = int(now.Month())
	day := map[int]string{now.Day(): now.Weekday().String()}

	c.Days = day

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
		<body style="background-color:#adebad">
		<div class="container">
		<iframe src="https://code2go.dev" style="border:none;"></iframe> 
		</div>
   		<div class="container" id="search" style="color:white; font-size:30px;">
		<form class="form-inline" role="form">
	   	<input class="form-control mr-sm-2" type="text" placeholder="Search" aria-label="Search" id ="find" name ="find">
	   	<button class="btn btn-outline-light my-2 my-sm-1" type="submit">Search</button><br>
		</form>
		</div>
		<br>
		<div class="container" id="nav" style="color:white;">
		` + time.Now().Format("2006") + `
		<br>
		`

	i := 1

	for i < 32 {

		d := now.AddDate(0, 0, i)

		m := int(d.Month())

		if m != c.Month {

			break

		}

		e := d.Day()

		c.Days[e] = d.Weekday().String()

		i++

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

	if n == 0 {

		str = str + `
		<button type="button" class="btn btn-light">` + time.Now().Format("Jan") + `
		 </button>
		 `

	} else {

		str = str + `
		<button type="button" class="btn btn-outline-dark" onclick="window.location.href='0'">` + time.Now().Format("Jan") + `
		 </button>
		 `

	}

	var p, q int

	if now.Month() == time.Now().Month() {

		p = now.Day()

	} else {

		p = 1

	}

	l := len(c.Days)

	for i := l; i >= p; i-- {

		q = i

	}

	c.Year = time.Now().Year()

	t := 1

	for t < 21 {

		y := time.Now().AddDate(0, t, 0).Year()

		if y > c.Year {

			if t == n {

				str = str + `
				<br>
				` + time.Now().AddDate(0, t, 0).Format("2006") + `
				<br>
				<button type="button" class="btn btn-light" onclick="window.location.href='` + strconv.Itoa(t) + `'">
				` + time.Now().AddDate(0, t, 0).Format("Jan") + `
				</button>
				`

			} else {

				str = str + `
				<br>
				` + time.Now().AddDate(0, t, 0).Format("2006") + `
				<br>
				<button type="button" class="btn btn-outline-dark" onclick="window.location.href='` + strconv.Itoa(t) + `'">
				` + time.Now().AddDate(0, t, 0).Format("Jan") + `
				</button>
				`

			}

			c.Year = y

		} else {

			if t == n {

				str = str + `
				
				<button type="button" class="btn btn-light" onclick="window.location.href='` + strconv.Itoa(t) + `'">
				` + time.Now().AddDate(0, t, 0).Format("Jan") + `
				</button>
				`

			} else {

				str = str + `
				
				<button type="button" class="btn btn-outline-dark" onclick="window.location.href='` + strconv.Itoa(t) + `'">
				` + time.Now().AddDate(0, t, 0).Format("Jan") + `
				</button>
				`

			}

		}

		/* 			str = str + `
				<br>
				` + time.Now().AddDate(0, t, 0).Format("2006") + `
				<br>
				<button type="button" class="btn btn-outline-dark" onclick="window.location.href='` + strconv.Itoa(t) + `'">
				` + time.Now().AddDate(0, t, 0).Format("Jan") + `
				</button>
				`
			c.Year = time.Now().AddDate(0, t, 0).Year()
		} else {
			str = str + `
				<button type="button" class="btn btn-outline-dark" onclick="window.location.href='` + strconv.Itoa(t) + `'">
				` + time.Now().AddDate(0, t, 0).Format("Jan") + `
				</button>
				`
		} */
				
		t++

	}

	c.Year = now.Year()

	str = str + `
		<br>
		</div>
		<br>
		<div class="container" id="data" style="color:white;">
		<form class="form-inline" role="form"  method="post">
		<ul class="list-group">
		`

	switch c.Days[q] {

	case "Monday":
		break
	case "Tuesday":
		str = str + `
			<br>
			<button type="button" class="btn btn-link" onclick="window.location.href='Monday'">
			<span class="badge badge-pill badge-dark">
			Monday
			</span>
			</button>
			`
	case "Wednesday":
		str = str + `
			<br>
			<button type="button" class="btn btn-link" onclick="window.location.href='Monday'">
			<span class="badge badge-pill badge-dark">
			Monday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Tuesday'">
			<span class="badge badge-pill badge-dark">
			Tuesday
			</span>
			</button>
			`
	case "Thursday":
		str = str + `
			<br>
			<button type="button" class="btn btn-link" onclick="window.location.href='Monday'">
			<span class="badge badge-pill badge-dark">
			Monday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Tuesday'">
			<span class="badge badge-pill badge-dark">
			Tuesday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Wednesday'">
			<span class="badge badge-pill badge-dark">
			Wednesday
			</span>
			</button>
			`
	case "Friday":
		str = str + `
			<br>
			<button type="button" class="btn btn-link" onclick="window.location.href='Monday'">
			<span class="badge badge-pill badge-dark">
			Monday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Tuesday'">
			<span class="badge badge-pill badge-dark">
			Tuesday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Wednesday'">
			<span class="badge badge-pill badge-dark">
			Wednesday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Thursday'">
			<span class="badge badge-pill badge-dark">
			Thursday
			</span>
			</button>
			`
	case "Saturday":
		str = str + `
			<br>
			<button type="button" class="btn btn-link" onclick="window.location.href='Monday'">
			<span class="badge badge-pill badge-dark">
			Monday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Tuesday'">
			<span class="badge badge-pill badge-dark">
			Tuesday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Wednesday'">
			<span class="badge badge-pill badge-dark">
			Wednesday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Thursday'">
			<span class="badge badge-pill badge-dark">
			Thursday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Friday'">
			<span class="badge badge-pill badge-dark">
			Friday
			</span>
			</button>
			`
	case "Sunday":
		str = str + `
			<br>
			<button type="button" class="btn btn-link" onclick="window.location.href='Monday'">
			<span class="badge badge-pill badge-dark">
			Monday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Tuesday'">
			<span class="badge badge-pill badge-dark">
			Tuesday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Wednesday'">
			<span class="badge badge-pill badge-dark">
			Wednesday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Thursday'">
			<span class="badge badge-pill badge-dark">
			Thursday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Friday'">
			<span class="badge badge-pill badge-dark">
			Friday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Saturday'">
			<span class="badge badge-pill badge-dark">
			Saturday
			</span>
			</button>
			`
	}

	fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

	x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database(now.Format("2006")), "role": "server"}))

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

	value := now.Format("2006-01-02")

	mo := fmt.Sprintf("%02d", c.Month)

	ye := strconv.Itoa(c.Year)

	v1 := make(map[string]interface{})

	for k := q; k < 32; k++ {

		value = ye + `-` + mo + `-` + fmt.Sprintf("%02d", k)

		var q struct {
			PostsByDate struct {
				Data []Post
			} `graphql:"postsByDate(date: $date iscommited: true)"`
		}

		v1["date"] = graphql.String(value)

		switch c.Days[k] {

		case "Monday":

			str = str + `
				<br>
				<button type="button" class="btn btn-link" onclick="window.location.href='` + c.Days[k] + `'">
				<span class="badge badge-pill badge-dark">
				` + c.Days[k] + `
				</span>
				</button>
				<button type="button" class="btn btn-light btn-link" onclick="window.location.href='entry#` + value + `'">
				<span class="badge badge-pill badge-light">
				<input readonly class="form-control-plaintext list-group-item-action" value="` + strconv.Itoa(k) + `" >
				</span>
				</button>
				`

			if err = call.Query(context.Background(), &q, v1); err != nil {
				fmt.Fprint(w, err)
			}

			l := len(q.PostsByDate.Data)

			if l > 0 {

				for _, p := range q.PostsByDate.Data {

					cache = append(cache, p.ID)

				}

				str = str + `
			<span style="text-align: inherit; color: #70db70" class="badge badge-pill badge-dark">
			` + strconv.Itoa(l) + `
			</span>
			</button>
			`

			} else {

				str = str + `
			</button>
			`
				break

			}

			/* s := `{"query":"query{` + dir + `(date:\"` + value + `\" iscommited: true){data{_id}}}"}`
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

			if i == nil {

				str = str + `
				</button>
				`

				break

			}

			a := i.(map[string]interface{})

			b := a["data"]

			if b == nil {

				str = str + `
				</button>
				`

				break

			}

			c := b.(map[string]interface{})

			d := c[dir]

			if d == nil {

				str = str + `
				</button>
				`

				break

			}

			e := d.(map[string]interface{})

			f := e["data"]

			if f == nil {

				str = str + `
				</button>
				`

				break

			}

			g := f.([]interface{})

			if g == nil {

				str = str + `
				</button>
				`
				break

			} else {

				h := make([]map[string]interface{}, len(g))

				for j := 0; j < len(g); j++ {

					h[j] = g[j].(map[string]interface{})

				}

				k := make(map[int]string, 0)

				for j := 0; j < len(g); j++ {

					k[j] = h[j]["_id"].(string)

					if k[j] != "" {

						cache = append(cache, k[j])

					}

				}

				if len(k) > 0 {

					str = str + `
				<span style="text-align: inherit; color: #70db70" class="badge badge-pill badge-dark">
				` + strconv.Itoa(len(k)) + `
				</span>
				</button>
				`

				} else {

					str = str + `
				</button>
				`
					break

				}

			} */

			/* if string(b[k-q]) != "0" {
				str = str + `
				<span class="badge badge-dark">
				` + string(b[k-q]) + `
				</span>
				</button>
				`
			} else {
				str = str + `
				</button>
				`
			} */

		case "Tuesday":

			str = str + `
				<br>
				<button type="button" class="btn btn-link" onclick="window.location.href='` + c.Days[k] + `'">
				<span class="badge badge-pill badge-dark">
				` + c.Days[k] + `
				</span>
				</button>
				<button type="button" class="btn btn-light btn-link" onclick="window.location.href='entry#` + value + `'">
				<span class="badge badge-pill badge-light">
				<input readonly class="form-control-plaintext list-group-item-action" value="` + strconv.Itoa(k) + `" >
				</span>
				</button>
				`

			if err = call.Query(context.Background(), &q, v1); err != nil {
				fmt.Fprint(w, err)
			}

			l := len(q.PostsByDate.Data)

			if l > 0 {

				for _, p := range q.PostsByDate.Data {

					cache = append(cache, p.ID)

				}

				str = str + `
			<span style="text-align: inherit; color: #70db70" class="badge badge-pill badge-dark">
			` + strconv.Itoa(l) + `
			</span>
			</button>
			`

			} else {

				str = str + `
			</button>
			`
				break

			}

		case "Wednesday":

			str = str + `
				<br>
				<button type="button" class="btn btn-link" onclick="window.location.href='` + c.Days[k] + `'">
				<span class="badge badge-pill badge-dark">
				` + c.Days[k] + `
				</span>
				</button>
				<button type="button" class="btn btn-light btn-link" onclick="window.location.href='entry#` + value + `'">
				<span class="badge badge-pill badge-light">
				<input readonly class="form-control-plaintext list-group-item-action" value="` + strconv.Itoa(k) + `" >
				</span>
				</button>
				`

			if err = call.Query(context.Background(), &q, v1); err != nil {
				fmt.Fprint(w, err)
			}

			l := len(q.PostsByDate.Data)

			if l > 0 {

				for _, p := range q.PostsByDate.Data {

					cache = append(cache, p.ID)

				}

				str = str + `
			<span style="text-align: inherit; color: #70db70" class="badge badge-pill badge-dark">
			` + strconv.Itoa(l) + `
			</span>
			</button>
			`

			} else {

				str = str + `
			</button>
			`
				break

			}

		case "Thursday":

			str = str + `
				<br>
				<button type="button" class="btn btn-link" onclick="window.location.href='` + c.Days[k] + `'">
				<span class="badge badge-pill badge-dark">
				` + c.Days[k] + `
				</span>
				</button>
				<button type="button" class="btn btn-light btn-link" onclick="window.location.href='entry#` + value + `'">
				<span class="badge badge-pill badge-light">
				<input readonly class="form-control-plaintext list-group-item-action" value="` + strconv.Itoa(k) + `" >
				</span>
				</button>
				`

			if err = call.Query(context.Background(), &q, v1); err != nil {
				fmt.Fprint(w, err)
			}

			l := len(q.PostsByDate.Data)

			if l > 0 {

				for _, p := range q.PostsByDate.Data {

					cache = append(cache, p.ID)

				}

				str = str + `
			<span style="text-align: inherit; color: #70db70" class="badge badge-pill badge-dark">
			` + strconv.Itoa(l) + `
			</span>
			</button>
			`

			} else {

				str = str + `
			</button>
			`
				break

			}

		case "Friday":

			str = str + `
				<br>
				<button type="button" class="btn btn-link" onclick="window.location.href='` + c.Days[k] + `'">
				<span class="badge badge-pill badge-dark">
				` + c.Days[k] + `
				</span>
				</button>
				<button type="button" class="btn btn-light btn-link" onclick="window.location.href='entry#` + value + `'">
				<span class="badge badge-pill badge-light">
				<input readonly class="form-control-plaintext list-group-item-action" value="` + strconv.Itoa(k) + `" >
				</span>
				</button>
				`

			if err = call.Query(context.Background(), &q, v1); err != nil {
				fmt.Fprint(w, err)
			}

			l := len(q.PostsByDate.Data)

			if l > 0 {

				for _, p := range q.PostsByDate.Data {

					cache = append(cache, p.ID)

				}

				str = str + `
			<span style="text-align: inherit; color: #70db70" class="badge badge-pill badge-dark">
			` + strconv.Itoa(l) + `
			</span>
			</button>
			`

			} else {

				str = str + `
			</button>
			`
				break

			}

		case "Saturday":

			str = str + `
				<br>
				<button type="button" class="btn btn-link" onclick="window.location.href='` + c.Days[k] + `'">
				<span class="badge badge-pill badge-dark">
				` + c.Days[k] + `
				</span>
				</button>
				<button type="button" class="btn btn-light btn-link" onclick="window.location.href='entry#` + value + `'">
				<span class="badge badge-pill badge-light">
				<input readonly class="form-control-plaintext list-group-item-action" value="` + strconv.Itoa(k) + `" >
				</span>
				</button>
				`

			if err = call.Query(context.Background(), &q, v1); err != nil {
				fmt.Fprint(w, err)
			}

			l := len(q.PostsByDate.Data)

			if l > 0 {

				for _, p := range q.PostsByDate.Data {

					cache = append(cache, p.ID)

				}

				str = str + `
			<span style="text-align: inherit; color: #70db70" class="badge badge-pill badge-dark">
			` + strconv.Itoa(l) + `
			</span>
			</button>
			`

			} else {

				str = str + `
			</button>
			`
				break

			}

		case "Sunday":

			str = str + `
				<br>
				<button type="button" class="btn btn-link" onclick="window.location.href='` + c.Days[k] + `'">
				<span class="badge badge-pill badge-dark">
				` + c.Days[k] + `
				</span>
				</button>
				<button type="button" class="btn btn-light btn-link" onclick="window.location.href='entry#` + value + `'">
				<span class="badge badge-pill badge-light">
				<input readonly class="form-control-plaintext list-group-item-action" value="` + strconv.Itoa(k) + `" >
				</span>
				</button>
				`

			if err = call.Query(context.Background(), &q, v1); err != nil {
				fmt.Fprint(w, err)
			}

			l := len(q.PostsByDate.Data)

			if l > 0 {

				for _, p := range q.PostsByDate.Data {

					cache = append(cache, p.ID)

				}

				str = str + `
			<span style="text-align: inherit; color: #70db70" class="badge badge-pill badge-dark">
			` + strconv.Itoa(l) + `
			</span>
			</button>
			`

			} else {

				str = str + `
			</button>
			`
				break

			}

		}

	}

	if cache != nil {

		var q struct {
			CacheByMonth struct {
				ID    graphql.ID       `graphql:"_id"`
				Month graphql.String   `graphql:"month"`
				Posts []graphql.String `graphql:"posts"`
			} `graphql:"cacheByMonth(month: $month)"`
		}

		v2 := map[string]interface{}{
			"month": graphql.String(ye + `-` + mo),
		}

		if err = call.Query(context.Background(), &q, v2); err != nil {
			fmt.Fprintf(w, "get cache error: %v", err)
		}

		result := q.CacheByMonth

		var posts []graphql.String = make([]graphql.String, 0)

		if result.Posts == nil {

			var m struct {
				CreateCache struct {
					ID    graphql.ID       `graphql:"_id"`
					Month graphql.String   `graphql:"month"`
					Posts []graphql.String `graphql:"posts"`
				} `graphql:"createCache(data:{month: $month, posts: $posts})"`
			}

			for _, c := range cache {

				posts = append(posts, c)
			}

			v3 := map[string]interface{}{
				"month": v2["month"],
				"posts": posts,
			}

			if err = call.Mutate(context.Background(), &m, v3); err != nil {
				fmt.Fprintf(w, "create cache error: %v", err)
			}

		} else {

			var m struct {
				UpdateCache struct {
					ID    graphql.ID       `graphql:"_id"`
					Month graphql.String   `graphql:"month"`
					Posts []graphql.String `graphql:"posts"`
				} `graphql:"updateCache(id: $id, data:{month: $month, posts: $posts})"`
			}

			for _, c := range cache {

				posts = append(posts, c)
			}

			//posts = append(posts, "253012617168159243")

			v4 := map[string]interface{}{
				"id":    result.ID,
				"month": v2["month"],
				"posts": posts,
			}

			if err = call.Mutate(context.Background(), &m, v4); err != nil {
				fmt.Fprintf(w, "update cache error %v", err)
			}

		}

	}

	switch c.Days[len(c.Days)] {

	case "Monday":
		str = str + `
			<button type="button" class="btn btn-link" onclick="window.location.href='Tuesday'">
			<span class="badge badge-pill badge-dark">
			Tuesday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Wednesday'">
			<span class="badge badge-pill badge-dark">
			Wednesday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Thursday'">
			<span class="badge badge-pill badge-dark">
			Thursday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Friday'">
			<span class="badge badge-pill badge-dark">
			Friday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Saturday'">
			<span class="badge badge-pill badge-dark">
			Saturday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Sunday'">
			<span class="badge badge-pill badge-dark">
			Sunday
			</span>
			</button>
			`

	case "Tuesday":
		str = str + `
			<button type="button" class="btn btn-link" onclick="window.location.href='Wednesday'">
			<span class="badge badge-pill badge-dark">
			Wednesday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Thursday'">
			<span class="badge badge-pill badge-dark">
			Thursday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Friday'">
			<span class="badge badge-pill badge-dark">
			Friday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Saturday'">
			<span class="badge badge-pill badge-dark">
			Saturday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Sunday'">
			<span class="badge badge-pill badge-dark">
			Sunday
			</span>
			</button>
			`

	case "Wednesday":
		str = str + `
			<button type="button" class="btn btn-link" onclick="window.location.href='Thursday'">
			<span class="badge badge-pill badge-dark">
			Thursday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Friday'">
			<span class="badge badge-pill badge-dark">
			Friday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Saturday'">
			<span class="badge badge-pill badge-dark">
			Saturday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Sunday'">
			<span class="badge badge-pill badge-dark">
			Sunday
			</span>
			</button>
			`

	case "Thursday":
		str = str + `
			<button type="button" class="btn btn-link" onclick="window.location.href='Friday'">
			<span class="badge badge-pill badge-dark">
			Friday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Saturday'">
			<span class="badge badge-pill badge-dark">
			Saturday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Sunday'">
			<span class="badge badge-pill badge-dark">
			Sunday
			</span>
			</button>
			`

	case "Friday":
		str = str + `
			<button type="button" class="btn btn-link" onclick="window.location.href='Saturday'">
			<span class="badge badge-pill badge-dark">
			Saturday
			</span>
			</button>
			<button type="button" class="btn btn-link" onclick="window.location.href='Sunday'">
			<span class="badge badge-pill badge-dark">
			Sunday
			</span>
			</button>
			`

	case "Saturday":
		str = str + `
			<button type="button" class="btn btn-link" onclick="window.location.href='Sunday'">
			<span class="badge badge-pill badge-dark">
			Sunday
			</span>
			</button>
			`

	case "Sunday":
		str = str + `
			</ul>
			</form>	
			</div>
			<br>
			`

		break

	}

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
