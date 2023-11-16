package src

import (
	"errors"
	"strconv"

	"github.com/shopspring/decimal"
)

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

func BuildFulfillments(orders []OrderParam) ([]Fulfillment, error) {
	if len(orders) < 2 {
		return nil, errors.New("fulfillment match mush more than two order")
	}

	fulfillments := []Fulfillment{}
	for i, offer := range orders {
		for j, consideration := range orders {
			offerIdx := strconv.FormatInt(int64(i), 10)
			considerationIdx := strconv.FormatInt(int64(j), 10)
			fulfillments = append(fulfillments, buildFulfillments(offerIdx, considerationIdx, offer.Offer, consideration.Consideration)...)
		}
	}

	return fulfillments, nil
}

func buildFulfillments(offerOrderIdx, considerationOrderIdx string, offer []OfferItem, consideration []ConsiderationItem) []Fulfillment {
	fulfillments := []Fulfillment{}
	for i, offerItem := range offer {
		for j, considerationItem := range consideration {
			if offerItem.ItemType == considerationItem.ItemType &&
				offerItem.Token == considerationItem.Token &&
				offerItem.IdentifierOrCriteria == considerationItem.IdentifierOrCriteria {
				fulfillments = append(fulfillments,
					Fulfillment{
						OfferComponents:         []FulfillmentComponent{{OrderIndex: offerOrderIdx, ItemIndex: strconv.Itoa(i)}},
						ConsiderationComponents: []FulfillmentComponent{{OrderIndex: considerationOrderIdx, ItemIndex: strconv.Itoa(j)}},
					},
				)
			}
		}
	}
	return fulfillments
}
