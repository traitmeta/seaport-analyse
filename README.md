# seaport-analyse

主要分析各种交易场景，交易单构建、签名、不同的交易函数、参数已经 event 分析；

## EIP712 订单结构解析

结构化签名数的如下，我们默认省略这些信息，只展示message在后面的分析中
```json
{
  "types": {
    "EIP712Domain": [
      {
        "name": "name",
        "type": "string"
      },
      {
        "name": "version",
        "type": "string"
      },
      {
        "name": "chainId",
        "type": "uint256"
      },
      {
        "name": "verifyingContract",
        "type": "address"
      }
    ],
    "OrderComponents": [
      {
        "name": "offerer",
        "type": "address"
      },
      {
        "name": "zone",
        "type": "address"
      },
      {
        "name": "offer",
        "type": "OfferItem[]"
      },
      {
        "name": "consideration",
        "type": "ConsiderationItem[]"
      },
      {
        "name": "orderType",
        "type": "uint8"
      },
      {
        "name": "startTime",
        "type": "uint256"
      },
      {
        "name": "endTime",
        "type": "uint256"
      },
      {
        "name": "zoneHash",
        "type": "bytes32"
      },
      {
        "name": "salt",
        "type": "uint256"
      },
      {
        "name": "conduitKey",
        "type": "bytes32"
      },
      {
        "name": "counter",
        "type": "uint256"
      }
    ],
    "OfferItem": [
      {
        "name": "itemType",
        "type": "uint8"
      },
      {
        "name": "token",
        "type": "address"
      },
      {
        "name": "identifierOrCriteria",
        "type": "uint256"
      },
      {
        "name": "startAmount",
        "type": "uint256"
      },
      {
        "name": "endAmount",
        "type": "uint256"
      }
    ],
    "ConsiderationItem": [
      {
        "name": "itemType",
        "type": "uint8"
      },
      {
        "name": "token",
        "type": "address"
      },
      {
        "name": "identifierOrCriteria",
        "type": "uint256"
      },
      {
        "name": "startAmount",
        "type": "uint256"
      },
      {
        "name": "endAmount",
        "type": "uint256"
      },
      {
        "name": "recipient",
        "type": "address"
      }
    ]
  },
  "primaryType": "OrderComponents",
  "domain": {
    "name": "Seaport",
    "version": "1.5",
    "chainId": "5",
    "verifyingContract": "0x00000000000000ADc04C56Bf30aC9d3c0aAF14dC"
  },
  "message": {}
}
```

## 普通交易模式

场景：拥有721NFT，想要出售这个 NFT 换取 ETH

链接：[普通交易](normal/normal.md)

## 指定买家模式

场景：拥有NFT，想要出售，并且指定了具体的买家
链接：[指定买家](normal/private.md)