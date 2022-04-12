package cache

import (
	"reflect"
	"testing"
)

type cacheEntity struct {
	Key   string
	Value string
}

func TestCache_Add(t *testing.T) {
	cases := []struct {
		limit         int64
		cacheEntities []cacheEntity
		validKeys     map[string]cacheEntity
		expectedSize  int
	}{
		{
			limit: 5,
			cacheEntities: []cacheEntity{
				{
					Key:   "foo",
					Value: "value0",
				},
				{
					Key:   "bar",
					Value: "value1",
				},
				{
					Key:   "baz",
					Value: "value2",
				},
				{
					Key:   "qux",
					Value: "value3",
				},
				{
					Key:   "moo",
					Value: "value4",
				},
			},
			validKeys: map[string]cacheEntity{
				"foo": {
					Key:   "foo",
					Value: "value0",
				},
				"bar": {
					Key:   "bar",
					Value: "value1",
				},
				"baz": {
					Key:   "baz",
					Value: "value2",
				},
				"qux": {
					Key:   "qux",
					Value: "value3",
				},
				"moo": {
					Key:   "moo",
					Value: "value4",
				},
			},
			expectedSize: 5,
		},
		{
			limit: 2,
			cacheEntities: []cacheEntity{
				{
					Key:   "bar",
					Value: "value1",
				},
				{
					Key:   "foo",
					Value: "value0",
				},
				{
					Key:   "baz",
					Value: "value2",
				},
				{
					Key:   "qux",
					Value: "value3",
				},
				{
					Key:   "moo",
					Value: "value4",
				},
			},
			validKeys: map[string]cacheEntity{
				"foo": {
					Key:   "foo",
					Value: "value0",
				},
				"bar": {
					Key:   "bar",
					Value: "value1",
				},
				"baz": {
					Key:   "baz",
					Value: "value2",
				},
				"qux": {
					Key:   "qux",
					Value: "value3",
				},
				"moo": {
					Key:   "moo",
					Value: "value4",
				},
			},
			expectedSize: 2,
		},
	}

	for _, c := range cases {
		cache := New(c.limit)

		for _, entity := range c.cacheEntities {
			cache.Add(entity.Key, entity)
		}

		count := 0
		cacheEntities := map[string]cacheEntity{}
		cache.caches.Range(func(key, value interface{}) bool {
			count++

			cacheEntities[key.(string)] = value.(cacheEntity)
			return true
		})

		if e, a := c.expectedSize, cache.size; int64(e) != a {
			t.Errorf("expected %v, but received %v", e, a)
		}

		if e, a := c.expectedSize, count; e != a {
			t.Errorf("expected %v, but received %v", e, a)
		}

		for k, ep := range cacheEntities {
			entity, ok := c.validKeys[k]
			if !ok {
				t.Errorf("unrecognized key %q in cache", k)
			}
			if e, a := entity, ep; !reflect.DeepEqual(e, a) {
				t.Errorf("expected %v, but received %v", e, a)
			}
		}
	}
}

func TestCache_Get(t *testing.T) {
	cases := []struct {
		limit         int64
		cacheEntities []cacheEntity
		validKeys     map[string]cacheEntity
	}{
		{
			limit: 5,
			cacheEntities: []cacheEntity{
				{
					Key:   "foo",
					Value: "value0",
				},
				{
					Key:   "bar",
					Value: "value1",
				},
				{
					Key:   "baz",
					Value: "value2",
				},
				{
					Key:   "qux",
					Value: "value3",
				},
				{
					Key:   "moo",
					Value: "value4",
				},
			},
			validKeys: map[string]cacheEntity{
				"foo": {
					Key:   "foo",
					Value: "value0",
				},
				"bar": {
					Key:   "bar",
					Value: "value1",
				},
				"baz": {
					Key:   "baz",
					Value: "value2",
				},
				"qux": {
					Key:   "qux",
					Value: "value3",
				},
				"moo": {
					Key:   "moo",
					Value: "value4",
				},
			},
		},
		{
			limit: 2,
			cacheEntities: []cacheEntity{
				{
					Key:   "bar",
					Value: "value1",
				},
				{
					Key:   "foo",
					Value: "value0",
				},
				{
					Key:   "baz",
					Value: "value2",
				},
				{
					Key:   "qux",
					Value: "value3",
				},
				{
					Key:   "moo",
					Value: "value4",
				},
			},
			validKeys: map[string]cacheEntity{
				"foo": {
					Key:   "foo",
					Value: "value0",
				},
				"bar": {
					Key:   "bar",
					Value: "value1",
				},
				"baz": {
					Key:   "baz",
					Value: "value2",
				},
				"qux": {
					Key:   "qux",
					Value: "value3",
				},
				"moo": {
					Key:   "moo",
					Value: "value4",
				},
			},
		},
	}

	for _, c := range cases {
		cache := New(c.limit)

		for _, entity := range c.cacheEntities {
			cache.Add(entity.Key, entity)
		}

		var keys []string
		cache.caches.Range(func(key, value interface{}) bool {
			a := value.(cacheEntity)
			e, ok := c.validKeys[key.(string)]
			if !ok {
				t.Errorf("unrecognized key %q in cache", key.(string))
			}

			if !reflect.DeepEqual(e, a) {
				t.Errorf("expected %v, but received %v", e, a)
			}

			keys = append(keys, key.(string))
			return true
		})

		for _, key := range keys {
			a, ok := cache.Get(key)
			if !ok {
				t.Errorf("expected key to be present: %q", key)
			}

			e := c.validKeys[key]
			if !reflect.DeepEqual(e, a) {
				t.Errorf("expected %v, but received %v", e, a)
			}
		}
	}
}

func TestCache_Delete(t *testing.T) {
	cases := []struct {
		limit         int64
		cacheEntities []cacheEntity
		deletedKeys   []string
		validKeys     map[string]cacheEntity
		expectedSize  int
	}{
		{
			limit: 5,
			cacheEntities: []cacheEntity{
				{
					Key:   "foo",
					Value: "value0",
				},
				{
					Key:   "bar",
					Value: "value1",
				},
				{
					Key:   "baz",
					Value: "value2",
				},
				{
					Key:   "qux",
					Value: "value3",
				},
				{
					Key:   "moo",
					Value: "value4",
				},
			},
			deletedKeys: []string{
				"foo", "bar",
			},
			validKeys: map[string]cacheEntity{
				"baz": {
					Key:   "baz",
					Value: "value2",
				},
				"qux": {
					Key:   "qux",
					Value: "value3",
				},
				"moo": {
					Key:   "moo",
					Value: "value4",
				},
			},
			expectedSize: 3,
		},
		{
			limit: 2,
			cacheEntities: []cacheEntity{
				{
					Key:   "bar",
					Value: "value1",
				},
				{
					Key:   "foo",
					Value: "value0",
				},
				{
					Key:   "baz",
					Value: "value2",
				},
				{
					Key:   "qux",
					Value: "value3",
				},
				{
					Key:   "moo",
					Value: "value4",
				},
			},
			deletedKeys: []string{
				"bar", "foo", "baz", "qux", "moo",
			},
			validKeys:    map[string]cacheEntity{},
			expectedSize: 0,
		},
	}

	for _, c := range cases {
		cache := New(c.limit)

		for _, entity := range c.cacheEntities {
			cache.Add(entity.Key, entity)
		}

		for _, key := range c.deletedKeys {
			cache.Delete(key)
		}

		count := 0
		var keys []string
		cache.caches.Range(func(key, value interface{}) bool {
			count++

			a := value.(cacheEntity)
			e, ok := c.validKeys[key.(string)]
			if !ok {
				t.Errorf("unrecognized key %q in cache", key.(string))
			}

			if !reflect.DeepEqual(e, a) {
				t.Errorf("expected %v, but received %v", e, a)
			}

			keys = append(keys, key.(string))
			return true
		})

		if e, a := c.expectedSize, cache.size; int64(e) != a {
			t.Errorf("expected %v, but received %v", e, a)
		}

		if e, a := c.expectedSize, count; e != a {
			t.Errorf("expected %v, but received %v", e, a)
		}

		for _, key := range keys {
			a, ok := cache.Get(key)
			if !ok {
				t.Errorf("expected key to be present: %q", key)
			}

			e := c.validKeys[key]
			if !reflect.DeepEqual(e, a) {
				t.Errorf("expected %v, but received %v", e, a)
			}
		}
	}
}

func TestCache_Keys(t *testing.T) {
	cases := []struct {
		limit         int64
		cacheEntities []cacheEntity
		validKeys     map[string]cacheEntity
	}{
		{
			limit: 5,
			cacheEntities: []cacheEntity{
				{
					Key:   "foo",
					Value: "value0",
				},
				{
					Key:   "bar",
					Value: "value1",
				},
				{
					Key:   "baz",
					Value: "value2",
				},
				{
					Key:   "qux",
					Value: "value3",
				},
				{
					Key:   "moo",
					Value: "value4",
				},
			},
			validKeys: map[string]cacheEntity{
				"foo": {
					Key:   "foo",
					Value: "value0",
				},
				"bar": {
					Key:   "bar",
					Value: "value1",
				},
				"baz": {
					Key:   "baz",
					Value: "value2",
				},
				"qux": {
					Key:   "qux",
					Value: "value3",
				},
				"moo": {
					Key:   "moo",
					Value: "value4",
				},
			},
		},
		{
			limit: 2,
			cacheEntities: []cacheEntity{
				{
					Key:   "bar",
					Value: "value1",
				},
				{
					Key:   "foo",
					Value: "value0",
				},
				{
					Key:   "baz",
					Value: "value2",
				},
				{
					Key:   "qux",
					Value: "value3",
				},
				{
					Key:   "moo",
					Value: "value4",
				},
			},
			validKeys: map[string]cacheEntity{
				"foo": {
					Key:   "foo",
					Value: "value0",
				},
				"bar": {
					Key:   "bar",
					Value: "value1",
				},
				"baz": {
					Key:   "baz",
					Value: "value2",
				},
				"qux": {
					Key:   "qux",
					Value: "value3",
				},
				"moo": {
					Key:   "moo",
					Value: "value4",
				},
			},
		},
	}

	for _, c := range cases {
		cache := New(c.limit)

		for _, entity := range c.cacheEntities {
			cache.Add(entity.Key, entity)
		}

		keys := cache.Keys()

		for _, key := range keys {
			a, ok := cache.Get(key)
			if !ok {
				t.Errorf("expected key to be present: %q", key.(string))
			}

			e, ok := c.validKeys[key.(string)]
			if !ok {
				t.Errorf("unrecognized key %q in cache", key.(string))
			}

			if !reflect.DeepEqual(e, a) {
				t.Errorf("expected %v, but received %v", e, a)
			}
		}
	}
}

func TestCache_UpdateLimit(t *testing.T) {
	cases := []struct {
		limit         int64
		cacheEntities []cacheEntity
		validKeys     map[string]cacheEntity
		expectedSize  int
		newLimit      int64
	}{
		{
			limit: 5,
			cacheEntities: []cacheEntity{
				{
					Key:   "foo",
					Value: "value0",
				},
				{
					Key:   "bar",
					Value: "value1",
				},
				{
					Key:   "baz",
					Value: "value2",
				},
				{
					Key:   "qux",
					Value: "value3",
				},
				{
					Key:   "moo",
					Value: "value4",
				},
			},
			validKeys: map[string]cacheEntity{
				"foo": {
					Key:   "foo",
					Value: "value0",
				},
				"bar": {
					Key:   "bar",
					Value: "value1",
				},
				"baz": {
					Key:   "baz",
					Value: "value2",
				},
				"qux": {
					Key:   "qux",
					Value: "value3",
				},
				"moo": {
					Key:   "moo",
					Value: "value4",
				},
			},
			newLimit:     5,
			expectedSize: 5,
		},
		{
			limit: 2,
			cacheEntities: []cacheEntity{
				{
					Key:   "bar",
					Value: "value1",
				},
				{
					Key:   "foo",
					Value: "value0",
				},
				{
					Key:   "baz",
					Value: "value2",
				},
				{
					Key:   "qux",
					Value: "value3",
				},
				{
					Key:   "moo",
					Value: "value4",
				},
			},
			validKeys: map[string]cacheEntity{
				"foo": {
					Key:   "foo",
					Value: "value0",
				},
				"bar": {
					Key:   "bar",
					Value: "value1",
				},
				"baz": {
					Key:   "baz",
					Value: "value2",
				},
				"qux": {
					Key:   "qux",
					Value: "value3",
				},
				"moo": {
					Key:   "moo",
					Value: "value4",
				},
			},
			newLimit:     2,
			expectedSize: 2,
		},
		{
			limit: 5,
			cacheEntities: []cacheEntity{
				{
					Key:   "foo",
					Value: "value0",
				},
				{
					Key:   "bar",
					Value: "value1",
				},
				{
					Key:   "baz",
					Value: "value2",
				},
				{
					Key:   "qux",
					Value: "value3",
				},
				{
					Key:   "moo",
					Value: "value4",
				},
			},
			validKeys: map[string]cacheEntity{
				"foo": {
					Key:   "foo",
					Value: "value0",
				},
				"bar": {
					Key:   "bar",
					Value: "value1",
				},
				"baz": {
					Key:   "baz",
					Value: "value2",
				},
				"qux": {
					Key:   "qux",
					Value: "value3",
				},
				"moo": {
					Key:   "moo",
					Value: "value4",
				},
			},
			newLimit:     2,
			expectedSize: 2,
		},
		{
			limit: 2,
			cacheEntities: []cacheEntity{
				{
					Key:   "bar",
					Value: "value1",
				},
				{
					Key:   "foo",
					Value: "value0",
				},
				{
					Key:   "baz",
					Value: "value2",
				},
				{
					Key:   "qux",
					Value: "value3",
				},
				{
					Key:   "moo",
					Value: "value4",
				},
			},
			validKeys: map[string]cacheEntity{
				"foo": {
					Key:   "foo",
					Value: "value0",
				},
				"bar": {
					Key:   "bar",
					Value: "value1",
				},
				"baz": {
					Key:   "baz",
					Value: "value2",
				},
				"qux": {
					Key:   "qux",
					Value: "value3",
				},
				"moo": {
					Key:   "moo",
					Value: "value4",
				},
			},
			newLimit:     5,
			expectedSize: 5,
		},
	}

	for i, c := range cases {
		cache := New(c.limit)

		for _, entity := range c.cacheEntities {
			cache.Add(entity.Key, entity)
		}

		cache.UpdateCacheLimit(c.newLimit)

		for _, entity := range c.cacheEntities {
			cache.Add(entity.Key, entity)
		}

		count := 0
		cacheEntities := map[string]cacheEntity{}
		cache.caches.Range(func(key, value interface{}) bool {
			count++

			cacheEntities[key.(string)] = value.(cacheEntity)
			return true
		})

		if e, a := c.expectedSize, cache.size; int64(e) != a {
			t.Errorf("case %d, expected %v, but received %v", i, e, a)
		}

		if e, a := c.expectedSize, count; e != a {
			t.Errorf("case %d, expected %v, but received %v", i, e, a)
		}

		for k, ep := range cacheEntities {
			entity, ok := c.validKeys[k]
			if !ok {
				t.Errorf("case %d, unrecognized key %q in cache", i, k)
			}
			if e, a := entity, ep; !reflect.DeepEqual(e, a) {
				t.Errorf("case %d, expected %v, but received %v", i, e, a)
			}
		}
	}
}
