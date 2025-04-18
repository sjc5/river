package framework

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/river-now/river/kit/htmlutil"
	"github.com/river-now/river/kit/mux"
)

type SSRInnerHTMLInput struct {
	RiverSymbolStr      string
	IsDev               bool
	BuildID             string
	ViteDevURL          string
	LoadersData         []any
	ImportURLs          []string
	ExportKeys          []string
	OutermostErrorIndex int
	SplatValues         SplatValues
	Params              mux.Params
	HasRootData         bool
	Deps                []string
	CSSBundles          []string
	PublicPathPrefix    string
}

// Sadly, must include the script tags so html/template parses this correctly.
// They are stripped off later in order to get the correct sha256 hash.
// Then they are added back via htmlutil.RenderElement.
const ssrInnerHTMLTmplStr = `<script>
	globalThis[Symbol.for("{{.RiverSymbolStr}}")] = {};
	const x = globalThis[Symbol.for("{{.RiverSymbolStr}}")];
	x.isDev = {{.IsDev}};
	x.viteDevURL = {{.ViteDevURL}};
	x.publicPathPrefix = "{{.PublicPathPrefix}}";
	x.buildID = {{.BuildID}};
	x.loadersData = {{.LoadersData}};
	x.importURLs = {{.ImportURLs}};
	x.exportKeys = {{.ExportKeys}};
	x.outermostErrorIndex = {{.OutermostErrorIndex}};
	x.splatValues = {{.SplatValues}};
	x.params = {{.Params}};
	x.hasRootData = {{.HasRootData}};
	if (!x.isDev) {
		const deps = {{.Deps}};
		deps.forEach((y) => {
			const link = document.createElement("link");
			link.rel = "modulepreload";
			link.href = x.publicPathPrefix + y;
			document.head.appendChild(link);
		});
		const cssBundles = {{.CSSBundles}};
		cssBundles.forEach((y) => {
			const link = document.createElement("link");
			link.rel = "stylesheet";
			link.href = x.publicPathPrefix + y;
			link.setAttribute("data-river-css-bundle", y);
			document.head.appendChild(link);
		});
	}
</script>`

var ssrInnerTmpl = template.Must(template.New("ssr").Parse(ssrInnerHTMLTmplStr))

type GetSSRInnerHTMLOutput struct {
	Script     *template.HTML
	Sha256Hash string
}

func (h *River) GetSSRInnerHTML(routeData *UIRouteOutput) (*GetSSRInnerHTMLOutput, error) {
	var htmlBuilder strings.Builder

	dto := SSRInnerHTMLInput{
		RiverSymbolStr:      RiverSymbolStr,
		IsDev:               h._isDev,
		BuildID:             h._buildID,
		ViteDevURL:          routeData.ViteDevURL,
		LoadersData:         routeData.LoadersData,
		ImportURLs:          routeData.ImportURLs,
		ExportKeys:          routeData.ExportKeys,
		OutermostErrorIndex: routeData.OutermostErrorIndex,
		SplatValues:         routeData.SplatValues,
		Params:              routeData.Params,
		HasRootData:         routeData.HasRootData,
		Deps:                routeData.Deps,
		CSSBundles:          routeData.CSSBundles,
		PublicPathPrefix:    h.Wave.GetPublicPathPrefix(),
	}
	if err := ssrInnerTmpl.Execute(&htmlBuilder, dto); err != nil {
		wrapped := fmt.Errorf("could not execute SSR inner HTML template: %w", err)
		Log.Error(wrapped.Error())
		return nil, wrapped
	}

	innerHTML := htmlBuilder.String()
	innerHTML = strings.TrimPrefix(innerHTML, "<script>")
	innerHTML = strings.TrimSuffix(innerHTML, "</script>")

	el := htmlutil.Element{
		Tag:       "script",
		InnerHTML: template.HTML(innerHTML),
	}

	sha256Hash, err := htmlutil.AddSha256HashInline(&el, true)
	if err != nil {
		wrapped := fmt.Errorf("could not handle CSP for SSR inner HTML: %w", err)
		Log.Error(wrapped.Error())
		return nil, wrapped
	}

	renderedEl, err := htmlutil.RenderElement(&el)
	if err != nil {
		wrapped := fmt.Errorf("could not render SSR inner HTML: %w", err)
		Log.Error(wrapped.Error())
		return nil, wrapped
	}

	return &GetSSRInnerHTMLOutput{Script: &renderedEl, Sha256Hash: sha256Hash}, nil
}
