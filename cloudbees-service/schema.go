package main

import "github.com/hashicorp/go-memdb"

var schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		"Train": &memdb.TableSchema{
			Name: "Train",
			Indexes: map[string]*memdb.IndexSchema{
				"id": &memdb.IndexSchema{
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "Id"},
				},
				"year": &memdb.IndexSchema{
					Name:    "year",
					Unique:  false,
					Indexer: &memdb.IntFieldIndex{Field: "Year"},
				},
				"month": &memdb.IndexSchema{
					Name:    "month",
					Unique:  false,
					Indexer: &memdb.IntFieldIndex{Field: "Month"},
				},
				"day": &memdb.IndexSchema{
					Name:    "day",
					Unique:  false,
					Indexer: &memdb.IntFieldIndex{Field: "Day"},
				},
			},
		},
		"Seat": &memdb.TableSchema{
			Name: "Seat",
			Indexes: map[string]*memdb.IndexSchema{
				"id": &memdb.IndexSchema{
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "Id"},
				},
				"trainId": &memdb.IndexSchema{
					Name:    "trainId",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "TrainId"},
				},
				"row": &memdb.IndexSchema{
					Name:    "row",
					Unique:  false,
					Indexer: &memdb.IntFieldIndex{Field: "Row"},
				},
				"car": &memdb.IndexSchema{
					Name:    "car",
					Unique:  false,
					Indexer: &memdb.StringFieldIndex{Field: "Car"},
				},
				"status": &memdb.IndexSchema{
					Name:    "status",
					Unique:  false,
					Indexer: &memdb.StringFieldIndex{Field: "Status"},
				},
				"passengerId": &memdb.IndexSchema{
					Name:    "passengerId",
					Unique:  false,
					Indexer: &memdb.StringFieldIndex{Field: "PassengerId"},
				},
			},
		},
	},
}
