package mybot

import (
	"testing"
)

func TestTwitterAction(t *testing.T) {
	a1 := TwitterAction{
		Retweet:     true,
		Favorite:    false,
		Follow:      true,
		Collections: []string{"col1", "col2"},
	}
	a2 := a1
	a3 := TwitterAction{
		Retweet:     true,
		Favorite:    true,
		Follow:      false,
		Collections: []string{"col1", "col3"},
	}

	a1.Add(&a3)
	if a1.Retweet != true {
		t.Fatalf("%v expected but %v found", false, a1.Retweet)
	}
	if a1.Favorite != true {
		t.Fatalf("%v expected but %v found", true, a1.Favorite)
	}
	if a1.Follow != true {
		t.Fatalf("%v expected but %v found", true, a1.Follow)
	}
	if len(a1.Collections) != 3 {
		t.Fatalf("%d expected but %d found", 3, len(a1.Collections))
	}

	a2.Sub(&a3)
	if a2.Retweet != false {
		t.Fatalf("%v expected but %v found", false, a2.Retweet)
	}
	if a2.Favorite != false {
		t.Fatalf("%v expected but %v found", false, a2.Favorite)
	}
	if a2.Follow != true {
		t.Fatalf("%v expected but %v found", true, a2.Follow)
	}
	if len(a2.Collections) != 1 {
		t.Fatalf("%d expected but %d found", 1, len(a2.Collections))
	}
}
