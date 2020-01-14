package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	//"strconv"
	//"strings"

	"golang.org/x/oauth2"

	//f "github.com/fauna/faunadb-go/faunadb"
	"github.com/plutov/paypal"
)

func Handler(w http.ResponseWriter, r *http.Request) {

		_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, err)
	}

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

	c.CreateOrder()

}
