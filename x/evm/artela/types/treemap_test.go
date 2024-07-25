package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/maps/treemap"
	pq "github.com/emirpasic/gods/queues/priorityqueue"
)

func TestPriorityQueue(t *testing.T) {
	/*	a := BondAspect{AspectId: "1111", priority: 1}
		b := BondAspect{AspectId: "22222", priority: -10}
		c := BondAspect{AspectId: "3333", priority: 3}
	*/
	ma := make(map[string]interface{}, 2)
	ma[AspectIDMapKey] = "1111"
	ma[PriorityMapKey] = 1
	mb := make(map[string]interface{}, 2)
	mb[AspectIDMapKey] = "2222"
	mb[PriorityMapKey] = -10
	mc := make(map[string]interface{}, 2)
	mc[AspectIDMapKey] = "3333"
	mc[PriorityMapKey] = 3

	queue := pq.NewWith(ByMapKeyPriority) // empty
	queue.Enqueue(ma)                     // {a 1}
	queue.Enqueue(mb)                     // {c 3}, {a 1}
	queue.Enqueue(mc)                     // {c 3}, {b 2}, {a 1}

	values := queue.Values()
	fmt.Println(values)
	toJSON, _ := queue.MarshalJSON()

	queue2 := pq.NewWith(ByMapKeyPriority)
	err := queue2.UnmarshalJSON(toJSON)
	if err != nil {
		fmt.Println(err)
	}
	v2 := queue2.Values()
	fmt.Println(v2)
	newQueue := pq.NewWith(ByMapKeyPriority)

	value, ok := queue.Dequeue()
	for ok {
		s, toStrOk := value.(map[string]interface{})["AspectId"].(string)
		if toStrOk && strings.EqualFold("3333", s) {
			break
		}
		newQueue.Enqueue(value)
		value, ok = queue.Dequeue()
	}
	fmt.Println(newQueue.Values())
}

func TestTreemap(_ *testing.T) {
	m := hashmap.New()
	m.Put("a", "1")
	m.Put("b", "2")
	m.Put("c", "3")

	bytes, err := json.Marshal(m) // Same as "m.ToJSON(m)"
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bytes))

	hm := hashmap.New()

	err2 := json.Unmarshal(bytes, &hm) // Same as "hm.FromJSON(bytes)"
	if err2 != nil {
		fmt.Println(err2)
	}

	fmt.Println(hm) // HashMap map[b:2 a:1]

	treem := treemap.NewWithIntComparator() // empty (keys are of type int)
	treem.Put(1000, "v1.0")                 // 1->x
	treem.Put(2000, "v2.0")                 // 1->x, 2->b (in order)
	treem.Put(90000000, "a")
	foundKey, foundValue := treem.Floor(4)
	foundKey2, foundValue2 := treem.Floor(5)
	foundKey3, foundValue3 := treem.Floor(6)
	foundKey4, foundValue4 := treem.Floor(90000000)
	foundKey5, foundValue5 := treem.Floor(90000001)

	bytes, err = json.Marshal(treem) // Same as "m.ToJSON(m)"
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bytes))

	comparator := treemap.NewWithIntComparator()
	err = json.Unmarshal(bytes, &comparator)
	if err != nil {
		return
	}
	foundKey7, foundValue7 := comparator.Floor("90")
	fmt.Print("%w,  %w \n", foundKey, foundValue)
	fmt.Print("%w,  %w \n", foundKey2, foundValue2)
	fmt.Print("%w,  %w \n", foundKey3, foundValue3)
	fmt.Print("%w,  %w \n", foundKey4, foundValue4)
	fmt.Print("%w,  %w \n", foundKey5, foundValue5)
	fmt.Print("%w,  %w \n", foundKey7, foundValue7)
	/*
		m := treemap.NewWithIntComparator() // empty (keys are of type int)
		m.Put(1, "x")                       // 1->x
		m.Put(4, "b")                       // 1->x, 2->b (in order)
		m.Put(90000000, "a")

		key, value := m.Floor(3)
		key1, value1 := m.Floor(50)
		fmt.Print("%w,  %w \n", key, value)
		fmt.Print("%w,  %w \n", key1, value1)
		bytes, err := json.Marshal(m)
		if err == nil {
			tree := treemap.NewWithIntComparator()
			_ = json.Unmarshal(bytes, tree)
			foundKey, foundValue := tree.Floor(90)
			fmt.Print("%w,  %w \n", foundKey, foundValue)
		}
	*/
}
