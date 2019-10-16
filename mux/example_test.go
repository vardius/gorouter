package mux_test

import (
	"fmt"

	"github.com/vardius/gorouter/v4/mux"
)

func Example() {
	root := mux.NewNode("GET", 0)

	lang := mux.NewNode("{lang:en|pl}", root.MaxParamsSize())
	blog := mux.NewNode("blog", lang.MaxParamsSize())
	search := mux.NewNode("search", blog.MaxParamsSize())

	page := mux.NewNode("page", blog.MaxParamsSize())
	pageID := mux.NewNode(`{pageId:[^/]+}`, page.MaxParamsSize())

	posts := mux.NewNode("posts", blog.MaxParamsSize())
	postsID := mux.NewNode(`{postsId:[^/]+}`, posts.MaxParamsSize())

	comments := mux.NewNode("comments", blog.MaxParamsSize())
	commentID := mux.NewNode(`{commentId:\d+}`, comments.MaxParamsSize())
	commentNew := mux.NewNode("new", commentID.MaxParamsSize())

	root.WithChildren(root.Tree().WithNode(lang))
	lang.WithChildren(lang.Tree().WithNode(blog))
	blog.WithChildren(blog.Tree().WithNode(search))
	blog.WithChildren(blog.Tree().WithNode(page))
	blog.WithChildren(blog.Tree().WithNode(posts))
	blog.WithChildren(blog.Tree().WithNode(comments))
	page.WithChildren(page.Tree().WithNode(pageID))
	posts.WithChildren(posts.Tree().WithNode(postsID))
	comments.WithChildren(comments.Tree().WithNode(commentID))
	commentID.WithChildren(commentID.Tree().WithNode(commentNew))

	fmt.Print(root.Tree().PrettyPrint())

	root.WithChildren(root.Tree().Compile())

	fmt.Printf("\n\n")
	fmt.Print(root.Tree().PrettyPrint())

	// Output:
	// {lang:^(?P<lang>en|pl)}
	// 	blog
	// 		comments
	// 		{commentId:^(?P<commentId>\d+)}
	// 		new
	// 	posts
	// 		{postsId:^(?P<postsId>[^/]+)}
	// 	page
	// 		{pageId:^(?P<pageId>[^/]+)}
	// 	search
	//
	// 	{lang-blog:^(?P<lang>en|pl)\/blog}
	// 		{comments-commentId-new:^comments\/(?P<commentId>\d+)\/new}
	// 		{posts-postsId:^posts\/(?P<postsId>[/]+)}
	// 		{page-pageId:^page\/(?P<pageId>[/]+)}
	// 		search
}
