package joins

// JoinType represents the type of join in a database query.
type JoinType int

const (
	// INNER_JOIN is used to specify an inner join in a database query.
	INNER_JOIN = JoinType(1)

	// LEFT_JOIN is used to specify a left join in a database query.
	LEFT_JOIN = JoinType(2)

	// RIGHT_JOIN is used to specify a right join in a database query.
	RIGHT_JOIN = JoinType(3)

	// CROSS_JOIN is used to specify a cross join in a database query.
	CROSS_JOIN = JoinType(4)
)

// Join represents a join operation in a database query, including the join type,
// the name of the table to join, and the join condition.
type Join struct {
	// Operator specifies the type of join, such as INNER_JOIN, LEFT_JOIN, RIGHT_JOIN, or CROSS_JOIN.
	Operator JoinType

	// TableName is the name of the table to join in the database query.
	TableName string

	// Condition is the join condition that specifies how the tables are related in the query.
	Condition string
}
