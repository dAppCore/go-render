package preload

import core "dappco.re/go"

func renderStoragePolyfills(pageURL string, canPersist bool) string {
	meta := map[string]any{
		"pageURL":       pageURL,
		"storageOrigin": storageOriginForPageURL(pageURL),
		"storeGroup":    "gui.preload.storage",
		"canPersist":    canPersist,
	}

	return core.Replace(
		storagePolyfillsAsset,
		"__CORE_PRELOAD_META__",
		core.JSONMarshalString(meta))

}
