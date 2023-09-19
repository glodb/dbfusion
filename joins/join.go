package joins

type JoinType int

const (
	INNER_JOIN = JoinType(1)
	LEFT_JOIN  = JoinType(2)
	RIGHT_JOIN = JoinType(3)
	CROSS_JOIN = JoinType(4)
)

type Join struct {
	Operator  JoinType
	TableName string
	Condition string
}
