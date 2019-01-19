package stores

import (
	"errors"
	"fmt"
	"log"
	"context"
	"cloud.google.com/go/firestore"

	"google.golang.org/api/iterator"

	"github.com/mchmarny/gauther/utils"


)

const (
	defaultCollectionName = "gauther"
)

var (
	coll   *firestore.CollectionRef
)



// InitDataStore initializes client
func InitDataStore() {

	projectID := utils.MustGetEnv("GCP_PROJECT_ID", "")
	collName := utils.MustGetEnv("FIRESTORE_COLL_NAME", defaultCollectionName)

	log.Printf("Initiating firestore client for %s collection in %s project",
		collName, projectID)

	// Assumes GOOGLE_APPLICATION_CREDENTIALS is set
	dbClient, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("Error while creating Firestore client: %v", err)
	}
	coll = dbClient.Collection(collName)
}


// GetAll retreaves all data for all data in the collection
func GetAll(ctx context.Context) (data []map[string]interface{}, err error) {

	list := make([]map[string]interface{}, 0)

	iter := coll.OrderBy("email", firestore.Asc).Documents(ctx)
	for {
		d, e := iter.Next()
		if e == iterator.Done {
			return list, nil
		}
		if e != nil {
			return nil, e
		}

		list = append(list, d.Data())
	}

}


// SaveData creates data or updates if exists with id
func SaveData(ctx context.Context, id string, data map[string]interface{}) error {

	if id == "" {
		return errors.New("Nil id")
	}

	_, err := coll.Doc(id).Set(ctx, data, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("Error on job save: %v", err)
	}

	return nil

}

// GetData retrieves saved sentiment request
func GetData(ctx context.Context, id string) (data map[string]interface{}, err error) {

	if id == "" {
		return nil, errors.New("Nil job ID parameter")
	}

	d, err := coll.Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	if d == nil {
		return nil, fmt.Errorf("No doc for ID: %s", id)
	}

	return d.Data(), nil


}


// DeleteData deletes data by id
func DeleteData(ctx context.Context, id string) error {

	if id == "" {
		return errors.New("Nil job ID parameter")
	}

	_, err := coll.Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("Error deleting data: %v", err)
	}

	return nil
}
