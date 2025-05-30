// Copyright 2018 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package templates contains functions for template processing of Resource objects.
package templates

import (
	"context"
	"fmt"

	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/internal"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/tpl/tplimpl"
)

// Client contains methods to perform template processing of Resource objects.
type Client struct {
	rs *resources.Spec
	t  tplimpl.TemplateStoreProvider
}

// New creates a new Client with the given specification.
func New(rs *resources.Spec, t tplimpl.TemplateStoreProvider) *Client {
	if rs == nil {
		panic("must provide a resource Spec")
	}
	if t == nil {
		panic("must provide a template provider")
	}
	return &Client{rs: rs, t: t}
}

type executeAsTemplateTransform struct {
	rs         *resources.Spec
	t          tplimpl.TemplateStoreProvider
	targetPath string
	data       any
}

func (t *executeAsTemplateTransform) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey("execute-as-template", t.targetPath)
}

func (t *executeAsTemplateTransform) Transform(ctx *resources.ResourceTransformationCtx) error {
	tplStr := helpers.ReaderToString(ctx.From)
	th := t.t.GetTemplateStore()
	ti, err := th.TextParse(ctx.InPath, tplStr)
	if err != nil {
		return fmt.Errorf("failed to parse Resource %q as Template:: %w", ctx.InPath, err)
	}
	ctx.OutPath = t.targetPath
	return th.ExecuteWithContext(ctx.Ctx, ti, ctx.To, t.data)
}

func (c *Client) ExecuteAsTemplate(ctx context.Context, res resources.ResourceTransformer, targetPath string, data any) (resource.Resource, error) {
	return res.TransformWithContext(ctx, &executeAsTemplateTransform{
		rs:         c.rs,
		targetPath: paths.ToSlashTrimLeading(targetPath),
		t:          c.t,
		data:       data,
	})
}
