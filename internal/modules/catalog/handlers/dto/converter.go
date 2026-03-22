package dto

import (
	"github.com/nlsnnn/berezhok/internal/modules/catalog/domain"
	"github.com/nlsnnn/berezhok/internal/modules/catalog/service"
)

func (r CreateBoxRequest) ToInput() service.CreateBoxInput {
	return service.CreateBoxInput{
		LocationID:      r.LocationID,
		Name:            r.Name,
		Description:     r.Description,
		OriginalPrice:   r.OriginalPrice,
		DiscountPrice:   r.DiscountPrice,
		PickupTimeStart: r.PickupTimeStart,
		PickupTimeEnd:   r.PickupTimeEnd,
		Quantity:        r.Quantity,
		Status:          r.Status,
		Image:           r.Image,
	}
}

func (r UpdateBoxRequest) ToInput(boxID string) service.UpdateBoxInput {
	return service.UpdateBoxInput{
		ID:              boxID,
		Name:            r.Name,
		Description:     r.Description,
		OriginalPrice:   r.OriginalPrice,
		DiscountPrice:   r.DiscountPrice,
		PickupTimeStart: r.PickupTimeStart,
		PickupTimeEnd:   r.PickupTimeEnd,
		Quantity:        r.Quantity,
		Status:          r.Status,
		Image:           r.Image,
	}
}

func BoxToResponse(box domain.SurpriseBox) BoxResponse {
	return BoxResponse{
		ID:            box.ID,
		LocationID:    box.LocationID,
		Name:          box.Name,
		Description:   box.Description,
		OriginalPrice: box.Price.Original,
		DiscountPrice: box.Price.Discount,
		PickupTime: PickupTimeResponse{
			Start: box.PickupTime.Start.Format("15:04"),
			End:   box.PickupTime.End.Format("15:04"),
		},
		Quantity:  box.Quantity,
		Image:     box.Image,
		Status:    string(box.Status),
		CreatedAt: box.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func BoxesToResponses(boxes []domain.SurpriseBox) []BoxResponse {
	responses := make([]BoxResponse, len(boxes))
	for i, box := range boxes {
		responses[i] = BoxToResponse(box)
	}
	return responses
}
