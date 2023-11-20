package utils

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

const EIP712DomainField = "EIP712Domain"

func EIP712Sign(prv *ecdsa.PrivateKey, data *apitypes.TypedData) (string, error) {
	var signature string
	dataHash, err := data.HashStruct(data.PrimaryType, data.Message)
	if err != nil {
		return signature, errors.New("calculate data hash falied")
	}

	domainSep, err := data.HashStruct(EIP712DomainField, data.Domain.Map())
	if err != nil {
		return signature, errors.New("calculate domain hash falied")
	}

	signData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSep), string(dataHash)))
	signDataDigest := ethcrypto.Keccak256Hash(signData)

	sig, err := ethcrypto.Sign(signDataDigest.Bytes(), prv)
	if err != nil {
		return signature, errors.New("EIP712 sign falied")
	}

	sig[64] += 27
	signature = hexutil.Encode(sig)
	return signature, err
}
