package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	/* 	"github.com/google/uuid"
	   	"github.com/gorilla/schema" */)

type Cal struct {
	Year  int
	Month int
	Days  map[int]string
}

/* type Form struct {
	Id    uuid.UUID
	Topic string
	Tag   string
	Event string
	Date  string
	Time  time.Time
} */

func Handler(w http.ResponseWriter, r *http.Request) {

	url := r.URL

	f := url.Fragment

	resp, _ := http.Get("http://example.com/" + f)

	b, _ := ioutil.ReadAll(resp.Body)

	resp.Body.Close()

	fmt.Fprint(w, string(b))

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
				<body style="background-color: #bcbcbc;">

				<div class="container" id="data" style="color:white; font-size:30px;">
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

		var p int
		var q int

		l := len(c.Days)

		p, _ = strconv.Atoi(time.Now().Format("02"))

		for i := l; i >= p; i-- {

			q = i

		}

		for k := q; k <= l; k++ {

			str = str + `

			<button type="button" class="list-group-item list-group-item-action" id="` + strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `" onclick="window.location.href='form#` + strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `'">` + strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `<span class="badge badge-pill badge-light">` + c.Days[k] + `</span>
	 <br>`
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

			l = len(c.Days)

			for k := 1; k <= l; k++ {

				str = str + `
				<button type="button" class="list-group-item list-group-item-action" id="` + strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `" onclick="window.location.href='form#` + strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `'">` + strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `<span class="badge badge-pill badge-light">` + c.Days[k] + `</span>
				<br>`
			}

		}

		str = str + `
 </ul>
						   </div>
						   <!-- Then Material JavaScript on top of Bootstrap JavaScript -->
<script src="https://assets.medienwerk.now.sh/material.min.js"></script>



					</body>
					</html>
					`

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(str)))
		w.Write([]byte(str))

	}

}
