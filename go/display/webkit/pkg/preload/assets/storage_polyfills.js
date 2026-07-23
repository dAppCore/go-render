(function () {
  if (globalThis.__corePreloadStorageInstalled) {
    return;
  }
  globalThis.__corePreloadStorageInstalled = true;

  const meta = __CORE_PRELOAD_META__;
  const pageURL = String(meta.pageURL || "");
  const storageOrigin = String(meta.storageOrigin || pageURL || "");
  const storeGroup = String(meta.storeGroup || "gui.preload.storage");
  const canPersist = !!meta.canPersist;

  const asPromise = (value) => (
    value && typeof value.then === "function" ? value : Promise.resolve(value)
  );

  const runCoreCall = (target, methodNames, name, payload) => {
    if (!target || typeof target !== "object") {
      return undefined;
    }
    for (const methodName of methodNames) {
      const method = target[methodName];
      if (typeof method !== "function") {
        continue;
      }
      try {
        const direct = method.call(target, name, payload);
        if (direct && typeof direct.Run === "function") {
          try {
            return direct.Run(payload);
          } catch (_) {
            return direct.Run();
          }
        }
        return direct;
      } catch (_) {
        try {
          const deferred = method.call(target, name);
          if (deferred && typeof deferred.Run === "function") {
            try {
              return deferred.Run(payload);
            } catch (_) {
              return deferred.Run();
            }
          }
          return deferred;
        } catch (_) {}
      }
    }
    return undefined;
  };

  const bridge = globalThis.__corePreloadBridge || (globalThis.__corePreloadBridge = {
    action(name, payload) {
      const candidates = [globalThis.c, globalThis.Core, globalThis.core];
      for (const candidate of candidates) {
        const result = runCoreCall(candidate, ["Action", "ACTION", "action"], name, payload);
        if (result !== undefined) {
          return asPromise(result);
        }
      }
      if (typeof globalThis.__CORE_GUI_INVOKE__ === "function") {
        return asPromise(globalThis.__CORE_GUI_INVOKE__(name, payload, { mode: "action" }));
      }
      return Promise.resolve(undefined);
    },
    query(name, payload) {
      const candidates = [globalThis.c, globalThis.Core, globalThis.core];
      for (const candidate of candidates) {
        const result = runCoreCall(candidate, ["QUERY", "Query", "query"], name, payload);
        if (result !== undefined) {
          return asPromise(result);
        }
      }
      if (typeof globalThis.__CORE_GUI_INVOKE__ === "function") {
        return asPromise(globalThis.__CORE_GUI_INVOKE__(name, payload, { mode: "query" }));
      }
      return Promise.resolve(undefined);
    }
  });

  const storageScopes = globalThis.__corePreloadStorageScopes || (globalThis.__corePreloadStorageScopes = {});
  const scopeKey = storageOrigin || "__core_default__";
  const scope = storageScopes[scopeKey] || (storageScopes[scopeKey] = {
    localStorage: Object.create(null),
    sessionStorage: Object.create(null),
    indexedDB: Object.create(null)
  });

  const persistKey = (bucket, key) => [storageOrigin, bucket, String(key ?? "")].join(":");

  const persistSet = (bucket, key, value) => {
    if (!canPersist) {
      return;
    }
    bridge.action("store.set", {
      group: storeGroup,
      key: persistKey(bucket, key),
      value: String(value ?? "")
    }).catch(() => undefined);
  };

  const persistDelete = (bucket, key) => {
    if (!canPersist) {
      return;
    }
    bridge.action("store.delete", {
      group: storeGroup,
      key: persistKey(bucket, key)
    }).catch(() => undefined);
  };

  const createStorage = (bucketName, bucket) => ({
    getItem(key) {
      const normalized = String(key ?? "");
      return Object.prototype.hasOwnProperty.call(bucket, normalized) ? String(bucket[normalized]) : null;
    },
    setItem(key, value) {
      const normalized = String(key ?? "");
      bucket[normalized] = String(value ?? "");
      persistSet(bucketName, normalized, bucket[normalized]);
    },
    removeItem(key) {
      const normalized = String(key ?? "");
      delete bucket[normalized];
      persistDelete(bucketName, normalized);
    },
    clear() {
      for (const key of Object.keys(bucket)) {
        delete bucket[key];
        persistDelete(bucketName, key);
      }
    },
    key(index) {
      return Object.keys(bucket)[Number(index)] ?? null;
    },
    get length() {
      return Object.keys(bucket).length;
    }
  });

  const queueTask = (callback) => {
    if (typeof queueMicrotask === "function") {
      queueMicrotask(callback);
      return;
    }
    Promise.resolve().then(callback).catch(() => undefined);
  };

  const createRequest = (result, upgrade) => {
    const request = { result, error: null, onsuccess: null, onerror: null, onupgradeneeded: null };
    queueTask(() => {
      if (upgrade) {
        request.onupgradeneeded?.({ target: request });
      }
      request.onsuccess?.({ target: request });
    });
    return request;
  };

  const serializeRecord = (value) => {
    if (typeof value === "string") {
      return value;
    }
    try {
      return JSON.stringify(value);
    } catch (_) {
      return String(value ?? "");
    }
  };

  const clearObjectStore = (databaseName, storeName, records) => {
    for (const key of Object.keys(records)) {
      persistDelete("indexeddb:" + databaseName + ":" + storeName, key);
    }
  };

  const createObjectStore = (databaseName, database, storeName) => ({
    put(value, key) {
      const resolvedKey = String(key ?? value?.id ?? Date.now());
      database.stores[storeName] = database.stores[storeName] || Object.create(null);
      database.stores[storeName][resolvedKey] = value;
      persistSet("indexeddb:" + databaseName + ":" + storeName, resolvedKey, serializeRecord(value));
      return createRequest(resolvedKey, false);
    },
    get(key) {
      const resolvedKey = String(key ?? "");
      return createRequest(database.stores?.[storeName]?.[resolvedKey], false);
    },
    getAll() {
      return createRequest(Object.values(database.stores?.[storeName] || {}), false);
    },
    delete(key) {
      const resolvedKey = String(key ?? "");
      if (database.stores?.[storeName]) {
        delete database.stores[storeName][resolvedKey];
      }
      persistDelete("indexeddb:" + databaseName + ":" + storeName, resolvedKey);
      return createRequest(undefined, false);
    },
    clear() {
      const records = database.stores?.[storeName] || Object.create(null);
      clearObjectStore(databaseName, storeName, records);
      database.stores[storeName] = Object.create(null);
      return createRequest(undefined, false);
    },
    createIndex() {
      return this;
    }
  });

  const createDatabase = (name, upgrade) => {
    const database = scope.indexedDB[name] || (scope.indexedDB[name] = { stores: Object.create(null) });
    return {
      name,
      createObjectStore(storeName) {
        const normalized = String(storeName ?? "default");
        database.stores[normalized] = database.stores[normalized] || Object.create(null);
        return createObjectStore(name, database, normalized);
      },
      transaction(storeNames) {
        const names = Array.isArray(storeNames) ? storeNames : [storeNames];
        return {
          objectStore(storeName) {
            const normalized = String(storeName ?? names[0] ?? "default");
            database.stores[normalized] = database.stores[normalized] || Object.create(null);
            return createObjectStore(name, database, normalized);
          }
        };
      },
      close() {}
    };
  };

  globalThis.core = globalThis.core || {};
  globalThis.core.storage = globalThis.core.storage || {};
  globalThis.core.storage.local = createStorage("localStorage", scope.localStorage);
  globalThis.core.storage.session = createStorage("sessionStorage", scope.sessionStorage);

  try {
    Object.defineProperty(globalThis, "localStorage", {
      configurable: true,
      enumerable: true,
      get() {
        return globalThis.core.storage.local;
      }
    });
  } catch (_) {}

  try {
    Object.defineProperty(globalThis, "sessionStorage", {
      configurable: true,
      enumerable: true,
      get() {
        return globalThis.core.storage.session;
      }
    });
  } catch (_) {}

  try {
    if (!globalThis.indexedDB) {
      globalThis.indexedDB = {
        open(name) {
          const normalized = String(name ?? "default");
          const upgrade = !scope.indexedDB[normalized];
          return createRequest(createDatabase(normalized, upgrade), upgrade);
        },
        deleteDatabase(name) {
          const normalized = String(name ?? "default");
          const database = scope.indexedDB[normalized];
          if (database && database.stores) {
            for (const [storeName, records] of Object.entries(database.stores)) {
              clearObjectStore(normalized, storeName, records);
            }
          }
          delete scope.indexedDB[normalized];
          return createRequest(undefined, false);
        }
      };
    }
  } catch (_) {}
})();
