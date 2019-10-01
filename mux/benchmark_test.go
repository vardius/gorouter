package mux

import (
	"testing"
)

func BenchmarkStatic1(b *testing.B) {
	root := NewNode("{lang:en|pl}", nil)
	blog := NewNode("blog", root)
	/* search := */ NewNode("search", blog)

	page := NewNode("page", blog)
	/* pageID := */ NewNode(`{pageId:[^/]++}`, page)

	posts := NewNode("posts", blog)
	/* postsID := */ NewNode(`{postsId:[^/]++}`, posts)

	comments := NewNode("comments", blog)
	commentID := NewNode(`{commentId:\d+}`, comments)
	commentNew := NewNode("new", commentID)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n, _, _ := root.FindByPath("/pl/blog/comments/123/new")

			if n == nil {
				b.Fatalf("%v", n)
			}

			if n.id != commentNew.id {
				b.Fatalf("%s != %s", n.id, commentNew.id)
			}
		}
	})
}
