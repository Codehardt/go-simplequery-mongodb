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
	// [{"Key":"$and","Value":[{"Key":"foo","Value":{"Key":"$lt","Value":5}},{"Key":"$or","Value":[{"Key":"$or","Value":[{"Key":"bar","Value":{"Key":"$gte","Value":10}},{"Key":"baz","Value":{"Key":"$eq","Value":"example"}}]},{"Key":"qux","Value":{"Key":"$eq","Value":{"Key":"$regex","Value":{"Pattern":"[a-z]{3,5}","Options":"i"}}}}]}]}]
	// [{"Key":"$and","Value":[{"Key":"foo","Value":{"Key":"$eq","Value":"bar"}},{"Key":"$nor","Value":[{"Key":"baz","Value":{"Key":"$eq","Value":""}}]}]}]
}
