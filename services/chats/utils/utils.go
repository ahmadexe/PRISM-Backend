package utils

import (
	"sort"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SortIDs(id1 primitive.ObjectID, id2 primitive.ObjectID) string {
	ids := []primitive.ObjectID{id1, id2}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i].Timestamp().Before(ids[j].Timestamp())
	})

	idsJoined := ids[0].Hex() + ids[1].Hex()
	return idsJoined
}
