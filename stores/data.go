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
	db   *firestore.Client
	coll string
)



// InitStore returns configured store
func InitStore(ctx context.Context) error {

	projectID := utils.MustGetEnv("GCP_PROJECT_ID", "")
	coll = utils.MustGetEnv("FIRESTORE_COLL_NAME", defaultCollectionName)

	log.Printf("Initiating firestore client for %s collection in %s project",
		coll, projectID)

	dbClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("Error while creating Firestore client: %v", err)
	}
	db = dbClient

	return nil

}

// CloseStore closes the DB connection
func CloseStore(){
	if db != nil {
		db.Close()
	}
}


// GetAllData retreaves all data in the collection
func GetAllData(ctx context.Context, data chan<- map[string]interface{},
									 err chan<- error,
									 done chan<- bool) {

	iter := db.Collection(coll).Documents(ctx)
	for {
		d, e := iter.Next()
		if e == iterator.Done {
			done <- true
			break
		}
		if e != nil {
			err <- e
		}

		data <- d.Data()
	}

}


// SaveData creates data or updates if exists with id
func SaveData(ctx context.Context, id string, data map[string]interface{}) error {

	if id == "" {
		return errors.New("Nil id")
	}

	_, err := db.Collection(coll).Doc(id).Set(ctx, data, firestore.MergeAll)
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

	d, err := db.Collection(coll).Doc(id).Get(ctx)
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

	_, err := db.Collection(coll).Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("Error deleting data: %v", err)
	}

	return nil
}
