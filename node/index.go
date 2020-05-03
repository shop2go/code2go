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

	var cat string

	node := strings.TrimPrefix(r.URL.Path, "/")

	id := r.Host

	id = strings.TrimSuffix(id, "shop2go.cloud")

	id = strings.TrimSuffix(id, ".")

	fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

	x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database(node), "role": "server"}))

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

		if products[i].ProductCat < products[j].ProductCat {
			return true
		}

		if products[i].ProductCat > products[j].ProductCat {
			return false
		}

		if products[i].ProductPrice < products[j].ProductPrice {
			return true
		}

		if products[i].ProductPrice > products[j].ProductPrice {
			return false
		}

		return products[i].ProductName < products[j].ProductName

	})

	if id != "" {

		var q struct {
			FindCartByID struct {
				CartEntry
			} `graphql:"findCartByID(id: $ID)"`
		}

		v := map[string]interface{}{
			"ID": graphql.ID(id),
		}

		if err = call.Query(context.Background(), &q, v); err != nil {
			fmt.Fprintf(w, "error with products: %v\n", err)
		}

		if q.FindCartByID.Products != nil {

			m := make(map[graphql.ID]struct{}, 0)

			for _, id := range q.FindCartByID.Products {

				var q struct {
					FindProductByID struct {
						ProductEntry
					} `graphql:"findProductByID(id: $ID)"`
				}

				v = map[string]interface{}{
					"ID": id,
				}

				if err = call.Query(context.Background(), &q, v); err != nil {
					fmt.Fprintf(w, "error with products: %v\n", err)
				}

				total = total + float64(q.FindProductByID.ProductPrice)

				m[id] = struct{}{}

			}

			for i := 0; i < len(products); i++ {

				if _, ok := m[products[i].ID]; ok {

					products[i] = ProductEntry{}

				}

			}

		} else {

			http.Redirect(w, r, "https://shop2go.cloud/"+node, http.StatusSeeOther)

		}

	}

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
	`

	if id == "" {

		content = content +

		`
		<li class="list-group-item">
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
			<li class="list-group-item">
			<div class="media">
			<img class="mr-3" src="https://assets.medienwerk.now.sh/love.svg" width="100" >
					
			<div class="media-body"><br><br>

			<h3>In Stadt Salzburg innerhalb eines Tages an ihrer Haustür.</h3>				
			<br><p>Lieferpauschale:<br><h2>€ 5</h2><br>Bestellsumme:<br><h2>€ ` + price + `</h2></p>
			<button type="button" class="btn btn-light" onclick="window.location.href='order'">Ware jetzt bestellen</button>
			</div></div>		
			</li>
			<br><br>
			`

		} else {

			content = content +

			`
			<li class="list-group-item">
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

	if products[0].ID != nil {

		cat = string(products[0].ProductCat)

		id := fmt.Sprintf("%s", products[0].ID)
		price := strconv.FormatFloat(float64(products[0].ProductPrice), 'f', 2, 64)
		pack := strconv.Itoa(int(products[0].PackSize))

		dim := strconv.Itoa(int(products[0].LogoDim))
		//quant := int(products[0].ProductQuant)

		content = content + 
		
		`		 
		<h1 id="` + cat + `">` + cat + `</h1>
		<br>

		<li class="list-group-item">
		<div class="media">
		<img class="mr-3" src="` + string(products[0].ProductImgLink) + `" width="100">

		<div class="media-body">

		<h2>` + string(products[0].ProductName) + `</h2>

		<h4>` + string(products[0].ProductInfo) + `</h4>

		<p><h2>€ ` + price + `</h2>` + pack + ` ` + string(products[0].PackUnit) + `<br><br>

		<form class="form-inline" role="form" method="POST">
				
		<label class="form-check-label" for="` + id + `" style="font-size:25px;">Mengenauswahl:</label>
		
		<select style="font-size:30px;" class="form-control" id="` + id + `" name="` + id + `">
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
		<a href="` + string(products[0].ProductLink) + `" target="_blank"><img class="mr-3" src="` + string(products[0].LogoLink) + `" width="` + dim + `">
		</a>
		
		</div>
		</div>
		</li>
		<br><br>
		`

	}

	for k := 1; k < len(products); k++ {

		if products[k].ID != nil {

			if string(products[k].ProductCat) == cat {

				id := fmt.Sprintf("%s", products[k].ID)
				price := strconv.FormatFloat(float64(products[k].ProductPrice), 'f', 2, 64)
				pack := strconv.Itoa(int(products[k].PackSize))
				dim := strconv.Itoa(int(products[k].LogoDim))
				//quant := int(products[0].ProductQuant)

				content = content +

				`
				<br>
				<li class="list-group-item">
				<div class="media">
				<img class="mr-3" src="` + string(products[k].ProductImgLink) + `" width="100">
				
				<div class="media-body">

				<h2>` + string(products[k].ProductName) + `</h2>

				<h4>` + string(products[k].ProductInfo) + `</h4>
	
				<p><h2>€ ` + price + `</h2>` + pack + ` ` + string(products[k].PackUnit) + `<br><br>
		
				<form class="form-inline" role="form" method="POST">
		
				<label class="form-check-label" for="` + id + `" style="font-size:25px;">Mengenauswahl:</label>
			
				<select style="font-size:30px;" class="form-control" id="` + id + `" name="` + id + `">
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
				<a href="` + string(products[k].ProductLink) + `" target="_blank"><img class="mr-3" src="` + string(products[k].LogoLink) + `" width="` + dim + `"></a>
	
				</div>
				</div>
				</li>
				<br><br>
				`

			} else {

				cat = string(products[k].ProductCat)

				id := fmt.Sprintf("%s", products[k].ID)
				price := strconv.FormatFloat(float64(products[k].ProductPrice), 'f', 2, 64)
				pack := strconv.Itoa(int(products[k].PackSize))
				dim := strconv.Itoa(int(products[k].LogoDim))
				//quant := int(products[0].ProductQuant)

				content = content +

				`
				<h1 id="` + cat + `">` + cat + `</h1>
				<br>

				<li class="list-group-item">
				<div class="media">
				<img class="mr-3" src="` + string(products[k].ProductImgLink) + `" width="100">
			
				<div class="media-body">
			
				<h2>` + string(products[k].ProductName) + `</h2>
			
				<h4>` + string(products[k].ProductInfo) + `</h4>
			
				<p><h2>€ ` + price + `</h2>` + pack + ` ` + string(products[k].PackUnit) + `<br><br>
				
				<form class="form-inline" role="form" method="POST">
				
				<label class="form-check-label" for="` + id + `" style="font-size:25px;">Mengenauswahl:</label>
				
				<select style="font-size:30px;" class="form-control" id="` + id + `" name="` + id + `">
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
				<a href="` + string(products[k].ProductLink) + `" target="_blank"><img class="mr-3" src="` + string(products[k].LogoLink) + `" width="` + dim + `"></a>
			
				</div>
				</div>
				</li>
				<br><br>
				`

			}

		}

		content = content +
		
		`
		</ul>
		`

	}

	content = content +

	`			   
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

		//node = ok
		var cart CartEntry

		cart.Products = make([]graphql.ID, 0)

		r.ParseForm()

		//form parsing
		for k := 0; k < len(products); k++ {

			id := fmt.Sprintf("%s", products[k].ID)

			cnt := r.Form.Get(id)

			count, _ := strconv.Atoi(cnt)

			if count == 0 {

				continue

			} else {

				for i := 0; i < count; i++ {

					cart.Products = append(cart.Products, products[k].ID)

				}

			}

		}

		if id != "" {

			cart.ID = graphql.ID(id)

			var q struct {
				FindCartByID struct {
					CartEntry
				} `graphql:"findCartByID(id: $ID)"`
			}

			v := map[string]interface{}{
				"ID": cart.ID,
			}

			if err = call.Query(context.Background(), &q, v); err != nil {
				fmt.Fprintf(w, "error with cart: %v\n", err)
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

			v = map[string]interface{}{
				"ID":       cart.ID,
				"Products": cart.Products,
			}

			if err = call.Mutate(context.Background(), &m, v); err != nil {
				fmt.Fprintf(w, "error with cart: %v\n", err)

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
				fmt.Fprintf(w, "error with cart: %v\n", err)

			}

			cart.ID = m.CreateCart.ID

			//node == shop
			x, err = fc.Query(f.CreateKey(f.Obj{"database": f.Database("shop"), "role": "server"}))

			if err != nil {

				fmt.Fprintf(w, "connection error: %v\n", err)

			}

			x.Get(&access)

			src = oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: access.Secret},
			)

			httpClient = oauth2.NewClient(context.Background(), src)

			call = graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

			var s struct {
				CreateSource struct {
					SourceEntry
				} `graphql:"createSource(data:{origin: $Origin, link: $Link})"`
			}

			v = map[string]interface{}{
				"Origin": graphql.String(node),
				"Link":   cart.ID,
			}

			if err = call.Mutate(context.Background(), &s, v); err != nil {
				fmt.Fprintf(w, "error with source: %v\n", err)

			}

		}

		id = fmt.Sprintf("%s", cart.ID)

		http.Redirect(w, r, "https://"+id+".shop2go.cloud/"+node, http.StatusSeeOther)

	}

}
