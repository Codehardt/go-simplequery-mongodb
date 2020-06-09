package simplequery

import (
	"encoding/json"
	"fmt"
)

func parseAndPrint(cond string) {
	filter, err := Parse(cond)
	if err != nil {
		panic(err)
	}
	b, err := json.Marshal(filter)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func ExampleParse() {
	parseAndPrint(`foo < 5 AND (bar >= 10 OR baz = "example" OR qux = /[a-z]{3,5}/i)`)
	parseAndPrint(`foo = "bar" AND NOT (baz = "")`)
	// Output:
	// {"$and":[{"foo":{"$lt":5}},{"$or":[{"$or":[{"bar":{"$gte":10}},{"baz":{"$eq":"example"}}]},{"qux":{"$eq":{"Key":"$regex","Value":{"Pattern":"[a-z]{3,5}","Options":"i"}}}}]}]}
	// {"$and":[{"foo":{"$eq":"bar"}},{"$nor":[{"baz":{"$eq":""}}]}]}
}
