package db

import (
	"context"
	"time"

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

		// 设置过期时间，单位秒
		o.Set("expire", func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			seconds := call.Argument(1).ToInteger()
			ctx := context.Background()
			bc := rdb.Expire(ctx, key, time.Duration(seconds))
			if bc.Err() != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}

			b, err := bc.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}
			return lib.MakeReturnValue(runtime, b)
		})

		// 设置过期时间，时间戳
		o.Set("expireAt", func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			unixTimestamp := call.Argument(1).ToInteger()
			ctx := context.Background()
			bc := rdb.ExpireAt(ctx, key, time.Unix(unixTimestamp, 0))
			if bc.Err() != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}

			b, err := bc.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}
			return lib.MakeReturnValue(runtime, b)
		})

		// 设置过期时间，单位毫秒
		o.Set("pexpire", func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			milliseconds := call.Argument(1).ToInteger()
			ctx := context.Background()
			bc := rdb.PExpire(ctx, key, time.Duration(milliseconds))
			if bc.Err() != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}

			b, err := bc.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}
			return lib.MakeReturnValue(runtime, b)
		})

		// 设置过期时间，毫秒级时间戳
		o.Set("pexpireAt", func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			milliseconds := call.Argument(1).ToInteger()
			ctx := context.Background()
			bc := rdb.PExpireAt(ctx, key, time.UnixMilli(milliseconds))
			if bc.Err() != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}

			b, err := bc.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}
			return lib.MakeReturnValue(runtime, b)
		})

		// 查找所有符合给定模式 pattern 的 key
		o.Set("keys", func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			ctx := context.Background()
			bc := rdb.Keys(ctx, key)
			if bc.Err() != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}

			b, err := bc.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}
			return lib.MakeReturnValue(runtime, b)
		})

		// 将当前数据库的 key 移动到给定的数据库 db 当中
		o.Set("move", func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			dbName := call.Argument(1).ToInteger()
			ctx := context.Background()
			bc := rdb.Move(ctx, key, int(dbName))
			if bc.Err() != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}

			b, err := bc.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}
			return lib.MakeReturnValue(runtime, b)
		})

		// 移除给定 key 的过期时间，使得 key 永不过期
		o.Set("persist", func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			ctx := context.Background()
			bc := rdb.Persist(ctx, key)
			if bc.Err() != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}

			b, err := bc.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}
			return lib.MakeReturnValue(runtime, b)
		})

		// 以毫秒为单位返回 key 的剩余过期时间
		o.Set("pttl", func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			ctx := context.Background()
			bc := rdb.PTTL(ctx, key)
			if bc.Err() != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}

			b, err := bc.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}
			return lib.MakeReturnValue(runtime, b)
		})

		// 以秒为单位返回 key 的剩余过期时间
		o.Set("ttl", func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			ctx := context.Background()
			bc := rdb.TTL(ctx, key)
			if bc.Err() != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}

			b, err := bc.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}
			return lib.MakeReturnValue(runtime, b)
		})

		// 从当前数据库中随机返回一个 key
		o.Set("randomKey", func(call goja.FunctionCall) goja.Value {
			ctx := context.Background()
			bc := rdb.RandomKey(ctx)
			if bc.Err() != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}

			b, err := bc.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}
			return lib.MakeReturnValue(runtime, b)
		})

		// 修改 key 的名称
		o.Set("rename", func(call goja.FunctionCall) goja.Value {
			oldName := call.Argument(0).String()
			newName := call.Argument(1).String()
			ctx := context.Background()
			bc := rdb.Rename(ctx, oldName, newName)
			if bc.Err() != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}

			b, err := bc.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}
			return lib.MakeReturnValue(runtime, b)
		})

		// 在新的 key 不存在时修改 key 的名称
		o.Set("renameNX", func(call goja.FunctionCall) goja.Value {
			oldName := call.Argument(0).String()
			newName := call.Argument(1).String()
			ctx := context.Background()
			bc := rdb.RenameNX(ctx, oldName, newName)
			if bc.Err() != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}

			b, err := bc.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}
			return lib.MakeReturnValue(runtime, b)
		})

		// 迭代数据库中的数据库键
		o.Set("scan", func(call goja.FunctionCall) goja.Value {
			cursor := call.Argument(0).ToInteger()
			match := call.Argument(1).String()
			count := call.Argument(2).ToInteger()
			ctx := context.Background()
			bc := rdb.Scan(ctx, uint64(cursor), match, count)
			if bc.Err() != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}

			keys, c, err := bc.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}

			// 返回数据
			value := struct {
				Keys   []string `json:"keys"`
				Cursor uint64   `json:"cursor"`
			}{
				Keys:   keys,
				Cursor: c,
			}

			return lib.MakeReturnValue(runtime, value)
		})

		// 返回 key 所储存的值的类型
		o.Set("type", func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			ctx := context.Background()
			bc := rdb.Type(ctx, key)
			if bc.Err() != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}

			b, err := bc.Result()
			if err != nil {
				return lib.MakeErrorValue(runtime, bc.Err())
			}
			return lib.MakeReturnValue(runtime, b)
		})
	})
}
