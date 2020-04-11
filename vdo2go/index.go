package main

import (
	//"context"
	//"encoding/base64"
	"fmt"
	"net/http"
	"os"
	//"sort"
	"strconv"
	//"strings"

	"github.com/muxinc/mux-go"
	/* 	f "github.com/fauna/faunadb-go/faunadb"
	   	"github.com/shurcooL/graphql" *///"golang.org/x/oauth2"
	//"github.com/plutov/paypal"
)

type Access struct {
	//Reference *f.RefV `fauna:"ref"`
	Timestamp int    `fauna:"ts"`
	Secret    string `fauna:"secret"`
	Role      string `fauna:"role"`
}

/* type ProductEntry struct {
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

type SourceEntry struct {
	ID     graphql.ID     `graphql:"_id"`
	Link   graphql.ID     `graphql:"link"`
	Origin graphql.String `graphql:"origin"`
} */

func Handler(w http.ResponseWriter, r *http.Request) {

	client := muxgo.NewAPIClient(
		muxgo.NewConfiguration(
			muxgo.WithBasicAuth(os.Getenv("MUX_ID"), os.Getenv("MUX_SECRET")),
		))

	assets, err := client.AssetsApi.ListAssets()

	if err != nil {

		fmt.Fprintf(w, "something went wrong", err)

	}

	for _, a := range assets.Data {

		fmt.Fprintf(w, "%v\n%v\n%v", a.Id, a.CreatedAt, a.PlaybackIds)
	}

}
