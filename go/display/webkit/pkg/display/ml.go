package display

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/chat"
)

type ModelRuntimeState struct {
	APIURL       string            `json:"api_url"`
	Loaded       []chat.ModelEntry `json:"loaded"`
	Available    []chat.ModelEntry `json:"available"`
	VRAMBytes    int64             `json:"vram_bytes"`
	Backend      string            `json:"backend"`
	InferenceURL string            `json:"inference_url"`
}

func (s *Service) modelState() ModelRuntimeState {
	apiURL := trimRight(core.Trim(core.Getenv("CORE_ML_API_URL")), "/")
	if apiURL == "" {
		apiURL = "http://localhost:8090"
	}
	models := s.chatModels()
	loaded := make([]chat.ModelEntry, 0, len(models))
	available := make([]chat.ModelEntry, 0, len(models))
	for _, model := range models {
		if model.Loaded {
			loaded = append(loaded, model)
			continue
		}
		available = append(available, model)
	}
	return ModelRuntimeState{
		APIURL:       apiURL,
		Loaded:       loaded,
		Available:    available,
		VRAMBytes:    estimateVRAM(models),
		Backend:      dominantBackend(models),
		InferenceURL: apiURL + "/v1/chat/completions",
	}
}

func (s *Service) chatModels() []chat.ModelEntry {
	if s == nil || s.ServiceRuntime == nil {
		return nil
	}
	result := s.Core().Action("gui.chat.models").Run(context.Background(), core.NewOptions())
	models, _ := result.Value.([]chat.ModelEntry)
	return models
}

func estimateVRAM(models []chat.ModelEntry) int64 {
	var total int64
	for _, model := range models {
		if model.Loaded {
			total += model.SizeBytes
		}
	}
	return total
}

func dominantBackend(models []chat.ModelEntry) string {
	for _, model := range models {
		if core.Trim(model.Backend) != "" {
			return model.Backend
		}
	}
	return "local"
}

func quoteJS(value string) string {
	escaped := core.Replace(value, `\`, `\\`)
	escaped = core.Replace(escaped, "\n", `\n`)
	escaped = core.Replace(escaped, "\u2028", `\u2028`)
	escaped = core.Replace(escaped, "\u2029", `\u2029`)
	escaped = core.Replace(escaped, `"`, `\"`)
	escaped = escapeClosingScriptTag(escaped)
	return `"` + escaped + `"`
}

func escapeClosingScriptTag(value string) string {
	for {
		index := indexString(core.Lower(value), "</script>")
		if index < 0 {
			return value
		}
		value = value[:index] + `<\/script>` + value[index+len("</script>"):]
	}
}
