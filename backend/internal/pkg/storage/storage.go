package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SavedFile struct {
	Key    string
	URL    string
	Size   int64
	SHA256 string
}

type LocalStore struct {
	root    string
	baseURL string
}

func NewLocalStore(root, baseURL string) *LocalStore {
	if root == "" {
		root = "./storage"
	}
	return &LocalStore{root: root, baseURL: strings.TrimRight(baseURL, "/")}
}

func (s *LocalStore) SaveUploaded(prefix string, file multipart.File, header *multipart.FileHeader) (*SavedFile, error) {
	if prefix == "" {
		prefix = "files"
	}
	ext := filepath.Ext(header.Filename)
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	key := filepath.ToSlash(filepath.Join(prefix, time.Now().Format("20060102"), name))
	fullPath := filepath.Join(s.root, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, err
	}
	out, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	hasher := sha256.New()
	size, err := io.Copy(io.MultiWriter(out, hasher), file)
	if err != nil {
		return nil, err
	}
	hash := "0x" + hex.EncodeToString(hasher.Sum(nil))
	return &SavedFile{Key: key, URL: s.publicURL(key), Size: size, SHA256: hash}, nil
}

func (s *LocalStore) SaveBytes(prefix, filename string, data []byte) (*SavedFile, error) {
	if prefix == "" {
		prefix = "files"
	}
	if filename == "" {
		filename = fmt.Sprintf("%d.bin", time.Now().UnixNano())
	}
	key := filepath.ToSlash(filepath.Join(prefix, time.Now().Format("20060102"), filename))
	fullPath := filepath.Join(s.root, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return nil, err
	}
	sum := sha256.Sum256(data)
	return &SavedFile{
		Key:    key,
		URL:    s.publicURL(key),
		Size:   int64(len(data)),
		SHA256: "0x" + hex.EncodeToString(sum[:]),
	}, nil
}

func (s *LocalStore) publicURL(key string) string {
	if s.baseURL == "" {
		return "/storage/" + filepath.ToSlash(key)
	}
	return s.baseURL + "/" + filepath.ToSlash(key)
}

func SHA256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return "0x" + hex.EncodeToString(sum[:])
}
