package main

import (
	"fmt"
	"net/http"
	//"os"
	"strconv"
	//"strings"

/* 	"golang.org/x/oauth2"

	//f "github.com/fauna/faunadb-go/faunadb"
	"github.com/plutov/paypal" */

)



func Handler(w http.ResponseWriter, r *http.Request) {

/* 	c, err := paypal.NewClient(os.Getenv("PP_ID"), os.Getenv("PP_SECRET"), paypal.APIBaseSandBox)

	if err != nil {

		fmt.Printf(w, err)
	} */

	switch r.Method {

	case "GET":

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
		src="https://www.paypal.com/sdk/js?client-id=AbBxx3BR2eA63A4i1g5rQduQ5K2LSqkybP7IdOAlTS65SoRfqwxqaEymvl5DHy183eUO1QQ8hqWwB9mE&currency=EUR">
	  	</script>

   		<div class="container" id="search" style="color:white; font-size:30px;">
		<form class="form-inline" role="form">
	   	<input class="form-control mr-sm-2" type="text" placeholder="Search" aria-label="Search" id ="find" name ="find">
	   	<button class="btn btn-outline-light my-2 my-sm-1" type="submit">Search</button><br>
		</div>
		<br>
		<div class="container" id="data" style="color:white;">
		<br>
		<form class="form-inline" role="form" method="POST">
		<input type="email" class="form-control" placeholder="name@example.com" aria-label="Email" id ="Email" name ="Email">
		<br>
		<input class="form-control mr-sm-2" type="password" placeholder="Secret" aria-label="Secret" id ="Secret" name ="Secret" value="">
		<input class="form-control mr-sm-2" type="text" placeholder="Title" aria-label="Title" id ="Title" name ="Title" required>
		<!--input class="form-control mr-sm-2" type="text" placeholder="Entry" aria-label="Entry" id ="Entry" name ="Entry" required-->
		<input class="form-control mr-sm-2" type="text" placeholder="Tags" aria-label="Tags" id ="Tags" name ="Tags">
		<input class="form-control mr-sm-2" tyoe="text" aria-label="Content" id ="Content" name ="Content" placeholder="Content"></textarea>
		<br>
		<button type="submit" class="btn btn-light">submit</button>
		</form>
		</div>
		`

		

		str = str + `
		<br>
		<br>
		<div class="container" id="paypal-button-container">
		</div>

	 	<script>
		 paypal.Buttons({
			createOrder: function(data, actions) {
			  // This function sets up the details of the transaction, including the amount and line item details.
			  return actions.order.create({
				purchase_units: [{
				  amount: {
					value: '55.00'
				  }
				}]
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
