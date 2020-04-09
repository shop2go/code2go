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
		<img class="mr-3" src="https://assets.medienwerk.now.sh/love.svg">
		<br>
		
			<p><h2>€ ` + sum + `</h2>Einkaufsumme<p>
			<br>
			<button type="button" class="btn btn-light" onclick="window.location.href='`+node+`'">Mit dem Einkauf fortfahren</button>

			<br><br>

		<form role="form" method="POST">

		<ul class="list-group">
		<li class="list-group-item">
		<button type="submit" class="btn btn-light" style="font-size:30px;">Bezahlen</button>
		<br><br>
		<h4>per eps, Sofort., Debit-/Kreditkarte oder PayPal.</h4>
		<p>Transaktionshandler:</p>
		<img src="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjMyIiB2aWV3Qm94PSIwIDAgMTAwIDMyIiBwcmVzZXJ2ZUFzcGVjdFJhdGlvPSJ4TWluWU1pbiBtZWV0IiB4bWxucz0iaHR0cDomI3gyRjsmI3gyRjt3d3cudzMub3JnJiN4MkY7MjAwMCYjeDJGO3N2ZyI+PHBhdGggZmlsbD0iIzAwMzA4NyIgZD0iTSAxMi4yMzcgMi40NDQgTCA0LjQzNyAyLjQ0NCBDIDMuOTM3IDIuNDQ0IDMuNDM3IDIuODQ0IDMuMzM3IDMuMzQ0IEwgMC4yMzcgMjMuMzQ0IEMgMC4xMzcgMjMuNzQ0IDAuNDM3IDI0LjA0NCAwLjgzNyAyNC4wNDQgTCA0LjUzNyAyNC4wNDQgQyA1LjAzNyAyNC4wNDQgNS41MzcgMjMuNjQ0IDUuNjM3IDIzLjE0NCBMIDYuNDM3IDE3Ljc0NCBDIDYuNTM3IDE3LjI0NCA2LjkzNyAxNi44NDQgNy41MzcgMTYuODQ0IEwgMTAuMDM3IDE2Ljg0NCBDIDE1LjEzNyAxNi44NDQgMTguMTM3IDE0LjM0NCAxOC45MzcgOS40NDQgQyAxOS4yMzcgNy4zNDQgMTguOTM3IDUuNjQ0IDE3LjkzNyA0LjQ0NCBDIDE2LjgzNyAzLjE0NCAxNC44MzcgMi40NDQgMTIuMjM3IDIuNDQ0IFogTSAxMy4xMzcgOS43NDQgQyAxMi43MzcgMTIuNTQ0IDEwLjUzNyAxMi41NDQgOC41MzcgMTIuNTQ0IEwgNy4zMzcgMTIuNTQ0IEwgOC4xMzcgNy4zNDQgQyA4LjEzNyA3LjA0NCA4LjQzNyA2Ljg0NCA4LjczNyA2Ljg0NCBMIDkuMjM3IDYuODQ0IEMgMTAuNjM3IDYuODQ0IDExLjkzNyA2Ljg0NCAxMi42MzcgNy42NDQgQyAxMy4xMzcgOC4wNDQgMTMuMzM3IDguNzQ0IDEzLjEzNyA5Ljc0NCBaIj48L3BhdGg+PHBhdGggZmlsbD0iIzAwMzA4NyIgZD0iTSAzNS40MzcgOS42NDQgTCAzMS43MzcgOS42NDQgQyAzMS40MzcgOS42NDQgMzEuMTM3IDkuODQ0IDMxLjEzNyAxMC4xNDQgTCAzMC45MzcgMTEuMTQ0IEwgMzAuNjM3IDEwLjc0NCBDIDI5LjgzNyA5LjU0NCAyOC4wMzcgOS4xNDQgMjYuMjM3IDkuMTQ0IEMgMjIuMTM3IDkuMTQ0IDE4LjYzNyAxMi4yNDQgMTcuOTM3IDE2LjY0NCBDIDE3LjUzNyAxOC44NDQgMTguMDM3IDIwLjk0NCAxOS4zMzcgMjIuMzQ0IEMgMjAuNDM3IDIzLjY0NCAyMi4xMzcgMjQuMjQ0IDI0LjAzNyAyNC4yNDQgQyAyNy4zMzcgMjQuMjQ0IDI5LjIzNyAyMi4xNDQgMjkuMjM3IDIyLjE0NCBMIDI5LjAzNyAyMy4xNDQgQyAyOC45MzcgMjMuNTQ0IDI5LjIzNyAyMy45NDQgMjkuNjM3IDIzLjk0NCBMIDMzLjAzNyAyMy45NDQgQyAzMy41MzcgMjMuOTQ0IDM0LjAzNyAyMy41NDQgMzQuMTM3IDIzLjA0NCBMIDM2LjEzNyAxMC4yNDQgQyAzNi4yMzcgMTAuMDQ0IDM1LjgzNyA5LjY0NCAzNS40MzcgOS42NDQgWiBNIDMwLjMzNyAxNi44NDQgQyAyOS45MzcgMTguOTQ0IDI4LjMzNyAyMC40NDQgMjYuMTM3IDIwLjQ0NCBDIDI1LjAzNyAyMC40NDQgMjQuMjM3IDIwLjE0NCAyMy42MzcgMTkuNDQ0IEMgMjMuMDM3IDE4Ljc0NCAyMi44MzcgMTcuODQ0IDIzLjAzNyAxNi44NDQgQyAyMy4zMzcgMTQuNzQ0IDI1LjEzNyAxMy4yNDQgMjcuMjM3IDEzLjI0NCBDIDI4LjMzNyAxMy4yNDQgMjkuMTM3IDEzLjY0NCAyOS43MzcgMTQuMjQ0IEMgMzAuMjM3IDE0Ljk0NCAzMC40MzcgMTUuODQ0IDMwLjMzNyAxNi44NDQgWiI+PC9wYXRoPjxwYXRoIGZpbGw9IiMwMDMwODciIGQ9Ik0gNTUuMzM3IDkuNjQ0IEwgNTEuNjM3IDkuNjQ0IEMgNTEuMjM3IDkuNjQ0IDUwLjkzNyA5Ljg0NCA1MC43MzcgMTAuMTQ0IEwgNDUuNTM3IDE3Ljc0NCBMIDQzLjMzNyAxMC40NDQgQyA0My4yMzcgOS45NDQgNDIuNzM3IDkuNjQ0IDQyLjMzNyA5LjY0NCBMIDM4LjYzNyA5LjY0NCBDIDM4LjIzNyA5LjY0NCAzNy44MzcgMTAuMDQ0IDM4LjAzNyAxMC41NDQgTCA0Mi4xMzcgMjIuNjQ0IEwgMzguMjM3IDI4LjA0NCBDIDM3LjkzNyAyOC40NDQgMzguMjM3IDI5LjA0NCAzOC43MzcgMjkuMDQ0IEwgNDIuNDM3IDI5LjA0NCBDIDQyLjgzNyAyOS4wNDQgNDMuMTM3IDI4Ljg0NCA0My4zMzcgMjguNTQ0IEwgNTUuODM3IDEwLjU0NCBDIDU2LjEzNyAxMC4yNDQgNTUuODM3IDkuNjQ0IDU1LjMzNyA5LjY0NCBaIj48L3BhdGg+PHBhdGggZmlsbD0iIzAwOWNkZSIgZD0iTSA2Ny43MzcgMi40NDQgTCA1OS45MzcgMi40NDQgQyA1OS40MzcgMi40NDQgNTguOTM3IDIuODQ0IDU4LjgzNyAzLjM0NCBMIDU1LjczNyAyMy4yNDQgQyA1NS42MzcgMjMuNjQ0IDU1LjkzNyAyMy45NDQgNTYuMzM3IDIzLjk0NCBMIDYwLjMzNyAyMy45NDQgQyA2MC43MzcgMjMuOTQ0IDYxLjAzNyAyMy42NDQgNjEuMDM3IDIzLjM0NCBMIDYxLjkzNyAxNy42NDQgQyA2Mi4wMzcgMTcuMTQ0IDYyLjQzNyAxNi43NDQgNjMuMDM3IDE2Ljc0NCBMIDY1LjUzNyAxNi43NDQgQyA3MC42MzcgMTYuNzQ0IDczLjYzNyAxNC4yNDQgNzQuNDM3IDkuMzQ0IEMgNzQuNzM3IDcuMjQ0IDc0LjQzNyA1LjU0NCA3My40MzcgNC4zNDQgQyA3Mi4yMzcgMy4xNDQgNzAuMzM3IDIuNDQ0IDY3LjczNyAyLjQ0NCBaIE0gNjguNjM3IDkuNzQ0IEMgNjguMjM3IDEyLjU0NCA2Ni4wMzcgMTIuNTQ0IDY0LjAzNyAxMi41NDQgTCA2Mi44MzcgMTIuNTQ0IEwgNjMuNjM3IDcuMzQ0IEMgNjMuNjM3IDcuMDQ0IDYzLjkzNyA2Ljg0NCA2NC4yMzcgNi44NDQgTCA2NC43MzcgNi44NDQgQyA2Ni4xMzcgNi44NDQgNjcuNDM3IDYuODQ0IDY4LjEzNyA3LjY0NCBDIDY4LjYzNyA4LjA0NCA2OC43MzcgOC43NDQgNjguNjM3IDkuNzQ0IFoiPjwvcGF0aD48cGF0aCBmaWxsPSIjMDA5Y2RlIiBkPSJNIDkwLjkzNyA5LjY0NCBMIDg3LjIzNyA5LjY0NCBDIDg2LjkzNyA5LjY0NCA4Ni42MzcgOS44NDQgODYuNjM3IDEwLjE0NCBMIDg2LjQzNyAxMS4xNDQgTCA4Ni4xMzcgMTAuNzQ0IEMgODUuMzM3IDkuNTQ0IDgzLjUzNyA5LjE0NCA4MS43MzcgOS4xNDQgQyA3Ny42MzcgOS4xNDQgNzQuMTM3IDEyLjI0NCA3My40MzcgMTYuNjQ0IEMgNzMuMDM3IDE4Ljg0NCA3My41MzcgMjAuOTQ0IDc0LjgzNyAyMi4zNDQgQyA3NS45MzcgMjMuNjQ0IDc3LjYzNyAyNC4yNDQgNzkuNTM3IDI0LjI0NCBDIDgyLjgzNyAyNC4yNDQgODQuNzM3IDIyLjE0NCA4NC43MzcgMjIuMTQ0IEwgODQuNTM3IDIzLjE0NCBDIDg0LjQzNyAyMy41NDQgODQuNzM3IDIzLjk0NCA4NS4xMzcgMjMuOTQ0IEwgODguNTM3IDIzLjk0NCBDIDg5LjAzNyAyMy45NDQgODkuNTM3IDIzLjU0NCA4OS42MzcgMjMuMDQ0IEwgOTEuNjM3IDEwLjI0NCBDIDkxLjYzNyAxMC4wNDQgOTEuMzM3IDkuNjQ0IDkwLjkzNyA5LjY0NCBaIE0gODUuNzM3IDE2Ljg0NCBDIDg1LjMzNyAxOC45NDQgODMuNzM3IDIwLjQ0NCA4MS41MzcgMjAuNDQ0IEMgODAuNDM3IDIwLjQ0NCA3OS42MzcgMjAuMTQ0IDc5LjAzNyAxOS40NDQgQyA3OC40MzcgMTguNzQ0IDc4LjIzNyAxNy44NDQgNzguNDM3IDE2Ljg0NCBDIDc4LjczNyAxNC43NDQgODAuNTM3IDEzLjI0NCA4Mi42MzcgMTMuMjQ0IEMgODMuNzM3IDEzLjI0NCA4NC41MzcgMTMuNjQ0IDg1LjEzNyAxNC4yNDQgQyA4NS43MzcgMTQuOTQ0IDg1LjkzNyAxNS44NDQgODUuNzM3IDE2Ljg0NCBaIj48L3BhdGg+PHBhdGggZmlsbD0iIzAwOWNkZSIgZD0iTSA5NS4zMzcgMi45NDQgTCA5Mi4xMzcgMjMuMjQ0IEMgOTIuMDM3IDIzLjY0NCA5Mi4zMzcgMjMuOTQ0IDkyLjczNyAyMy45NDQgTCA5NS45MzcgMjMuOTQ0IEMgOTYuNDM3IDIzLjk0NCA5Ni45MzcgMjMuNTQ0IDk3LjAzNyAyMy4wNDQgTCAxMDAuMjM3IDMuMTQ0IEMgMTAwLjMzNyAyLjc0NCAxMDAuMDM3IDIuNDQ0IDk5LjYzNyAyLjQ0NCBMIDk2LjAzNyAyLjQ0NCBDIDk1LjYzNyAyLjQ0NCA5NS40MzcgMi42NDQgOTUuMzM3IDIuOTQ0IFoiPjwvcGF0aD48L3N2Zz4=">
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

	<h1>Einkauf</h1>
	

	<br>
	
		<p><h2>€ ` + sum + `</h2>Gesamtpreis<p>
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
