# 拍卖交易

目前测试发现 opensea 并不会主动帮用户撮合拍卖的订单；需要用户自己通过接受报价的方式成交；所以和普通的接受报价单的交易一样；

拍卖单和普通挂单的区别在于，后端的记录以及时间的展示不同；

## 挂单签名数据的构造

使用 EIP712 结构化签名

```JSON
{
 "message": {
  "offerer": "0xDdA9c09736c7B36316b758734cfdf4aeBCC2968A",
  "offer": [
   {
    "itemType": "2", // 721
    "token": "0xB87B37475C67761Fa36F312966A89E5A191ce069",
    "identifierOrCriteria": "2",
    "startAmount": "1",
    "endAmount": "1"
   }
  ],
  "consideration": [
   {
    "itemType": "1",
    "token": "0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6",
    "identifierOrCriteria": "0",
    "startAmount": "9750000000000000",
    "endAmount": "9750000000000000",
    "recipient": "0xDdA9c09736c7B36316b758734cfdf4aeBCC2968A"
   }
  ],
  "startTime": "1690345170",
  "endTime": "1690953568",
  "orderType": "2", // FULL_RESTRICTED
  "zone": "0x9B814233894Cd227f561B78Cc65891AA55C62Ad2",
  "zoneHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "salt": "24446860302761739304752683030156737591518664810215442929800467309395918027993",
  "conduitKey": "0x0000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000",
  "totalOriginalConsiderationItems": "1",
  "counter": "0"
 }
}
```

1. 注意，这里我们可以看到 startAmount 和 endAmount 是相同的；
2. 注意，orderType 是 2，如果 msg.sender 不是 offerer，那么就要通过 ZONE 进行验证
3. conduitKey：的目的是指定哪写通道可以对用户授权的资产进行转账
4. Zone 是一个合约，用于自定义验证规则；目前看到 seaport 用的合约是 signedZone；

## 用户提供报价单

```json
{
  "message": {
    "offerer": "0x560D5cdC0CdcA5496606Ed40EB6a9F886B768960",
    "offer": [
      {
        "itemType": "1",
        "token": "0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6",
        "identifierOrCriteria": "0",
        "startAmount": "10000000000000000",
        "endAmount": "10000000000000000"
      }
    ],
    "consideration": [
      {
        "itemType": "2",
        "token": "0xB87B37475C67761Fa36F312966A89E5A191ce069",
        "identifierOrCriteria": "2",
        "startAmount": "1",
        "endAmount": "1",
        "recipient": "0x560D5cdC0CdcA5496606Ed40EB6a9F886B768960"
      },
      {
        "itemType": "1",
        "token": "0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6",
        "identifierOrCriteria": "0",
        "startAmount": "250000000000000",
        "endAmount": "250000000000000",
        "recipient": "0x0000a26b00c1F0DF003000390027140000fAa719"
      }
    ],
    "startTime": "1690345516",
    "endTime": "1690607968",
    "orderType": "2",
    "zone": "0x000000e7Ec00e7B300774b00001314B8610022b8",
    "zoneHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "salt": "24446860302761739304752683030156737591518664810215442929817238018060826374234",
    "conduitKey": "0x0000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000",
    "totalOriginalConsiderationItems": "2",
    "counter": "0"
  }
}
```

- 和上面拍卖挂单类似，orderType 也是 2，也就是需要通过 ZONE 合约进行额外的验证；

## 用户购买 NFT 的时候调用的合约方法

1. 使用 matchAdvancedOrders 函数进行撮合交易
2. 参数解析；
   - AdvancedOrder 是订单列表，也就是上面两个订单 + advance 数据; [参考代码](../src/order/advanced_order.go)
     1. numerator 和 denominator 表示分子分母，这里是 721，只有一个 NFT，所有 1 / 1 = 1；
     2. signature 是订单的签名数据
     3. extraData： [参考代码](../src/zone/extra_data_builder.go)
   - CriteriaResolver： 这里没用到暂时填写空数组
   - Fulfillment： 我们在介绍指定用户购买的交易中介绍过
   - address：指定 recipent，主要是 seaport 不会主动帮助拍卖单子撮合，最后还是要用户接受报价单；所以两边会有差价，差价是要给这个地址。

```javascript
    function matchAdvancedOrders(
        AdvancedOrder[] calldata orders,
        CriteriaResolver[] calldata criteriaResolvers,
        Fulfillment[] calldata fulfillments,
        address recipient
    ) external payable returns (Execution[] memory executions);

    struct AdvancedOrder {
        OrderParameters parameters;
        uint120 numerator;
        uint120 denominator;
        bytes signature;
        bytes extraData;
    }

    struct OrderParameters {
        address offerer; // 0x00
        address zone; // 0x20
        OfferItem[] offer; // 0x40
        ConsiderationItem[] consideration; // 0x60
        OrderType orderType; // 0x80
        uint256 startTime; // 0xa0
        uint256 endTime; // 0xc0
        bytes32 zoneHash; // 0xe0
        uint256 salt; // 0x100
        bytes32 conduitKey; // 0x120
        uint256 totalOriginalConsiderationItems; // 0x140
    }

    struct OfferItem {
        ItemType itemType;
        address token;
        uint256 identifierOrCriteria;
        uint256 startAmount;
        uint256 endAmount;
    }

    struct ConsiderationItem {
        ItemType itemType;
        address token;
        uint256 identifierOrCriteria;
        uint256 startAmount;
        uint256 endAmount;
        address payable recipient;
    }

    struct Fulfillment {
        FulfillmentComponent[] offerComponents;
        FulfillmentComponent[] considerationComponents;
    }

    struct FulfillmentComponent {
        uint256 orderIndex;
        uint256 itemIndex;
    }
```

### AdvancedOrder

- order 1: 列表中第一个元素，这里是报价订单，其中构造了 extraData，还有 signature

```javascript
struct AdvancedOrder {
    struct OrderParameters {
		    address offerer; // 0x560D5cdC0CdcA5496606Ed40EB6a9F886B768960
		    address zone; // 0x000000e7Ec00e7B300774b00001314B8610022b8
		    OfferItem[] offer;
					struct OfferItem {
					    ItemType itemType; // 1
					    address token; // 0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6
					    uint256 identifierOrCriteria; // 0
					    uint256 startAmount;  // 10000000000000000
					    uint256 endAmount; // 10000000000000000
				  }
		    ConsiderationItem[] consideration;
						struct ConsiderationItem { //
						    ItemType itemType; // 2
						    address token; // 0xB87B37475C67761Fa36F312966A89E5A191ce069
						    uint256 identifierOrCriteria; // 2
						    uint256 startAmount; // 1
						    uint256 endAmount; // 1
						    address payable recipient; // 0x560D5cdC0CdcA5496606Ed40EB6a9F886B768960
						}
						struct ConsiderationItem { //
						    ItemType itemType; // 1
						    address token; // 0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6
						    uint256 identifierOrCriteria; // 0
						    uint256 startAmount; // 250000000000000
						    uint256 endAmount; // 250000000000000
						    address payable recipient; // 0x0000a26b00c1F0DF003000390027140000fAa719
						}
		    OrderType orderType; // 2
		    uint256 startTime; // 1690345516
		    uint256 endTime; // 1690607968
		    bytes32 zoneHash; // 0x0000000000000000000000000000000000000000000000000000000000000000
		    uint256 salt; // 24446860302761739304752683030156737591518664810215442929817238018060826374234
		    bytes32 conduitKey; // 0x0000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000
		    uint256 totalOriginalConsiderationItems; // 2
		}
    uint120 numerator; // 1
    uint120 denominator; // 1
    bytes signature; // 0x2e17ca884bfab390baf01d361734e74eac398cbcb86516ccdd89eea4d80b07a38abfdb1c4ea2ab1cb24644fd51f5f5acca91aa5b858dc13e7f9c74d8a2827217
    bytes extraData; // 0x00dda9c09736c7b36316b758734cfdf4aebcc2968a0000000064c0a1e66572e51ecc255ad1c736ca96ef5bc2d0a9c60772773f3614d074a3e0657e3742f00f8084ff750e4c31d020fe6ef765c0dd68b1e54650cc2688b6dc57c165f49c000000000000000000000000000000000000000000000000000000000000000002
}
```

- order2：订单列表第二个元素，也是拍卖订单；其中 signature 和 extraData 都是空，原因是 opensea 不会主动去撮合匹配的拍卖单子，需要用户自己选择报价单然后发送交易，所以用户是发起交易者不需要签名，也不需要 ZONE 验证；

```javascript
struct AdvancedOrder {
    struct OrderParameters {
		    address offerer; // 0xDdA9c09736c7B36316b758734cfdf4aeBCC2968A // 拍卖商家
		    address zone; // 0x0000000000000000000000000000000000000000
		    OfferItem[] offer;
					struct OfferItem {
					    ItemType itemType; // 2
					    address token; // 0xB87B37475C67761Fa36F312966A89E5A191ce069
					    uint256 identifierOrCriteria; // 2
					    uint256 startAmount;  // 1
					    uint256 endAmount; // 1
				  }
		    ConsiderationItem[] consideration;
						struct ConsiderationItem { //
						    ItemType itemType; // 1
						    address token; // 0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6
						    uint256 identifierOrCriteria; // 0
						    uint256 startAmount; // 9750000000000000
						    uint256 endAmount; // 9750000000000000
						    address payable recipient; // 0xDdA9c09736c7B36316b758734cfdf4aeBCC2968A
						}
		    OrderType orderType; // 0
		    uint256 startTime; // 1690345516
		    uint256 endTime; // 1690607968
		    bytes32 zoneHash; // 0x0000000000000000000000000000000000000000000000000000000000000000
		    uint256 salt; // 24446860302761739304752683030156737591518664810215442929818213209080753270031
		    bytes32 conduitKey; // 0x0000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000
		    uint256 totalOriginalConsiderationItems; // 1
		}
    uint120 numerator; // 1
    uint120 denominator; // 1
    bytes signature; // 0x
    bytes extraData; // 0x
}
```

### fulfillment

直接介绍过如何构造fulfillment，这里简单的看下参数对照下就行；

```javascript
	struct Fulfillment {
	    FulfillmentComponent[] offerComponents{
					struct FulfillmentComponent {
					    uint256 orderIndex; // 1
					    uint256 itemIndex;  // 0
					}
					struct FulfillmentComponent {
					    uint256 orderIndex; // 0
					    uint256 itemIndex;  // 0
					}
					struct FulfillmentComponent {
					    uint256 orderIndex; // 0
					    uint256 itemIndex;  // 0
					}
			}
	    FulfillmentComponent[] considerationComponents{
					struct FulfillmentComponent {
					    uint256 orderIndex; // 0
					    uint256 itemIndex;  // 0
					}
					struct FulfillmentComponent {
					    uint256 orderIndex; // 0
					    uint256 itemIndex;  // 1
					}
					struct FulfillmentComponent {
					    uint256 orderIndex; // 1
					    uint256 itemIndex;  // 0
					}
			}
	}
```

## Event 事件

和前面说的指定买家的撮合类似；
- 两个 OrderFulfilled 事件，应为有两个订单
- 一个 OrderMatched 事件，返回的是所有订单的 OrderHash
- 多个 Transfer 事件，表示 TOKEN 或者 NFT 的转移

## 总结

1. opensea的订单都是中心化的，拍卖也是中心化的，没有放到合约中进行；
2. 拍卖需要用户自己接受报价，后台不会帮助撮合上链

所以，如果要再拍卖结束之后帮用户提交撮合交易，需要后端实现一个Task去维护
