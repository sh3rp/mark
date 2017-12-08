package main

import (
	"fmt"
	"os"

	"github.com/sh3rp/mark"
)

func main() {
	svc, err := mark.NewGCPBookmarkService()

	if err != nil {
		panic(err)
	}

	if len(os.Args) == 1 {
		bookmarks, err := svc.GetBookmarks()
		if err != nil {
			panic(err)
		}
		for _, bookmark := range bookmarks {
			fmt.Printf("%s\n", bookmark.URL)
		}
	} else {

		theURL := os.Args[1]

		bookmark, err := mark.BookmarkFromLink(theURL)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if err := svc.SaveBookmark(bookmark); err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println("OK")
		}
	}
}
