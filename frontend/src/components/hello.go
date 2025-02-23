package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type Hello struct {
	app.Compo

	wailsApp app.Value
	name     string
	greet    string
}

func (h *Hello) Render() app.UI {
	return app.Div().ID("App").Body(
		app.Img().Src("assets/images/logo-universal.png").Alt("logo").ID("logo"),
		app.Div().ID("result").Class("result").Body(
			app.If(h.name != "", func() app.UI { return app.Text(h.greet) }),
		),
		app.Div().ID("input").Class("input-box").Body(
			app.Input().ID("name").Class("input text-black").AutoComplete(false).
				Name("input").Type("text").OnChange(h.ValueTo(&h.name)),
			app.Button().Class("btn").OnClick(h.Greet).Text("Greet"),
		),
	)
}

func (h *Hello) OnMount(ctx app.Context) {
	h.wailsApp = app.Window().Get("go").Get("main").Get("App")
}

func (h *Hello) Greet(ctx app.Context, e app.Event) {
	res := h.wailsApp.Call("Greet", h.name)
	res.Then(func(v app.Value) {
		h.greet = v.String()
		ctx.Update()
	})
}
