package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	f "github.com/fauna/faunadb-go/faunadb"
	//"github.com/muxinc/mux-go"
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

	//m := make(map[graphql.ID]string, 0)
	var content string

	var assets []AssetEntry

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

		httpClient := oauth2.NewClient(context.Background(), src)

		caller := graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

		/* client := muxgo.NewAPIClient(
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

		} */

		var q struct {
			AllAssets struct {
				Data []AssetEntry
			}
		}

		if err = caller.Query(context.Background(), &q, nil); err != nil {
			fmt.Fprintf(w, "error with products: %v\n", err)
		}

		assets = q.AllAssets.Data

		sort.Slice(assets, func(i, j int) bool {

			if assets[i].Category < assets[j].Category {
				return true
			}

			if assets[i].Category > assets[i].Category {
				return false
			}

			return assets[i].Title < assets[i].Title

		})

	
	case "public":

	case "signed":

	default:
	}

	switch r.Method {

	case "GET":

		switch id {

		case "":

			content =

				`
			<!DOCTYPE html>
			<html lang="en">
			<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<title>vdo2go</title>
			<!-- CSS -->
			<!-- Add Material font (Roboto) and Material icon as needed -->
			<link href="https://fonts.googleapis.com/css?family=Roboto:300,300i,400,400i,500,500i,700,700i|Roboto+Mono:300,400,700|Roboto+Slab:300,400,700" rel="stylesheet">
			<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">

			<!-- Add Material CSS, replace Bootstrap CSS -->
			<link href="https://assets.medienwerk.now.sh/material.min.css" rel="stylesheet">
			</head>
			<body style="background-color: #a1b116;">

			<div class="container" id="videos" style="color:rgb(255, 255, 255); font-size:30px;">

			<ul class="list-group">
			<br>
			<br>

			<h1>Selection</h1>

			<form role="form" method="POST">

			<li class="list-group-item">

			`

			if assets[0].ID != nil {

				c := string(assets[0].Category)

				id := fmt.Sprintf("%s", assets[0].ID)
				//price := strconv.FormatFloat(float64(products[0].Price), 'f', 2, 64)
				//pack := strconv.Itoa(int(products[0].Pack))
				//dim := strconv.Itoa(int(products[0].LinkDIM))

				content = content + ` 
	
				<br>
				<h1>` + c + `</h1>
	
				<li class="list-group-item">
	
				<div class="media">
				<img class="mr-3" src=https://image.mux.com/` + string(assets[0].PbID) + `/thumbnail.jpg?width=214&height=121&fit_mode=pad" width="100">
	
				<div class="media-body">
	
				<h2>` + string(assets[0].Title) + `</h2>
	
				<h4>` + string(assets[0].Content) + `</h4>
	
				<p><a href="mailto:` + string(assets[0].Email) + `?subject=`+string(assets[0].AssetID)+`" target="_top">content provider</a><br><br>

				
				
				<label class="form-check-label" for="` + id + `" style="font-size:25px;">Mengenauswahl:</label>

				<input type="checkbox" id="Access" name="Access" value="select">
 				<label for="Access">select content</label><br>

				</p>

				`

				for k := 1; k < len(assets); k++ {


				}

			}

		
		case "public":

		case "signed":
	
		default:

		}

	case "POST":

	case "":

	case "public":

	case "signed":

	default:
	

	}

}
