package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type Publication struct {
    ID         primitive.ObjectID `bson:"_id,omitempty"`
    UserID     primitive.ObjectID `bson:"user_id"`
    Title      string             `bson:"title"`
    Content    string             `bson:"content"`
    Category   string             `bson:"category,omitempty"`
    Image      string             `bson:"image,omitempty"`
    CreatedAt  time.Time          `bson:"created_at"`
    Likes      int                `bson:"likes"`        // Nombre de likes
    Dislikes   int                `bson:"dislikes"`     // Nombre de dislikes
}
