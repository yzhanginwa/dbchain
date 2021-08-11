## dbchain 部署及初始化文档

### 一、dbchain部署

dbchain采用的BFT 共识算法，部署节点应当最少4个。下面以4个为例介绍部署文档

现在假设有4台服务器，分别为node1，node2，node3，node4

#### 1、在每个节点上初始化,后续操作都是在容器进行的

node1

```shell
~$ sudo docker run --net dbcnet --ip node1_ip --user 1000:1000 -it -v /dbcpool/chain_data/docker_store:/home/dbchain dbchain bash
	In the container's bash session:
	#初始化
	~$ dbchaind init node1 --chain-id testnet
	#创世用户
	~$ dbchaincli keys add node1
	
```

node2

```shell
~$ sudo docker run --net dbcnet --ip node2_ip --user 1000:1000 -it -v /dbcpool/chain_data/docker_store:/home/dbchain dbchain bash
	In the container's bash session:
	$ dbchaind init node2 --chain-id testnet
	$ dbchaincli keys add node2
	
```

node3

```shell
~$ sudo docker run --net dbcnet --ip node3_ip --user 1000:1000 -it -v /dbcpool/chain_data/docker_store:/home/dbchain dbchain bash
	In the container's bash session:
	$ dbchaind init node3 --chain-id testnet
	$ dbchaincli keys add node3
	
```

node4

```shell
~$ sudo docker run --net dbcnet --ip node4_ip --user 1000:1000 -it -v /dbcpool/chain_data/docker_store:/home/dbchain dbchain bash
	In the container's bash session:
	$ dbchaind init node4 --chain-id testnet
	$ dbchaincli keys add node4
	
```

#### 2、**add-genesis-account将创世用户添加到创始文件**

node1

```shell
#stake 用于验证节点staking，dbctoken可用于充值gas，数量可以自定义
~$ dbchaind add-genesis-account node1Address   10000000000dbctoken,10000000000stake  
~$ dbchaind add-genesis-account node2Address   10000000000dbctoken,10000000000stake
~$ dbchaind add-genesis-account node3Address   10000000000dbctoken,10000000000stake
~$ dbchaind add-genesis-account node4Address   10000000000dbctoken,10000000000stake

~$ dbchaind add-genesis-admin-account node1Address
~$ dbchaind add-genesis-admin-account node2Address 
~$ dbchaind add-genesis-admin-account node3Address
~$ dbchaind add-genesis-admin-account node4Address
```

node2 

```shell
~$ dbchaind add-genesis-account node2Address   10000000000dbctoken,10000000000stake
~$ dbchaind add-genesis-admin-account node2Address 
```

node3

```shell
~$ dbchaind add-genesis-account node3Address   10000000000dbctoken,10000000000stake
~$ dbchaind add-genesis-admin-account node3Address 
```

node4

```shell
~$ dbchaind add-genesis-account node4Address   10000000000dbctoken,10000000000stake
~$ dbchaind add-genesis-admin-account node4Address 
```

（注：其中一个节点需要添加所有账号，其他节点只需要添加自己即可）

#### 3、**创建gentx**，**collect-gentxs收集创世交易到创世文件**

node1

```shell
#创建gentx
~$ dbchaind gentx --name node1
#收集创世交易
~$ dbchaind collect-gentxs
```

node2

```shell
~$ dbchaind gentx --name node2
~$ dbchaind collect-gentxs
```

node3

```shell
~$ dbchaind gentx --name node3
~$ dbchaind collect-gentxs
```

node4

```shell
~$ dbchaind gentx --name node4
~$ dbchaind collect-gentxs
```

#### 4、**拷贝创世交易**

将其他节点的创世文件里的创世交易手动拷贝到第一个节点的创世文件中的"gentxs": []里



##### 4.1 如node1的genesis.json(~/docker_store/.dbchaind/config/genesis.json) 内容如下

```json
{
    ...
    "genutil": {
      "gentxs": [
        {
          "type": "cosmos-sdk/StdTx",
          "value": {
            "msg": [
              {
                "type": "cosmos-sdk/MsgCreateValidator",
                "value": {
                  "description": {
                    "moniker": "node1",
                    "identity": "",
                    "website": "",
                    "security_contact": "",
                    "details": ""
                  },
                  "commission": {
                    "rate": "0.100000000000000000",
                    "max_rate": "0.200000000000000000",
                    "max_change_rate": "0.010000000000000000"
                  },
                  "min_self_delegation": "1",
                  "delegator_address": "cosmos1w539arkjhd9ceppje6q4jfvwwtcm2ducy3wy5z",
                  "validator_address": "cosmosvaloper1w539arkjhd9ceppje6q4jfvwwtcm2ducp963c3",
                  "pubkey": "cosmosvalconspub1ulx45dfpqg6vtvzexe02axd2je73q2xl8mhzf0lsqzlm6t95ylz3a7kn4k5sk0ngzn9",
                  "value": {
                    "denom": "stake",
                    "amount": "100000000"
                  }
                }
              }
            ],
            "fee": {
              "amount": [],
              "gas": "200000"
            },
            "signatures": [
              {
                "pub_key": {
                  "type": "tendermint/PubKeySm2",
                  "value": "Aqmm3ZNydd9l+8J6GOU5jrH0aZ2ydGhTBJS5YDDDcNvI"
                },
                "signature": "lAN2t+usO8fwzgZl/oaDYCsiGqITDM5colZ/kFtka42Py196AweiXPv6x/5chXMDtSDGXZyfg+WpYosscaMQ0Q=="
              }
            ],
            "memo": "45ce89c3804dacf760f60635fa3b971fe568de52@172.20.0.101:26656"
          }
        }
    ...
}
```

##### 4.2 如node2的 内容如下

```json
{
    ...
    "genutil": {
      "gentxs": [
        {
          "type": "cosmos-sdk/StdTx",
          "value": {
            "msg": [
              {
                "type": "cosmos-sdk/MsgCreateValidator",
                "value": {
                  "description": {
                    "moniker": "node1",
                    "identity": "",
                    "website": "",
                    "security_contact": "",
                    "details": ""
                  },
                  "commission": {
                    "rate": "0.100000000000000000",
                    "max_rate": "0.200000000000000000",
                    "max_change_rate": "0.010000000000000000"
                  },
                  "min_self_delegation": "1",
                  "delegator_address": "cosmos1hctlp06x9zx9uk2nznzwy68xadnvqs5wphwnmn",
                  "validator_address": "cosmosvaloper1hctlp06x9zx9uk2nznzwy68xadnvqs5wyr6xhq",
                  "pubkey": "cosmosvalconspub1ulx45dfpqfmfzyulg936am7p56e9x5dmj3m86c725gytslzhre0yxcy45kpfsaeqtyf",
                  "value": {
                    "denom": "stake",
                    "amount": "100000000"
                  }
                }
              }
            ],
            "fee": {
              "amount": [],
              "gas": "200000"
            },
            "signatures": [
              {
                "pub_key": {
                  "type": "tendermint/PubKeySm2",
                  "value": "A8PRjVTAYzPdA9vcAht4cZuJA7oZd5OxaLlSmF7X1dJS"
                },
                "signature": "ZI6eJPkwyTg7kQvEil8OU37ztsHprxZj2M5BYNmFiNQ8ScDCAoWB3ZjE7LyT7CVU4ZTtaMhnUMho2ske4Lg+8g=="
              }
            ],
            "memo": "368524d422cf3a4530cd1006d2c019ce565a3f15@172.20.0.102:26656"
          }
        }
    ...
}
```

##### 4.3 将node2的gentx拷贝到node1中，合并如下

```json
{
    ...
    "genutil": {
      "gentxs": [
        {
          "type": "cosmos-sdk/StdTx",
          "value": {
            "msg": [
              {
                "type": "cosmos-sdk/MsgCreateValidator",
                "value": {
                  "description": {
                    "moniker": "node1",
                    "identity": "",
                    "website": "",
                    "security_contact": "",
                    "details": ""
                  },
                  "commission": {
                    "rate": "0.100000000000000000",
                    "max_rate": "0.200000000000000000",
                    "max_change_rate": "0.010000000000000000"
                  },
                  "min_self_delegation": "1",
                  "delegator_address": "cosmos1w539arkjhd9ceppje6q4jfvwwtcm2ducy3wy5z",
                  "validator_address": "cosmosvaloper1w539arkjhd9ceppje6q4jfvwwtcm2ducp963c3",
                  "pubkey": "cosmosvalconspub1ulx45dfpqg6vtvzexe02axd2je73q2xl8mhzf0lsqzlm6t95ylz3a7kn4k5sk0ngzn9",
                  "value": {
                    "denom": "stake",
                    "amount": "100000000"
                  }
                }
              }
            ],
            "fee": {
              "amount": [],
              "gas": "200000"
            },
            "signatures": [
              {
                "pub_key": {
                  "type": "tendermint/PubKeySm2",
                  "value": "Aqmm3ZNydd9l+8J6GOU5jrH0aZ2ydGhTBJS5YDDDcNvI"
                },
                "signature": "lAN2t+usO8fwzgZl/oaDYCsiGqITDM5colZ/kFtka42Py196AweiXPv6x/5chXMDtSDGXZyfg+WpYosscaMQ0Q=="
              }
            ],
            "memo": "45ce89c3804dacf760f60635fa3b971fe568de52@172.20.0.101:26656"
          }
        },
        {
          "type": "cosmos-sdk/StdTx",
          "value": {
            "msg": [
              {
                "type": "cosmos-sdk/MsgCreateValidator",
                "value": {
                  "description": {
                    "moniker": "node1",
                    "identity": "",
                    "website": "",
                    "security_contact": "",
                    "details": ""
                  },
                  "commission": {
                    "rate": "0.100000000000000000",
                    "max_rate": "0.200000000000000000",
                    "max_change_rate": "0.010000000000000000"
                  },
                  "min_self_delegation": "1",
                  "delegator_address": "cosmos1hctlp06x9zx9uk2nznzwy68xadnvqs5wphwnmn",
                  "validator_address": "cosmosvaloper1hctlp06x9zx9uk2nznzwy68xadnvqs5wyr6xhq",
                  "pubkey": "cosmosvalconspub1ulx45dfpqfmfzyulg936am7p56e9x5dmj3m86c725gytslzhre0yxcy45kpfsaeqtyf",
                  "value": {
                    "denom": "stake",
                    "amount": "100000000"
                  }
                }
              }
            ],
            "fee": {
              "amount": [],
              "gas": "200000"
            },
            "signatures": [
              {
                "pub_key": {
                  "type": "tendermint/PubKeySm2",
                  "value": "A8PRjVTAYzPdA9vcAht4cZuJA7oZd5OxaLlSmF7X1dJS"
                },
                "signature": "ZI6eJPkwyTg7kQvEil8OU37ztsHprxZj2M5BYNmFiNQ8ScDCAoWB3ZjE7LyT7CVU4ZTtaMhnUMho2ske4Lg+8g=="
              }
            ],
            "memo": "368524d422cf3a4530cd1006d2c019ce565a3f15@172.20.0.102:26656"
          }
        }
    ...
}
```

##### 4.4 同样的方式将node3、node4中的gentx拷贝到node1中

#### 5、将最终的node1 中的genesis.json**覆盖其他节点的genesis.json**

#### 6、**配置每个节点的config.toml**

​	配置config.toml(~/docker_store/.dbchaind/config/config.toml)中的persistent_peers，值为"node1-id@node1-ip:26656,node2-id@node2-ip:26656,node3-id@node3-ip:26656,node4-id@node4-ip:26656,"，laddr = "tcp://127.0.0.1:26657"改为laddr = "tcp://0.0.0.0:26657"

​	示例配置如下

```toml
...
[rpc]

# TCP or UNIX socket address for the RPC server to listen on
laddr = "tcp://0.0.0.0:26657"
...


...
# Comma separated list of nodes to keep persistent connections to
persistent_peers = "45ce89c3804dacf760f60635fa3b971fe568de52@172.20.0.101:26656,368524d422cf3a4530cd1006d2c019ce565a3f15@172.20.0.102:26656,39b7c8fc33dad4789aae2681d61600d4043fb897@172.20.0.103:26656,e89b76750b53a4a852de65919729599f2e7f042f@172.20.0.104:26656"
...

```

node-id 是在第一步初始化的时候生成的，如下

```shell
~$ dbchaind init node4 --chain-id testnet
{"app_message":{"auth":{"accounts":[],"params":{"max_memo_characters":"256","sig_verify_cost_Sm2":"1000","sig_verify_cost_ed25519":"590","sig_verify_cost_secp256k1":"1000","tx_sig_limit":"7","tx_size_cost_per_byte":"10"}},"bank":{"send_enabled":true},"dbchain":{"admin_addresses":null},"distribution":{"delegator_starting_infos":[],"delegator_withdraw_infos":[],"fee_pool":{"community_pool":[]},"outstanding_rewards":[],"params":{"base_proposer_reward":"0.010000000000000000","bonus_proposer_reward":"0.040000000000000000","community_tax":"0.020000000000000000","withdraw_addr_enabled":true},"previous_proposer":"","validator_accumulated_commissions":[],"validator_current_rewards":[],"validator_historical_rewards":[],"validator_slash_events":[]},"genutil":{"gentxs":[]},"params":null,"slashing":{"missed_blocks":{},"params":{"downtime_jail_duration":"600000000000","min_signed_per_window":"0.500000000000000000","signed_blocks_window":"100","slash_fraction_double_sign":"0.050000000000000000","slash_fraction_downtime":"0.010000000000000000"},"signing_infos":{}},"staking":{"delegations":null,"exported":false,"last_total_power":"0","last_validator_powers":null,"params":{"bond_denom":"stake","historical_entries":0,"max_entries":7,"max_validators":100,"unbonding_time":"1814400000000000"},"redelegations":null,"unbonding_delegations":null,"validators":null},"supply":{"supply":[]}},"chain_id":"testnet","gentxs_dir":"","moniker":"node4","node_id":"e89b76750b53a4a852de65919729599f2e7f042f"}

```

#### 7、退出容器，启动节点,完成部署

```shell
~$ sudo docker run --net dbcnet --ip node1-ip --user 1000:1000 -it -d -v ~/docker_store:/home/dbchain dbchain
~$ sudo docker run --net dbcnet --ip node2-ip --user 1000:1000 -it -d -v ~/docker_store:/home/dbchain dbchain
~$ sudo docker run --net dbcnet --ip node3-ip --user 1000:1000 -it -d -v ~/docker_store:/home/dbchain dbchain
~$ sudo docker run --net dbcnet --ip node4-ip --user 1000:1000 -it -d -v ~/docker_store:/home/dbchain dbchain
```

#### 8、如果要保持容器时间地区与宿主机一致，启动容器时需要添加下面参数

```
-v /etc/localtime:/etc/localtime
```



### 二、客户端设置

#### 1、进入任意一个节点容器中，以node1为例

```shell
#配置cli
~$ dbchaincli config chain-id testnet
~$ dbchaincli config output json
~$ dbchaincli config indent true
~$ dbchaincli config trust-node true

#配置oracle
~$ dbchainoracle config chain-id testnet
~$ dbchainoracle config output json
~$ dbchainoracle config indent true
~$ dbchainoracle config trust-node true

#生成oracle密钥
~$ dbchainoracle query oracle-info
	"oracle-encrypted-key: H9WPAb8xjGMs58PECwVTpTsnCrJQDuCHPohiufexmTDE"
```

#### 2、将oracle密钥配置在config.toml中

在../.dbchainoracle/config/config.toml 中添加oracle密钥，示例如下

```toml
chain-id = "testnet"
indent = true
output = "json"
trust-node = true
secret_key = "asdfqweasd123456"
oracle-encrypted-key = "H9WPAb8xjGMs58PECwVTpTsnCrJQDuCHPohiufexmTDE"
```

#### 3、上一步的config.toml中配置secret_key， 示例如下

```shell
chain-id = "testnet"
indent = true
output = "json"
trust-node = true
secret_key = "asdfqweasd123456"
oracle-encrypted-key = "H9WPAb8xjGMs58PECwVTpTsnCrJQDuCHPohiufexmTDE"
secret_key = "asdfqweasd123456" #任意的长度为16的字符串
```

#### 4、获取oracle地址，并为oracle地址发送token

再次输入步骤1中的命令

```shell
~$ dbchainoracle query oracle-info
"Address: cosmos14jjve5l4yx7mnt2c3e9sxsxul7avsz66ldtwx6"  #这就是第一步生成的密钥对应的地址

#获取创世用户地址
~$ dbchaincli keys show node1
Enter keyring passphrase:
{
  "name": "node1",
  "type": "local",
  "address": "cosmos1w539arkjhd9ceppje6q4jfvwwtcm2ducy3wy5z",
  "pubkey": "cosmospub1ulx45dfpq256dhvnwf6a7e0mcfap3efe36clg6vakf6xs5cyjjukqvxrwrdusj5pq5r"
}
#发送积分,积分用于oracle创建交易，当oracle积分消耗光，oracle创建的交易会失败
#oracle 主要用于将用户托管密钥上链
~$ dbchaincli tx bank send cosmos1w539arkjhd9ceppje6q4jfvwwtcm2ducy3wy5z cosmos14jjve5l4yx7mnt2c3e9sxsxul7avsz66ldtwx6 100000dbctoken

```

#### 4、将../.dbchaincli 和 ../.dbchainoracle 目录拷贝到其他上个节点



### 三、链初始设置

#### 1、还是在容器中，配置token keeper，用户给普通用户充值积分使用，node1节点为例

```shell
# 将node1 自己添加为token keeper
# 获取地址的命令 dbchaincli keys show name 
~$ dbchaincli tx dbchain modify-token-keeper add node1Address --from node1
#将其他三个节点的创世账号也添加为 token keeper（建议）
~$ dbchaincli tx dbchain modify-token-keeper add node2Address --from node1
~$ dbchaincli tx dbchain modify-token-keeper add node3Address --from node1
~$ dbchaincli tx dbchain modify-token-keeper add node4Address --from node1
```

#### 2、设置限制p2p转账

```shell
#此操作只有token keeper 有权限
~$ dbchaincli tx dbchain  modify-p2p-transfer-limit true --from node1
#可以通过如下命令，取消限制
#dbchaincli tx dbchain  modify-p2p-transfer-limit false --from node1
```

