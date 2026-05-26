// registry service
package registry

import "sync"

type Registry struct {
	mu sync.RWMutex // thread safety , multiple workers may deploy or work simultan. wihtout race condition

	functions map[string]Function
}

func NewRegistry() *Registry {

	return &Registry{
		functions: make(
			map[string]Function,
		),
	}
}

func (r *Registry) Register(
	fn Function,
) {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.functions[fn.Name] = fn
}

func (r *Registry) Get(
	name string,
) (Function, bool) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	fn, exists := r.functions[name]

	return fn, exists
}

func (r *Registry) List() []Function {

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(
		[]Function,
		0,
		len(r.functions),
	)

	for _, fn := range r.functions {
		result = append(result, fn)
	}

	return result
}

func (r *Registry) IncrementInvocation(
	name string,
) {

	r.mu.Lock()
	defer r.mu.Unlock()

	fn, exists := r.functions[name]

	if !exists {
		return
	}

	fn.InvocationCount++

	r.functions[name] = fn
}
