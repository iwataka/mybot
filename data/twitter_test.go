package data_test

import (
	. "github.com/iwataka/mybot/data"

	"testing"
)

func TestTwitterAction(t *testing.T) {
	a1 := TwitterAction{
		Collections: []string{"col1", "col2"},
	}
	a1.Retweet = true
	a2 := TwitterAction{
		Collections: []string{"col1", "col3"},
	}
	a2.Retweet = true
	a2.Favorite = true

	result1 := a1.Add(a2)
	if result1.Retweet != true {
		t.Fatalf("%v expected but %v found", false, result1.Retweet)
	}
	if result1.Favorite != true {
		t.Fatalf("%v expected but %v found", true, result1.Favorite)
	}
	if len(result1.Collections) != 3 {
		t.Fatalf("%d expected but %d found", 3, len(result1.Collections))
	}

	result2 := a1.Sub(a2)
	if result2.Retweet != false {
		t.Fatalf("%v expected but %v found", false, result2.Retweet)
	}
	if result2.Favorite != false {
		t.Fatalf("%v expected but %v found", false, result2.Favorite)
	}
	if len(result2.Collections) != 1 {
		t.Fatalf("%d expected but %d found", 1, len(result2.Collections))
	}
}
