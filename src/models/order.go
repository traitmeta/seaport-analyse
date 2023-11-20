package models

import "github.com/shopspring/decimal"

type OfferItem struct {
	ItemType             int8            `json:"item_type"`              // 资产类型
	Token                string          `json:"token"`                  // 资产地址
	IdentifierOrCriteria string          `json:"identifier_or_criteria"` // Token为 "0"，NFT为tokenId
	StartAmount          decimal.Decimal `json:"start_amount"`           // (起始)数量/金额
	EndAmount            decimal.Decimal `json:"end_amount"`             // (结束)数量/金额
}

type ConsiderationItem struct {
	ItemType             int8            `json:"item_type"`              // 资产类型
	Token                string          `json:"token"`                  // 资产地址
	IdentifierOrCriteria string          `json:"identifier_or_criteria"` // Token为 `"0"`，NFT为tokenId
	StartAmount          decimal.Decimal `json:"start_amount"`           // (起始)数量/金额
	EndAmount            decimal.Decimal `json:"end_amount"`             // (结束)数量/金额
	Recipient            string          `json:"recipient"`              // 资产接收者地址
}

type Fulfillment struct {
	OfferComponents         []FulfillmentComponent `json:"offer_components"`
	ConsiderationComponents []FulfillmentComponent `json:"consideration_components"`
}

type FulfillmentComponent struct {
	OrderIndex string `json:"order_index"`
	ItemIndex  string `json:"item_index"`
}

type OrderParam struct {
	Offerer                         string              `json:"offerer" `
	StartTime                       int64               `json:"start_time"`
	EndTime                         int64               `json:"end_time"`
	OrderType                       int8                `json:"order_type"`
	Zone                            string              `json:"zone"`
	ZoneHash                        string              `json:"zone_hash"`
	Salt                            string              `json:"salt"`
	ConduitKey                      string              `json:"conduit_key"`
	Counter                         string              `json:"counter"`
	Offer                           []OfferItem         `json:"offer"`
	Consideration                   []ConsiderationItem `json:"consideration"`
	TotalOriginalConsiderationItems string              `json:"total_original_consideration_items"`
}

type Order struct {
	OrderParam
	OrderHash string `json:"order_hash"`
	Signature string `json:"signature"`
}

type AdvancedOrder struct {
	Parameters  OrderParam `json:"parameters"`
	Numerator   string     `json:"numerator"`
	Denominator string     `json:"denominator"`
	Signature   string     `json:"signature"`
	ExtraData   string     `json:"extra_data"`
}
