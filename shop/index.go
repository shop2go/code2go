package main

import (
	"context"
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

	"github.com/google/uuid"
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

func Handler(w http.ResponseWriter, r *http.Request) {

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
		fmt.Fprintf(w, "error with getting products: %v\n", err)
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

	/* 	c := make([]string, 0)

	var s string */

	switch r.Method {

	case "GET":

		str := `<!DOCTYPE html>
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
						   <br><br>
	
						<div class="container" id="data" style="color:rgb(255, 255, 255); font-size:30px;"></div>

						<h1>Zustellung</h1>
						<li class="list-group-item">
						<div class="media">
				  <img class="mr-3" src="https://assets.medienwerk.now.sh/love.svg">
				  <div class="media-body"><br><br>
					  <h3>	
						In 5020 Salzburg innert eines Tages an die Haustür
					</h3>
					<h4></h4>
						<p><h2>€ 6</h2>Pauschal<br>Mindesteinkaufsumme: € 14</p>
						
				</div>
				
						</li>
<br>
`

		f := strconv.FormatFloat(float64(products[0].Price),'f',10,64)
		
		//c = append(c, s)

		str = str + 
		`
		<h1>` + string(products[0].Cat) + `</h1>
		<li class="list-group-item">
	<div class="media">
	<img class="mr-3" src="`+string(products[0].ImgURL) +`" width="200">
	<div class="media-body">
	<h2>`+string(products[0].Product)+`</h2>
	<h4>`+ string(products[0].Info)+ `</h4>
	<p><h2>€ `+ f +`</h2>`+strconv.Itoa(int(products[0].Pack+` Gramm<br><br>
		<form class="form-inline" role="form" method="POST">
		
			<label class="form-check-label" for="`+string(products[0].Product)+`" style="font-size:25px;">Mengenauswahl:   _____________   +	</label>
			
			<select style="font-size:30px;" class="form-control" id="`+string(products[0].Product)+`">
				<option>0</option>
				<option>1</option>
				<option>2</option>
				<option>3</option>
				<option>4</option>
				<option>5</option>
				<option>6</option>
				<option>7</option>
				<option>8</option>
				<option>9</option>
			  </select><br>
			  <button type="submit" class="btn btn-light">Einkaufliste</button>

			</form>
			</p>
	
			<a href="`+string(products[0].InfoURL)+`" target="_blank"><img class="mr-3" src="`+string(products[0].LinkURL)+`" width="`+strconv.Itoa(int(products[0].LinkDIM))+`"></a>
	</div>
	</div>
	
	</li>
	<br> 
	`
	
	     for i := 1; i < len(products); i++ {

			if string(products[i].Cat) == s {

				f = strconv.FormatFloat(float64(products[i].Price),'f',10,64)

				str = str +
				`
		<li class="list-group-item">
	<div class="media">
	<img class="mr-3" src="`+string(products[i].ImgURL) +`" width="200">
	<div class="media-body">
	<h2>`+string(products[i].Product)+`</h2>
	<h4>`+ string(products[i].Info)+ `</h4>
	<p><h2>€ `+ f +`</h2>`+strconv.Itoa(int(products[i].Pack+` Gramm<br><br>
		<form class="form-inline" role="form" method="POST">
		
			<label class="form-check-label" for="`+string(products[i].Product)+`" style="font-size:25px;">Mengenauswahl:   _____________   +	</label>
			
			<select style="font-size:30px;" class="form-control" id="`+string(products[i].Product)+`">
				<option>0</option>
				<option>1</option>
				<option>2</option>
				<option>3</option>
				<option>4</option>
				<option>5</option>
				<option>6</option>
				<option>7</option>
				<option>8</option>
				<option>9</option>
			  </select><br>
			  <button type="submit" class="btn btn-light">Einkaufliste</button>

			</form>
			</p>
	
			<a href="`+string(products[i].InfoURL)+`" target="_blank"><img class="mr-3" src="`+string(products[i].LinkURL)+`" width="`+strconv.Itoa(int(products[i].LinkDIM))+`"></a>
	</div>
	</div>
	
	</li>
	<br>
	`
				continue

			} else {

				f = strconv.FormatFloat(float64(products[i].Price),'f',10,64)

				str = str + 
				`
				<h1>` + string(products[i].Cat) + `</h1>
				<li class="list-group-item">
			<div class="media">
			<img class="mr-3" src="`+string(products[i].ImgURL) +`" width="200">
			<div class="media-body">
			<h2>`+string(products[i].Product)+`</h2>
			<h4>`+ string(products[i].Info)+ `</h4>
			<p><h2>€ `+ f +`</h2>`+strconv.Itoa(int(products[i].Pack+` Gramm<br><br>
				<form class="form-inline" role="form" method="POST">
				
					<label class="form-check-label" for="`+string(products[i].Product)+`" style="font-size:25px;">Mengenauswahl:   _____________   +	</label>
					
					<select style="font-size:30px;" class="form-control" id="`+string(products[i].Product)+`">
						<option>0</option>
						<option>1</option>
						<option>2</option>
						<option>3</option>
						<option>4</option>
						<option>5</option>
						<option>6</option>
						<option>7</option>
						<option>8</option>
						<option>9</option>
					  </select><br>
					  <button type="submit" class="btn btn-light">Einkaufliste</button>
					  
					</form>
					</p>
			
					<a href="`+string(products[i].InfoURL)+`" target="_blank"><img class="mr-3" src="`+string(products[i].LinkURL)+`" width="`+strconv.Itoa(int(products[i].LinkDIM))+`"></a>
			</div>
			</div>
			
			</li>
			<br>
			`

				continue

			}

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

	}

}
