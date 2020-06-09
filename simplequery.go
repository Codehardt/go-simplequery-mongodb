package simplequery

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Codehardt/go-simplequery-mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Codehardt/go-simplequery-parser"
	"go.mongodb.org/mongo-driver/bson"
)

func Parse(condition string) (bson.D, error) {
	root, err := simplequery.Parse(condition)
	if err != nil {
		return nil, err
	}
	e, err := parse(root)
	return bson.D{e}, err
}

var rgx = regexp.MustCompile(`^/(.*)/([gimy]*)$`)

func parse(node simplequery.Node) (bson.E, error) {
	if node == nil {
		return bson.E{}, nil
	}
	c1, c2 := node.Children()
	d1, err := parse(c1)
	if err != nil {
		return bson.E{}, err
	}
	d2, err := parse(c2)
	if err != nil {
		return bson.E{}, err
	}
	switch node.(type) {
	case simplequery.AND:
		return bson.E{
			Key:   "$and",
			Value: bson.D{d1, d2},
		}, nil
	case simplequery.OR:
		return bson.E{
			Key:   "$or",
			Value: bson.D{d1, d2},
		}, nil
	case simplequery.NOT:
		return bson.E{
			Key:   "$nor",
			Value: bson.D{d1},
		}, nil
	case simplequery.EQ, simplequery.NE, simplequery.GT, simplequery.GTE, simplequery.LT, simplequery.LTE:
		op := "$" + strings.ToLower(reflect.TypeOf(node).Name()) // $eq, $ne, $gt, ...
		return bson.E{
			Key: d1.Key,
			Value: bson.M{
				op: d2.Value,
			},
		}, nil
	case simplequery.ID:
		return bson.E{Key: node.Value()}, nil
	case simplequery.VAL:
		var v interface{}
		str := node.Value()
		if strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"") {
			// trim " prefix and suffix for strings
			v = str[1 : len(str)-1]
		} else if matches := rgx.FindStringSubmatch(str); len(matches) == 3 {
			// parse regex /<pattern>/<options>
			v = bson.E{Key: "$regex", Value: primitive.Regex{Pattern: matches[1], Options: matches[2]}}
		} else {
			// parse integer
			v, err = strconv.Atoi(str)
			if err != nil {
				return bson.E{}, err
			}
		}
		return bson.E{Value: v}, nil
	default:
		return bson.E{}, fmt.Errorf("unknown node type %T", node)
	}
}
