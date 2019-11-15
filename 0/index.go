package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	//"github.com/mschneider82/problem"

	f "github.com/fauna/faunadb-go/faunadb"
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

func response(w http.ResponseWriter, success bool, message string, method string) {
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

}

func Handler(w http.ResponseWriter, r *http.Request) {

	//var id f.RefV
	url := strings.TrimPrefix(r.URL.Path, "/")

	n, err := strconv.Atoi(url)

	if err != nil {

		fmt.Fprint(w, "... an error occured ... please refresh browser window ...")

		return

	}

	now := time.Now().AddDate(0, n, 0)

	/* 	if r.Method == http.MethodOptions {
		response(w, true, "", r.Method)
		return
	} */

	var c Cal

	c.Year = now.Year()

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

	month, _ := strconv.Atoi(now.Format("01"))
	c.Month = month
	day := map[int]string{now.Day(): now.Weekday().String()}

	c.Days = day

	i := 1

	for i < 32 {

		d := now.AddDate(0, 0, i)

		//m, _ := strconv.Atoi(d.Format("01"))

		if d.Month() != now.Month() {

			break

		}

		e, _ := strconv.Atoi(d.Format("02"))

		c.Days[e] = d.Weekday().String()

		i++

	}

	j := 1

	for j > 0 {

		d := now.AddDate(0, 0, -j)

		//m, _ := strconv.Atoi(d.Format("01"))

		if d.Month() != now.Month() {

			break

		}

		e, _ := strconv.Atoi(d.Format("02"))

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

	l := len(c.Days)

	if now.Month() == time.Now().Month() {

		p, _ = strconv.Atoi(time.Now().Format("02"))

	} else {

		p = 1

	}

	for i := l; i >= p; i-- {

		q = i

	}

	t := 0

	for t < 20 {

		t++

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

		c.Year = y

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

	x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database(now.Format("2006")), "role": "server-readonly"}))

	if err != nil {

		fmt.Fprint(w, "... an error occured ... please refresh browser window ...")

		return

	}

	var access *Access

	x.Get(&access)
	//TODO:make cert client

	/* 		resp, err := http.Get("https://"+ip+"/" + url)
	   		if err != nil {
	   			problem.New(problem.Type("https://"+ip+"/404"), problem.Status(404)).WriteTo(w)
	   			os.Exit(2)
	   		}
	   		b, _ := ioutil.ReadAll(resp.Body)
			   resp.Body.Close() */

	dir := "messagesByDate"
	value := now.Format("2006-01-02")

	for k := q; k < 32; k++ {

		s := fmt.Sprintf("%02d", k)

		value = strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + s

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

			s := `{"query":"query{` + dir + `(date:\"` + value + `\" appended: true){data{_id}}}"}`
			body := strings.NewReader(s)
			req, err := http.NewRequest("POST", "https://graphql.fauna.com/graphql", body)

			req.Header.Set("Authorization", "Bearer "+access.Secret)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
			req.Header.Set("X-Schema-Preview", "partial-update-mutation")

			resp, err := http.DefaultClient.Do(req)

			if err != nil {

				fmt.Fprint(w, "... an error occured.... please refresh browser window...")

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

			}

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

			s := `{"query":"query{` + dir + `(date:\"` + value + `\" appended: true){data{_id}}}"}`
			body := strings.NewReader(s)
			req, err := http.NewRequest("POST", "https://graphql.fauna.com/graphql", body)

			req.Header.Set("Authorization", "Bearer "+access.Secret)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
			req.Header.Set("X-Schema-Preview", "partial-update-mutation")

			resp, err := http.DefaultClient.Do(req)

			if err != nil {

				fmt.Fprint(w, "... an error occured.... please refresh browser window...")

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

			}

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

			s := `{"query":"query{` + dir + `(date:\"` + value + `\" appended: true){data{_id}}}"}`
			body := strings.NewReader(s)
			req, err := http.NewRequest("POST", "https://graphql.fauna.com/graphql", body)

			req.Header.Set("Authorization", "Bearer "+access.Secret)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
			req.Header.Set("X-Schema-Preview", "partial-update-mutation")

			resp, err := http.DefaultClient.Do(req)

			if err != nil {

				fmt.Fprint(w, "... an error occured.... please refresh browser window...")

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

			}

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

			s := `{"query":"query{` + dir + `(date:\"` + value + `\" appended: true){data{_id}}}"}`
			body := strings.NewReader(s)
			req, err := http.NewRequest("POST", "https://graphql.fauna.com/graphql", body)

			req.Header.Set("Authorization", "Bearer "+access.Secret)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
			req.Header.Set("X-Schema-Preview", "partial-update-mutation")

			resp, err := http.DefaultClient.Do(req)

			if err != nil {

				fmt.Fprint(w, "... an error occured.... please refresh browser window...")

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

			}

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

			s := `{"query":"query{` + dir + `(date:\"` + value + `\" appended: true){data{_id}}}"}`
			body := strings.NewReader(s)
			req, err := http.NewRequest("POST", "https://graphql.fauna.com/graphql", body)

			req.Header.Set("Authorization", "Bearer "+access.Secret)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
			req.Header.Set("X-Schema-Preview", "partial-update-mutation")

			resp, err := http.DefaultClient.Do(req)

			if err != nil {

				fmt.Fprint(w, "... an error occured.... please refresh browser window...")

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

			}

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

			s := `{"query":"query{` + dir + `(date:\"` + value + `\" appended: true){data{_id}}}"}`
			body := strings.NewReader(s)
			req, err := http.NewRequest("POST", "https://graphql.fauna.com/graphql", body)

			req.Header.Set("Authorization", "Bearer "+access.Secret)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
			req.Header.Set("X-Schema-Preview", "partial-update-mutation")

			resp, err := http.DefaultClient.Do(req)

			if err != nil {

				fmt.Fprint(w, "... an error occured.... please refresh browser window...")

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

			}

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

			s := `{"query":"query{` + dir + `(date:\"` + value + `\" appended: true){data{_id}}}"}`
			body := strings.NewReader(s)
			req, err := http.NewRequest("POST", "https://graphql.fauna.com/graphql", body)

			req.Header.Set("Authorization", "Bearer "+access.Secret)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
			req.Header.Set("X-Schema-Preview", "partial-update-mutation")

			resp, err := http.DefaultClient.Do(req)

			if err != nil {

				fmt.Fprint(w, "... an error occured.... please refresh browser window...")

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

			}

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

	//case "POST":

	/* client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, "localhost:5000/bolt/users/user2", Reader(s))
	if err != nil {
		fmt.Fprint(w, err)
	}
	_, err = client.Do(req)
	if err != nil {
		fmt.Fprint(w, err)
	} */

	//}

}
