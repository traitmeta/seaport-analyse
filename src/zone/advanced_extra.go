package zone

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/holiman/uint256"

	"github.com/traitmeta/seaport-analyse/src/utils"
)

const SignedZoneExpirationSeconds = 60

var BizAdvancedExtra *AdvancedExtraBuilder

// BestOrderBiz Asset's best ask/bid 相关 biz
type AdvancedExtraBuilder struct {
	contract  string
	chainId   string
	fulfiller string
}

func NewAdvancedExtraBuilder(fulfiller, chainId, contract string) *AdvancedExtraBuilder {
	return &AdvancedExtraBuilder{
		contract:  contract,
		chainId:   chainId,
		fulfiller: fulfiller,
	}
}

func (b *AdvancedExtraBuilder) BuildExtraData(priv, orderHash, context string) (string, error) {
	expiration := time.Now().Unix() + int64(SignedZoneExpirationSeconds)
	return b.buildExtraData(priv, orderHash, context, expiration)
}

func (b *AdvancedExtraBuilder) buildExtraData(priv, orderHash, context string, expiration int64) (string, error) {
	signaData := b.buildSignData(orderHash, context, expiration)
	signature, err := b.SignForZone(priv, signaData)
	if err != nil {
		return "", fmt.Errorf("failed to sign typed data: %v", err)
	}

	compactSignature := b.convertSignatureToEIP2098(signature)
	extraData := fmt.Sprintf("0x00%s%s%s%s", b.fulfiller[2:], b.toPaddedBytes(uint64(expiration), 8), compactSignature[2:], context[2:])

	return extraData, nil
}

func (b *AdvancedExtraBuilder) convertSignatureToEIP2098(signature string) string {
	if len(signature) == 130 {
		return signature
	}

	if len(signature) != 132 {
		fmt.Println("Error")
	}
	sigBytes, err := hexutil.Decode(signature)
	if err != nil {
		fmt.Println("Failed to decode signature:", err)
		return ""
	}
	r := sigBytes[:32]
	s, ok := uint256.FromBig(big.NewInt(0).SetBytes(sigBytes[32:64]))
	if ok {
		return ""
	}

	v := sigBytes[64]
	if v-27 == 1 {
		max := uint256.NewInt(1).Lsh(uint256.NewInt(1), 255)
		s = uint256.NewInt(0).Or(max, s)
	}

	compactSignature := append(r, s.Bytes()...)
	return hexutil.Encode(compactSignature)
}

func (b *AdvancedExtraBuilder) toPaddedBytes(value uint64, numBytes int) string {
	hexVal := hexutil.EncodeUint64(value)
	valueBytes := common.LeftPadBytes(hexutil.MustDecode(hexVal), numBytes)

	return hexutil.Encode(valueBytes)[2 : numBytes*2+2]
}

func (b *AdvancedExtraBuilder) BuildZoneContext(considerFirstItemIdentifier string) (string, bool) {
	var valueBytes []byte
	s := considerFirstItemIdentifier
	if len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		valueBytes = common.LeftPadBytes(common.FromHex(s), 32)
	} else {
		bigIdentifier, success := big.NewInt(0).SetString(considerFirstItemIdentifier, 10)
		if !success {
			return "", false
		}
		valueBytes = common.LeftPadBytes(bigIdentifier.Bytes(), 32)
	}

	return "0x00" + hexutil.Encode(valueBytes)[2:32*2+2], true
}

// contract is signedZone contract address
// fulfiller is activedSigner for signedZone contract
// orderHash is order hash
// context is extend of first consideration identify
func (b *AdvancedExtraBuilder) buildSignData(orderHash, context string, expiration int64) string {
	var signFormat = `
	{
		"types": {
		  "SignedOrder": [
			{ "name": "fulfiller", "type": "address" },
			{ "name": "expiration", "type": "uint64" },
			{ "name": "orderHash", "type": "bytes32" },
			{ "name": "context", "type": "bytes" }
		  ],
		  "EIP712Domain": [
			{ "name": "name", "type": "string" },
			{ "name": "version", "type": "string" },
			{ "name": "chainId", "type": "uint256" },
			{ "name": "verifyingContract", "type": "address" }
		  ]
		},
		"domain": {
		  "name": "SignedZone",
		  "version": "1.0",
		  "chainId": "%s",
		  "verifyingContract": "%s"
		},
		"primaryType": "SignedOrder",
		"message": {
		  "fulfiller": "%s",
		  "expiration": "%d",
		  "orderHash": "%s",
		  "context": "%s"
		}
	  }
	`
	return fmt.Sprintf(signFormat, b.chainId, b.contract, b.fulfiller, expiration, orderHash, context)
}

func (b *AdvancedExtraBuilder) SignForZone(priv string, signData string) (signature string, err error) {
	var td apitypes.TypedData
	if err = json.Unmarshal([]byte(signData), &td); err != nil {
		return "", err
	}

	prv, err := crypto.HexToECDSA(priv)
	if err != nil {
		return "", err
	}

	return utils.EIP712Sign(prv, &td)
}
