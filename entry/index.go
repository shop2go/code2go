package main

import (
	//"encoding/gob"
	//"log"
	//"net"
	"net/http"
	"strconv"
	//"strings"
	"time"

	//"github.com/aerogo/packet"
	//"github.com/mmaedel/code2go/pb"
)

type Cal struct {
	Year  int
	Month int
	Days  map[int]string
}

func Handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":

		/* store := make([]pb.ReqPost, 0)

		query := strings.TrimPrefix(r.URL.Path, "/entry#")

		qu := strings.SplitN(query, "-", -1)

		cue, err := strconv.Atoi(qu[2])

		if err != nil {

			log.Println(err)
			cue = 0

		}

		//persistence layer

		conn, err := net.Dial("tcp", "localhost:80")

		if err != nil {

			log.Println(err)

		}

		// Create a stream
		stream := packet.NewStream(1024)

		stream.SetConnection(conn)

		// Send a message
		stream.Outgoing <- packet.New(byte(cue), []byte(query))

		//the response gob from conn

		dec := gob.NewDecoder(conn)

		dec.Decode(&store)

		numberOfEntries := len(store) */

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

		<body style="background-color: #bcbcbc;">

		<div class="container" id="data" style="color:white; font-size:30px;">
		<form class="form-inline" role="form" method="POST">
		<ul class="list-group">
		`

		now := time.Now()

		var c Cal

		c.Year = now.Year()
		month, _ := strconv.Atoi(now.Format("01"))
		c.Month = month
		day := map[int]string{now.Day(): now.Weekday().String()}

		c.Days = day

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

		j := 1

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

		var p, q int

		l := len(c.Days)

		p, _ = strconv.Atoi(time.Now().Format("02"))

		for i := l; i >= p; i-- {

			q = i

		}

		//expose the anchor of specified date++; list apropriate entries for this month form persitence layer

		for k := q; k <= l; k++ {

			schedule := strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k)

			str = str + `
			<button type="button" class="btn btn-link" onclick="window.location.href='` + c.Days[k] + `'">
			<span class="badge badge-pill badge-dark">
			` + c.Days[k] + `
			</span>
			</button>
			<button class="btn btn-light" type="button" data-toggle="collapse" data-target="#entry" aria-expanded="false" aria-controls="entry">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" id="` + schedule + `" value="` + schedule + `" placeholder="` + schedule + `">
			<div class="collapse" id="entry">
			<div class="card card-body">
			<input readonly class="form-control-plaintext list-group-item-action" id="Schedule" aria-label="Schedule" name ="Schedule" value="` + schedule + `" placeholder="` + schedule + `" required>>
			<input class="form-control mr-sm-2" type="text" placeholder="topic" aria-label="Topic" id ="Topic" name ="Topic" required><br>
			<input class="form-control mr-sm-2" type="text" placeholder="entry" aria-label="Event" id ="Event" name ="Event" required><br>
			<input class="form-control mr-sm-2" type="text" placeholder="tags" aria-label="Tag" id ="Tag" name ="Tag"><br>
			<button type="submit" class="btn btn-light">submit
			</button>
			</div>
			</div>
			</span>
			</button>
			`

			/* if numberOfEntries > 0 {

				for n := 0; n < numberOfEntries; n++ {

					switch string(store[n].Schedule) {

					case schedule:

						str = str + `
						<input readonly class="form-control-plaintext list-group-item-action" id="` + string(store[n].PostId) + `" value="` + string(store[n].PostId) + `" placeholder="` + string(store[n].Tags) + `">
						`

					}

				}

			} */

		}

		for o := 1; o < 21; o++ {

			now = time.Now().AddDate(0, o, 0)

			c.Year = now.Year()
			month, _ = strconv.Atoi(now.Format("01"))
			c.Month = month
			day = map[int]string{now.Day(): now.Weekday().String()}

			c.Days = day

			i = 1

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

			for k := 1; k <= l; k++ {

				schedule := strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k)

				str = str + `
				<button type="button" class="btn btn-link" onclick="window.location.href='` + c.Days[k] + `'">
				<span class="badge badge-pill badge-dark">
				` + c.Days[k] + `
				</span>
				</button>
				<button type="submit" class="btn btn-light">
				<span class="badge badge-pill badge-light">
				<input readonly class="form-control-plaintext list-group-item-action" id="` + schedule + `" value="` + schedule + `" placeholder="` + schedule + `">
				</span>
				</button>	
				`

			}

		}

		str = str + `
 		</ul>
		</form>
		</div>
		
		<!-- Then Material JavaScript on top of Bootstrap JavaScript -->
		<script src="https://assets.medienwerk.now.sh/material.min.js"></script>

		</body>
		</html>
		`

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(str)))
		w.Write([]byte(str))

		/* var req pb.ReqPost

		url := strings.TrimPrefix(r.URL.Path, "/entry#")

		r.ParseForm()

		for k, v := range r.Form {

			switch k {

			case "Topic":

				s := strings.Join(v, " ")

				req.Topic = []byte(s)

			case "Entry":

				s := strings.Join(v, " ")

				req.Entry = []byte(s)

			case "Schedule":

				s := strings.Join(v, " ")

				req.Schedule = []byte(s)

			case "Tags":

				s := strings.Join(v, " ")

				req.Tags = []byte(s)

			default:

				continue

			}

		}

		//persistence layer

		conn, err := net.Dial("tcp", "localhost:80")

		if err != nil {

			log.Println(err)

		}

		// Create a stream
		stream := packet.NewStream(1024)

		stream.SetConnection(conn)

		// Send a message
		stream.Outgoing <- packet.New('P', []byte(req)) */

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

	}

}
