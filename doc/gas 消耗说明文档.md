### gas 消耗说明文档

bsn 版本的dbchain由于会使用gas 扣取token，所以对gas消耗机制做了一些调整

主要的调整有两点：

​	1、普通的交易去掉了读取数据消耗gas的机制，因为插入数据的时候，会因过滤器，触发器等脚本对数据库实行读取操作，当表的数据太多后，会导致消耗gas激增，所以取消了读取消耗gas机制

```go
//cosmos-sdk/store/types/gas.go
func KVGasConfig() GasConfig {
	return GasConfig{
		HasCost:          0, //置0
		DeleteCost:       1000,
		ReadCostFlat:     0,//置0
		ReadCostPerByte:  0,//置0
		WriteCostFlat:    2000,
		WriteCostPerByte: 30,
		IterNextCostFlat: 0,//置0
	}
}
```

​	2、当我们把数据添加了索引，每次插入数据都会更新索引，更新索引也会消耗大量的gas，所以去掉了索引相关操作消耗积分的机制

```go
//dbchain/x/dbchain/internal/keeper/store.go
func DbChainStoreWithOutGas(ctx sdk.Context,storeKey sdk.StoreKey) *SafeStore{
	gasConfig := sdk.GasConfig {
		HasCost:          0,
		DeleteCost:       0,
		ReadCostFlat:     0,
		ReadCostPerByte:  0,
		WriteCostFlat:    0,
		WriteCostPerByte: 0,
		IterNextCostFlat: 0,
	}
	rawStore := gaskv.NewStore(ctx.MultiStore().GetKVStore(storeKey), ctx.GasMeter(), gasConfig)
	return NewSafeStore(rawStore)
}
```

