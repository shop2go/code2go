package main

import (
	"context"
	//"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"sort"
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

	//var ID []byte

	u := r.Host

	u = strings.TrimSuffix(u, "code2go.dev")

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

	call := graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

	var q struct {
		AllProducts struct {
			Data []ProductEntry
		}
	}

	if err = call.Query(context.Background(), &q, nil); err != nil {
		fmt.Fprintf(w, "error with products: %v\n", err)
	}

	products := q.AllProducts.Data

	sort.Slice(products, func(i, j int) bool {

		if products[i].Cat < products[j].Cat {
			return true
		}

		if products[i].Cat > products[j].Cat {
			return false
		}

		if products[i].Price < products[j].Price {
			return true
		}

		if products[i].Price > products[j].Price {
			return false
		}

		return products[i].Product < products[j].Product

	})

	if u != "" {

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

		if q.FindCartByID.Products != nil {

			m := make(map[graphql.ID]struct{}, 0)

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

				m[id] = struct{}{}

			}

			for i := 0; i < len(products); i++ {

				if _, ok := m[products[i].ID]; ok {

					products[i] = ProductEntry{}

				}

			}

		} else {

			http.Redirect(w, r, "https://code2go.dev/shop", http.StatusSeeOther)

		}

	}

	var s string

	switch r.Method {

	case "GET":

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

		<div class="container" id="shop" style="color:rgb(255, 255, 255); font-size:30px;">

		<ul class="list-group">
		<br>
		<br>

		<h1>Zustellung</h1>

		<li class="list-group-item">

		`

		if u != "" {

			if total >= 14 {

				price := strconv.FormatFloat(total, 'f', 2, 64)

				str = str +

					`				
							
				<h3>In Stadt Salzburg innerhalb eines Tages an ihrer Haustür.</h3>				
				<p>Pauschal<br><h2>€ 5</h2> + Bestellsumme: <h2>€ ` + price + `</h2></p>
				<button type="button" class="btn btn-light" onclick="window.location.href='order'">Ware jetzt bestellen</button>
						
				</li>
				<br><br>
				`

			} else {

				str = str +

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

		} else {

			str = str +

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

		}

		if products[0].ID != nil {

			s = string(products[0].Cat)

			price := strconv.FormatFloat(float64(products[0].Price), 'f', 2, 64)
			pack := strconv.Itoa(int(products[0].Pack))
			dim := strconv.Itoa(int(products[0].LinkDIM))

			str = str + ` 

		<br>
		<h1>` + s + `</h1>

		<li class="list-group-item">

		<div class="media">
		<img class="mr-3" src="` + string(products[0].ImgURL) + `" width="100">

		<div class="media-body">

		<h2>` + string(products[0].Product) + `</h2>

		<h4>` + string(products[0].Info) + `</h4>

		<p><h2>€ ` + price + `</h2>` + pack + ` Gramm<br><br>

		<form class="form-inline" role="form" method="POST">
				
		<label class="form-check-label" for="` + string(products[0].Product) + `" style="font-size:25px;">Mengenauswahl:</label>
		
		<select style="font-size:30px;" class="form-control" id="` + string(products[0].Product) + `" name="` + string(products[0].Product) + `">

			<option>1</option>
			<option>2</option>
			<option>3</option>
			<option>4</option>
			<option>5</option>
			<option>6</option>
			<option>7</option>
			<option>8</option>
			<option>9</option>
		</select>

		<button type="submit" class="btn btn-light">nehmen</button>
		  
		</form>
		</p>
		<br>
		<a href="` + string(products[0].InfoURL) + `" target="_blank"><img class="mr-3" src="` + string(products[0].LinkURL) + `" width="` + dim + `">
		</a>
		
		</div>
		</div>
		</li>
		<br><br>
		`

		}

		for k := 1; k < len(products); k++ {

			if products[k].ID != nil {

				if string(products[k].Cat) == s {

					price := strconv.FormatFloat(float64(products[k].Price), 'f', 2, 64)
					pack := strconv.Itoa(int(products[k].Pack))
					dim := strconv.Itoa(int(products[k].LinkDIM))

					str = str +
						`
				
				<li class="list-group-item">

				<div class="media">
				<img class="mr-3" src="` + string(products[k].ImgURL) + `" width="100">
				
				<div class="media-body">

				<h2>` + string(products[k].Product) + `</h2>

				<h4>` + string(products[k].Info) + `</h4>
	
				<p><h2>€ ` + price + `</h2>` + pack + ` Gramm<br><br>
		
				<form class="form-inline" role="form" method="POST">
		
				<label class="form-check-label" for="` + string(products[k].Product) + `" style="font-size:25px;">Mengenauswahl:</label>
			
				<select style="font-size:30px;" class="form-control" id="` + string(products[k].Product) + `" name="` + string(products[k].Product) + `">

					<option>1</option>
					<option>2</option>
					<option>3</option>
					<option>4</option>
					<option>5</option>
					<option>6</option>
					<option>7</option>
					<option>8</option>
					<option>9</option>
				</select>

			  	<button type="submit" class="btn btn-light">nehmen</button>

				</form>
				</p>
				<br>
				<a href="` + string(products[k].InfoURL) + `" target="_blank"><img class="mr-3" src="` + string(products[k].LinkURL) + `" width="` + dim + `"></a>
	
				</div>
				</div>
				</li>
				<br><br>
				`

				} else {

					s = string(products[k].Cat)

					price := strconv.FormatFloat(float64(products[k].Price), 'f', 2, 64)
					pack := strconv.Itoa(int(products[k].Pack))
					dim := strconv.Itoa(int(products[k].LinkDIM))

					str = str +
						`
				<br>
				<h1>` + s + `</h1>

				<li class="list-group-item">

				<div class="media">
				<img class="mr-3" src="` + string(products[k].ImgURL) + `" width="100">
			
				<div class="media-body">
			
				<h2>` + string(products[k].Product) + `</h2>
			
				<h4>` + string(products[k].Info) + `</h4>
			
				<p><h2>€ ` + price + `</h2>` + pack + ` Gramm<br><br>
				
				<form class="form-inline" role="form" method="POST">
				
				<label class="form-check-label" for="` + string(products[k].Product) + `" style="font-size:25px;">Mengenauswahl:</label>
				
				<select style="font-size:30px;" class="form-control" id="` + string(products[k].Product) + `" name="` + string(products[k].Product) + `">

					<option>1</option>
					<option>2</option>
					<option>3</option>
					<option>4</option>
					<option>5</option>
					<option>6</option>
					<option>7</option>
					<option>8</option>
					<option>9</option>
				</select>
					 
				<button type="submit" class="btn btn-light">nehmen</button>
					  
				</form>
				</p>
				<br>
				<a href="` + string(products[k].InfoURL) + `" target="_blank"><img class="mr-3" src="` + string(products[k].LinkURL) + `" width="` + dim + `"></a>
			
				</div>
				</div>
				</li>
				<br><br>
				`

				}

			}

			str = str + `

		</ul>

		`

		}

		str = str + `
					   
		<script src="https://assets.medienwerk.now.sh/material.min.js"></script>
		</body>
		</html>
		`

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(str)))
		w.Write([]byte(str))

	case "POST":

		var cart CartEntry

		cart.Products = make([]graphql.ID, 0)

		r.ParseForm()

		//form parsing
		for k := 0; k < len(products); k++ {

			cnt := r.Form.Get(string(products[k].Product))

			count, _ := strconv.Atoi(cnt)

			if count == 0 {

				continue

			} else {

				for l := 0; l < count; l++ {

					cart.Products = append(cart.Products, products[k].ID)

				}

			}

		}

		//if len(cart.Products) == 0 {

		if u != "" {

			cart.ID = graphql.ID(u)

			var q struct {
				FindCartByID struct {
					CartEntry
				} `graphql:"findCartByID(id: $ID)"`
			}

			doc := map[string]interface{}{
				"ID": cart.ID,
			}

			if err = call.Query(context.Background(), &q, doc); err != nil {
				fmt.Fprintf(w, "error with products: %v\n", err)
			}

			// appending additional products
			for _, p := range q.FindCartByID.Products {

				cart.Products = append(cart.Products, p)

			}

			var m struct {
				UpdateCart struct {
					CartEntry
				} `graphql:"updateCart(id: $ID, data:{products: $Products})"`
			}

			v := map[string]interface{}{
				"ID":       cart.ID,
				"Products": cart.Products,
			}

			if err = call.Mutate(context.Background(), &m, v); err != nil {
				fmt.Fprintf(w, "error with products: %v\n", err)

			}

		} else {

			var m struct {
				CreateCart struct {
					CartEntry
				} `graphql:"createCart(data:{products: $Products})"`
			}

			v := map[string]interface{}{
				"Products": cart.Products,
			}

			if err = call.Mutate(context.Background(), &m, v); err != nil {
				fmt.Fprintf(w, "error with products: %v\n", err)

			}

			cart.ID = m.CreateCart.ID

		}

		s = fmt.Sprintf("%s", cart.ID)

		//code := base64.StdEncoding.EncodeToString([]byte(s))

		http.Redirect(w, r, "https://"+s+".code2go.dev/shop", http.StatusSeeOther)

	}

}
