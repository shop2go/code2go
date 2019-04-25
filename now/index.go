package main

import (
	//"fmt"

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

	/* 	var start time.Time
	   	var end time.Time */

	now := time.Now()
	/* year, _ := strconv.Atoi(now.Format("2006"))

	m := time.Date(year, 04, 01, 00, 00, 00, 0, time.UTC)
	o := time.Date(year, 10, 01, 00, 00, 00, 0, time.UTC)

	switch {

	case m.AddDate(0, 0, -1).Weekday() == 0:

		start = m.AddDate(0, 0, -1)

	case m.AddDate(0, 0, -2).Weekday() == 0:

		start = m.AddDate(0, 0, -2)

	case m.AddDate(0, 0, -3).Weekday() == 0:

		start = m.AddDate(0, 0, -3)

	case m.AddDate(0, 0, -4).Weekday() == 0:

		start = m.AddDate(0, 0, -4)

	case m.AddDate(0, 0, -5).Weekday() == 0:

		start = m.AddDate(0, 0, -5)

	case m.AddDate(0, 0, -6).Weekday() == 0:

		start = m.AddDate(0, 0, -6)

	case m.AddDate(0, 0, -7).Weekday() == 0:

		start = m.AddDate(0, 0, -7)

	}

	switch {

	case o.AddDate(0, 0, -1).Weekday() == 0:

		end = o.AddDate(0, 0, -1)

	case o.AddDate(0, 0, -2).Weekday() == 0:

		end = o.AddDate(0, 0, -2)

	case o.AddDate(0, 0, -3).Weekday() == 0:

		end = o.AddDate(0, 0, -3)

	case o.AddDate(0, 0, -4).Weekday() == 0:

		end = o.AddDate(0, 0, -4)

	case o.AddDate(0, 0, -5).Weekday() == 0:

		end = o.AddDate(0, 0, -5)

	case o.AddDate(0, 0, -6).Weekday() == 0:

		end = o.AddDate(0, 0, -6)

	case o.AddDate(0, 0, -7).Weekday() == 0:

		end = o.AddDate(0, 0, -7)

	}
	*/

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

	if time.Now().Month() == now.Month() {

		p, _ = strconv.Atoi(time.Now().Format("02"))

	} else {

		p = 1

	}

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
					<body style="background-color: #bcbcbc;">



					<div class="container" id="date" style="color:white; font-size:30px;">
					<form class="form-inline" role="form" method="post">
    <input class="form-control mr-sm-2" type="text" placeholder="Search" aria-label="Search" id ="find" name ="find">
    <button class="btn btn-outline-light my-2 my-sm-1" type="submit">Search</button>
  </form><br>
  <form class="form-inline" role="form" method="get">
					` + strconv.Itoa(c.Year) + `/` + now.Month().String() + `<button type="button" class="btn btn-link" onclick="window.location.href='next'">` + now.AddDate(0, +1, 0).Month().String() + `</button>
					</form>
					</div>
					<div class="container" id="data" style="color:white;">
	<br>

		`

	switch c.Days[q] {

	case "Monday":
		break
	case "Tuesday":
		str = str + `

		<div class="row">
		<span class="border">
		<div class="col-sm">
		<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Monday</span></button></div></span>
		`
	case "Wednesday":
		str = str + `
					<div class="row">

					<span class="border">
					<div class="col-sm">
					<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Monday</span></button>s</div></span>
					<span class="border">
					<div class="col-sm">
					<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Tuesday</span></button></div></span>`
	case "Thursday":
		str = str + `
									<div class="row">

									<span class="border">
									<div class="col-sm">
									<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Monday</span></button></div></span>
									<span class="border">
									<div class="col-sm">
									<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Tuesday</span></button></div></span>
									<span class="border">
									<div class="col-sm">
									<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Wednesday</span></button></div></span>`
	case "Friday":
		str = str + `
									<div class="row">

									<span class="border">
									<div class="col-sm">
									<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Monday</span></button></div></span>
									<span class="border">
									<div class="col-sm">
									<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Tuesday</span></button></div></span>
									<span class="border">
									<div class="col-sm">
									<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Wednesday</span></button></div></span>
									<span class="border">
									<div class="col-sm">
									<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Thursday</span></button></div></span>`
	case "Saturday":
		str = str + `
																																	<div class="row">

																																	<span class="border">
																																	<div class="col-sm">
																																	<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Monday</span></button></div></span>
																																	<span class="border">
																																	<div class="col-sm">
																																	<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Tuesday</span></button></div></span>
																																	<span class="border">
																																	<div class="col-sm">
																																	<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Wednesday</span></button></div></span>
																																	<span class="border">
																																	<div class="col-sm">
																																	<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Thursday</span></button></div></span>
																																	<span class="border">
																																	<div class="col-sm">
																																	<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Friday</span></button></div></span>
																																	`
	case "Sunday":
		str = str + `
																																																																	<div class="row">

																																																																	<span class="border">
																																																																	<div class="col-sm">
																																																																	<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Monday</span></button></div></span>
																																																																	<span class="border">
																																																																	<div class="col-sm">
																																																																	<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Tuesday</span></button></div></span>
																																																																	<span class="border">
																																																																	<div class="col-sm">
																																																																	<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Wednesday</span></button></div></span>
																																																																	<span class="border">
																																																																	<div class="col-sm">
																																																																	<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Thursday</span></button></div></span>
																																																																	<span class="border">
																																																																	<div class="col-sm">
																																																																	<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Friday</span></button></div></span>
																																																																	<span class="border">
																																																																	<div class="col-sm">
																																																																	<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Saturday</span></button></div></span>
																																																																	`
	}

	for k := q; k < 32; k++ {

		switch c.Days[k] {

		case "Monday":

			str = str + `
			</div><div class="row">

			<span class="border">

			<div class="col-sm">
			<button type="button" class="btn btn-link" onclick="window.location.href='entry#` + c.Days[k] + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `'"><span class="badge badge-pill badge-light">` + c.Days[k] + `<br>` + strconv.Itoa(k) + `</span></button>
				</div> </span>

				`

		case "Tuesday":

			str = str + `

			<span class="border">

			<div class="col-sm">
			<button type="button" class="btn btn-link" onclick="window.location.href='entry#` + c.Days[k] + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `'"><span class="badge badge-pill badge-light">` + c.Days[k] + `<br>` + strconv.Itoa(k) + `</span></button>
				</div> </span>

			`

		case "Wednesday":

			str = str + `

			<span class="border">

			<div class="col-sm">
			<button type="button" class="btn btn-link" onclick="window.location.href='entry#` + c.Days[k] + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `'"><span class="badge badge-pill badge-light">` + c.Days[k] + `<br>` + strconv.Itoa(k) + `</span></button>
				</div> </span>

				`

		case "Thursday":

			str = str + `

			<span class="border">

			<div class="col-sm">
			<button type="button" class="btn btn-link" onclick="window.location.href='entry#` + c.Days[k] + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `'"><span class="badge badge-pill badge-light">` + c.Days[k] + `<br>` + strconv.Itoa(k) + `</span></button>
				</div> </span>

				`

		case "Friday":

			str = str + `

			<span class="border">

			<div class="col-sm">
			<button type="button" class="btn btn-link" onclick="window.location.href='entry#` + c.Days[k] + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `'"><span class="badge badge-pill badge-light">` + c.Days[k] + `<br>` + strconv.Itoa(k) + `</span></button>
				</div> </span>

				`

		case "Saturday":

			str = str + `

			<span class="border">

			<div class="col-sm">
			<button type="button" class="btn btn-link" onclick="window.location.href='entry#` + c.Days[k] + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `'"><span class="badge badge-pill badge-light">` + c.Days[k] + `<br>` + strconv.Itoa(k) + `</span></button>
				</div> </span>

				`

		case "Sunday":

			str = str + `

			<span class="border">

			<div class="col-sm">
			<button type="button" class="btn btn-link" onclick="window.location.href='entry#` + c.Days[k] + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `'"><span class="badge badge-pill badge-light">` + c.Days[k] + `<br>` + strconv.Itoa(k) + `</span></button>
				</div> </span>

				`

		}

	}

	switch c.Days[len(c.Days)] {

	case "Monday":
		str = str + `
		<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Tuesday</span></button></div></span>

			<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Wednesday</span></button></div></span>

			   <span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Thursday</span></button></div></span>

			   <span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Friday</span></button></div></span>

			   <span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Saturday</span></button></div></span>

			   <span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Sunday</span></button></div></span>


`
	case "Tuesday":
		str = str + `
		<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Wednesday</span></button></div></span>

			<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Thursday</span></button></div></span>

			<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Friday</span></button></div></span>

			<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Saturday</span></button></div></span>

			<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Sunday</span></button></div></span>


`
	case "Wednesday":
		str = str + `
		<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Thursday</span></button></div></span>

			<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Friday</span></button></div></span>

			<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Saturday</span></button></div></span>

			<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Sunday</span></button></div></span>


`
	case "Thursday":
		str = str + `
		<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Friday</span></button></div></span>

			<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Saturday</span></button></div></span>

			<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Sunday</span></button></div></span>


`
	case "Friday":
		str = str + `
		<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Saturday</span></button></div></span>

			<span class="border">
			<div class="col-sm">
			<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Sunday</span></button></div></span>


`
	case "Saturday":
		str = str + `
		<span class="border">

		<div class="col-sm">
		<button type="button" class="btn btn-link"><span class="badge badge-pill badge-light">Sunday</span></button></div></span>


`

	case "Sunday":
		str = str + `


		</div>

`
		break

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
