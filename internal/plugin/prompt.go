package plugin

import (
	"context"
	"seall/internal/extension"
	"seall/internal/extension_repo/prompt"
	"seall/internal/util/result"
	"strings"
)

func newPromptCache() *result.Cache[string, bool] {
	return result.NewCache[string, bool]()
}

func promptKey(parts ...string) string {
	ret := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		ret = append(ret, part)
	}

	return strings.Join(ret, ":")
}

func (a *AppContextImpl) ask(ext *extension.Extension, opts prompt.Options) error {
	manager, ok := a.promptManager.Get()
	if !ok || manager == nil {
		return prompt.ErrUnavailable
	}
	if opts.AllowLabel == "" {
		opts.AllowLabel = "Allow"
	}
	if opts.DenyLabel == "" {
		opts.DenyLabel = "Don't Allow"
	}

	return manager.Ask(context.Background(), ext, opts)
}
