package conditions

type ConditionType int

const (
	EQUAL              = ConditionType(1)
	NOT_EQUAL          = ConditionType(2)
	GREATER_THAN       = ConditionType(3)
	GREATER_THAN_EQUAL = ConditionType(4)
	LESSER_THAN        = ConditionType(5)
	LESSER_THAN_EQUAL  = ConditionType(6)
	LIKE               = ConditionType(7)
	IN                 = ConditionType(8)
	NOT_IN             = ConditionType(9)
	IS_NULL            = ConditionType(10)
	IS_NOT_NULL        = ConditionType(11)
	NOT_LIKE           = ConditionType(12)
)

type ConditionBuilder interface {
	And(...ConditionBuilder) ConditionBuilder
	Or(...ConditionBuilder) ConditionBuilder
	Add(...ConditionBuilder) ConditionBuilder
	GroupConditions(ConditionBuilder) ConditionBuilder
	getKey() string
	getOperator() ConditionType
	getValue() interface{}
	isInner() bool
	getInnerQuery() DBFusionData
	GetQuery(conditions ConditionBuilder) DBFusionData
}
