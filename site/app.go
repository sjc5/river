package app

import (
	"embed"
	"net/http"
	"path"

	"github.com/river-now/river"
	"github.com/river-now/river/kit/colorlog"
	"github.com/river-now/river/kit/headels"
	"github.com/river-now/river/kit/htmlutil"
	"github.com/river-now/river/kit/theme"
	"github.com/river-now/river/wave"
)

const (
	Domain          = "river.now"
	Origin          = "https://" + Domain
	SiteTitle       = "river.now"
	SiteDescription = "River is a Go / TypeScript meta-framework with first-class support for React, Solid, and Preact – built on Vite."
)

var River = &river.River{
	Wave: Wave,

	GetHeadElUniqueRules: func() *headels.HeadEls {
		e := river.NewHeadEls(7)

		e.Meta(e.Property("og:title"))
		e.Meta(e.Property("og:description"))
		e.Meta(e.Property("og:type"))
		e.Meta(e.Property("og:image"))
		e.Meta(e.Property("og:url"))
		e.Meta(e.Name("twitter:card"))
		e.Meta(e.Name("twitter:site"))

		return e
	},

	GetDefaultHeadEls: func(r *http.Request) ([]*htmlutil.Element, error) {
		root := "https://river.now"
		currentURL := path.Join(root, r.URL.Path)
		imgURL := path.Join(root, Wave.GetPublicURL("river-banner.webp"))

		e := river.NewHeadEls()

		e.Title(SiteTitle)
		e.Description(SiteDescription)

		e.Meta(e.Property("og:title"), e.Content(SiteTitle))
		e.Meta(e.Property("og:description"), e.Content(SiteDescription))
		e.Meta(e.Property("og:type"), e.Content("website"))
		e.Meta(e.Property("og:image"), e.Content(imgURL))
		e.Meta(e.Property("og:url"), e.Content(currentURL))

		e.Meta(e.Name("twitter:card"), e.Content("summary_large_image"))
		e.Meta(e.Name("twitter:site"), e.Content("@riverframework"))

		e.Link(e.Attr("rel", "icon"), e.Attr("href", Wave.GetPublicURL("favicon.svg")))

		return e.Collect(), nil
	},

	GetRootTemplateData: func(r *http.Request) (map[string]any, error) {
		return map[string]any{
			"HTMLClass":                   theme.GetThemeData(r).HTMLClass,
			"SystemThemeScript":           theme.SystemThemeScript,
			"SystemThemeScriptSha256Hash": theme.SystemThemeScriptSha256Hash,
		}, nil
	},
}

//go:embed wave.config.json
var configBytes []byte

//go:embed all:__dist/static
var staticFS embed.FS

var Wave = wave.New(&wave.Config{
	ConfigBytes:            configBytes,
	StaticFS:               staticFS,
	StaticFSEmbedDirective: "all:__dist/static",
})

var Log = colorlog.New("app server")
