package order

import (
	"errors"
	"strconv"

	"github.com/traitmeta/seaport-analyse/src/models"
)

func BuildFulfillments(orders []models.OrderParam) ([]models.Fulfillment, error) {
	if len(orders) < 2 {
		return nil, errors.New("fulfillment match mush more than two order")
	}

	fulfillments := []models.Fulfillment{}
	for i, offer := range orders {
		for j, consideration := range orders {
			offerIdx := strconv.FormatInt(int64(i), 10)
			considerationIdx := strconv.FormatInt(int64(j), 10)
			fulfillments = append(fulfillments, buildFulfillments(offerIdx, considerationIdx, offer.Offer, consideration.Consideration)...)
		}
	}

	return fulfillments, nil
}

func buildFulfillments(offerOrderIdx, considerationOrderIdx string, offer []models.OfferItem, consideration []models.ConsiderationItem) []models.Fulfillment {
	fulfillments := []models.Fulfillment{}
	for i, offerItem := range offer {
		for j, considerationItem := range consideration {
			if offerItem.ItemType == considerationItem.ItemType &&
				offerItem.Token == considerationItem.Token &&
				offerItem.IdentifierOrCriteria == considerationItem.IdentifierOrCriteria {
				fulfillments = append(fulfillments,
					models.Fulfillment{
						OfferComponents:         []models.FulfillmentComponent{{OrderIndex: offerOrderIdx, ItemIndex: strconv.Itoa(i)}},
						ConsiderationComponents: []models.FulfillmentComponent{{OrderIndex: considerationOrderIdx, ItemIndex: strconv.Itoa(j)}},
					},
				)
			}
		}
	}
	return fulfillments
}
