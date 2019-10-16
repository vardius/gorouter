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
