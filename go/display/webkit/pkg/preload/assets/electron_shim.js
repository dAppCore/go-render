(function () {
  if (globalThis.__corePreloadElectronInstalled) {
    return;
  }
  globalThis.__corePreloadElectronInstalled = true;

  const meta = __CORE_PRELOAD_META__;
  if (!meta.allow) {
    return;
  }

  const bridge = globalThis.__corePreloadBridge;
  if (!bridge) {
    return;
  }

  const listeners = new Map();
  const eventName = (channel) => "__core_preload_electron__:" + String(channel ?? "");

  const ipcRenderer = {
    send(channel, ...args) {
      const normalized = String(channel ?? "");
      globalThis.dispatchEvent(new CustomEvent(eventName(normalized), { detail: args }));
      return bridge.action(normalized, { channel: normalized, args }).then(() => undefined);
    },
    invoke(channel, ...args) {
      const normalized = String(channel ?? "");
      return bridge.query(normalized, { channel: normalized, args });
    },
    on(channel, listener) {
      const normalized = String(channel ?? "");
      const handler = (event) => listener(event, ...(event.detail || []));
      listeners.set(listener, handler);
      globalThis.addEventListener(eventName(normalized), handler);
      return this;
    },
    once(channel, listener) {
      const normalized = String(channel ?? "");
      const onceListener = (event, ...args) => {
        ipcRenderer.removeListener(normalized, listener);
        listener(event, ...args);
      };
      return ipcRenderer.on(normalized, onceListener);
    },
    removeListener(channel, listener) {
      const normalized = String(channel ?? "");
      const handler = listeners.get(listener);
      if (handler) {
        globalThis.removeEventListener(eventName(normalized), handler);
        listeners.delete(listener);
      }
      return this;
    }
  };

  const remote = {
    getGlobal(name) {
      return bridge.query("electron.remote.getGlobal", { name: String(name ?? "") });
    },
    app: {
      getPath(name) {
        return bridge.query("electron.app.getPath", { name: String(name ?? "") });
      }
    }
  };

  const shell = {
    openExternal(url) {
      return bridge.action("browser.openURL", { url: String(url ?? "") }).then(() => undefined);
    },
    openPath(path) {
      return bridge.action("browser.openFile", { path: String(path ?? "") }).then(() => "");
    }
  };

  const contextBridge = {
    exposeInMainWorld(name, api) {
      globalThis[name] = api;
    }
  };

  const processShim = globalThis.process || {
    env: {},
    platform: "wails",
    type: "renderer",
    versions: {}
  };
  processShim.versions = processShim.versions || {};
  processShim.versions.electron = processShim.versions.electron || "wails-shim";

  const electron = {
    ipcRenderer,
    remote,
    shell,
    contextBridge
  };

  globalThis.process = processShim;
  globalThis.electron = electron;
  globalThis.require = globalThis.require || ((name) => (name === "electron" ? electron : undefined));
})();
