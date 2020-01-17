package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func Handler(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		//email := r.FormValue("Email")
		count := r.Form.Get("Count")
		//price := r.FormValue("Price")

		i, err := strconv.Atoi(count)

		if err != nil {

			fmt.Fprint(w, err)

		}

		price := i * 50

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
	   <br>
	   <br>
	<div class="container" id="paypal-button-container">
	</div>

	<script>
paypal.Buttons({
createOrder: function(data, actions) {
  return actions.order.create({
	purchase_units: [{
	  amount: {
		  value: '` + strconv.Itoa(price) + `'
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

		/* var payer *paypal.CreateOrderPayer

		if email != "" {

			payer.EmailAddress = email

		}

		var purchase []paypal.PurchaseUnitRequest = make([]paypal.PurchaseUnitRequest, 1)

		purchase[0].Amount.Currency = "EUR"
		purchase[0].Amount.Value = strconv.Itoa(price)

		c, err := paypal.NewClient("AbBxx3BR2eA63A4i1g5rQduQ5K2LSqkybP7IdOAlTS65SoRfqwxqaEymvl5DHy183eUO1QQ8hqWwB9mE", os.Getenv("PP_SECRET"), paypal.APIBaseSandBox)

		if err != nil {

			fmt.Fprint(w, err)

		}

		accessToken, err := c.GetAccessToken()

		if err != nil {

			fmt.Fprint(w, err)

		}

		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken.Token},
		)

		httpClient := oauth2.NewClient(context.Background(), src)

		c.SetHTTPClient(httpClient)

		order, err := c.CreateOrder("CAPTURE", purchase, payer, nil)

		if err != nil {

			fmt.Fprint(w, err)

		}

		fmt.Fprint(w, order.ID)
		*/
	

}
