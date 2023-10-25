package db

const (
	Asc  = 1
	Desc = -1
)

type UpdateResult struct {
	MatchedCount  int64       // The number of documents matched by the filter.
	ModifiedCount int64       // The number of documents modified by the operation.
	UpsertedCount int64       // The number of documents upserted by the operation.
	UpsertedID    interface{} // The _id field of the upserted document, or nil if no upsert was done.
}

// DeleteResult is the result type returned by DeleteOne and DeleteMany operations.
type DeleteResult struct {
	DeletedCount int64 `bson:"n"` // The number of documents deleted.
}

// pagination 라이브러리
type PaginationData struct {
	Total     int64 `json:"total"`
	Page      int64 `json:"page"`
	PerPage   int64 `json:"perPage"`
	Prev      int64 `json:"prev"`
	Next      int64 `json:"next"`
	TotalPage int64 `json:"totalPage"`
}
