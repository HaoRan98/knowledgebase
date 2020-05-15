package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"net/http"
	"strings"
	"time"
)

// Initialize a client with the configuration object.
func NewEsClient() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://127.0.0.1:9200",
		},
		//Username: "foo",
		//Password: "bar",
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Println("Error creating the client: ", err)
		return nil, err
	}

	log.Println(elasticsearch.Version)
	log.Println(es.Info())
	return es, nil
}

// Index documents concurrently
func EsDocument(index, id, jsonStr string) (string, string) {
	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       strings.NewReader(jsonStr),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return "", fmt.Sprintf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return "", fmt.Sprintf("[%s] Error indexing document ID=%s", res.Status(), id)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return "", fmt.Sprintf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			return fmt.Sprintf("[%s] %s; version=%d",
				res.Status(),
				r["result"],
				int(r["_version"].(float64)),
			), ""
		}
	}
}

// Search for the indexed documents
func EsSearch(index, key string) ([]interface{}, string) {

	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithQuery(key),
		es.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Sprintf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Sprintf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			return nil, fmt.Sprintf("search response:[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var results []interface{}
	// Build the request body.
	var r map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Sprintf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	total := int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	took := int(r["took"].(float64))
	log.Printf("[%s] %d hits; took: %dms", res.Status(), total, took)

	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"]
		// debug
		//id := hit.(map[string]interface{})["_id"]
		//log.Printf(" * ID=%s, %s", id, source)
		results = append(results, source)
	}
	return results, ""
}

//根据id获取document
func EsGetSource(index, id string) (string, error) {
	res, err := es.GetSource(
		index,
		id,
		es.GetSource.WithPretty(),
	)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println(res.String())
	beg := strings.LastIndex(res.String(), "]")
	log.Println("_source:\n", res.String()[beg+1:])
	source := res.String()[beg+1:]
	return source, nil
}

//删除索引
func EsDeleteIndex(index string) error {
	req := esapi.IndicesDeleteRequest{
		Index: []string{index},
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	log.Println(res.String())
	if res.StatusCode == 200 || res.StatusCode == 404 {
		return nil
	}
	return nil
}
