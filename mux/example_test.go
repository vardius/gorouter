package mux

import "fmt"

func Example() {
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

	root.WithChildren(root.Tree().withNode(lang).sort())
	lang.WithChildren(lang.Tree().withNode(blog).sort())
	blog.WithChildren(blog.Tree().withNode(search).sort())
	blog.WithChildren(blog.Tree().withNode(page).sort())
	blog.WithChildren(blog.Tree().withNode(posts).sort())
	blog.WithChildren(blog.Tree().withNode(comments).sort())
	search.WithChildren(search.Tree().withNode(searchAuthor).sort())
	page.WithChildren(page.Tree().withNode(pageID).sort())
	posts.WithChildren(posts.Tree().withNode(postsID).sort())
	comments.WithChildren(comments.Tree().withNode(commentID).sort())
	commentID.WithChildren(commentID.Tree().withNode(commentNew).sort())

	fmt.Printf("Raw tree:\n")
	fmt.Print(root.Tree().PrettyPrint())

	root.WithChildren(root.Tree().Compile())

	fmt.Printf("Compiled tree:\n")
	fmt.Print(root.Tree().PrettyPrint())

	// Output:
	// Raw tree:
	// 	{lang:en|pl}
	// 		blog
	// 		page
	// 		{pageId:[^/]+}
	// 	posts
	// 		{postsId:[^/]+}
	// 	search
	// 		author
	// 	comments
	// 		{commentId:\d+}
	// 		new
	// Compiled tree:
	// 	{lang:en|pl}
	// 		blog
	// 		page
	// 		{pageId:[^/]+}
	// 	posts
	// 		{postsId:[^/]+}
	// 	search/author
	// 	comments
	// 		{commentId:\d+}
	// 		new
}
