package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	catalogDomain "github.com/nlsnnn/berezhok/internal/modules/catalog/domain"
	catalogErrors "github.com/nlsnnn/berezhok/internal/modules/catalog/errors"
	orderErrors "github.com/nlsnnn/berezhok/internal/modules/order/errors"
	orderService "github.com/nlsnnn/berezhok/internal/modules/order/service"
)

type catalogBoxReader interface {
	GetBoxByID(ctx context.Context, id string) (*catalogDomain.SurpriseBox, error)
}

type CatalogBoxProvider struct {
	boxes catalogBoxReader
}

func NewCatalogBoxProvider(boxes catalogBoxReader) *CatalogBoxProvider {
	return &CatalogBoxProvider{boxes: boxes}
}

func (p *CatalogBoxProvider) GetBoxForOrder(ctx context.Context, boxID uuid.UUID) (*orderService.BoxForOrder, error) {
	box, err := p.boxes.GetBoxByID(ctx, boxID.String())
	if err != nil {
		if errors.Is(err, catalogErrors.ErrBoxNotFound) {
			return nil, orderErrors.ErrBoxNotAvailable
		}
		return nil, fmt.Errorf("get catalog box for order: %w", err)
	}

	return &orderService.BoxForOrder{
		LocationID: box.LocationID,
		Amount:     box.Price.Discount,
		PickupTime: box.PickupTime,
		Available:  box.IsAvailable(),
	}, nil
}
