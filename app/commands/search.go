package commands

import (
	"bytes"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

type SearchCmd struct {
	Query string `arg name:"query" help:"Query be used in the search." type:"string"`
}

func (r *SearchCmd) Run(*Context) error {
	jsonData, err := r.query()
	if err != nil {
		return err
	}

	jsonParser := gjson.Parse(jsonData.String())

	products := jsonParser.Get(`#(title="Products").values`)
	for _, product := range products.Array() {
		id := product.Get(`project_id`).String()
		title := product.Get(`title`).String()
		category := product.Get(`cat_title`).String()
		username := product.Get(`username`).String()
		fmt.Printf("appimagehub:%s %s (%s) by %s\n", id, title, category, username)
	}

	return nil
}

func (r *SearchCmd) query() (bytes.Buffer, error) {
	query := "https://www.pling.com/json/search/p/" + r.Query + "/s/AppImageHub.com"
	resp, err := http.Get(query)
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer resp.Body.Close()

	var json bytes.Buffer
	_, err = io.Copy(&json, resp.Body)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return json, nil
}
