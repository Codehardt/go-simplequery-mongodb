package simplequery

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Codehardt/go-simplequery-parser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Parse(condition string) (bson.M, error) {
	root, err := simplequery.Parse(condition)
	if err != nil {
		return nil, err
	}
	e, _, _, err := parse(root)
	return e, err
}

var rgx = regexp.MustCompile(`^/(.*)/([gimy]*)$`)

func parse(node simplequery.Node) (bson.M, string, interface{}, error) {
	if node == nil {
		return nil, "", nil, nil
	}
	c1, c2 := node.Children()
	d1, k1, _, err := parse(c1)
	if err != nil {
		return nil, "", nil, err
	}
	d2, _, v2, err := parse(c2)
	if err != nil {
		return nil, "", nil, err
	}
	switch node.(type) {
	case simplequery.AND:
		return bson.M{"$and": []interface{}{d1, d2}}, "", nil, nil
	case simplequery.OR:
		return bson.M{"$or": []interface{}{d1, d2}}, "", nil, nil
	case simplequery.NOT:
		return bson.M{"$nor": []interface{}{d1}}, "", nil, nil
	case simplequery.EQ, simplequery.NE, simplequery.GT, simplequery.GTE, simplequery.LT, simplequery.LTE:
		op := "$" + strings.ToLower(reflect.TypeOf(node).Name()) // $eq, $ne, $gt, ...
		return bson.M{k1: bson.M{op: v2}}, "", nil, nil
	case simplequery.ID:
		return nil, node.Value(), nil, nil
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
				return nil, "", nil, err
			}
		}
		return nil, "", v, nil
	default:
		return nil, "", nil, fmt.Errorf("unknown node type %T", node)
	}
}
