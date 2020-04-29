package db

import (
	"github.com/srgyrn/mdns-todo-api/storage/db/suitetest"
	"os"
	"testing"
)

func TestGateway(t *testing.T) {
	gw := BoltDB{
		dbName:         "test.db",
		rootBucketName: "DB",
		bucketName:     "ITEMS",
		db:             nil,
	}

	gwt := suitetest.GatewayTest{
		GW:     &gw,
		Before: func() {
			err := initBoltDB(&gw)

			if err != nil {
				t.Fatalf("init db error: %v", err)
			}
		},
		After:  func() {
			err := os.Remove("test.db")

			if err != nil {
				t.Errorf("failed to delete db: %v", err)
			}
		},
	}
	suitetest.Gateway(t, gwt)
}
