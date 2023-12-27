package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
)

type Local struct {
	ids   map[string]string
	store map[string][]byte
}

func NewLocal() *Local {
	return &Local{
		ids:   map[string]string{},
		store: map[string][]byte{},
	}
}

func (l *Local) GetFileNameFromURL(string) string {
	return ""
}

func (l *Local) Upload(ctx context.Context, file *File) (string, error) {
	id := uuid.NewString()
	buf, err := io.ReadAll(file.Content)
	if err != nil {
		return "", err
	}
	defer file.Content.Close()
	l.ids[id] = file.Name
	l.store[file.Name] = buf
	return id, nil
}

func (l *Local) Get(ctx context.Context, fileName string) (io.ReadCloser, error) {
	b, ok := l.store[fileName]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return io.NopCloser(bytes.NewBuffer(b)), nil
}

func (l *Local) GetURLFromFileName(name string) string {
	return "www.google.com/image.jpg"
}
