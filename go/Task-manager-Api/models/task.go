package models

type Tasks struct {
	ID    string `bson:"_id,omitempty" json:"id"` // MongoDB will set this
	Title string `bson:"title" json:"title"`
}
