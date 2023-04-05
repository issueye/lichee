package compiler

import (
	_ "embed"
	"sync"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
)

//go:embed lib/babel.min.js
var babelSrc string

var (
	DefaultOpts = map[string]interface{}{
		// "presets": []string{"latest"},
		"plugins": []interface{}{
			// es2015 https://github.com/babel/babel/blob/v6.26.0/packages/babel-preset-es2015/src/index.js
			// in goja
			// []interface{}{"transform-es2015-template-literals", map[string]interface{}{"loose": false, "spec": false}},
			// "transform-es2015-literals", // in goja
			// "transform-es2015-function-name", // in goja
			// []interface{}{"transform-es2015-arrow-functions", map[string]interface{}{"spec": false}}, // in goja
			// "transform-es2015-block-scoped-functions", // in goja
			// []interface{}{"transform-es2015-classes", map[string]interface{}{"loose": false}}, // in goja
			// "transform-es2015-object-super", // in goja
			// "transform-es2015-shorthand-properties", // in goja
			// "transform-es2015-duplicate-keys", // in goja
			// []interface{}{"transform-es2015-computed-properties", map[string]interface{}{"loose": false}}, // in goja
			// "transform-es2015-for-of", // in goja
			// "transform-es2015-sticky-regex", // in goja
			// "transform-es2015-unicode-regex", // in goja
			// "check-es2015-constants", // in goja
			// []interface{}{"transform-es2015-spread", map[string]interface{}{"loose": false}}, // in goja
			// "transform-es2015-parameters", // in goja
			// []interface{}{"transform-es2015-destructuring", map[string]interface{}{"loose": false}}, // in goja
			// "transform-es2015-block-scoping", // in goja
			// "transform-es2015-typeof-symbol", // in goja
			// all the other module plugins are just dropped
			[]interface{}{"transform-es2015-modules-commonjs", map[string]interface{}{"loose": false}},
			// "transform-regenerator", // Doesn't really work unless regeneratorRuntime is also added

			// es2016 https://github.com/babel/babel/blob/v6.26.0/packages/babel-preset-es2016/src/index.js
			// "transform-exponentiation-operator",

			// es2017 https://github.com/babel/babel/blob/v6.26.0/packages/babel-preset-es2017/src/index.js
			// "syntax-trailing-function-commas", // in goja
			// "transform-async-to-generator", // Doesn't really work unless regeneratorRuntime is also added
		},
		"ast":           false,
		"sourceMaps":    false,
		"babelrc":       false,
		"compact":       false,
		"retainLines":   true,
		"highlightCode": false,
	}

	onceBabelCode      sync.Once
	globalBabelCode    *goja.Program
	errGlobalBabelCode error
	globalBabel        *babel
	onceBabel          sync.Once
)

type CompatibilityMode uint8

const (
	// CompatibilityModeExtended achieves ES6+ compatibility with Babel
	CompatibilityModeExtended CompatibilityMode = iota + 1
	// CompatibilityModeBase is standard goja ES5.1+
	CompatibilityModeBase
)

type Options struct {
	CompatibilityMode CompatibilityMode
	SourceMapLoader   func(string) ([]byte, error)
	Strict            bool
}

type babel struct {
	vm        *goja.Runtime
	this      goja.Value
	transform goja.Callable
	m         sync.Mutex
}

func newBabel() (*babel, error) {
	onceBabelCode.Do(func() {
		globalBabelCode, errGlobalBabelCode = goja.Compile("<lib/babel.min.js>", babelSrc, false)
	})
	if errGlobalBabelCode != nil {
		return nil, errGlobalBabelCode
	}
	vm := goja.New()
	_, err := vm.RunProgram(globalBabelCode)
	if err != nil {
		return nil, err
	}

	this := vm.Get("Babel")
	bObj := this.ToObject(vm)
	result := &babel{vm: vm, this: this}
	if err = vm.ExportTo(bObj.Get("transform"), &result.transform); err != nil {
		return nil, err
	}

	return result, err
}

func (b *babel) transformImpl(src, filename string) (string, error) {
	b.m.Lock()
	defer b.m.Unlock()
	opts := make(map[string]interface{})
	for k, v := range DefaultOpts {
		opts[k] = v
	}
	opts["filename"] = filename

	// startTime := time.Now()
	v, err := b.transform(b.this, b.vm.ToValue(src), b.vm.ToValue(opts))
	if err != nil {
		return "", err
	}

	vO := v.ToObject(b.vm)
	var code string
	if err = b.vm.ExportTo(vO.Get("code"), &code); err != nil {
		return code, err
	}

	return code, nil
}

type Compiler struct {
	babel   *babel
	Options Options
}

func New() *Compiler {
	return &Compiler{}
}

// initializeBabel
// 初始化 babel
func (c *Compiler) initializeBabel() error {
	var err error
	if c.babel == nil {
		c.babel, err = newBabel()
	}
	return err
}

func (c *Compiler) Transform(src, filename string) (code string, err error) {
	if c.babel == nil {
		onceBabel.Do(func() {
			globalBabel, err = newBabel()
		})
		c.babel = globalBabel
	}
	if err != nil {
		return
	}

	code, err = c.babel.transformImpl(src, filename)
	return
}

func (c *Compiler) Compile(src, filename string, isESM bool) (*goja.Program, string, error) {
	return c.compileImpl(src, filename, !isESM, c.Options.CompatibilityMode)
}

func (c *Compiler) compileImpl(src, filename string, wrap bool, compatibilityMode CompatibilityMode) (*goja.Program, string, error) {
	code := src
	if wrap { // the lines in the sourcemap (if available) will be fixed by increaseMappingsByOne
		code = "(function(module, exports){\n" + code + "\n})\n"
	}
	opts := parser.WithDisableSourceMaps
	ast, err := parser.ParseFile(nil, filename, code, 0, opts)

	if err != nil {
		if compatibilityMode == CompatibilityModeExtended {
			code, err = c.Transform(src, filename)
			if err != nil {
				return nil, code, err
			}
			// the compatibility mode "decreases" here as we shouldn't transform twice
			return c.compileImpl(code, filename, wrap, CompatibilityModeBase)
		}
		return nil, code, err
	}
	pgm, err := goja.CompileAST(ast, c.Options.Strict)
	return pgm, code, err
}

// Pool is a pool of compilers so it can be used easier in parallel tests as they have their own babel.
type Pool struct {
	c chan *Compiler
}

// NewPool creates a Pool that will be using the provided logger and will preallocate (in parallel)
// the count of compilers each with their own babel.
func NewPool(count int) *Pool {
	c := &Pool{
		c: make(chan *Compiler, count),
	}
	go func() {
		for i := 0; i < count; i++ {
			go func() {
				co := New()
				err := co.initializeBabel()
				if err != nil {
					panic(err)
				}
				c.Put(co)
			}()
		}
	}()

	return c
}

// Get a compiler from the pool.
func (c *Pool) Get() *Compiler {
	return <-c.c
}

// Put a compiler back in the pool.
func (c *Pool) Put(co *Compiler) {
	c.c <- co
}
