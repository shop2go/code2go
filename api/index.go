package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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

type CacheEntry struct {
	Month graphql.String   `graphql:"month"`
	Posts []graphql.String `graphql:"posts"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	url := strings.TrimPrefix(r.URL.Path, "/")

	n, err := strconv.Atoi(url)

	now := time.Now()

	var c Cal

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

	var t, y int

	for t < 20 {

		t++

		y = time.Now().AddDate(0, t, 0).Year()

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

	mo := fmt.Sprintf("%02d", c.Month)

	ye := strconv.Itoa(c.Year)

	var value string

	var query struct {
		CacheByMonth struct {
			CacheEntry
		} `graphql:"cacheByMonth(month: $Month)"`
	}

	v1 := map[string]interface{}{
		"Month": graphql.String(ye + `-` + mo),
	}

	if err = call.Query(context.Background(), &query, v1); err != nil {
		fmt.Fprintf(w, "get cache error: %v\n", err)
	}

	result := query.CacheByMonth.Posts

	hits := make(map[string]int, len(result))

	if result != nil {

		for _, v := range result {

			m := strings.Split(string(v), ":")

			if c, ok := hits[m[1]]; ok {

				hits[m[1]] = c + 1

			} else {

				hits[m[1]] = 1

			}

		}

	}

	for k := q; k < 32; k++ {

		d := fmt.Sprintf("%02d", k)

		value = ye + "-" + mo + "-" + d

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
				`

			if l, ok := hits[d]; ok {

				str = str + `
				<span style="text-align: inherit; color: #70db70" class="badge badge-pill badge-dark">
				` + strconv.Itoa(l) + `
				</span>
				`

			}

			str = str + `
			</button>
			`

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
			`

			if l, ok := hits[d]; ok {

				str = str + `
			<span style="text-align: inherit; color: #70db70" class="badge badge-pill badge-dark">
			` + strconv.Itoa(l) + `
			</span>
			`

			}

			str = str + `
		</button>
		`

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
			`

			if l, ok := hits[d]; ok {

				str = str + `
			<span style="text-align: inherit; color: #70db70" class="badge badge-pill badge-dark">
			` + strconv.Itoa(l) + `
			</span>
			`

			}

			str = str + `
		</button>
		`

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
			`

			if l, ok := hits[d]; ok {

				str = str + `
			<span style="text-align: inherit; color: #70db70" class="badge badge-pill badge-dark">
			` + strconv.Itoa(l) + `
			</span>
			`

			}

			str = str + `
		</button>
		`

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
			`

			if l, ok := hits[d]; ok {

				str = str + `
			<span style="text-align: inherit; color: #70db70" class="badge badge-pill badge-dark">
			` + strconv.Itoa(l) + `
			</span>
			`

			}

			str = str + `
		</button>
		`

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
			`

			if l, ok := hits[d]; ok {

				str = str + `
			<span style="text-align: inherit; color: #70db70" class="badge badge-pill badge-dark">
			` + strconv.Itoa(l) + `
			</span>
			`

			}

			str = str + `
		</button>
		`

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
			`

			if l, ok := hits[d]; ok {

				str = str + `
			<span style="text-align: inherit; color: #70db70" class="badge badge-pill badge-dark">
			` + strconv.Itoa(l) + `
			</span>
			`

			}

			str = str + `
		</button>
		`
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
