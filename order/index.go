package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	//"time"

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

type CostumerEntry struct {
	ID         graphql.ID      `graphql:"_id"`
	First      graphql.String  `graphql:"first"`
	Last       graphql.String  `graphql:"last"`
	Email      graphql.String  `graphql:"email"`
	Phone      graphql.String  `graphql:"phone"`
	Registered graphql.Boolean `graphql:"registered"`
}

type AddressEntry struct {
	ID     graphql.ID     `graphql:"_id"`
	Street graphql.String `graphql:"street"`
	Number graphql.String `graphql:"number"`
	Door   graphql.String `graphql:"door"`
	City   graphql.String `graphql:"city"`
	Zip    graphql.String `graphql:"zip"`
}

type StatusEntry struct {
	ID       graphql.ID      `graphql:"_id"`
	Order    graphql.ID      `graphql:"order"`
	Payment  graphql.Boolean `graphql:"payment"`
	Delivery graphql.Boolean `graphql:"delivery"`
}

type ProductEntry struct {
	ID      graphql.ID     `graphql:"_id"`
	ImgURL  graphql.String `graphql:"imgURL"`
	Product graphql.String `graphql:"product"`
	Cat     graphql.String `graphql:"cat"`
	Info    graphql.String `graphql:"info"`
	Price   graphql.Float  `graphql:"price"`
	Pack    graphql.Int    `graphql:"pack"`
	InfoURL graphql.String `graphql:"infoURL"`
	LinkURL graphql.String `graphql:"linkURL"`
	LinkDIM graphql.Int    `graphql:"linkDIM"`
}

type CartEntry struct {
	ID       graphql.ID   `graphql:"_id"`
	Order    graphql.ID   `graphql:"order"`
	Products []graphql.ID `graphql:"products"`
}

type OrderEntry struct {
	ID       graphql.ID     `graphql:"_id"`
	Date     graphql.String `graphql:"date"`
	Costumer graphql.ID     `graphql:"costumer"`
	Cart     graphql.ID     `graphql:"cart"`
	Amount   graphql.Float  `graphql:"amount"`
	Status   graphql.ID     `graphql:"status"`
}

type SourceEntry struct {
	ID     graphql.ID     `graphql:"_id"`
	Link   graphql.ID     `graphql:"link"`
	Origin graphql.String `graphql:"origin"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	var total float64

	var node string

	m := make(map[string]float64, 0)

	id := r.Host

	id = strings.TrimSuffix(id, "code2go.dev")

	id = strings.TrimSuffix(id, ".")

	if id == "" {

		http.Redirect(w, r, "https://code2go.dev/", http.StatusSeeOther)

	} else {

		fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

		x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database("shop"), "role": "server"}))

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

		//ID, _ := base64.StdEncoding.DecodeString(u)

		var q struct {
			SourceByLink struct {
				SourceEntry
			} `graphql:"sourceByLink(link: $Link)"`
		}

		v := map[string]interface{}{
			"Link": graphql.ID(id),
		}

		if err = call.Query(context.Background(), &q, v); err != nil {
			fmt.Fprintf(w, "error with source: %v\n", err)
		}

		if q.SourceByLink.ID == nil {

			http.Redirect(w, r, "https://code2go.dev/", http.StatusSeeOther)

		} else {

			node = string(q.SourceByLink.Origin)

			x, err = fc.Query(f.CreateKey(f.Obj{"database": f.Database(node), "role": "server"}))

			if err != nil {

				fmt.Fprintf(w, "connection error: %v\n", err)

			}

			x.Get(&access)

			src = oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: access.Secret},
			)

			httpClient = oauth2.NewClient(context.Background(), src)

			call = graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

			var q struct {
				FindCartByID struct {
					CartEntry
				} `graphql:"findCartByID(id: $ID)"`
			}

			v = map[string]interface{}{
				"ID": graphql.ID(id),
			}

			if err = call.Query(context.Background(), &q, v); err != nil {
				fmt.Fprintf(w, "error with cart: %v\n", err)
			}

			if q.FindCartByID.Products != nil {

				for _, id := range q.FindCartByID.Products {

					var q struct {
						FindProductByID struct {
							ProductEntry
						} `graphql:"findProductByID(id: $ID)"`
					}

					v := map[string]interface{}{
						"ID": id,
					}

					if err = call.Query(context.Background(), &q, v); err != nil {
						fmt.Fprintf(w, "error with products: %v\n", err)
					}

					total = total + float64(q.FindProductByID.Price)

					if l, ok := m[string(q.FindProductByID.Product)]; ok {

						m[string(q.FindProductByID.Product)] = l + float64(q.FindProductByID.Price)

					} else {

						m[string(q.FindProductByID.Product)] = float64(q.FindProductByID.Price)

					}

				}

			} else {

				http.Redirect(w, r, "https://"+id+".code2go.dev/"+node, http.StatusSeeOther)

			}

		}

	}

	switch r.Method {

	case "GET":

		sum := strconv.FormatFloat(total+5.00, 'f', 2, 64)

		str :=

			`
		<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>shop2go</title>
		<!-- CSS -->
		<!-- Add Material font (Roboto) and Material icon as needed -->
		<link href="https://fonts.googleapis.com/css?family=Roboto:300,300i,400,400i,500,500i,700,700i|Roboto+Mono:300,400,700|Roboto+Slab:300,400,700" rel="stylesheet">
		<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">

		<!-- Add Material CSS, replace Bootstrap CSS -->
		<link href="https://assets.medienwerk.now.sh/material.min.css" rel="stylesheet">
		</head>
		<body style="background-color: #a1b116;">

		<div class="container" id="order" style="color:rgb(255, 255, 255); font-size:30px;">

		<br>
		<br>

		<h1>Einkauf</h1>
		<br>
		
			<p><h2>€ ` + sum + `</h2>Einkaufsumme<p>
			<br>
			<button type="button" class="btn btn-light" onclick="window.location.href='`+node+`'">Mit dem Einkauf fortfahren</button>

			<br><br>

		<form role="form" method="POST">

		<ul class="list-group">
		<li class="list-group-item">
		<button type="submit" class="btn btn-light" style="font-size:30px;">Bezahlen</button>
		</li><br>
		`

		for prod, flo := range m {

			//prod := fmt.Sprintf("%s", pro)
			price := strconv.FormatFloat(flo, 'f', 2, 64)

			str = str +

				`

			<li class="list-group-item">

			<label class="form-check-label" for="` + prod + `" style="font-size:25px;">` + prod + `</label>

			<input readonly="true" class="form-control-plaintext" id="` + prod + `" aria-label="` + prod + `" name ="` + prod + `" value="€ ` + price + `" style="font-size:30px;">
			<br>
			<button type="button" class="btn btn-light" onclick="window.location.href='product'">Produkt ändern</button>
			</li><br>

			`

		}

		str = str + `
			
			
			</ul>
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

	case "POST":

		//var registered bool

		//var c CostumerEntry

		//var gqlid graphql.ID

		/* r.ParseForm()

		first := r.Form.Get("first")
		last := r.Form.Get("last")
		email := r.Form.Get("email")
		phone := r.Form.Get("phone")
		street := r.Form.Get("street")
		number := r.Form.Get("number")
		door := r.Form.Get("door")

		if email != "" {

			registered = true

		} */

		/* fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

		x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database("shop"), "role": "server"}))

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
 */
		/* var q struct {
			CostumersByName struct {
				Data []CostumerEntry
			} `graphql:"costumersByName(last: $Last)"`
		}

		v := map[string]interface{}{
			"Last": graphql.String(last),
		}

		if err := call.Query(context.Background(), &q, v); err != nil {
			fmt.Fprintf(w, "error with costumer: %v\n", err)

		}

		if q.CostumersByName.Data != nil {

			for _, c = range q.CostumersByName.Data {

				if string(c.Phone) == phone {

					gqlid = c.ID

					break

				}

			}

		} else {

			var a struct {
				CreateAddress struct {
					AddressEntry
				} `graphql:"createAddress(data:{street: $Street, number: $Number, door: $Door})"`
			}

			x := map[string]interface{}{
				"Street": graphql.String(street),
				"Number": graphql.String(number),
				"Door":   graphql.String(door),
			}

			if err := call.Mutate(context.Background(), &a, x); err != nil {
				fmt.Fprintf(w, "error with address: %v\n", err)

			}

			var m struct {
				CreateCostumer struct {
					CostumerEntry
				} `graphql:"createCostumer(data:{first: $First, last: $Last, email: $Email, phone: $Phone, address: $Address, registered: $Registered})"`
			}

			z := map[string]interface{}{
				"First":      graphql.String(first),
				"Last":       graphql.String(last),
				"Email":      graphql.String(email),
				"Phone":      graphql.String(phone),
				"Address":    a.CreateAddress.ID,
				"Registered": graphql.Boolean(registered),
			}

			if err := call.Mutate(context.Background(), &m, z); err != nil {
				fmt.Fprintf(w, "error with costumer: %v\n", err)

			}
			gqlid = m.CreateCostumer.CostumerEntry.ID

		}

		var m1 struct {
			CreateOrder struct {
				OrderEntry
			} `graphql:"createOrder(data:{date: $Date, costumer: $Costumer, cart: $Cart, amount: $Amount})"`
		}

		x1 := map[string]interface{}{
			"Date":     graphql.String(time.Now().UTC().Format("2006-01-02")),
			"Costumer": gqlid,
			"Cart":     graphql.ID(node),
			"Amount":   graphql.Float(total),
		}

		if err := call.Mutate(context.Background(), &m1, x1); err != nil {
			fmt.Fprintf(w, "error with order: %v\n", err)

		}

		var m2 struct {
			CreateStatus struct {
				StatusEntry
			} `graphql:"createStatus(data:{order: $Order, payment: $Payment, delivery: $Delivery})"`
		}

		x2 := map[string]interface{}{
			"Order":    m1.CreateOrder.ID,
			"Payment":  graphql.Boolean(false),
			"Delivery": graphql.Boolean(false),
		}

		if err := call.Mutate(context.Background(), &m2, x2); err != nil {
			fmt.Fprintf(w, "error with status: %v\n", err)

		}
 */
		sum := strconv.FormatFloat(total+5.00, 'f', 2, 64)

		str :=

			`
		
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

	<div class="container" id="order" style="color:rgb(255, 255, 255); font-size:30px;">

	<br>
	<br>

	<h1>Einkauf</h1>´
	um<br>

	<br>
	
		<p><h2>€ ` + sum + `</h2>Einkaufsumme<p>
		<br>

	</div>

	<script
	src="https://www.paypal.com/sdk/js?client-id=` + os.Getenv("PP_CLIENT_ID") + `&currency=EUR">
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
		  value: "` + sum + `"
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
