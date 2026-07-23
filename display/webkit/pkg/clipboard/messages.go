// pkg/clipboard/messages.go
package clipboard

const MaxImageBytes = 16 << 20

// QueryText reads the clipboard. Result: ClipboardContent
type QueryText struct{}

// QueryImage reads image data from the clipboard. Result: ImageContent
type QueryImage struct{}

// TaskSetText writes text to the clipboard. Result: bool (success)
type TaskSetText struct{ Text string }

// TaskSetImage writes image data to the clipboard. Result: bool (success)
type TaskSetImage struct{ Data []byte }

// TaskClear clears the clipboard. Result: bool (success)
type TaskClear struct{}

// ClipboardImageContent contains clipboard image data encoded for transport.
type ClipboardImageContent struct {
	Base64     string `json:"base64"`
	MimeType   string `json:"mimeType"`
	HasContent bool   `json:"hasContent"`
}
