package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	f "github.com/fauna/faunadb-go/faunadb"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
	//"github.com/plutov/paypal"
)

type Access struct {
	//Reference *f.RefV `fauna:"ref"`
	Timestamp int    `fauna:"ts"`
	Secret    string `fauna:"secret"`
	Role      string `fauna:"role"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	var result map[string]int = make(map[string]int, 0)

	u := r.Host

	u = strings.TrimSuffix(u, "code2go.dev")

	token := "test"

	s := r.Cookies()

	for _, c := range s {

		if c.Name == "code2go.dev" {

			token = c.Value

		}

	}

	fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

	x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database("tickets"), "role": "server"}))

	if err != nil {

		fmt.Fprintf(w, "connection error: %v\n", err)

	}

	var access *Access

	x.Get(&access)

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: access.Secret},
	)

	httpClient := oauth2.NewClient(context.Background(), src)

	call := graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

	if u != "" {

		u = strings.TrimSuffix(u, ".")

		//fmt.Fprintf(w, "%v\n", u)

		var q2 struct {
			EventByName struct {
				ID        graphql.ID      `graphql:"_id"`
				Date      graphql.String  `graphql:"date"`
				Confirmed graphql.Boolean `graphql:"confirmed"`
				Host      struct {
					ID       graphql.ID     `graphql:"_id"`
					Username graphql.String `graphql:"username"`
					Email    graphql.String `graphql:"email"`
				} `graphql:"host"`
				Tickets struct {
					Data []struct {
						Total graphql.Int `graphql:"total"`
						Cat   struct {
							Data []struct {
								Category graphql.String `graphql:"category"`
								Price    graphql.Float  `graphql:"price"`
								Issued   graphql.Int    `graphql:"issued"`
							} `graphql:"data"`
						} `graphql:"cat"`
					} `graphql:"data"`
				} `graphql:"tickets"`
			} `graphql:"eventByName(name: $name)"`
		}

		v1 := map[string]interface{}{
			"name": graphql.String(u),
		}

		if err := call.Query(context.Background(), &q2, v1); err != nil {

			fmt.Fprintf(w, "%v\n", err)

		}

		//fmt.Fprintf(w, "%v\n", q2)

		r := q2.EventByName.Tickets

		if len(r.Data) > 0 {

			for _, v := range r.Data {

				//for _, y := range v {

				for _, x := range v.Cat.Data {

					var i int

					if j, ok := result[string(x.Category)+":"+strconv.FormatFloat(float64(x.Price), 'f', 2, 64)]; ok {

						i = j

					}

					//if v.Event.Name == r.Name {

					result[string(x.Category)+":"+strconv.FormatFloat(float64(x.Price), 'f', 2, 64)] = i + int(x.Issued)

					//}

				}

			}
		}

	}

	/* 	c, err := paypal.NewClient(os.Getenv("PP_ID"), os.Getenv("PP_SECRET"), paypal.APIBaseSandBox)

	   	if err != nil {

			   fmt.Printf(w, err)

		} */

	switch r.Method {

	case "GET":

		switch token {

		default:

			var q1 struct {
				UserByToken struct {
					ID         graphql.ID      `graphql:"_id"`
					Username   graphql.String  `graphql:"username"`
					registered graphql.Boolean `graphql:"isregistered`
					Email      graphql.String  `graphql:"email`
					Token      graphql.String  `graphql:"token`
				} `graphql:"userByToken(token: $token)"`
			}

			v1 := map[string]interface{}{
				"token": graphql.String(token),
			}

			if err := call.Query(context.Background(), &q1, v1); err != nil {

				fmt.Fprintf(w, "%v\n", err)

			}

			//result1 := string(q1.UserByToken.Email)

			str := `
			<!DOCTYPE html>
			<html lang="en">
			<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<title>` + u + `</title>
			<link href="https://fonts.googleapis.com/css?family=Roboto:300,300i,400,400i,500,500i,700,700i|Roboto+Mono:300,400,700|Roboto+Slab:300,400,700" rel="stylesheet">
			<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
			   <link href="https://assets.medienwerk.now.sh/material.min.css" rel="stylesheet">
			</head>
			<body style="background-color: #bcbcbc;">
			   
			<div class="container" id="data" style="color:white;">
			<br>
			<form class="form-inline" role="form" method="POST">
			<input type="email" class="form-control" value="` + string(q1.UserByToken.Email) + `" aria-label="Email" id ="Email" name ="Email">
			<br>
			`

			for k, v := range result {

				catprice := strings.SplitN(k, ":", -1)

				str = str + `

				<span><br>
				Category: ` + catprice[0] + `</span>
				<br><br>
				<input readonly="true" class="form-control-plaintext" id="Ticket` + k + `" aria-label="Ticket` + k + `" name ="Ticket` + k + `" value="` + strconv.Itoa(v) + `">
				<input class="form-control-plaintext" id="Count` + k + `" aria-label="Count` + k + `" name ="Count` + k + `" value="0">
				<input readonly="true" class="form-control-plaintext" id="Price` + k + `" aria-label="Price` + k + `" name ="Price` + k + `" value="` + catprice[1] + `">
				<br>

				`

			}

			str = str + `
			
			<button type="submit" class="btn btn-light">checkout</button>
			</form>
			</div>
			<br>
			<br>
	
			   
			<script src="https://assets.medienwerk.now.sh/material.min.js">
			</script>
			</body>
			</html>
			`

			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("Content-Length", strconv.Itoa(len(str)))
			w.Write([]byte(str))

		case "":

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
			   
			<div class="container" id="data" style="color:white;">
			<br>
			<form class="form-inline" role="form" method="POST">
			<input type="email" class="form-control" placeholder="name@example.com" aria-label="Email" id ="Email" name ="Email">
			<br>
			<input readonly="true" class="form-control-plaintext" id="Ticket" aria-label="Ticket" name ="Ticket" value="Ticket">
			<input class="form-control-plaintext" id="Count" aria-label="Count" name ="Count" placeholder="" value="1">
			<input readonly="true" class="form-control-plaintext" id="Price" aria-label="Price" name ="Price" value="50">
			<br>
			<button type="submit" class="btn btn-light">checkout</button>
			</form>
			</div>
			<br>
			<br>
	
			   
			<script src="https://assets.medienwerk.now.sh/material.min.js">
			</script>
			</body>
			</html>
			`

			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("Content-Length", strconv.Itoa(len(str)))
			w.Write([]byte(str))

		}

	case "POST":

		var i int
		var j, sum float64

		r.ParseForm()

		for k, _ := range result {

			//	http.Redirect(w, r, "/transaction", http.StatusFound)

			count := r.Form.Get("Count" + k)
			price := r.Form.Get("Price" + k)

			i, err = strconv.Atoi(count)

			if err != nil {

				i = 0

			}

			j, err = strconv.ParseFloat(price, 64)

			if i != 0 || j != 0 {

				sum = float64(i) * j

			}

		}

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
	<script
	src="https://www.paypal.com/sdk/js?client-id=` +
			os.Getenv("PP_CLIENT_ID") + `&currency=EUR">
	  </script>
	   <br>
	   <br>
	<div class="container" id="paypal-button-container">
	</div>

	<script>
paypal.Buttons({
createOrder: function(data, actions) {
  return actions.order.create({
	"intent": "CAPTURE", 
	purchase_units: [{
	  amount: {
		"currency_code": "EUR",
		  value: '` + strconv.FormatFloat(sum, 'f', 2, 64) + `'
	  }
	}]
  });
},
onApprove: function(data, actions) {
  return actions.order.capture().then(function(details) {
	alert('Transaction completed by ' + details.payer.name.given_name);
  });
}
}).render('#paypal-button-container');
</script>
	   
	<script src="https://assets.medienwerk.now.sh/material.min.js">
	</script>
	</body>
	</html>
	`

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(str)))
		w.Write([]byte(str))

	}

}
