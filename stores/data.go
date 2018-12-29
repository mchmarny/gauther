package stores

import (
	"errors"
	"fmt"
	"log"
	"sort"
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



// InitDataStore initializes client
func InitDataStore() {

	projectID := utils.MustGetEnv("GCP_PROJECT_ID", "")
	coll = utils.MustGetEnv("FIRESTORE_COLL_NAME", defaultCollectionName)

	log.Printf("Initiating firestore client for %s collection in %s project",
		coll, projectID)

	dbClient, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("Error while creating Firestore client: %v", err)
	}
	db = dbClient
}

// CloseStore closes the DB connection
func CloseStore(){
	if db != nil {
		db.Close()
	}
}


// GetAllEmails retreaves Email for all data in the collection
func GetAllEmails(ctx context.Context) (data []string, err error) {

	list := []string{}

	iter := db.Collection(coll).OrderBy("email", firestore.Asc).Documents(ctx)
	for {
		d, e := iter.Next()
		if e == iterator.Done {
			sort.Strings(list)
			return list, nil
		}
		if e != nil {
			return nil, e
		}

		m := d.Data()

		log.Printf("[%v] %v", m["email"], m["id"])

		list = append(list, m["email"].(string))
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
