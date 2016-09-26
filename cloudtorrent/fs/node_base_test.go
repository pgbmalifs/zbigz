package fs

import "testing"

func TestBaseNode1(t *testing.T) {
	b := NewBaseNode("root")
	foo := NewBaseNode("foo")
	b.Upsert("foo", foo)
	foo.Upsert("bazz", NewBaseNode("bazz"))
	b.Upsert("bar", NewBaseNode("bar"))
	c := b.GetChildren()
	if l := len(c); l != 2 {
		t.Fatalf("should be 2, got %d", l)
	}
}

func TestBaseNode2(t *testing.T) {
	b := NewBaseNode("root")
	b.Upsert("/foo/bazz", NewBaseNode("bazz"))
	if l := len(b.GetChildren()); l != 1 {
		t.Fatalf("root should have 1, got %d", l)
	}
	f := b.Get("/foo")
	if l := len(f.GetChildren()); l != 1 {
		t.Fatalf("foo should have 1, got %d", l)
	}
}

func TestBaseNode3(t *testing.T) {
	b := NewBaseNode("root")
	b.Upsert("/zip", NewBaseNode("zip"))
	b.Upsert("/foo/bazz", NewBaseNode("bazz"))
	b.Upsert("/zop", NewBaseNode("zop"))
	if l := len(b.GetChildren()); l != 3 {
		t.Fatalf("root should have 3, got %d", l)
	}
	f := b.Get("/foo")
	if l := len(f.GetChildren()); l != 1 {
		t.Fatalf("foo should have 1, got %d", l)
	}
	//
	if !b.Delete("/zop") {
		t.Fatalf("root failed to delete zop")
	}
	if l := len(b.GetChildren()); l != 2 {
		t.Fatalf("root should have 2, got %d", l)
	}
	if b.Delete("/zop") {
		t.Fatalf("root double deleted zop")
	}
}
