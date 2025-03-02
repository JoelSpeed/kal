/*
nomaps provides a linter to ensure that structs do not contain a map field.

Map fields are discouraged in Kubernetes API types. Use a list of key-value pairs.

The following struct isn't recommended

type Spec struct {
	Pairs map[string]string
}

Instead, use a list of key-value pairs

type Spec struct {
	Pairs []Pair
}

type Pair struct {
	Key   string
	Value string
}
*/

package nomaps
