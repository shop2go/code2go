package main

import (
	"context"
	//"encoding/base64"
	"fmt"
	"net/http"
	"os"
	//"sort"
	//"strconv"
	//"strings"

	f "github.com/fauna/faunadb-go/faunadb"
	"github.com/muxinc/mux-go"
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

type AssetEntry struct {
	ID         graphql.ID     `graphql:"_id"`
	AssetID    graphql.String `graphql:"assetId"`
	PlaybackID graphql.String `graphql:"playbackId"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

	x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database("2020"), "role": "server"}))

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

	/* 	var q struct {
	   		AllAssets struct {
	   			Data []AssetEntry
	   		}
	   	}

	   	if err = call.Query(context.Background(), &q, nil); err != nil {
	   		fmt.Fprintf(w, "error with assets: %v\n", err)
	   	} */

	mux := muxgo.NewAPIClient(
		muxgo.NewConfiguration(
			muxgo.WithBasicAuth(os.Getenv("MUX_ID"), os.Getenv("MUX_SECRET")),
		))

	assets, err := mux.AssetsApi.ListAssets()

	if err != nil {

		fmt.Fprintf(w, "something went wrong... %v\n", err)

	}

	for _, a := range assets.Data {

		var q struct {
			AssetByAssetID struct {
				AssetEntry
			} `graphql:"assetById(assetId: $AssetID)"`
		}

		v := map[string]interface{}{
			"AssetID": graphql.String(a.Id),
		}

		if err = call.Mutate(context.Background(), &q, v); err != nil {
			fmt.Fprintf(w, "error with asset: %v\n", err)
		}

		if q.AssetByAssetID.ID == nil {

			req := muxgo.CreatePlaybackIdRequest{muxgo.SIGNED}
			res, err := mux.AssetsApi.CreateAssetPlaybackId(a.Id, req)

			if err != nil {

				fmt.Fprintf(w, "something went wrong... %v\n", err)

			}

			var m struct {
				CreateAsset struct {
					AssetEntry
				} `graphql:"createAsset(assetId: $AssetID, playbackId: $PlaybackID)"`
			}

			v = map[string]interface{}{
				"AssetID":    graphql.String(a.Id),
				"PlaybackID": graphql.String(res.Data.Id),
			}

			if err = call.Query(context.Background(), &m, v); err != nil {
				fmt.Fprintf(w, "error with asset: %v\n", err)
			}

		}

	}

	/* 	k, _ := mux.URLSigningKeysApi.CreateUrlSigningKey()
	   	kg, _ := mux.URLSigningKeysApi.GetUrlSigningKey(k.Data.Id)

	   	fmt.Fprintf(w, "%v\n\n", kg.Data.PrivateKey) */

}
