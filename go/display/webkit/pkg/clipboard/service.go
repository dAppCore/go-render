// pkg/clipboard/service.go
package clipboard

import (
	"context"
	"encoding/base64"

	core "dappco.re/go"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// Register(p) binds the clipboard service to a Core instance.
// c.WithService(clipboard.Register(wailsClipboard))
func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, OK: true}
	}
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	s.Core().RegisterQuery(s.handleQuery)
	setText := func(_ context.Context, opts core.Options) core.Result {
		success := s.platform.SetText(opts.String("text"))
		return core.Result{Value: success, OK: true}
	}
	setImage := func(_ context.Context, opts core.Options) core.Result {
		imgPlatform, ok := s.platform.(ImagePlatform)
		if !ok {
			return core.Result{Value: false, OK: true}
		}
		data, err := clipboardImageData(opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		if len(data) == 0 || len(data) > MaxImageBytes {
			return core.Result{Value: false, OK: true}
		}
		success := imgPlatform.SetImage(data)
		return core.Result{Value: success, OK: true}
	}
	clear := func(_ context.Context, _ core.Options) core.Result {
		oldText, hadText := s.platform.Text()
		var oldImage []byte
		var hadImage bool
		if imgPlatform, ok := s.platform.(ImagePlatform); ok {
			oldImage, hadImage = imgPlatform.Image()
			oldImage = append([]byte(nil), oldImage...)
			if !imgPlatform.SetImage(nil) {
				return core.Result{Value: false, OK: true}
			}
			if !s.platform.SetText("") {
				if hadImage {
					if !imgPlatform.SetImage(oldImage) {
						return core.Result{Value: false, OK: true}
					}
				}
				return core.Result{Value: false, OK: true}
			}
			return core.Result{Value: true, OK: true}
		}
		success := s.platform.SetText("")
		if !success && hadText {
			if !s.platform.SetText(oldText) {
				return core.Result{Value: false, OK: true}
			}
		}
		return core.Result{Value: success, OK: true}
	}
	read := func(_ context.Context, _ core.Options) core.Result {
		text, ok := s.platform.Text()
		return core.Result{Value: ClipboardContent{Text: text, HasContent: ok && text != ""}, OK: true}
	}
	s.Core().Action("clipboard.set_text", setText)
	s.Core().Action("gui.clipboard.write", setText)
	s.Core().Action("clipboard.set_image", setImage)
	s.Core().Action("gui.clipboard.writeImage", setImage)
	s.Core().Action("gui.clipboard.readImage", func(_ context.Context, _ core.Options) core.Result {
		imgPlatform, ok := s.platform.(ImagePlatform)
		if !ok {
			return core.Result{Value: ImageContent{}, OK: true}
		}
		data, hasImage := imgPlatform.Image()
		return core.Result{Value: ImageContent{Data: append([]byte(nil), data...), HasImage: hasImage && len(data) > 0}, OK: true}
	})
	s.Core().Action("clipboard.clear", clear)
	s.Core().Action("gui.clipboard.clear", clear)
	s.Core().Action("gui.clipboard.read", read)
	return core.Result{OK: true}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch q.(type) {
	case QueryText:
		text, ok := s.platform.Text()
		return core.Result{Value: ClipboardContent{Text: text, HasContent: ok && text != ""}, OK: true}
	case QueryImage:
		imgPlatform, ok := s.platform.(ImagePlatform)
		if !ok {
			return core.Result{Value: ImageContent{}, OK: true}
		}
		data, hasImage := imgPlatform.Image()
		return core.Result{Value: ImageContent{Data: append([]byte(nil), data...), HasImage: hasImage && len(data) > 0}, OK: true}
	default:
		return core.Result{}
	}
}

// clipboardImageData normalizes clipboard image inputs from MCP, preload bridge, and WS callers.
// Use: bytes, err := clipboardImageData(core.NewOptions(core.Option{Key: "data", Value: "iVBORw0KGgo..."}))
func clipboardImageData(opts core.Options) ([]byte, resultFailure) {
	if raw, ok := opts.Get("data").Value.([]byte); ok && len(raw) > 0 {
		return append([]byte(nil), raw...), nil
	}
	encoded := opts.String("data")
	if encoded == "" {
		return nil, nil
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, core.E("clipboard.imageData", "invalid base64 image data", err)
	}
	return data, nil
}
