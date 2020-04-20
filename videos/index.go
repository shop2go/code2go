package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	f "github.com/fauna/faunadb-go/faunadb"
	"github.com/muxinc/mux-go"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

type Access struct {
	//Reference *f.RefV `fauna:"ref"`
	Timestamp int    `fauna:"ts"`
	Secret    string `fauna:"secret"`
	Role      string `fauna:"role"`
}

type AssetEntry struct {
	ID       graphql.ID     `graphql:"_id"`
	SourceID graphql.String `graphql:"sourceID"`
	AssetID  graphql.String `graphql:"assetID"`
	PbID     graphql.String `graphql:"pbID"`
	First    graphql.String `graphql:"first"`
	Last     graphql.String `graphql:"last"`
	Email    graphql.String `graphql:"email"`
	Title    graphql.String `graphql:"title"`
	Category graphql.String `graphql:"category"`
	Content  graphql.String `graphql:"content"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	m := make(map[graphql.ID]string, 0)

	id := r.Host

	id = strings.TrimSuffix(id, "code2go.dev")

	id = strings.TrimSuffix(id, ".")

	switch id {

	case "":

		fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

		x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database("assets"), "role": "server"}))

		if err != nil {

			fmt.Fprintf(w, "a connection error occured: %v\n", err)

		}

		var access *Access

		x.Get(&access)

		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: access.Secret},
		)

		access = &Access{}

		httpClient := oauth2.NewClient(context.Background(), src)

		caller := graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

		client := muxgo.NewAPIClient(
			muxgo.NewConfiguration(
				muxgo.WithBasicAuth(os.Getenv("MUX_ID"), os.Getenv("MUX_SECRET")),
			))

		assets, err := client.AssetsApi.ListAssets()

		if err != nil {

			fmt.Fprintf(w, "%v", err)

		}

		for _, a := range assets.Data {

			var q struct {
				AssetByAssetID struct {
					AssetEntry
				} `graphql:"assetByAssetID(assetID: $AssetID)"`
			}

			v := map[string]interface{}{
				"AssetID": graphql.String(a.Id),
			}

			if err := caller.Query(context.Background(), &q, v); err != nil {
				fmt.Fprintf(w, "error with asset source: %v\n", err)
			}



			m[q.AssetByAssetID.ID] = string(q.AssetByAssetID.Category)

		}

	}

	switch r.Method {

	case "GET":

		switch id {

		case "":

			fmt.Fprint(w, m)

		}

	}

}
