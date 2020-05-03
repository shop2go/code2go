package main

import (
	"context"
	//"encoding/base64"
	"fmt"
	"net/http"
	"os"
	//"sort"
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

type ProductEntry struct {
	ID             graphql.ID     `graphql:"_id"`
	ProductImgLink graphql.String `graphql:"productImgLink"`
	ProductName    graphql.String `graphql:"productName"`
	ProductCat     graphql.String `graphql:"productCat"`
	ProductInfo    graphql.String `graphql:"productInfo"`
	ProductPrice   graphql.Float  `graphql:"productPrice"`
	PackSize       graphql.Int    `graphql:"packSize"`
	PackUnit       graphql.String `graphql:"packUnit"`
	ProductLink    graphql.String `graphql:"productLink"`
	LogoLink       graphql.String `graphql:"logoLink"`
	LogoDim        graphql.Int    `graphql:"logoDim"`
	ProductQuant   graphql.Int    `graphql:"productQuant"`
}

type CartEntry struct {
	ID       graphql.ID   `graphql:"_id"`
	Products []graphql.ID `graphql:"products"`
}

type SourceEntry struct {
	ID     graphql.ID     `graphql:"_id"`
	Link   graphql.ID     `graphql:"link"`
	Origin graphql.String `graphql:"origin"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	var total float64

	m := make(map[graphql.ID]int, 0)

	id := r.Host

	id = strings.TrimSuffix(id, "shop2go.cloud")

	id = strings.TrimSuffix(id, ".")

	if id == "" {

		http.Redirect(w, r, "https://shop2go.cloud", http.StatusSeeOther)

	}

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

	node := string(q.SourceByLink.Origin)

	if node != "" {

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

		var p struct {
			FindCartByID struct {
				CartEntry
			} `graphql:"findCartByID(id: $ID)"`
		}

		v := map[string]interface{}{
			"ID": graphql.ID(id),
		}

		if err = call.Query(context.Background(), &p, v); err != nil {
			fmt.Fprintf(w, "error with cart: %v\n", err)
		}

		if p.FindCartByID.Products != nil {

			//products = p.FindCartByID.Products

			for _, id := range p.FindCartByID.Products {

				if _, ok := m[id]; ok {

					m[id] = m[id] + 1

				} else {

					m[id] = 1

				}

				var n struct {
					FindProductByID struct {
						ProductEntry
					} `graphql:"findProductByID(id: $ID)"`
				}

				x := map[string]interface{}{
					"ID": id,
				}

				if err = call.Query(context.Background(), &n, x); err != nil {
					fmt.Fprintf(w, "error with products: %v\n", err)
				}

				total = total + float64(n.FindProductByID.ProductPrice)

			}

			/* for _, id := range p.FindCartByID.Products {

				var n struct {
					FindProductByID struct {
						ProductEntry
					} `graphql:"findProductByID(id: $ID)"`
				}

				x := map[string]interface{}{
					"ID": id,
				}

				if err = call.Query(context.Background(), &n, x); err != nil {
					fmt.Fprintf(w, "error with products: %v\n", err)
				}

				total = total + float64(n.FindProductByID.Price)

				if _, ok := m[id]; ok {

					m[id] = m[id] + 1

				} else {

					m[id] = 1

				}

				//products = append(products, n.FindProductByID.ID)

			} */

			/* for i := 0; i < len(products); i++ {

				if _, ok := m[products[i]]; !ok {

					products[i] = ProductEntry{}

				}

			} */

		} else {

			http.Redirect(w, r, "https://shop2go.cloud/"+node, http.StatusSeeOther)

		}

	} else {

		http.Redirect(w, r, "https://shop2go.cloud/"+node, http.StatusSeeOther)

	}

	var cat string

	content :=

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

			<div class="container" id="shop" style="color:rgb(255, 255, 255); font-size:30px;">

		<ul class="list-group">
		<br>
		<br>

		<h1>Bestellung</h1>

		<li class="list-group-item">

		`

	if id == "" {

		content = content +

			`
			<div class="media">
			<img class="mr-3" src="https://assets.medienwerk.now.sh/love.svg" width="100" >
					
			<div class="media-body"><br><br>
					
			<h3>In Stadt Salzburg innerhalb eines Tages an ihrer Haustür.</h3>
					
			<p><h2>€ 5</h2>Pauschal<br>Mindesteinkaufsumme: € 14</p>
			</div>
			</div>	
			</li>
			<br><br>
			`

	} else {

		if total >= 14 {

			price := strconv.FormatFloat(total, 'f', 2, 64)

			content = content +

				`				
					<div class="media">
					<img class="mr-3" src="https://assets.medienwerk.now.sh/love.svg" width="100" >
							
					<div class="media-body"><br><br>			
				<h3>In Stadt Salzburg innerhalb eines Tages an ihrer Haustür.</h3>				
				<p><br><br><h2>€ 5</h2>Bestellsumme: <h2>€ ` + price + `</h2></p>
				<button type="button" class="btn btn-light" onclick="window.location.href='order'">Ware jetzt bestellen</button>
				</div></div>	
				</li>
				<br><br>
				`

		} else {

			content = content +

				`
			<div class="media">
			<img class="mr-3" src="https://assets.medienwerk.now.sh/love.svg" width="100" >
					
			<div class="media-body"><br><br>
					
			<h3>In Stadt Salzburg innerhalb eines Tages an ihrer Haustür.</h3>
					
			<p><h2>€ 5</h2>Pauschal<br><h2>Mindesteinkaufsumme: € 14<h2></p>
			</div>
			</div>	
			</li>
			<br><br>
			`

		}

	}

	for k, i := range m {

		var q struct {
			FindProductByID struct {
				ProductEntry
			} `graphql:"findProductByID(id: $ID)"`
		}

		v := map[string]interface{}{
			"ID": k,
		}

		if err = call.Query(context.Background(), &q, v); err != nil {
			fmt.Fprintf(w, "error with products: %v\n", err)
		}

		if string(q.FindProductByID.ProductCat) != cat {

			cat = string(q.FindProductByID.ProductCat)

			content = content + ` 

		<br>
		<h1>` + cat + `</h1>

		`

		}

		//if products[0].ID != nil {
		id := fmt.Sprintf("%s", q.FindProductByID.ID)
		price := strconv.FormatFloat(float64(q.FindProductByID.ProductPrice), 'f', 2, 64)
		pack := strconv.Itoa(int(q.FindProductByID.PackSize))
		dim := strconv.Itoa(int(q.FindProductByID.LogoDim))

		content = content + ` 

		<li class="list-group-item">

		<div class="media">
		<img class="mr-3" src="` + string(q.FindProductByID.ProductImgLink) + `" width="100">

		<div class="media-body">

		<h2>` + string(q.FindProductByID.ProductName) + `</h2>

		<h4>` + string(q.FindProductByID.ProductInfo) + `</h4>

		<p><h2>€ ` + price + `</h2>` + pack + ` ` + string(q.FindProductByID.PackUnit) + `<br><br>

		<form class="form-inline" role="form" method="POST">
				
		<label class="form-check-label" for="` + id + `" style="font-size:25px;">Mengenauswahl:</label>
		
		<select style="font-size:30px;" class="form-control" id="` + id + `" name="` + id + `">
			
		`
		//if j, ok := m[q.FindProductByID.ID]; ok {

		i = i + 3

		for i > 0 {

			o := strconv.Itoa(i)

			content = content + `
				
				
				<option>` + o + `</option>
				
				`
			i--

		}

		content = content + `
		</select>

		<button type="submit" class="btn btn-light">ändern</button>
		
		</form>
		</p>
		<br>
		<a href="` + string(q.FindProductByID.ProductLink) + `" target="_blank"><img class="mr-3" src="` + string(q.FindProductByID.LogoLink) + `" width="` + dim + `">
		</a>
		
		</div>
		</div>
		</li>
		<br><br>
		`

	}

	content = content + `

		</ul>
		</div>
					   
		<script src="https://assets.medienwerk.now.sh/material.min.js"></script>
		</body>
		</html>
		`

	switch r.Method {

	case "GET":

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(content)))
		w.Write([]byte(content))

	case "POST":

		products := make([]graphql.ID, 0)

		r.ParseForm()

		for k, i := range m {

			id := fmt.Sprintf("%s", k)

			cnt := r.Form.Get(id)

			if cnt != "" {

				count, _ := strconv.Atoi(cnt)

				m[k] = count

				for count > 0 {

					products = append(products, k)

					count--

				}

				continue

			}

			if cnt != "0" {

				for i > 0 {

					products = append(products, k)

					i--
				}

			}

		}

		var m struct {
			UpdateCart struct {
				CartEntry
			} `graphql:"updateCart(id: $ID, data:{products: $Products})"`
		}

		v := map[string]interface{}{
			"ID":       graphql.ID(id),
			"Products": products,
		}

		if err = call.Mutate(context.Background(), &m, v); err != nil {
			fmt.Fprintf(w, "error with products: %v\n", err)

		}

		http.Redirect(w, r, "https://"+id+".shop2go.cloud/order", http.StatusSeeOther)

	}

}
