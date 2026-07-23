// pkg/contextmenu/wails.go
package contextmenu

import (
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// WailsPlatform implements Platform via Wails v3's ContextMenuManager.
//
// Frontend trigger pattern: a DOM element gets the CSS custom property
//
//	--custom-contextmenu: <menu-name>;
//
// (optionally with --custom-contextmenu-data: <string>; for per-element
// payload). On right-click Wails reads the CSS, opens the registered
// menu by name, and calls our OnClick handlers when the user picks an
// item — those handlers forward to the Platform's onItemClick callback,
// which the gui contextmenu service translates into ActionItemClicked
// dispatches the consumer can listen for.
//
//	wp := contextmenu.NewWailsPlatform(app)
//	core.WithService(contextmenu.Register(wp))
type WailsPlatform struct {
	app *application.App

	mu       sync.RWMutex
	registry map[string]*registeredMenu
}

// registeredMenu keeps both the wails handle (for Destroy on Remove)
// and the original definition (for Get / GetAll).
type registeredMenu struct {
	wails *application.ContextMenu
	def   ContextMenuDef
}

func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app, registry: make(map[string]*registeredMenu)}
}

// Add (re)registers a context menu by name. If the same name was already
// registered, the previous wails ContextMenu is destroyed first so the
// definition is fully replaced — matches the "update" semantics the gui
// service exposes via contextmenu.update.
func (wp *WailsPlatform) Add(name string, menu ContextMenuDef, onItemClick func(menuName, actionID, data string)) resultFailure {
	if wp == nil || wp.app == nil {
		return nil
	}
	wp.mu.Lock()
	if existing, ok := wp.registry[name]; ok && existing.wails != nil {
		existing.wails.Destroy()
	}
	wp.mu.Unlock()

	cm := application.NewContextMenu(name)
	for _, item := range menu.Items {
		buildItem(cm.Menu, item, name, onItemClick)
	}
	cm.Update()

	wp.mu.Lock()
	wp.registry[name] = &registeredMenu{wails: cm, def: menu}
	wp.mu.Unlock()
	return nil
}

// Remove unregisters a context menu by name. No-op if absent.
func (wp *WailsPlatform) Remove(name string) resultFailure {
	if wp == nil {
		return nil
	}
	wp.mu.Lock()
	defer wp.mu.Unlock()
	entry, ok := wp.registry[name]
	if !ok {
		return nil
	}
	if entry.wails != nil {
		entry.wails.Destroy()
	}
	delete(wp.registry, name)
	return nil
}

// Get returns the definition we registered, not the live wails state —
// the gui service uses this for round-tripping definitions back to MCP
// callers.
func (wp *WailsPlatform) Get(name string) (*ContextMenuDef, bool) {
	if wp == nil {
		return nil, false
	}
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	entry, ok := wp.registry[name]
	if !ok {
		return nil, false
	}
	def := entry.def
	return &def, true
}

func (wp *WailsPlatform) GetAll() map[string]ContextMenuDef {
	if wp == nil {
		return map[string]ContextMenuDef{}
	}
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	out := make(map[string]ContextMenuDef, len(wp.registry))
	for name, entry := range wp.registry {
		out[name] = entry.def
	}
	return out
}

// buildItem walks a single MenuItemDef into the wails Menu. Handles
// separator / submenu / checkbox / radio / normal item shapes; recurses
// for submenus. Each non-separator item gets an OnClick handler that
// captures its ActionID and forwards to the consumer's onItemClick
// alongside the per-click data payload Wails extracts from
// --custom-contextmenu-data on the right-clicked element.
func buildItem(parent *application.Menu, def MenuItemDef, menuName string, onItemClick func(menuName, actionID, data string)) {
	switch def.Type {
	case "separator":
		parent.AddSeparator()
		return
	case "submenu":
		sub := parent.AddSubmenu(def.Label)
		for _, child := range def.Items {
			buildItem(sub, child, menuName, onItemClick)
		}
		return
	}

	var item *application.MenuItem
	switch def.Type {
	case "checkbox":
		item = parent.AddCheckbox(def.Label, def.Checked)
	case "radio":
		item = parent.AddRadio(def.Label, def.Checked)
	default:
		item = parent.Add(def.Label)
	}

	if def.Accelerator != "" {
		item.SetAccelerator(def.Accelerator)
	}
	if def.Enabled != nil {
		item.SetEnabled(*def.Enabled)
	}

	if onItemClick != nil && def.ActionID != "" {
		actionID := def.ActionID
		item.OnClick(func(ctx *application.Context) {
			data := ""
			if ctx != nil {
				data = ctx.ContextMenuData()
			}
			onItemClick(menuName, actionID, data)
		})
	}
}
