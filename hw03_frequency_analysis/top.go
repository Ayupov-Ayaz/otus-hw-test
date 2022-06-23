package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type field struct {
	str   string
	count int
}

type bucket struct {
	fields         []field
	bucketCapacity int
}

func newBucket(bucketLen int) *bucket {
	return &bucket{
		fields:         make([]field, bucketLen),
		bucketCapacity: bucketLen,
	}
}

func (b *bucket) sortMinToMax() {
	sort.Slice(b.fields, func(i, j int) bool {
		a := b.fields[i]
		c := b.fields[j]

		if a.count == c.count {
			return a.str < c.str
		}

		return a.count < c.count
	})
}

func (b *bucket) sortMaxToMin() {
	sort.Slice(b.fields, func(i, j int) bool {
		a := b.fields[i]
		c := b.fields[j]

		if a.count == c.count {
			return a.str < c.str
		}

		return a.count > c.count
	})
}

func (b *bucket) add(str string, count int) {
	for i, f := range b.fields {
		if f.count < count {
			b.fields[i].count = count
			b.fields[i].str = str
			b.sortMinToMax()
			break
		}
	}
}

func (b *bucket) getTop() []string {
	b.sortMaxToMin()

	top := make([]string, len(b.fields))

	for i, f := range b.fields {
		top[i] = f.str
	}

	return top
}

func getStringMap(str string) map[string]int {
	fields := strings.Fields(str)

	results := make(map[string]int)

	for _, field := range fields {
		results[field]++
	}

	return results
}

func getTop10(strMap map[string]int) []string {
	bucket := newBucket(10)

	for str, count := range strMap {
		bucket.add(str, count)
	}

	return bucket.getTop()
}

func Top10(str string) []string {
	if len(str) == 0 {
		return nil
	}

	stringMap := getStringMap(str)

	return getTop10(stringMap)
}
