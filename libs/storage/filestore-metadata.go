package storage

import (
	"context"
	"fmt"
	"log/slog"
	"path"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type FilestoreMetadata struct {
	projectID        string
	filestoreLatest  string
	thumbnailsPrefix string
}

// NewMetadateStorageOptions for MetadateStorage.
type NewFilestoreOptions struct {
	ProjectID       string
	FilestoreLatest string
	ThumbnailPerfix string
}

// NewFilestore with the given options.
// If no logger is provided, logs are discarded.
func NewFilestore(opts NewFilestoreOptions) (*FilestoreMetadata, error) {
	if opts.ProjectID == "" {
		return nil, fmt.Errorf("ProjectID must be provided")
	}
	if opts.FilestoreLatest == "" {
		return nil, fmt.Errorf("FilestoreLatest must be provided")
	}
	if opts.ThumbnailPerfix == "" {
		return nil, fmt.Errorf("ThumbnailPerfix must be provided")
	}
	return &FilestoreMetadata{
		projectID:        opts.ProjectID,
		thumbnailsPrefix: opts.ThumbnailPerfix,
		filestoreLatest:  opts.FilestoreLatest,
	}, nil
}

func (f *FilestoreMetadata) StorePath(ctx context.Context, path string, timestamp int64) error {
	slog.Info("Storing metadata Path")
	client, err := firestore.NewClient(ctx, f.projectID)
	if err != nil {
		slog.Error("firestore client failed", slog.String("error", err.Error()))
		return err
	}
	defer client.Close()

	// Store photo name under user.
	_, err = client.Doc(path).Set(ctx, map[string]interface{}{
		"uploaded_time": timestamp,
	})
	if err != nil {
		return fmt.Errorf("firestore create failed for %s: %v", path, err)
	}
	return nil
}

func (f *FilestoreMetadata) StorePathWithUser(ctx context.Context, user, path string, timestamp int64) error {
	slog.Info("Storing metadata Path with user")
	client, err := firestore.NewClient(ctx, f.projectID)
	if err != nil {
		slog.Error("firestore client failed", slog.String("error", err.Error()))
		return err
	}
	defer client.Close()

	// Store photo name under user.
	_, err = client.Doc(path).Set(ctx, map[string]interface{}{
		"uploaded_time": timestamp,
		"user":          user,
	})
	if err != nil {
		return fmt.Errorf("firestore create failed for %s: %v", path, err)
	}
	return nil
}

func (f *FilestoreMetadata) StoreLatest(ctx context.Context, index uint32, latest, objName string) error {
	slog.Info("Storing metadata Latest")
	client, err := firestore.NewClient(ctx, f.projectID)
	if err != nil {
		slog.Error("firestore client failed", slog.String("error", err.Error()))
		return err
	}
	defer client.Close()

	id := path.Join(latest, strconv.Itoa(int(index)))
	_, err = client.Doc(id).Set(ctx, map[string]interface{}{
		"obj_name": objName,
	})
	if err != nil {
		return fmt.Errorf("firestore create failed for %s: %v", objName, err)
	}
	return nil
}

func getData(doc *firestore.DocumentSnapshot) string {
	data := doc.Data()
	return data["obj_name"].(string)
}

func getID(doc *firestore.DocumentSnapshot) string {
	return doc.Ref.ID
}

func (f *FilestoreMetadata) GetThumbnails(
	ctx context.Context,
	thumbnailCount int,
	keyword string) ([]string, error) {

	client, err := firestore.NewClient(ctx, f.projectID)
	if err != nil {
		slog.Error("firestore client failed", slog.String("error", err.Error()))
		return nil, err
	}
	defer client.Close()

	filter := getID
	if keyword == "" {
		keyword = f.filestoreLatest
		filter = getData
	}
	q := client.Collection(keyword)
	if q == nil {
		return nil, nil
	}
	iter := q.Documents(ctx)

	var urls []string
	for i := 0; i < thumbnailCount; i++ {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		urls = append(urls, f.thumbnailsPrefix+filter(doc))
	}

	return urls, nil

}
