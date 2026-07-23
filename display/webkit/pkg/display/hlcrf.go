package display

import (
	"regexp"

	core "dappco.re/go"
)

var hlcrfSlotPattern = regexp.MustCompile(`\{\{\s*slot\s+"([^"]*)"\s*\}\}`)

func (s *Service) buildHLCRFComponents(pageURL string) (string, resultFailure) {
	loaded, err := s.loadManifestForOrigin(pageURL)
	if err != nil || loaded == nil {
		if err != nil && core.Contains(err.Error(), "view manifest not found") {
			return "", nil
		}
		return "", err
	}
	var scripts []string
	for _, component := range loaded.Manifest.HLCRF {
		templateBody := core.Trim(component.Template)
		if templateBody == "" && core.Trim(component.Name) != "" {
			resolvedPath, pathErr := safeManifestRelativePath(loaded.BaseDir, component.Name, "hlcrf component path")
			if pathErr != nil {
				if isMissingManifestPath(pathErr) {
					continue
				}
				return "", pathErr
			}
			body, readErr := coreReadFile(resolvedPath)
			if readErr != nil {
				if isMissingManifestPath(readErr) {
					continue
				}
				return "", readErr
			}
			templateBody = string(body)
		}
		if templateBody == "" {
			continue
		}
		tag := core.Trim(component.Tag)
		if tag == "" {
			tag = defaultHLCRFTag(component.Name)
		}
		scripts = append(scripts, renderHLCRFComponent(tag, templateBody))
	}
	return core.Join("\n", scripts...), nil
}

func isMissingManifestPath(err error) bool {
	if err == nil {
		return false
	}
	return core.IsNotExist(err) || core.Contains(err.Error(), "does not exist") || core.Contains(err.Error(), "no such file")
}

func renderHLCRFComponent(tag, templateBody string) string {
	templateBody = compileHLCRFTemplate(templateBody)
	return `(function(){if(customElements.get(` + quoteJS(tag) + `)){return;}const tpl=document.createElement('template');tpl.innerHTML=` +
		quoteJS(templateBody) +
		`;class CoreHLCRFElement extends HTMLElement{connectedCallback(){if(this.shadowRoot){return;}const root=this.attachShadow({mode:'open'});root.appendChild(tpl.content.cloneNode(true));}}customElements.define(` +
		quoteJS(tag) +
		`,CoreHLCRFElement);})();`
}

func compileHLCRFTemplate(templateBody string) string {
	return hlcrfSlotPattern.ReplaceAllStringFunc(templateBody, func(source string) string {
		match := hlcrfSlotPattern.FindStringSubmatch(source)
		if len(match) < 2 {
			return source
		}
		slotName := core.Trim(match[1])
		if slotName == "" || equalFold(slotName, "default") {
			return "<slot></slot>"
		}
		return `<slot name="` + slotName + `"></slot>`
	})
}

func defaultHLCRFTag(name string) string {
	name = core.Trim(core.Lower(name))
	name = core.TrimSuffix(name, core.PathExt(name))
	name = core.Replace(name, "_", "-")
	name = core.Replace(name, " ", "-")
	if !core.Contains(name, "-") {
		name = "core-" + name
	}
	return name
}
