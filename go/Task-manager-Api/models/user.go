package models

type User struct {
	ID       string `bson:"_id,omitempty" json:"id"`
	UserName string `bson:"username" json:"username" binding:"required"`
	Password string `bson:"password" json:"password,omitempty" binding:"required"`
	Role     string `bson:"role" json:"role,omitempty"`
}
