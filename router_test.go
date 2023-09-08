package morphadon_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/mrparano1d/morphadon"
)

type testPage struct {
	morphadon.Component
}

func NewTestPage() *testPage {
	return &testPage{
		Component: morphadon.NewDefaultComponent(),
	}
}

func TestNewRouter(t *testing.T) {
	r := morphadon.NewRouter()

	r.AddRoute(
		morphadon.NewRoute(
			"home",
			"/",
			NewTestPage(),
			morphadon.NewRoute(
				"about",
				"/about",
				NewTestPage(),
			),
		),
	)
	r.AddRoute(
		morphadon.NewRoute(
			"blog",
			"/blog",
			NewTestPage(),
			morphadon.NewRoute(
				"post",
				"/{id}",
				NewTestPage(),
			),
		),
	)

	if len(r.Routes()) != 4 {
		t.Errorf("Expected 4 routes, got %d", len(r.Routes()))
	}
}

func TestCurrentRoute(t *testing.T) {
	r := morphadon.NewRouter()

	r.AddRoute(
		morphadon.NewRoute(
			"home",
			"/",
			NewTestPage(),
			morphadon.NewRoute(
				"about",
				"/about",
				NewTestPage(),
			),
		),
	)
	r.AddRoute(
		morphadon.NewRoute(
			"blog",
			"/blog",
			NewTestPage(),
			morphadon.NewRoute(
				"post",
				"/{id}",
				NewTestPage(),
			),
		),
	)

	ctx := morphadon.NewContext()
	ctx.Req = &http.Request{
		URL: &url.URL{
			Path: "/blog/1",
		},
	}

	currentRoute := r.CurrentRoute(ctx)
	if currentRoute == nil {
		t.Errorf("Expected current route to be non-nil")
	} else if currentRoute.Name() != "post" {
		t.Errorf("Expected current route to be post, got %s", currentRoute.Name())
	}

	if currentRoute.Parent() == nil {
		t.Errorf("Expected parent route to be non-nil")
	} else if currentRoute.Parent().Name() != "blog" {
		t.Errorf("Expected parent route to be blog, got %s", currentRoute.Parent().Name())
	}
}
