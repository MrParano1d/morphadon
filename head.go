package engine

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

type HeadAttributes map[string]string

type Head struct {
	Meta      []HeadAttributes
	Title     string
	Link      []HeadAttributes
	Base      HeadAttributes
	Scripts   []HeadAttributes
	HtmlAttrs HeadAttributes
	BodyAttrs HeadAttributes
}

func UseHead() *Head {
	return &Head{
		Meta:      []HeadAttributes{},
		Link:      []HeadAttributes{},
		Scripts:   []HeadAttributes{},
		Base:      HeadAttributes{},
		HtmlAttrs: HeadAttributes{},
		BodyAttrs: HeadAttributes{},
		Title:     "",
	}
}

func (h *Head) baseNode() g.Node {
	var attributes []g.Node

	for k, v := range h.Base {
		attributes = append(attributes, g.Attr(k, v))
	}

	return html.Base(g.Group(attributes))
}

func (h *Head) bodyAttributes() g.Node {
	var attributes []g.Node

	for k, v := range h.BodyAttrs {
		attributes = append(attributes, g.Attr(k, v))
	}

	return g.Group(attributes)
}

func mergeHeads(newHead *Head, currentHead *Head) *Head {
	currentHead.Meta = append(currentHead.Meta, newHead.Meta...)
	currentHead.Link = append(currentHead.Link, newHead.Link...)
	currentHead.Scripts = append(currentHead.Scripts, newHead.Scripts...)

	for k, v := range newHead.Base {
		currentHead.Base[k] = v
	}

	for k, v := range newHead.HtmlAttrs {
		currentHead.HtmlAttrs[k] = v
	}

	for k, v := range newHead.BodyAttrs {
		currentHead.BodyAttrs[k] = v
	}

	if newHead.Title != "" {
		currentHead.Title = newHead.Title
	}

	return currentHead
}
