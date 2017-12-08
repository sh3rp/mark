package mark

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

const BOOKMARK_BUCKET = "bookmark"

type Bookmark struct {
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	CreatedOn   time.Time `json:"created_on"`
}

type BookmarkService interface {
	SaveBookmark(Bookmark) error
	GetBookmarks() ([]Bookmark, error)
	GetBookmarksByTag(string) ([]Bookmark, error)
}

type GCPBookmarkService struct {
	bucket *storage.BucketHandle
	ctx    context.Context
}

func NewGCPBookmarkService() (BookmarkService, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return GCPBookmarkService{}, err
	}

	bucket := client.Bucket(BOOKMARK_BUCKET)

	return GCPBookmarkService{bucket, ctx}, nil
}

func (svc GCPBookmarkService) SaveBookmark(b Bookmark) error {
	hasher := sha256.New()
	hasher.Write([]byte(b.URL))
	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	obj := svc.bucket.Object(hash)
	writer := obj.NewWriter(svc.ctx)

	if err := json.NewEncoder(writer).Encode(&b); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}

func (svc GCPBookmarkService) GetBookmarks() ([]Bookmark, error) {
	var bookmarks []Bookmark
	objs := svc.bucket.Objects(svc.ctx, nil)

	for {
		objAttr, err := objs.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			return bookmarks, err
		}

		rawObj := svc.bucket.Object(objAttr.Name)
		reader, err := rawObj.NewReader(svc.ctx)

		if err != nil {
			return bookmarks, err
		}

		bookmark := &Bookmark{}

		err = json.NewDecoder(reader).Decode(bookmark)

		if err != nil {
			return bookmarks, err
		}

		bookmarks = append(bookmarks, *bookmark)
	}

	return bookmarks, nil
}

func (svc GCPBookmarkService) GetBookmarksByTag(tag string) ([]Bookmark, error) {
	return nil, nil
}
