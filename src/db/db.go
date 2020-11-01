package db

import (
	"log"

	"github.com/hashicorp/go-memdb"
)

//New create new empty memdb for keeping track of notifications
func New() (*memdb.MemDB, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"alerts": &memdb.TableSchema{
				Name: "alerts",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "RuleName"},
					},
					"alerting": &memdb.IndexSchema{
						Name:    "alerting",
						Unique:  true,
						Indexer: &memdb.BoolFieldIndex{Field: "Alerting"},
					},
				},
			},
		},
	}
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		log.Panic(err)
	}
	return db, err
}
