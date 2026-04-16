package registry

import "errors"

type Registry[K comparable, V any] map[K]V

func (r Registry[K, V]) Get(k K) (V, error) {
	v, ok := r[k]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}
