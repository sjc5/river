package router

import (
	"fmt"
	"net/http"
	"testers/routes-lit/internal/platform"

	"github.com/go-chi/chi/v5"
	"github.com/sjc5/hwy"
)

var HwyInstance = hwy.Hwy{}

type strMap map[string]string

func init() {
	defaultHeadBlocks := []hwy.HeadBlock{{Title: "JACOB"}}
	dataFuncs := hwy.DataFunctionMap{
		"/dashboard/customers/$customer_id/orders": hwy.Loader[any](
			func(ctx hwy.LoaderCtx[any]) {
				ctx.Res.Data = map[string]string{
					"message": "</script><script>alert('Hello, Bob!')</script>",
				}
				// res.Redirect("/login", http.StatusFound)
				ctx.Res.Headers.Set("bob3", "bob3")
				ctx.Res.Cookies = append(ctx.Res.Cookies, &http.Cookie{Name: "bob3", Value: "bob3"})
				ctx.Res.HeadBlocks = []*hwy.HeadBlock{
					{
						Tag:        "meta",
						Attributes: map[string]string{"name": "description", "content": "parent"},
					},
					{Title: "JACOB1"},
				}
			},
		),
		"/dashboard/customers/$customer_id/orders/$order_id": hwy.Loader[strMap](
			func(ctx hwy.LoaderCtx[strMap]) {
				ctx.Res.Data = strMap{"message": "kjbkjbkjbkjbkjbk"}
				ctx.Res.Headers.Set("bob3", "bob4")
				ctx.Res.Cookies = append(ctx.Res.Cookies, &http.Cookie{Name: "bob3", Value: "bob4"})
				ctx.Res.HeadBlocks = []*hwy.HeadBlock{
					{
						Tag:        "meta",
						Attributes: strMap{"name": "description", "content": "child"},
					},
					{Title: "JACOB2"},
				}
			},
		),
	}

	privateFS, err := platform.Kiruna.GetPrivateFS()
	if err != nil {
		panic(fmt.Sprintf("Error loading private FS: %v", err))
	}

	HwyInstance = hwy.Hwy{
		DefaultHeadBlocks:    defaultHeadBlocks,
		FS:                   privateFS,
		RootTemplateLocation: "templates/index.go.html",
		RootTemplateData: map[string]any{
			"Kiruna":         platform.Kiruna,
			"ClientEntryURL": platform.Kiruna.GetPublicURL("hwy_client_entry.js"),
		},
		Loaders: dataFuncs,
	}
	err = HwyInstance.Initialize()
	if err != nil {
		fmt.Println(err)
		panic("Error initializing Hwy")
	} else {
		fmt.Println("Hwy initialized")
	}
}

func Init() *chi.Mux {
	r := chi.NewRouter()
	r.Handle("/public/*", platform.Kiruna.GetServeStaticHandler("/public/", true))
	r.Handle("/*", HwyInstance.GetRootHandler())
	return r
}
