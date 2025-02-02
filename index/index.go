package index

import (
	"time"

	"github.com/ostafen/clover/v2/store"
)

type Type int

const (
	SingleField Type = iota
)

type Info struct {
	Field string
	Type  Type
}

type Index interface {
	Add(docId string, v interface{}, ttl time.Duration) error
	Remove(docId string, v interface{}) error
	Iterate(reverse bool, onValue func(docId string) error) error
	Drop() error
	Type() Type
	Collection() string
	Field() string
}

type indexBase struct {
	collection, field string
}

func (idx *indexBase) Collection() string {
	return idx.collection
}

func (idx *indexBase) Field() string {
	return idx.field
}

type Query interface {
	Run(onValue func(docId string) error) error
}

func CreateIndex(collection, field string, idxType Type, tx store.Tx) Index {
	indexBase := indexBase{collection: collection, field: field}
	switch idxType {
	case SingleField:
		return &rangeIndex{
			indexBase: indexBase,
			tx:        tx,
		}
	}
	return nil
}
