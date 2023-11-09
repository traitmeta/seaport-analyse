# 普通交易模式

## 挂单签名数据的构造

使用 EIP712 结构化签名

```JSON
{
 "message": {
  "offerer": "0x560D5cdC0CdcA5496606Ed40EB6a9F886B768960",
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
    "startAmount": "20475000000000000",
    "endAmount": "20475000000000000",
    "recipient": "0x560D5cdC0CdcA5496606Ed40EB6a9F886B768960"
   },
   {
    "itemType": "0",
    "token": "0x0000000000000000000000000000000000000000",
    "identifierOrCriteria": "0",
    "startAmount": "525000000000000",
    "endAmount": "525000000000000",
    "recipient": "0x0000a26b00c1F0DF003000390027140000fAa719"
   }
  ],
  "startTime": "1689923450",
  "endTime": "1689927048",
  "orderType": "0",
  "zone": "0x004C00500000aD104D7DBd00e3ae0A5C00560C00",
  "zoneHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "salt": "24446860302761739304752683030156737591518664810215442929816108075358245614181",
  "conduitKey": "0x0000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000",
  "totalOriginalConsiderationItems": "2",
  "counter": "0"
 }
}
```

1. 注意，这里我们可以看到 startAmount 和 endAmount 是相同的；
2. 注意，orderType 是 0，表示只要用户进行了签名，任何用户都可以将交易发送到链上
3. conduitKey：的目的是指定哪写通道可以对用户授权的资产进行转账
4. Zone 会在拍卖场景详细说目前就是 ZERO

## 用户购买 NFT 的时候调用的合约方法

1. 使用 fulfillBasicOrder_efficient_6GL6yc/fulfillBasicOrder 都可以完成交易，efficient 更节省 gas
2. 参数是上面的挂单信息进行转变而来；
   - consideration 开头的是上面 ConsiderationItem 的第一个 Item 的信息
   - additionalRecipients 是 ConsiderationItem 排除第一个 Item 之后剩余的信息
   - offer 开头的是上面 offerItem 的信息
   - basicOrderType 是需要转换的，[可以参考](https://github.com/cryptochou/seaport-analysis)
   - 签名是上面挂单用户的签名
   - 其他信息可以照搬下来

```rust
    function fulfillBasicOrder_efficient_6GL6yc(
            BasicOrderParameters calldata parameters
        ) external payable returns (bool fulfilled);

    struct BasicOrderParameters {
        // calldata offset
        address considerationToken; // 0x0000000000000000000000000000000000000000
        uint256 considerationIdentifier; // 0
        uint256 considerationAmount; // 20475000000000000
        address payable offerer; // 0x560D5cdC0CdcA5496606Ed40EB6a9F886B768960
        address zone; // 0x004C00500000aD104D7DBd00e3ae0A5C00560C00
        address offerToken; // 0xD5835369d4F691094D7509296cFC4dA19EFe4618
        uint256 offerIdentifier; // 84946
        uint256 offerAmount; // 1
        BasicOrderType basicOrderType; // 0
        uint256 startTime; // 1689923450
        uint256 endTime; // 1689927048
        bytes32 zoneHash; // 0x0000000000000000000000000000000000000000000000000000000000000000
        uint256 salt; // 24446860302761739304752683030156737591518664810215442929816108075358245614181
        bytes32 offererConduitKey; // 0x0000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000
        bytes32 fulfillerConduitKey; // 0x0000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000
        uint256 totalOriginalAdditionalRecipients; // 1
        AdditionalRecipient[] additionalRecipients; // [525000000000000,0x0000a26b00c1F0DF003000390027140000fAa719]
        bytes signature; // 0xc529f5de19907fbbd858d9f740e373e1c755dd6f499434052054464fcd9af2c7dcb72f5c4d28c36d6831711bf1329dfd0e0fb3c927b0a87f11e8118d9bd69db2
    }
```

## Event 事件

主要涉及到下面两个, [链接](https://goerli.etherscan.io/tx/0x9acd4ccb2337223b159f8668914269447fd88136951bef40a21c0fac5c8b1226#eventlog)

- `OrderFulfilled (bytes32 orderHash, index address offerer, index address zone, address recipient, tuple[] offer, tuple[] consideration)`
- `Transfer (index address from, index address to, uint256 tokens)`

## 总结

1. 首先用户想要出售 NFT，调用 opensea 后端构造买单数据，进行签名，之后将签名数据回传到 opensea 后端；
2. 买家要买 NFT 的时候，直接构造 fulfillBasicOrder_efficient_6GL6yc 交易，填充交易参数 BasicOrderParameters；同样的这些数据是 opensea 后端构造好的
3. 买家签名交易信息，发送到链上

注意：这里面只有卖家对订单进行了签名，而买家只对交易签名; seaport 合约会恢复出卖家的订单，并验证签名；
