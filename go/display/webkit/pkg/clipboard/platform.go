package clipboard

// Platform abstracts the system clipboard backend.
type Platform interface {
	Text() (string, bool)
	SetText(text string) bool
}

// ImagePlatform is an optional extension for clipboard backends that support images.
type ImagePlatform interface {
	Image() ([]byte, bool)
	SetImage(data []byte) bool
}

// ClipboardContent is the result of QueryText.
type ClipboardContent struct {
	Text       string `json:"text"`
	HasContent bool   `json:"hasContent"`
}

// ImageContent is the result of QueryImage.
type ImageContent struct {
	Data     []byte `json:"data"`
	HasImage bool   `json:"hasImage"`
}
