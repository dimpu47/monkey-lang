package object

// NewEnvironment constructs a new Environment object to hold bindings
// of identifiers to their names
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// Environment is an object that holds a mapping of names to bound objets
type Environment struct {
	store  map[string]Object
	parent *Environment
}

// Clone returns a new Environment with the parent set to the current
// environment (enclosing environment)
func (e *Environment) Clone() *Environment {
	env := NewEnvironment()
	env.parent = e
	return env
}

// Get returns the object bound by name
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.parent != nil {
		obj, ok = e.parent.Get(name)
	}
	return obj, ok
}

// Set stores the object with the given name
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
