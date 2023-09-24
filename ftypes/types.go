package ftypes

import "go.mongodb.org/mongo-driver/bson/primitive"

// DBTypes represents various types of databases, such as SQL and NoSQL.
type DBTypes int

// QMap is a shorthand for a map with string keys and interface{} values.
type QMap map[string]interface{}

// DMap is a shorthand for a BSON primitive.D, used for MongoDB queries.
type DMap primitive.D

// ConditionType represents different types of conditions for database queries.
type ConditionType int
