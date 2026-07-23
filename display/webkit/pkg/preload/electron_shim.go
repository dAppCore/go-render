package preload

import core "dappco.re/go"

func renderElectronShim(pageURL string) string {
	meta := map[string]any{
		"allow":   true,
		"pageURL": pageURL,
	}

	return core.Replace(
		electronShimAsset,
		"__CORE_PRELOAD_META__",
		core.JSONMarshalString(meta))

}
