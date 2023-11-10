# 指定交易用户模式

## 挂单签名数据的构造

使用 EIP712 结构化签名

```JSON
{
  "message": {
  "offerer": "0xDdA9c09736c7B36316b758734cfdf4aeBCC2968A",
  "offer": [
   {
    "itemType": "2",
    "token": "0xD5835369d4F691094D7509296cFC4dA19EFe4618",
    "identifierOrCriteria": "84946",
    "startAmount": "1",
    "endAmount": "1"
   }
  ],
  "consideration": [
   {
    "itemType": "0",
    "token": "0x0000000000000000000000000000000000000000",
    "identifierOrCriteria": "0",
    "startAmount": "19500000000000000",
    "endAmount": "19500000000000000",
    "recipient": "0xDdA9c09736c7B36316b758734cfdf4aeBCC2968A"
   },
   {
    "itemType": "0",
    "token": "0x0000000000000000000000000000000000000000",
    "identifierOrCriteria": "0",
    "startAmount": "500000000000000",
    "endAmount": "500000000000000",
    "recipient": "0x0000a26b00c1F0DF003000390027140000fAa719" // Fees
   },
   {
    "itemType": "2",
    "token": "0xD5835369d4F691094D7509296cFC4dA19EFe4618",
    "identifierOrCriteria": "84946",
    "startAmount": "1",
    "endAmount": "1",
    "recipient": "0x560D5cdC0CdcA5496606Ed40EB6a9F886B768960" // 指定买家
   }
  ],
  "startTime": "1689907839",
  "endTime": "1689911434",
  "orderType": "0",
  "zone": "0x004C00500000aD104D7DBd00e3ae0A5C00560C00",
  "zoneHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "salt": "24446860302761739304752683030156737591518664810215442929818138293242152910413",
  "conduitKey": "0x0000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000",
  "totalOriginalConsiderationItems": "3",
  "counter": "0"
 }
}
```

1. consideration 有三个 Item，第一个是卖家收到的金额，第二个是平台手续费，第三个是指定买家收到 NFT；
2. 其他数据和之前类似

## 用户购买 NFT 的时候调用的合约方法

1. 使用 matchOrders 进行撮合，matchOrders 有两个参数，一个是订单列表，另一个是 fulfillment，我们通过具体的例子来说明参数的含义

   ```rust
       function matchOrders(
           Order[] calldata orders,
           Fulfillment[] calldata fulfillments
       ) external payable returns (Execution[] memory executions);

       struct Fulfillment {
           FulfillmentComponent[] offerComponents;
           FulfillmentComponent[] considerationComponents;
       }

       struct FulfillmentComponent {
           uint256 orderIndex;
           uint256 itemIndex;
       }

       struct Order {
           OrderParameters parameters;
           bytes signature;
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
           // offer.length                          // 0x160
       }
   ```

2. 订单列表第一个元素分析, 上面的代码中可以看到每一个订单都包含两个字段，OrderParameters 和订单的签名；解析后可以看到第一个订单就是挂出售 NFT 的订单和签名信息；

   ```rust
   struct OrderParameters {
       address offerer; // 0xDdA9c09736c7B36316b758734cfdf4aeBCC2968A
       address zone; // 0x004C00500000aD104D7DBd00e3ae0A5C00560C00
       OfferItem[] offer;
               struct OfferItem {
                   ItemType itemType; // 2
                   address token; // 0xD5835369d4F691094D7509296cFC4dA19EFe4618
                   uint256 identifierOrCriteria; // 84946
                   uint256 startAmount;  // 1
                   uint256 endAmount; // 1
           }
       ConsiderationItem[] consideration;
                   struct ConsiderationItem { //
                       ItemType itemType; // 0
                       address token; // 0x0000000000000000000000000000000000000000
                       uint256 identifierOrCriteria; // 0
                       uint256 startAmount; // 19500000000000000
                       uint256 endAmount; // 19500000000000000
                       address payable recipient; // 0xDdA9c09736c7B36316b758734cfdf4aeBCC2968A
                   }
                   struct ConsiderationItem { //
                       ItemType itemType; // 0
                       address token; // 0x0000000000000000000000000000000000000000
                       uint256 identifierOrCriteria; // 0
                       uint256 startAmount; // 500000000000000
                       uint256 endAmount; // 500000000000000
                       address payable recipient; // 0x0000a26b00c1F0DF003000390027140000fAa719
                   }
                   struct ConsiderationItem { //
                       ItemType itemType; // 2
                       address token; // 0xD5835369d4F691094D7509296cFC4dA19EFe4618
                       uint256 identifierOrCriteria; // 84946
                       uint256 startAmount; // 1
                       uint256 endAmount; // 1
                       address payable recipient; // 0x560D5cdC0CdcA5496606Ed40EB6a9F886B768960
                   }
       OrderType orderType; // 0
       uint256 startTime; // 1689907839
       uint256 endTime; // 1689911434
       bytes32 zoneHash; // 0x0000000000000000000000000000000000000000000000000000000000000000
       uint256 salt; // 24446860302761739304752683030156737591518664810215442929818138293242152910413
       bytes32 conduitKey; // 0x0000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000
       uint256 totalOriginalConsiderationItems; // 3
       // offer.length                          // 0x160
   }

   0	.signature	bytes	0xa6ce64e20d5b3d8e27e727f66e4b601c1b052bf5d2b7b9795558cb0c35041fa21e709666b8eeadbf25b9c1624447b5537c243223cbe063c5ca993a309e42a698
   ```

3. 订单列表第二个元素分析：下面的参数可以看到，我们构造了一个对手订单，这个订单只填写了 Offer 信息，也就是买家愿意支付的金额；并且这个订单是没有签名信息，为什么没有订单的签名呢？**Opensea 的代码中会进行判断，如果发送交易的用户就是订单的 offerer，那么就会跳过签名验证；再加上用户已经对交易进行了签名，表示用户同意这笔交易，所以就不需要二次签名；**

   ```rust
       struct OrderParameters {
           address offerer; // 0x560D5cdC0CdcA5496606Ed40EB6a9F886B768960
           address zone; // 0x0000000000000000000000000000000000000000
           OfferItem[] offer;
                   struct OfferItem {
                       ItemType itemType; // 0
                       address token; // 0x0000000000000000000000000000000000000000
                       uint256 identifierOrCriteria; // 0
                       uint256 startAmount;  // 20000000000000000
                       uint256 endAmount; // 20000000000000000
               }
           ConsiderationItem[] consideration; // 空
           OrderType orderType; // 0
           uint256 startTime; // 1689907839
           uint256 endTime; // 1689911434
           bytes32 zoneHash; // 0x0000000000000000000000000000000000000000000000000000000000000000
           uint256 salt; // 24446860302761739304752683030156737591518664810215442929808117352648472516547
           bytes32 conduitKey; // 0x0000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000
           uint256 totalOriginalConsiderationItems; // 0
       },

       1	.signature	bytes	0x
   ```

4.fulfillment 参数的构建，fulfillment是对订单以及订单中Offer和Consideration的索引，所以我们可以看到一个FulFillment包含offerComponents和considerationComponents两个属性；
- 我们拆解一下FulfillmentComponent参数
    - orderIndex：表示再Orders中的索引，从0开始；
    - itemIndex：表示在Offer列表或者Consideration列表中的索引，从0开始； 
- 手工画了一个图帮助理解 [img](fulfillment.jpg)

    ```RUST
        Fulfillment[] calldata fulfillments = [
            struct Fulfillment {
                FulfillmentComponent[] offerComponents{
                    struct FulfillmentComponent {
                        uint256 orderIndex; // 0
                        uint256 itemIndex;  // 0
                    }
                }
                FulfillmentComponent[] considerationComponents{
                    struct FulfillmentComponent {
                        uint256 orderIndex; // 0
                        uint256 itemIndex;  // 2
                    }
                }
            },
            struct Fulfillment {
                FulfillmentComponent[] offerComponents{
                    struct FulfillmentComponent {
                        uint256 orderIndex; // 1
                        uint256 itemIndex;  // 0
                    }
                }
                FulfillmentComponent[] considerationComponents{
                    struct FulfillmentComponent {
                        uint256 orderIndex; // 0
                        uint256 itemIndex;  // 0
                    }
                }
            },
            struct Fulfillment {
                FulfillmentComponent[] offerComponents{
                    struct FulfillmentComponent {
                        uint256 orderIndex; // 1
                        uint256 itemIndex;  // 0
                    }
                }
                FulfillmentComponent[] considerationComponents{
                    struct FulfillmentComponent {
                        uint256 orderIndex; // 0
                        uint256 itemIndex;  // 1
                    }
                }
            },
        ]
    ```

## Event 事件

- 两个 OrderFulfilled 事件，应为有两个订单
- 一个 OrderMatched 事件，返回的是所有订单的 OrderHash
- 多个 Transfer 事件，表示 TOKEN 或者 NFT 的转移
- [链接](https://goerli.etherscan.io/tx/0x39006958ce86c456585ab5581119839697d2ca501876127d7e1c204022c77d9d)

1. `OrderFulfilled (bytes32 orderHash, index_topic_1 address offerer, index_topic_2 address zone, address recipient, tuple[] offer, tuple[] consideration)`
2. `OrderFulfilled (bytes32 orderHash, index_topic_1 address offerer, index_topic_2 address zone, address recipient, tuple[] offer, tuple[] consideration)`
3. `OrdersMatched (bytes32[] orderHashes)`
4. `Transfer (index_topic_1 address from, index_topic_2 address to, uint256 tokens)`

## 总结

1. 首先用户想要出售 NFT，调用 opensea 后端构造卖单数据，进行签名，之后将签名数据回传到 opensea 后端；**卖单的 consideration 中会增加一条购买人的信息**；
2. 买家要买 NFT 的时候，调用 opensea 后端，后端会根据卖单的订单信息，生成一个没签名信息的对手单，并且构建好 fulfillment；之后将构造好的交易返给前端让用户签名
3. 买家签名交易信息，发送到链上

**注意：买家订单不需要签名**
