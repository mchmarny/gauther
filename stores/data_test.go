package stores

import (
	"context"
	"log"
	"testing"
)

func TestJobData(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestJobData")
	}

	InitDataStore()

	ctx := context.Background()

	list, err := GetAllI(ctx)
	if err != nil {
		t.Errorf("Error on get all emails: %v", err)
	}

	for i, v := range list {
		log.Printf("[%d] %s - %s", i, v["id"], v["email"])
	}

	// configInitializer("test-data")

	// termReq := newRequest("test")

	// err := saveJob(termReq)

	// req, err := getJob(termReq.ID)

	// if err != nil {
	// 	t.Errorf("Error on job read: %v", err)
	// }

	// if req.ID != termReq.ID {
	// 	t.Errorf("Got invalid job: %v", req)
	// }

}
