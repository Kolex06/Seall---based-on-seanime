package plugin

import (
	"errors"
	"seall/internal/extension"
	"seall/internal/extension_repo/prompt"
	"seall/internal/goja/goja_bindings"
	gojautil "seall/internal/util/goja"

	"github.com/dop251/goja"
	"github.com/rs/zerolog"
)

func (a *AppContextImpl) bindAuthObj(vm *goja.Runtime, ext *extension.Extension, scheduler *gojautil.Scheduler) *goja.Object {
	authObj := vm.NewObject()
	cache := newPromptCache()

	_ = authObj.Set("login", func(token string) goja.Value {
		return a.authAction(vm, scheduler, ext, prompt.Options{
			Kind:     "auth",
			Action:   "log in to SIMKL",
			Resource: "SIMKL account",
			Message:  "Allow \"" + ext.Name + "\" to log in to SIMKL?",
			Cache:    cache,
			CacheKey: promptKey("auth", "login", "simkl"),
		}, func() error {
			if token == "" {
				return errors.New("token is empty")
			}
			if a.auth.Login == nil {
				return errors.New("auth login not set")
			}
			return a.auth.Login(token)
		})
	})

	_ = authObj.Set("logout", func() goja.Value {
		return a.authAction(vm, scheduler, ext, prompt.Options{
			Kind:     "auth",
			Action:   "log out of SIMKL",
			Resource: "SIMKL account",
			Message:  "Allow \"" + ext.Name + "\" to log out of SIMKL?",
			Cache:    cache,
			CacheKey: promptKey("auth", "logout", "simkl"),
		}, func() error {
			if a.auth.Logout == nil {
				return errors.New("auth logout not set")
			}
			return a.auth.Logout()
		})
	})

	return authObj
}

func (a *AppContextImpl) BindAuthToContextObj(vm *goja.Runtime, obj *goja.Object, logger *zerolog.Logger, ext *extension.Extension, scheduler *gojautil.Scheduler) {
	_ = logger
	_ = obj.Set("auth", a.bindAuthObj(vm, ext, scheduler))
}

func (a *AppContextImpl) authAction(vm *goja.Runtime, scheduler *gojautil.Scheduler, ext *extension.Extension, opts prompt.Options, run func() error) goja.Value {
	promise, resolve, reject := vm.NewPromise()

	go func() {
		err := a.ask(ext, opts)
		if err == nil {
			err = run()
		}

		scheduler.ScheduleAsync(func() error {
			if err != nil {
				reject(goja_bindings.NewErrorString(vm, err.Error()))
				return nil
			}
			resolve(vm.ToValue(true))
			return nil
		})
	}()

	return vm.ToValue(promise)
}
