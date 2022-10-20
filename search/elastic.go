package search

import (
	"bytes"
	"context"
	"cqrs/models"
	"encoding/json"
	"errors"

	elastic "github.com/elastic/go-elasticsearch/v7"
)

/*
	concret implentation for elastic-search
*/

type ElasticSearchRepo struct {
	client *elastic.Client
}

// new connection for elastic - search
func NewElastic(url string) (*ElasticSearchRepo, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{url},
	})

	if err != nil {
		return nil, err
	}

	return &ElasticSearchRepo{client: client}, nil
}

func (r *ElasticSearchRepo) Close() {
	//
}

// find one for feed
func (r *ElasticSearchRepo) IndexFeed(ctx context.Context, feed *models.Feed) error {
	body, _ := json.Marshal(feed) //return bytes of feed

	_, err := r.client.Index(
		"feeds",
		bytes.NewReader(body), //process data for elastic search
		r.client.Index.WithDocumentID(feed.ID),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithRefresh("wait_for"),
	)
	return err
}

func (r *ElasticSearchRepo) SearchFeed(ctx context.Context, query string) (results []models.Feed, err error) {
	var buff bytes.Buffer

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":            query,
				"fields":           []string{"title", "description"},
				"fuzziness":        3,
				"cutoff_frequency": 0.0001,
			},
		},
	}

	if err = json.NewEncoder(&buff).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("feeds"),
		r.client.Search.WithBody(&buff),
		r.client.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			results = nil
		}
	}()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var eREs map[string]interface{} // presentation json on golang

	if err := json.NewDecoder(res.Body).Decode(&eREs); err != nil {
		return nil, err
	}

	var feeds []models.Feed

	for _, hit := range eREs["hits"].(map[string]interface{})["hits"].([]interface{}) {
		feed := models.Feed{}
		source := hit.(map[string]interface{})["_source"] // hit -> source

		marshal, err := json.Marshal(source) //parse json - bytes
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(marshal, &feed); err == nil { //ERROR IS NIL BECAU  SE ALL IS CORRECT, YEAH?
			feeds = append(feeds, feed)
		}
	}
	return feeds, nil
}
