package main

import (
	"encoding/json"
	"errors"
	"fmt"

	//"log"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	f "github.com/fauna/faunadb-go/faunadb"
	/* 	"github.com/mschneider82/problem"
	   	"github.com/aerogo/packet"
	   	"github.com/mmaedel/code2go/pb" */)

type Cal struct {
	Year  int
	Month int
	Days  map[int]string
}

type Cache struct {
	Month string
	Posts []Post
}

type Post struct {
	ID       string
	Date     string
	Password string
	Title    string
	Content  interface{}
}

type Access struct {
	//Reference *f.RefV `fauna:"ref"`
	Timestamp int    `fauna:"ts"`
	Secret    string `fauna:"secret"`
	Role      string `fauna:"role"`
}

var errOnData = errors.New("error on data")

func getCache(a *Access) ([]Cache, error) {

	var result []Cache = make([]Cache, 0)

	dir := "allCaches"

	s := `{"query":"query{` + dir + `{data{month posts{_id date password title content}}}}"}`
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

					p := o[k].(map[string]interface{})

					resultP := make([]Post, len(p))

					resultP[k].ID = p["_id"].(string)

					resultP[k].Date = p["date"].(string)

					resultP[k].Password = p["password"].(string)

					resultP[k].Title = p["title"].(string)

					resultP[k].Content = p["content"]

					cache[j].Posts = append(cache[j].Posts, resultP[k])

				}

				result = append(result, cache[j])

			}

		}

	}

	return result, nil

}

func Handler(w http.ResponseWriter, r *http.Request) {

	var access *Access

	var c Cal

	var posts []Post

	now := time.Now()

	c.Year = now.Year()
	month, _ := strconv.Atoi(now.Format("01"))
	c.Month = month
	day := map[int]string{now.Day(): now.Weekday().String()}

	c.Days = day

	fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

	x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database(now.Format("2006")), "role": "server-readonly"}))

	if err != nil {

		fmt.Fprint(w, err)

		return

	}

	if err := x.Get(&access); err != nil {

		fmt.Fprint(w, err)

		return

	}

	result, err := getCache(access)

	if err != nil {

		fmt.Fprint(w, err)

		return

	}

	//fmt.Fprint(w, result)

	/* addr := &net.TCPAddr{net.ParseIP(a), 8080, "UTC"}

	//switch r.Method {

	//case "GET":

		store := make([]pb.ReqPost, 0)

		query := strings.TrimPrefix(r.URL.Path, "/entry#")

		qu := strings.SplitN(query, "-", -1)

		cue, err := strconv.Atoi(qu[2])

		if err != nil {

			//log.Println(err)
			cue = 0

		}

		//persistence layer

		conn, err := net.DialTCP("tcp", nil, addr)

		if err != nil {

			log.Fatal(err)

		}

		// Create a stream

		stream := packet.NewStream(1024)

		stream.SetConnection(conn)

		// Send a message

		stream.Outgoing <- packet.New(byte(cue), []byte(query))

		//conn.CloseWrite()

		//the response gob from conn

		dec := gob.NewDecoder(conn)

		dec.Decode(&store)

		//conn.CloseRead()
	*/

	z := 1

	for z < 32 {

		d := now.AddDate(0, 0, z)

		m, _ := strconv.Atoi(d.Format("01"))

		if m != c.Month {

			break

		}

		e, _ := strconv.Atoi(d.Format("02"))

		c.Days[e] = d.Weekday().String()

		z++

	}

	j := 1

	for j > 0 {

		d := now.AddDate(0, 0, -j)

		m, _ := strconv.Atoi(d.Format("01"))

		if m != c.Month {

			break

		}

		e, _ := strconv.Atoi(d.Format("02"))

		c.Days[e] = d.Weekday().String()

		j++

	}

	var p, q int

	l := len(c.Days)

	p, _ = strconv.Atoi(now.Format("02"))

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

	//expose the anchor of specified date++; list apropriate entries for that date whithin the actual month from persitence layer

	for k := q; k <= l; k++ {

		n := fmt.Sprintf("%02d", k)

		schedule := strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + n

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

			<div class="container" id="threads` + schedule + `">
			<!--form class="form-inline" role="form"-->
			<input readonly class="form-control-plaintext list-group-item-action" id="thread` + schedule + `" value="new thread" placeholder="new thread" onclick="window.location.href='https://` + schedule + `.code2go.dev/public'">
						
			`

		/* <form class="form-inline" role="form" method="POST">
		<input readonly="true" class="form-control-plaintext" id="Schedule" aria-label="Schedule" name ="Schedule" value="` + schedule + `" type="hidden">
		<input class="form-control mr-sm-2" type="text" placeholder="Title" aria-label="Title" id ="Title" name ="Title" required>
		<!--input class="form-control mr-sm-2" type="text" placeholder="entry" aria-label="Entry" id ="Entry" name ="Entry" required-->
		<input class="form-control mr-sm-2" type="text" placeholder="Tags" aria-label="Tags" id ="Tags" name ="Tags">
		<textarea class="form-control  mr-sm-2" id="Content" rows="2" placeholder="Content"></textarea>
		<br>
		<button type="submit" class="btn btn-light">submit</button>
		</form>
		</div>
		` */

		if result != nil {

			for _, v := range result {

				if v.Month == strconv.Itoa(c.Year)+`-`+strconv.Itoa(c.Month) {

					posts = v.Posts

				}

			}

			m := len(posts)

			if m > 0 {

				for n := 0; n < m; n++ {

					switch posts[n].Date {

					case schedule:

						if posts[n].Password == "" {

							str = str + `
						<input readonly class="form-control-plaintext list-group-item-action" id="` + posts[n].ID + `" value="` + posts[n].Title + `" placeholder="` + posts[n].Title + `" onclick="window.location.href='https://` + posts[n].ID + `.code2go.dev/public'">
						`

						} else {

							str = str + `
						<input readonly class="form-control-plaintext list-group-item-action" id="` + posts[n].ID + `" value="` + posts[n].Title + `" placeholder="password protected" onclick="window.location.href='https://` + posts[n].ID + `.code2go.dev/password'">
						`

						}

					}

				}

			}

		}

	}

	str = str + `
		<!--/form-->
		</div>
		`

	y := c.Year

	for o := 1; o < 21; o++ {

		now = time.Now().AddDate(0, o, 0)

		c.Year = now.Year()
		month, _ = strconv.Atoi(now.Format("01"))
		c.Month = month
		day = map[int]string{now.Day(): now.Weekday().String()}

		c.Days = day

	LOOP:

		if c.Year == y {

			i := 1

			for i < 32 {

				d := now.AddDate(0, 0, i)

				m, _ := strconv.Atoi(d.Format("01"))

				if m != month {

					break

				}

				e, _ := strconv.Atoi(d.Format("02"))

				c.Days[e] = d.Weekday().String()

				i++

			}

			j = 1

			for j > 0 {

				d := now.AddDate(0, 0, -j)

				m, _ := strconv.Atoi(d.Format("01"))

				if m != month {

					break

				}

				e, _ := strconv.Atoi(d.Format("02"))

				c.Days[e] = d.Weekday().String()

				j++

			}

			//all following months without entries

			//store = nil

			l = len(c.Days)

			for k := q; k <= l; k++ {

				n := fmt.Sprintf("%02d", k)

				schedule := strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + n

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
	
				<div class="container" id="threads` + schedule + `">
				<!--form class="form-inline" role="form"-->
				<input readonly class="form-control-plaintext list-group-item-action" id="thread` + schedule + `" value="new thread" placeholder="new thread" onclick="window.location.href=''https://` + schedule + `.code2go.dev/public'">
							
				`

				/* <form class="form-inline" role="form" method="POST">
				<input readonly="true" class="form-control-plaintext" id="Schedule" aria-label="Schedule" name ="Schedule" value="` + schedule + `" type="hidden">
				<input class="form-control mr-sm-2" type="text" placeholder="Title" aria-label="Title" id ="Title" name ="Title" required>
				<!--input class="form-control mr-sm-2" type="text" placeholder="entry" aria-label="Entry" id ="Entry" name ="Entry" required-->
				<input class="form-control mr-sm-2" type="text" placeholder="Tags" aria-label="Tags" id ="Tags" name ="Tags">
				<textarea class="form-control  mr-sm-2" id="Content" rows="2" placeholder="Content"></textarea>
				<br>
				<button type="submit" class="btn btn-light">submit</button>
				</form>
				</div>
				` */

				if result != nil {

					for _, v := range result {

						if v.Month == strconv.Itoa(c.Year)+`-`+strconv.Itoa(c.Month) {

							posts = v.Posts

						}

					}

					m := len(posts)

					if m > 0 {

						for n := 0; n < m; n++ {

							switch posts[n].Date {

							case schedule:

								if posts[n].Password == "" {

									str = str + `
							<input readonly class="form-control-plaintext list-group-item-action" id="` + posts[n].ID + `" value="` + posts[n].Title + `" placeholder="` + posts[n].Title + `" onclick="window.location.href='https://` + posts[n].ID + `.code2go.dev/public'">
							`

								} else {

									str = str + `
							<input readonly class="form-control-plaintext list-group-item-action" id="` + posts[n].ID + `" value="` + posts[n].Title + `" placeholder="password protected" onclick="window.location.href='https://` + posts[n].ID + `.code2go.dev/password'">
							`

								}

							}

						}

					}

				}

			}

			str = str + `
				<!--/form-->
				</div>
				`

		} else {

			fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

			x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database(now.Format("2006")), "role": "server-readonly"}))

			if err != nil {

				fmt.Fprint(w, err)

				return

			}

			if err := x.Get(&access); err != nil {

				fmt.Fprint(w, err)

				return

			}

			result, err = getCache(access)

			if err != nil {

				fmt.Fprint(w, err)

				return

			}

			y = c.Year

			goto LOOP

		}

	}

	str = str + `
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

	//}

	//var s string
	/* var req pb.ReqPost

	r.ParseForm()

	for k, v := range r.Form {

		switch k {

		case "Topic":

			//s = s + k + ": " + strings.Join(v, " ") + "\n\r"

			s := strings.Join(v, " ")

			req.Topic = []byte(s)

		case "Entry":

			sb := "\x00" + strings.Join(v, "\x20\x00") // x20 = space and x00 = null

			//s = s + k + ": " + strings.Join(v, " ") + "\n\r"

			req.Entry = []byte(sb)

		case "Schedule":

			//s = s + k + ": " + strings.Join(v, " ") + "\n\r"
			s := strings.Join(v, " ")

			req.Schedule = []byte(s)

		case "Tags":

			//s = s + k + ": " + strings.Join(v, " ") + "\n\r"
			s := strings.Join(v, "#")

			req.Tags = []byte(s)

		default:

			continue

		}

	}

	if req.Topic != nil {

		//Post data

		// Create a stream

		stream = packet.NewStream(1024)

		stream.SetConnection(conn)

		// Send data

		stream.Outgoing <- packet.New('T', req.Topic)

		stream.Outgoing <- packet.New('E', req.Entry)

		stream.Outgoing <- packet.New('S', req.Schedule)

		stream.Outgoing <- packet.New('#', req.Tags)

		w.Write(req.Schedule)
	*/
	//} else {
	/*
			w.Header().Set("Content-Type", "text/html")
		   	w.Header().Set("Content-Length", strconv.Itoa(len(str)))
		   	w.Write([]byte(str))
	*/
	//}

	/* client := &http.Client{}

	req, err := http.NewRequest("POST", "http://localhost/"+url, bytes.NewBuffer([]byte(s)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value") // This makes it work
	if err != nil {
		log.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Write(body) */

	//}

}
