package main

import (
	"context"
	//"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	Products []graphql.ID `graphql:"products"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	var total float64

	m := make(map[string]float64, 0)

	var call *graphql.Client

	u := r.Host

	u = strings.TrimSuffix(u, "code2go.dev")

	if u == "" {

		http.Redirect(w, r, "https://code2go.dev/shop", http.StatusSeeOther)

	} else {

		u = strings.TrimSuffix(u, ".")

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

		call = graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

		//ID, _ := base64.StdEncoding.DecodeString(u)

		var q struct {
			FindCartByID struct {
				CartEntry
			} `graphql:"findCartByID(id: $ID)"`
		}

		v1 := map[string]interface{}{
			"ID": graphql.ID(u),
		}

		if err = call.Query(context.Background(), &q, v1); err != nil {
			fmt.Fprintf(w, "error with products: %v\n", err)
		}

		if len(q.FindCartByID.Products) > 0 {

			for _, id := range q.FindCartByID.Products {

				var p struct {
					FindProductByID struct {
						ProductEntry
					} `graphql:"findProductByID(id: $ID)"`
				}

				v2 := map[string]interface{}{
					"ID": id,
				}

				if err = call.Query(context.Background(), &p, v2); err != nil {
					fmt.Fprintf(w, "error with products: %v\n", err)
				}

				total = total + float64(p.FindProductByID.Price)

				if l, ok := m[string(p.FindProductByID.Product)]; ok {

					m[string(p.FindProductByID.Product)] = l + float64(p.FindProductByID.Price)

				} else {

					m[string(p.FindProductByID.Product)] = float64(p.FindProductByID.Price)

				}

			}

		} else {

			http.Redirect(w, r, "https://"+u+".code2go.dev/shop", http.StatusSeeOther)

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
			<button type="button" class="btn btn-light" onclick="window.location.href='shop'">Mit dem Einkauf fortfahren</button>

			<br><br>

		<form class="form-inline" role="form" method="POST">
		<input type="text" class="form-control" value="" aria-label="first" id ="first" name ="first" placeholder="Vorname" required>
		<input type="text" class="form-control" value="" aria-label="last" id ="last" name ="last" placeholder="Nachname" required>
		<input type="email" class="form-control" value="" aria-label="email" id ="email" name ="email" placeholder="Email">
		<input type="number" class="form-control" value="" aria-label="phone" id ="phone" name ="phone" placeholder="Telefon" required>
		<input type="text" class="form-control" value="" aria-label="street" id ="street" name ="street" placeholder="Straße" required>
		<input type="text" class="form-control" value="" aria-label="number" id ="number" name ="number" placeholder="Hausnummer" required>
		<input type="text" class="form-control" value="" aria-label="door" id ="Email" door ="door" placeholder="Türnummer"required>
		<input readonly="true" type="text" class="form-control" value="Salzburg" aria-label="city" id ="city" name ="city">
		<input readonly="true" type="text" class="form-control" value="5020" aria-label="zip" id ="zip" name ="zip">
		
		<br>
		<br>
		
		<ul class="list-group">
		`

		for pro, flo := range m {

			price := strconv.FormatFloat(flo, 'f', 2, 64)

			str = str +

				`

			<li class="list-group-item">

			<label class="form-check-label" for="` + pro + `" style="font-size:25px;">` + pro + `</label>

			<input readonly="true" class="form-control-plaintext" id="` + pro + `" aria-label="` + pro + `" name ="` + pro + `" value="€ ` + price + `" style="font-size:30px;">
			<br>
			<button type="button" class="btn btn-light" onclick="window.location.href='product'">Produkt ändern</button>
			</li><br>

			`

		}

		str = str + `
			
			<button type="submit" class="btn btn-light" style="font-size:30px;">Bezahlen</button>
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

		var registered bool
		var costumerID graphql.ID

		r.ParseForm()

		first := r.Form.Get("first")
		last := r.Form.Get("last")
		email := r.Form.Get("email")
		phone := r.Form.Get("phone")
		street := r.Form.Get("street")
		number := r.Form.Get("number")
		door := r.Form.Get("door")

		if email != "" {

			registered = true

		}

		var q struct {
			CostumerByName struct {
				ID    graphql.ID     `graphql:"_id"`
				Phone graphql.String `graphql:"phone"`
			} `graphql:"costumerByName(name: $Name})"`
		}

		z := map[string]interface{}{
			"Name": graphql.String(last),
		}

		if err := call.Mutate(context.Background(), &q, z); err != nil {
			fmt.Fprintf(w, "error with costumer: %v\n", err)

		}

		if phone != string(q.CostumerByName.Phone) || q.CostumerByName.ID == nil {

			var a struct {
				CreateAddress struct {
					ID     graphql.ID     `graphql:"_id"`
					Street graphql.String `graphql:"street"`
					Number graphql.String `graphql:"number"`
					Door   graphql.String `graphql:"door"`
				} `graphql:"createAddress(data:{street: $Street, number: $Number, door: $Door})"`
			}

			y := map[string]interface{}{
				"Street": graphql.String(street),
				"Number": graphql.String(number),
				"Door":   graphql.String(door),
			}

			if err := call.Mutate(context.Background(), &a, y); err != nil {
				fmt.Fprintf(w, "error with address: %v\n", err)

			}

			var m struct {
				CreateCostumer struct {
					ID         graphql.ID      `graphql:"_id"`
					First      graphql.String  `graphql:"first"`
					Last       graphql.String  `graphql:"last"`
					Email      graphql.String  `graphql:"email"`
					Phone      graphql.String  `graphql:"phone"`
					Address    graphql.ID      `graphql:"address"`
					Registered graphql.Boolean `graphql:"registered"`
				} `graphql:"createOrder(data:{first: $First, last: $Last, email: $Email, phone: $Phone, address: $Address, registered: $Registered})"`
			}

			v := map[string]interface{}{
				"First":      graphql.String(first),
				"Last":       graphql.String(last),
				"Email":      graphql.String(email),
				"Phone":      graphql.String(phone),
				"Address":    a.CreateAddress.ID,
				"Registered": graphql.Boolean(registered),
			}

			if err := call.Mutate(context.Background(), &m, v); err != nil {
				fmt.Fprintf(w, "error with costumer: %v\n", err)

			}

			costumerID = m.CreateCostumer.ID

		} else {

			costumerID = q.CostumerByName.ID

		}

		var m1 struct {
			CreateOrder struct {
				ID       graphql.ID     `graphql:"_id"`
				Date     graphql.String `graphql:"date"`
				Costumer graphql.ID     `graphql:"costumer"`
				Cart     graphql.ID     `graphql:"cart"`
				Amount   graphql.Float  `graphql:"amount"`
			} `graphql:"createOrder(data:{date: $Date, costumer: $Costumer, cart: $Cart, amount: $Amount})"`
		}

		v1 := map[string]interface{}{
			"Date":     graphql.String(time.Now().UTC().Format("2006-01-02")),
			"Costumer": costumerID,
			"Cart":     graphql.ID(u),
			"Amount":   graphql.Float(total),
		}

		if err := call.Mutate(context.Background(), &m1, v1); err != nil {
			fmt.Fprintf(w, "error with order: %v\n", err)

		}

		var m2 struct {
			CreateStatus struct {
				ID       graphql.ID      `graphql:"_id"`
				Order    graphql.ID      `graphql:"order"`
				Payment  graphql.Boolean `graphql:"payment"`
				Delivery graphql.Boolean `graphql:"delivery"`
			} `graphql:"createStatus(data:{order: $Order, payment: $Payment, delivery: $Delivery})"`
		}

		v2 := map[string]interface{}{
			"Order":    m1.CreateOrder.ID,
			"payment":  graphql.Boolean(false),
			"delivery": graphql.Boolean(false),
		}

		if err := call.Mutate(context.Background(), &m2, v2); err != nil {
			fmt.Fprintf(w, "error with status: %v\n", err)

		}

		sum := strconv.FormatFloat(total, 'f', 2, 64)

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
