// SPDX-Licence-Identifier: EUPL-1.2

package ctml

import (
	"testing"

	core "dappco.re/go"
	html "dappco.re/go/html/engine/html"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubcommandList_Good(t *testing.T) {
	c := core.New()
	require.True(t, c.Command("deploy", core.Command{Description: "Deploy the application"}).OK)
	require.True(t, c.Command("status", core.Command{}).OK)

	tree := SubcommandList(c, "", []string{"deploy", "status"})
	out := html.Render(tree, html.NewContext())

	assert.Contains(t, out, "<ul>")
	assert.Contains(t, out, "Deploy the application")
	assert.Contains(t, out, "cmd.status.description", "a command with no Description falls back to its derived I18nKey")
}

func TestSubcommandList_Good_NestedRoot(t *testing.T) {
	c := core.New()
	require.True(t, c.Command("deploy/to/homelab", core.Command{Description: "Deploy to homelab"}).OK)
	require.True(t, c.Command("deploy/to/qa", core.Command{Description: "Deploy to QA"}).OK)

	paths := []string{"deploy", "deploy/to", "deploy/to/homelab", "deploy/to/qa"}
	tree := SubcommandList(c, "deploy", paths)
	out := html.Render(tree, html.NewContext())

	assert.Contains(t, out, "cmd.deploy.to.description", "the direct child \"deploy/to\" is listed")
	assert.NotContains(t, out, "Deploy to homelab", "a grandchild is not listed at this root")
	assert.NotContains(t, out, "Deploy to QA")
}

func TestSubcommandList_Good_HiddenSkipped(t *testing.T) {
	c := core.New()
	require.True(t, c.Command("public", core.Command{Description: "Public command"}).OK)
	require.True(t, c.Command("secret", core.Command{Description: "Secret command", Hidden: true}).OK)

	tree := SubcommandList(c, "", []string{"public", "secret"})
	out := html.Render(tree, html.NewContext())

	assert.Contains(t, out, "Public command")
	assert.NotContains(t, out, "Secret command")
}

func TestSubcommandList_Bad_UnregisteredPathSkipped(t *testing.T) {
	c := core.New()
	tree := SubcommandList(c, "", []string{"ghost"})
	assert.Equal(t, "<ul></ul>", html.Render(tree, html.NewContext()))
}

func TestSubcommandList_Ugly_EmptyPaths(t *testing.T) {
	c := core.New()
	tree := SubcommandList(c, "", nil)
	assert.Equal(t, "<ul></ul>", html.Render(tree, html.NewContext()))
}

func TestSubcommandList_Ugly_NilCoreNeverDereferencedWhenPathsEmpty(t *testing.T) {
	var c *core.Core
	tree := SubcommandList(c, "", nil)
	assert.Equal(t, "<ul></ul>", html.Render(tree, html.NewContext()))
}

func TestSubcommandList_RealCommandsCall(t *testing.T) {
	c := core.New()
	require.True(t, c.Command("build", core.Command{Description: "Build the project"}).OK)
	require.True(t, c.Command("test", core.Command{Description: "Run tests"}).OK)

	tree := SubcommandList(c, "", c.Commands())
	out := html.Render(tree, html.NewContext())
	assert.Contains(t, out, "Build the project")
	assert.Contains(t, out, "Run tests")
}

func TestDirectChildPaths_Good(t *core.T) {
	got := directChildPaths("deploy", []string{"deploy/to", "deploy/to/homelab", "other"})
	core.AssertEqual(t, []string{"deploy/to"}, got)
}

func TestDirectChildPaths_Bad(t *core.T) {
	got := directChildPaths("deploy", []string{"other", "deployX"})
	core.AssertEqual(t, 0, len(got))
}

func TestDirectChildPaths_Ugly(t *core.T) {
	got := directChildPaths("", nil)
	core.AssertEqual(t, 0, len(got))
}
