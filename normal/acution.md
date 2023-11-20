# 拍卖交易

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

1. 使用 matchAdvancedOrders函数进行撮合交易
2. 参数解析；
   - AdvancedOrder 是订单列表，也就是上面两个订单 + advance数据
      1. numerator和denominator表示分子分母，这里是721，只有一个NFT，所有 1 / 1 = 1；
      2. signature 是订单的签名数据
      3. extraData： 
   - CriteriaResolver： 这里没用到暂时填写空数组
   - Fulfillment： 我们在介绍指定用户购买的交易中介绍过
   - address：指定recipent，主要是seaport不会主动帮助拍卖单子撮合，最后还是要用户接受报价单；所以两边会有差价，差价是要给这个地址。

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

## Event 事件


## 总结

