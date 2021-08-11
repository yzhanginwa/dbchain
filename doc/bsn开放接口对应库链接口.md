# bsn开放接口对应库链接口



### 1. 密钥托管创建链账户

| 平台         | 地址                                    |
| ------------ | --------------------------------------- |
| bsn 接口     | /api/v1/account/apply(POST)             |
| dbchain 接口 | /dbchain/oracle/bsn/account/apply(POST) |

请求数据

| 序号 | 字段名 | 字段 | 类型 | 必填 | 备注 |
| ---- | ------ | ---- | ---- | ---- | ---- |
| 1    |        |      |      |      |      |

响应数据

| 序号 | 字段名   | 字段       | 类型   | 备注 |
| ---- | -------- | ---------- | ------ | ---- |
| 1    | 公钥     | publicKey  | String |      |
| 2    | 私钥     | privateKey | String |      |
| 3    | 账户地址 | address    | String |      |
| 4    | 助记词   | mnemonic   | String |      |

相应示例数据 ：

```json
{
    "address": "cosmos1gu6yn0gx3eutq0zh8nwzjslr6pusuc9tnf0yc5",
    "mnemonic": "curious offer wine reform series spare merit meadow pony embark quiz dilemma",
    "privateKey": "e1b0f79b20eaf586375171a25e65555f9e9a9534e28a63d4ba5956161851647b24b045da39",
    "publicKey": "eb5ae98721032dab88be43b4096e846cffd43ea0c3be36c37d424f7c2ced65ab206cc0da1339"
}
```



### 2. 公钥上传创建链账户

| 平台         | 地址                                              |
| ------------ | ------------------------------------------------- |
| bsn 接口     | /api/v1/account/apply/publicKey(POST)             |
| dbchain 接口 | /dbchain/oracle/bsn/account/apply/publicKey(POST) |

请求数据

| 序号 | 字段名 | 字段      | 类型   | 必填 | 备注 |
| ---- | ------ | --------- | ------ | ---- | ---- |
| 1    | 公钥   | publicKey | String | Y    |      |

响应数据

| 序号 | 字段名     | 字段      | 类型   | 备注 |
| ---- | ---------- | --------- | ------ | ---- |
| 1    | 公钥       | publicKey | String |      |
| 2    | 私账户地址 | address   | String |      |

请求示例数据

```json
{
    "publicKey" : "eb5ae98721023ff53cba69bdb83047be13a492bedbe12dcd17709e2564d4a0463f7c9be82b9a"
}
```

响应示例数据

```json
{
    "address": "cosmos1dx9h5yl0gxsc3snn9888cmm79xzqwa9lv9pnss",
    "publicKey": "eb5ae98721023ff53cba69bdb83047be13a492bedbe12dcd17709e2564d4a0463f7c9be82b9a"
}
```



### 3. 链账户充值

| 平台         | 地址                                       |
| ------------ | ------------------------------------------ |
| bsn 接口     | /api/v1/account/recharge(POST)             |
| dbchain 接口 | /dbchain/oracle/bsn/account/recharge(POST) |

请求数据

| 序号 | 字段名           | 字段               | 类型   | 必填 | 备注                                         |
| ---- | ---------------- | ------------------ | ------ | ---- | -------------------------------------------- |
| 1    | BSN 管理账户地址 | bsnAddress         | String | Y    | BSN 管理账户需要采用密钥托管模式，由节点代签 |
| 2    | 用户账户地址     | userAccountAddress | String | Y    |                                              |
| 3    | 充值 GAS         | rechargeGas        | String | Y    |                                              |

响应数据

| 序号 | 字段名    | 字段    | 类型   | 备注                      |
| ---- | --------- | ------- | ------ | ------------------------- |
| 1    | 交易 HASH | txHash  | String |                           |
| 2    | 状态      | state   | int    | 2:充值成功<br/>3:充值失败 |
| 3    | 备注      | remarks | String |                           |



请求示例数据

```json
{
		"bsnAddress":"cosmos1k5nt3qtpvjyfetma4krcwlnsnk9l87kvw75yyg",
        "userAccountAddress":"cosmos16fqf5vaf6cf224w8zddlzad2yd4zns9ua2sc28",
        "rechargeGas":"1dbctoken" //数量+币种
}
```

响应示例数据

```json
{
    "remarks": "",//成功时为空，出错时会返回相应的错误信息
    "state": 2,
    "txHash": "b37948cc1c28e27d3bc5617337e9af86c784684bb7f33b06240cde75437f1716"
}
```

### 4. 链账户能量值查询

| 平台         | 地址                                     |
| ------------ | ---------------------------------------- |
| bsn 接口     | /api/v1/account/gas(POST)                |
| dbchain 接口 | /bank/balances/{userAccountAddress}(GET) |

请求数据

| 序号 | 字段名 | 字段 | 类型 | 必填 | 备注 |
| ---- | ------ | ---- | ---- | ---- | ---- |
| 1    |        |      |      |      |      |

响应数据

| 序号 | 字段名 | 字段   | 类型   | 备注 |
| ---- | ------ | ------ | ------ | ---- |
| 1    | 名称   | denom  | String |      |
| 2    | 数量   | amount | String |      |

响应示例数据

```json
{
    "height": "138",
    "result": [
        {
            "denom": "dbctoken",
            "amount": "101"
        }
    ]
}
```

注： 这是底层标准接口，未做封装处理



### 5. 链账户交易查询

| 平台         | 地址                                 |
| ------------ | ------------------------------------ |
| bsn 接口     | /api/v1/account/tx(POST)             |
| dbchain 接口 | /dbchain/oracle/bsn/account/tx(POST) |

请求数据

| 序号 | 字段名       | 字段               | 类型   | 必填 | 备注                |
| ---- | ------------ | ------------------ | ------ | ---- | ------------------- |
| 1    | 用户账户地址 | userAccountAddress | String | Y    |                     |
| 2    | 起始时间     | startDate          | String | N    | yyyy-MM-dd HH:mm:ss |
| 3    | 结束时间     | endDate            | String | N    | yyyy-MM-dd HH:mm:ss |
| 4    | 起始块高     | startBlockHeight   | String | N    |                     |
| 5    | 结束块高     | endBlockHeight     | String | N    |                     |

响应数据

| 序号 | 字段名       | 字段               | 类型   | 备注                      |
| ---- | ------------ | ------------------ | ------ | ------------------------- |
| 1    | 用户账户地址 | userAccountAddress | String |                           |
| 2    | 交易 HASH    | txHash             | String |                           |
| 3    | 交易时间     | txTime             | String | yyyy-MM-dd HH:mm:ss       |
| 4    | 交易能量值   | gas                | String |                           |
| 5    | 交易状态     | state              | String | 2:处理成功<br/>3:处理失败 |
| 6    | 块高         | blockHeight        | int    |                           |



请求示例数据(按时间查询)

```json
{
    "userAccountAddress" : "cosmos1k5nt3qtpvjyfetma4krcwlnsnk9l87kvw75yyg",
    "startDate" : "2021-08-02 14:00:00",
    "endDate" : "2021-08-02 15:00:00"
}
```

响应示例数据

```json
[
    {
        "blockHeight": 147,
        "gas": "7dbctoken",
        "state": 2,
        "txHash": "a2a6c0b9f8c038ad52419476742f70092b29db5bc949b1f26fc02869c1a3f6c8",
        "txTime": "2021-08-02 14:47:33",
        "userAccountAddress": "cosmos1k5nt3qtpvjyfetma4krcwlnsnk9l87kvw75yyg"
    },
    {
        "blockHeight": 127,
        "gas": "",
        "state": 2,
        "txHash": "b37948cc1c28e27d3bc5617337e9af86c784684bb7f33b06240cde75437f1716",
        "txTime": "2021-08-02 14:15:58",
        "userAccountAddress": "cosmos1k5nt3qtpvjyfetma4krcwlnsnk9l87kvw75yyg"
    },
    {
        "blockHeight": 127,
        "gas": "", //为空表示未消耗手续费
        "state": 2,
        "txHash": "5a3171786291f535e3a812911a39af3f8c7007ab6c13616e093ab9107cadcfd0",
        "txTime": "2021-08-02 14:14:44",
        "userAccountAddress": "cosmos1k5nt3qtpvjyfetma4krcwlnsnk9l87kvw75yyg"
    }
]
```

如果想要按块高查询，请求数据如下(两个条件都写，优先时间)：

```json
{
    "userAccountAddress" : "cosmos1j8tdge2fev0z7g46rp5tkdgcjq3y652k6eflaa",
    "startBlockHeight" : "1",
    "endBlockHeight" : "150"
}
```



### 6. 查询有转账权限的管理员

| 平台         | 地址                                      |
| ------------ | ----------------------------------------- |
| dbchain 接口 | /dbchain/token_keepers/{accessToken}(GET) |

请求数据

| 序号 | 字段名 | 字段 | 类型 | 必填 | 备注 |
| ---- | ------ | ---- | ---- | ---- | ---- |
|      |        |      |      |      |      |

响应数据

| 序号 | 字段名 | 字段   | 类型 | 备注 |
| ---- | ------ | ------ | ---- | ---- |
|      | 结果   | result | 集合 |      |

响应示例数据

```json
{
    "height": "0",
    "result": [ //管理员集合
        "cosmos1dgk6nmqfm3yz5k00f075h5u4gt2e5jtwg9zc3f",
        "cosmos1dsq9pnsn59fgzucxf33y99hehp5z9g6ljydnyg"
    ]
}
```

注 ：只有管理员才有查询权限，普通用户查询结果为空；js-sdk有生成accessToken的API



### 7.查询当前是否限制P2P转账

| 平台         | 地址                                                  |
| ------------ | ----------------------------------------------------- |
| dbchain 接口 | /dbchain/limit_p2p_transfer_status/{accessToken}(GET) |

请求数据

| 序号 | 字段名 | 字段 | 类型 | 必填 | 备注 |
| ---- | ------ | ---- | ---- | ---- | ---- |
|      |        |      |      |      |      |

响应数据

| 序号 | 字段名 | 字段   | 类型 | 备注 |
| ---- | ------ | ------ | ---- | ---- |
|      | 结果   | result | bool |      |

响应示例数据

```json
{
    "height": "0",
    "result": true   //true 或者 false
}
```



### 8.查询当前链最小gas费

| 平台         | 地址                                       |
| ------------ | ------------------------------------------ |
| dbchain 接口 | /dbchain/min_gas_prices/{accessToken}(GET) |

请求数据

| 序号 | 字段名 | 字段 | 类型 | 必填 | 备注 |
| ---- | ------ | ---- | ---- | ---- | ---- |
|      |        |      |      |      |      |

响应数据

| 序号 | 字段名 | 字段   | 类型 | 备注 |
| ---- | ------ | ------ | ---- | ---- |
|      | 结果   | result | 数组 |      |

响应示例数据

```json
{
    "height": "0",
    "result": [
        {
            "denom": "dbctoken",
            "amount": "0.000001000000000000"
        }
    ]
}
```

