package order

import (
	"errors"

	"github.com/traitmeta/seaport-analyse/src/models"
	"github.com/traitmeta/seaport-analyse/src/utils"
	"github.com/traitmeta/seaport-analyse/src/zone"
)

// first element is bid order, second is ask order
func BuildAdvancedOrders(orders []models.Order) ([]models.AdvancedOrder, error) {
	counterAdvancedOrder, err := buildAdvancedOrder(orders[0])
	if err != nil {
		return nil, err
	}

	offererAdvancedOrder, err := buildAdvancedOrder(orders[1])
	if err != nil {
		return nil, err
	}

	return []models.AdvancedOrder{*counterAdvancedOrder, *offererAdvancedOrder}, nil
}

func buildAdvancedOrder(order models.Order) (*models.AdvancedOrder, error) {
	extraData := "0x"
	var err error
	if order.OrderType == utils.OrderTypeFullRestricted || order.OrderType == utils.OrderTypePartialRestricted {
		zoneCtx, success := zone.ZoneExtraDataBuilder.BuildZoneContext(order.Consideration[0].IdentifierOrCriteria)
		if !success {
			return nil, errors.New("cannot convert context from consideration")
		}
		extraData, err = zone.ZoneExtraDataBuilder.BuildExtraData(zone.ZoneActivedSignerPK, order.OrderHash, zoneCtx)
		if err != nil {
			return nil, err
		}
	}
	advancedOrder := &models.AdvancedOrder{
		Parameters:  order.OrderParam,
		Numerator:   "1",
		Denominator: "1",
		Signature:   order.Signature,
		ExtraData:   extraData,
	}

	return advancedOrder, nil
}
