package mux

import (
	"testing"
)

func TestTreeMatch(t *testing.T) {
	root := NewNode("GET", 0)

	lang := NewNode("{lang:en|pl}", root.MaxParamsSize())
	blog := NewNode("blog", lang.MaxParamsSize())

	search := NewNode("search", blog.MaxParamsSize())
	searchAuthor := NewNode("author", search.MaxParamsSize())

	page := NewNode("page", blog.MaxParamsSize())
	pageID := NewNode(`{pageId:[^/]+}`, page.MaxParamsSize())

	posts := NewNode("posts", blog.MaxParamsSize())
	postsID := NewNode(`{postsId:[^/]+}`, posts.MaxParamsSize())

	comments := NewNode("comments", blog.MaxParamsSize())
	commentID := NewNode(`{commentId:\d+}`, comments.MaxParamsSize())
	commentNew := NewNode("new", commentID.MaxParamsSize())

	root.WithChildren(root.Tree().withNode(lang))
	lang.WithChildren(lang.Tree().withNode(blog))
	blog.WithChildren(blog.Tree().withNode(search))
	blog.WithChildren(blog.Tree().withNode(page))
	blog.WithChildren(blog.Tree().withNode(posts))
	blog.WithChildren(blog.Tree().withNode(comments))
	search.WithChildren(search.Tree().withNode(searchAuthor))
	page.WithChildren(page.Tree().withNode(pageID))
	posts.WithChildren(posts.Tree().withNode(postsID))
	comments.WithChildren(comments.Tree().withNode(commentID))
	commentID.WithChildren(commentID.Tree().withNode(commentNew))

	root.WithChildren(root.Tree().Compile())

	n, _, _ := root.Tree().Match("pl/blog/comments/123/new")

	if n == nil {
		t.Fatalf("%v", n)
	}

	if n.Name() != commentNew.Name() {
		t.Fatalf("%s != %s", n.Name(), commentNew.Name())
	}
}
