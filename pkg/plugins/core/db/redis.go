package db

import (
	"context"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/issueye/lichee/pkg/plugins/core/lib"
	redis "github.com/redis/go-redis/v9"
)

// RegisterRedis
// 由外部传入redis 客户端
func RegisterRedis(moduleName string, rdb *redis.Client) {
	require.RegisterNativeModule(moduleName, func(runtime *goja.Runtime, module *goja.Object) {
		o := module.Get("exports").(*goja.Object)
		// key

		// 删除键
		o.Set("del", func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			ctx := context.Background()
			ic := rdb.Del(ctx, key)
			if ic.Err() != nil {
				return lib.MakeErrorValue(runtime, ic.Err())
			}
			return nil
		})

		// 序列化数据
		o.Set("dump", func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			ctx := context.Background()
			sc := rdb.Dump(ctx, key)
			if sc.Err() != nil {
				return lib.MakeErrorValue(runtime, sc.Err())
			}

			// 获取数据
			s, err := sc.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}

			return lib.MakeReturnValue(runtime, s)
		})

		// 判断键是否存在
		o.Set("exists", func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			ctx := context.Background()
			ic := rdb.Exists(ctx, key)
			if ic.Err() != nil {
				return lib.MakeErrorValue(runtime, ic.Err())
			}

			i, err := ic.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, err)
			}

			return lib.MakeReturnValue(runtime, i)
		})
	})
}
