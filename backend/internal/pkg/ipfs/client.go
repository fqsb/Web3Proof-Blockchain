package ipfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Client struct {
	pinataJWT string
	storageDir string
}

type PinResult struct {
	CID string
	URL string
}

func New(pinataJWT, storageDir string) *Client {
	if storageDir == "" {
		storageDir = "./storage/ipfs"
	}
	return &Client{pinataJWT: pinataJWT, storageDir: storageDir}
}

func (c *Client) PinJSON(name string, data interface{}) (*PinResult, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	if c.pinataJWT != "" {
		return c.pinToPinata(name, body, "application/json")
	}
	return c.pinLocal(name, body)
}

func (c *Client) PinFile(name string, file multipart.File, header *multipart.FileHeader) (*PinResult, error) {
	body, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if c.pinataJWT != "" {
		return c.pinFileToPinata(name, header.Filename, body, header.Header.Get("Content-Type"))
	}
	return c.pinLocal(name+"-"+header.Filename, body)
}

func (c *Client) pinToPinata(name string, body []byte, contentType string) (*PinResult, error) {
	payload := map[string]interface{}{
		"pinataContent": json.RawMessage(body),
		"pinataMetadata": map[string]string{
			"name": name,
		},
	}
	raw, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "https://api.pinata.cloud/pinning/pinJSONToIPFS", bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.pinataJWT)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("pinata error: %s", string(b))
	}
	var result struct {
		IpfsHash string `json:"IpfsHash"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &PinResult{
		CID: result.IpfsHash,
		URL: "https://gateway.pinata.cloud/ipfs/" + result.IpfsHash,
	}, nil
}

func (c *Client) pinFileToPinata(name, filename string, body []byte, contentType string) (*PinResult, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("pinataMetadata", fmt.Sprintf(`{"name":"%s"}`, name))
	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(body); err != nil {
		return nil, err
	}
	w.Close()

	req, err := http.NewRequest(http.MethodPost, "https://api.pinata.cloud/pinning/pinFileToIPFS", &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.pinataJWT)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("pinata file error: %s", string(b))
	}
	var result struct {
		IpfsHash string `json:"IpfsHash"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &PinResult{CID: result.IpfsHash, URL: "https://gateway.pinata.cloud/ipfs/" + result.IpfsHash}, nil
}

func (c *Client) pinLocal(name string, body []byte) (*PinResult, error) {
	cid := strings.TrimPrefix(ContentHashHex(body), "0x")
	cid = "bafyDEV" + cid[:52]
	if err := os.MkdirAll(c.storageDir, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(c.storageDir, cid+".json")
	meta := map[string]interface{}{
		"name":      name,
		"cid":       cid,
		"content":   json.RawMessage(body),
		"pinned_at": time.Now().UTC().Format(time.RFC3339),
	}
	raw, _ := json.MarshalIndent(meta, "", "  ")
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return nil, err
	}
	return &PinResult{CID: cid, URL: "local://" + strings.ReplaceAll(path, "\\", "/")}, nil
}
