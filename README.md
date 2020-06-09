# go-simplequery-mongodb

Convert a condition query to MongoDB query syntax

## Example

```golang
filter, _ := Parse(`foo < 5 AND (bar >= 10 OR baz = "example" OR qux = /[a-z]{3,5}/i)`)
```

generates a nested `bson.D` struct with the following nested structure:

```json
[{"Key":"$and","Value":[{"Key":"foo","Value":{"Key":"$lt","Value":5}},{"Key":"$or","Value":[{"Key":"$or","Value":[{"Key":"bar","Value":{"Key":"$gte","Value":10}},{"Key":"baz","Value":{"Key":"$eq","Value":"example"}}]},{"Key":"qux","Value":{"Key":"$eq","Value":{"Key":"$regex","Value":{"Pattern":"[a-z]{3,5}","Options":"i"}}}}]}]}]
```

that can be used in MongoDB's `.Find()`, `.Aggregate()`, `Count()`, `...` functions.
