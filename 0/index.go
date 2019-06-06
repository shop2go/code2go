package main

import (
	"io"
	"io/ioutil"
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

	switch r.Method {

	case "GET":

		url := strings.TrimPrefix(r.URL.Path, "/")

		n, _ := strconv.Atoi(url)

		now := time.Now().AddDate(0, n, 0)

		/* 	var start time.Time
		var end time.Time */
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
		<body style="background-color: #bcbcbc;">
   		<div class="container" id="search" style="color:white; font-size:30px;">
		<form class="form-inline" role="form" method="post">
	   	<input class="form-control mr-sm-2" type="text" placeholder="Search" aria-label="Search" id ="find" name ="find">
	   	<button class="btn btn-outline-light my-2 my-sm-1" type="submit">Search</button><br>
		 </div>
		 <br>
		 <div class="container" id="nav" style="color:white;">
		 ` + strconv.Itoa(time.Now().Year()) + `
		 <br>
		 `

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

		for t := 0; t < n; t++ {

			if time.Now().AddDate(0, t, 0).Year() != c.Year {

				str = str + `
				<br> 
				` + strconv.Itoa(time.Now().AddDate(0, t, 0).Year()) + `
				<br>
				<button type="button" class="btn btn-outline-dark" onclick="window.location.href='` + strconv.Itoa(t) + `'">` + time.Now().AddDate(0, t, 0).Month().String() + `
				</button>
				`

				c.Year = time.Now().AddDate(0, t, 0).Year()

			} else {

				str = str + ` 
				 <button type="button" class="btn btn-outline-dark" onclick="window.location.href='` + strconv.Itoa(t) + `'">` + time.Now().AddDate(0, t, 0).Month().String() + `
				 </button>
				 `

			}

		}

		str = str + `
		<button type="button" class="btn btn-light">` + now.Month().String() + `
		 </button>
		 `

		for t := n + 1; t < 21; t++ {

			if time.Now().AddDate(0, t, 0).Year() != c.Year {

				str = str + `
				<br>
				` + strconv.Itoa(time.Now().AddDate(0, t, 0).Year()) + `
				<br>
				<button type="button" class="btn btn-outline-dark" onclick="window.location.href='` + strconv.Itoa(t) + `'">` + time.Now().AddDate(0, t, 0).Month().String() + `
				</button>
				`

				c.Year = time.Now().AddDate(0, t, 0).Year()

			} else {

				str = str + `
				<button type="button" class="btn btn-outline-dark" onclick="window.location.href='` + strconv.Itoa(t) + `'">` + time.Now().AddDate(0, t, 0).Month().String() + `
				</button>
				`

			}
		}

		str = str + `
		</form>
		<br>
		</div>
		<br>
		<div class="container" id="data" style="color:white;">
		<form class="form-inline" role="form">
		<ul class="list-group">
		`

		c.Year = now.Year()

		switch c.Days[q] {

		case "Monday":
			break
		case "Tuesday":
			str = str + `
			<br>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Monday">
			</span>
			</button>
			`
		case "Wednesday":
			str = str + `
			<br>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Monday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Tuesday">
			</span>
			</button>
			`

		case "Thursday":
			str = str + `
			<br>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Monday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Tuesday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Wednesday">
			</span>
			</button>
			`
		case "Friday":
			str = str + `
			<br>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Monday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Tuesday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Wednesday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Thursday">
			</span>
			</button>
			`
		case "Saturday":
			str = str + `
			<br>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Monday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Tuesday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Wednesday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Thursday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Friday">
			</span>
			</button>
			`
		case "Sunday":
			str = str + `
			<br>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Monday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Tuesday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Wednesday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Thursday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Friday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Saturday">
			</span>
			</button>
			`
		}

		//make cert client
		resp, _ := http.Get("http://example.org/0" + url)

		b, _ := ioutil.ReadAll(resp.Body)

		resp.Body.Close()

		for k := q; k < 32; k++ {

			switch c.Days[k] {

			case "Monday":

				str = str + `
				<br>
				<button type="button" class="btn btn-link" onclick="window.location.href='entry#` + strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `'">
				<span class="badge badge-pill badge-dark">
				` + c.Days[k] + `
				</span>
				<span class="badge badge-pill badge-light">
				<input readonly class="form-control-plaintext list-group-item-action" value="` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `" placeholder="` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `">
				</span>
				`

				if string(b[k-q]) != "0" {

					str = str + `
					<span class="badge badge-pill badge-dark">
					` + string(b[k-q]) + `
					</span>
					</button>
					`

				} else {

					str = str + `
					</button>
					`

				}

			default:

				str = str + `
				<button type="button" class="btn btn-link" onclick="window.location.href='entry#` + strconv.Itoa(c.Year) + `-` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `'">
				<span class="badge badge-pill badge-dark">
				` + c.Days[k] + `
				</span>
				<span class="badge badge-pill badge-light">
				<input readonly class="form-control-plaintext list-group-item-action" value="` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `" placeholder="` + strconv.Itoa(c.Month) + `-` + strconv.Itoa(k) + `">
				</span>
				`

				if string(b[k-q]) != "0" {

					str = str + `
					<span class="badge badge-pill badge-dark">
					` + string(b[k-q]) + `
					</span>
					</button>
					`

				} else {

					str = str + `
					</button>
					`

				}

			}

		}

		switch c.Days[len(c.Days)] {

		case "Monday":
			str = str + `
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Tuesday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Wednesday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Thursday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Friday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Saturday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Sunday">
			</span>
			</button>
			`

		case "Tuesday":
			str = str + `
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Wednesday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Thursday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Friday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Saturday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Sunday">
			</span>
			</button>
			`

		case "Wednesday":
			str = str + `
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Thursday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Friday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Saturday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Sunday">
			</span>
			</button>
			`

		case "Thursday":
			str = str + `
			button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Friday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Saturday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Sunday">
			</span>
			</button>
			`

		case "Friday":
			str = str + `
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Saturday">
			</span>
			</button>
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Sunday">
			</span>
			</button>
			`

		case "Saturday":
			str = str + `
			<button type="button" class="btn btn-dark">
			<span class="badge badge-pill badge-light">
			<input readonly class="form-control-plaintext list-group-item-action" placeholder="Sunday">
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

	case "POST":

		r.ParseForm()

		//s := strings.Join(r.Form["find"], " ")

		/* client := &http.Client{}
		req, err := http.NewRequest(http.MethodPut, "localhost:5000/bolt/users/user2", Reader(s))
		if err != nil {
			fmt.Fprint(w, err)
		}
		_, err = client.Do(req)
		if err != nil {
			fmt.Fprint(w, err)
		} */

	}

}

func Reader(s string) io.Reader {

	return strings.NewReader(s)
}
