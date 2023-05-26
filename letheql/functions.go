package letheql

var AtModifierUnsafeFunctions = map[string]struct{}{
	// Step invariant functions.
	"days_in_month": {}, "day_of_month": {}, "day_of_week": {}, "day_of_year": {},
	"hour": {}, "minute": {}, "month": {}, "year": {},
	"predict_linear": {}, "time": {},
	// Uses timestamp of the argument for the result,
	// hence unsafe to use with @ modifier.
	"timestamp": {},
}
