package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Cal struct {
	Year  int
	Month int
	Days  map[int]string
}

func Handler(w http.ResponseWriter, r *http.Request) {

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
	 <li class="list-group-item" id="` + strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `">` + strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + ` - ` + c.Days[k] + `

	 </li><br>
	 `
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
	 <li class="list-group-item" id="` + strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `">` + strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + ` - ` + c.Days[k] + `

	 </li><br>
	 `
		}

		str = str + `
 </ul>
						   </div>
						   <!-- Then Material JavaScript on top of Bootstrap JavaScript -->
<script src="https://assets.medienwerk.now.sh/material.min.js"></script>



					</body>
					</html>
					`

	}

	if r.Method == "GET" {

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(str)))
		w.Write([]byte(str))

	} else {

		r.ParseForm()

		s := strings.Join(r.Form["find"], " ")

		fmt.Fprint(w, s)

	}

}
