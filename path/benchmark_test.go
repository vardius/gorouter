package path

import (
	"testing"
)

func BenchmarkTrimSlash(b *testing.B) {
	path := "/pl/blog/comments/123/new/"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p := TrimSlash(path)

			if p != "pl/blog/comments/123/new" {
				b.Fatalf("%s", p)
			}
		}
	})
}

func BenchmarkGetPart(b *testing.B) {
	path := "pl/blog/comments/123/new"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p, _ := GetPart(path)

			if p != "pl" {
				b.Fatalf("%s", p)
			}
		}
	})
}
