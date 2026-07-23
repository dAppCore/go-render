// pkg/webview/diagnostics.go
package webview

import (
	corego "dappco.re/go"
)

func jsQuote(v string) string {
	return corego.JSONMarshalString(v)
}

func computedStyleScript(selector string) string {
	sel := jsQuote(selector)
	return corego.Sprintf(`(function(){
  const el = document.querySelector(%s);
  if (!el) return null;
  const style = window.getComputedStyle(el);
  const out = {};
  for (let i = 0; i < style.length; i++) {
    const key = style[i];
    out[key] = style.getPropertyValue(key);
  }
  return out;
})()`, sel)
}

func ComputedStyleScript(selector string) string {
	return computedStyleScript(selector)
}

func highlightScript(selector, colour string) string {
	sel := jsQuote(selector)
	if colour == "" {
		colour = "#ff9800"
	}
	col := jsQuote(colour)
	return corego.Sprintf(`(function(){
  const el = document.querySelector(%s);
  if (!el) return false;
  if (el.__coreHighlightOrigOutline === undefined) {
    el.__coreHighlightOrigOutline = el.style.outline || "";
  }
  el.style.outline = "3px solid " + %s;
  el.style.outlineOffset = "2px";
  try { el.scrollIntoView({block: "center", inline: "center", behavior: "smooth"}); } catch (e) {}
  return true;
})()`, sel, col)
}

func HighlightScript(selector, colour string) string {
	return highlightScript(selector, colour)
}

func performanceScript() string {
	return `(function(){
  const nav = performance.getEntriesByType("navigation")[0] || {};
  const paints = performance.getEntriesByType("paint");
  const firstPaint = paints.find((entry) => entry.name === "first-paint");
  const firstContentfulPaint = paints.find((entry) => entry.name === "first-contentful-paint");
  const memory = performance.memory || {};
  return {
    navigationStart: nav.startTime || 0,
    domContentLoaded: nav.domContentLoadedEventEnd || 0,
    loadEventEnd: nav.loadEventEnd || 0,
    firstPaint: firstPaint ? firstPaint.startTime : 0,
    firstContentfulPaint: firstContentfulPaint ? firstContentfulPaint.startTime : 0,
    usedJSHeapSize: memory.usedJSHeapSize || 0,
    totalJSHeapSize: memory.totalJSHeapSize || 0
  };
})()`
}

func PerformanceScript() string {
	return performanceScript()
}

func resourcesScript() string {
	return `(function(){
  return performance.getEntriesByType("resource").map((entry) => ({
    name: entry.name,
    entryType: entry.entryType,
    initiatorType: entry.initiatorType,
    startTime: entry.startTime,
    duration: entry.duration,
    transferSize: entry.transferSize || 0,
    encodedBodySize: entry.encodedBodySize || 0,
    decodedBodySize: entry.decodedBodySize || 0
  }));
})()`
}

func ResourcesScript() string {
	return resourcesScript()
}

func networkInitScript() string {
	return `(function(){
  if (window.__coreNetworkLog) return true;
  window.__coreNetworkLog = [];
  const log = (entry) => { window.__coreNetworkLog.push(entry); };
  const originalFetch = window.fetch;
  if (originalFetch) {
    window.fetch = async function(input, init) {
      const request = typeof input === "string" ? input : (input && input.url) ? input.url : "";
      const method = (init && init.method) || (input && input.method) || "GET";
      const started = Date.now();
      try {
        const response = await originalFetch.call(this, input, init);
        log({
          url: response.url || request,
          method: method,
          status: response.status,
          ok: response.ok,
          resource: "fetch",
          timestamp: started
        });
        return response;
      } catch (error) {
        log({
          url: request,
          method: method,
          error: String(error),
          resource: "fetch",
          timestamp: started
        });
        throw error;
      }
    };
  }
  const originalOpen = XMLHttpRequest.prototype.open;
  const originalSend = XMLHttpRequest.prototype.send;
  XMLHttpRequest.prototype.open = function(method, url) {
    this.__coreMethod = method;
    this.__coreUrl = url;
    return originalOpen.apply(this, arguments);
  };
  XMLHttpRequest.prototype.send = function(body) {
    const started = Date.now();
    this.addEventListener("loadend", () => {
      log({
        url: this.__coreUrl || "",
        method: this.__coreMethod || "GET",
        status: this.status || 0,
        ok: this.status >= 200 && this.status < 400,
        resource: "xhr",
        timestamp: started
      });
    });
    return originalSend.apply(this, arguments);
  };
  return true;
})()`
}

func NetworkInitScript() string {
	return networkInitScript()
}

func networkClearScript() string {
	return `(function(){
  window.__coreNetworkLog = [];
  return true;
})()`
}

func NetworkClearScript() string {
	return networkClearScript()
}

func networkLogScript(limit int) string {
	if limit <= 0 {
		return `((window.__coreNetworkLog && window.__coreNetworkLog.length)
  ? window.__coreNetworkLog
  : performance.getEntriesByType("resource").map((entry) => ({
      url: entry.name,
      method: "GET",
      status: 0,
      ok: true,
      resource: entry.initiatorType || entry.entryType,
      timestamp: entry.startTime
    })))`
	}
	return corego.Sprintf(`(function(){
  const log = (window.__coreNetworkLog && window.__coreNetworkLog.length)
    ? window.__coreNetworkLog
    : performance.getEntriesByType("resource").map((entry) => ({
        url: entry.name,
        method: "GET",
        status: 0,
        ok: true,
        resource: entry.initiatorType || entry.entryType,
        timestamp: entry.startTime
      }));
  return log.slice(-%d);
})()`, limit)
}

func NetworkLogScript(limit int) string {
	return networkLogScript(limit)
}

func normalizeWhitespace(s string) string {
	return corego.Trim(s)
}
