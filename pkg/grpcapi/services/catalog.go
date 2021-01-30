package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/oam-dev/velacp/pkg/datastore"
	"github.com/oam-dev/velacp/pkg/datastore/model"
	"github.com/oam-dev/velacp/pkg/proto/catalogservice"
)

type CatalogService struct {
	Store datastore.CatalogStore
}

func (c *CatalogService) AddCatalog(ctx context.Context, request *catalogservice.AddCatalogRequest) (*catalogservice.AddCatalogResponse, error) {
	catalog := &model.Catalog{
		Id:   uuid.New().String(),
		Name: request.Name,
		Desc: request.Desc,
	}

	now := time.Now().Unix()
	if catalog.CreatedAt == 0 {
		catalog.CreatedAt = now
	}
	catalog.UpdatedAt = now

	err := c.Store.PutCatalog(ctx, catalog)
	if err != nil {
		return nil, err
	}
	return &catalogservice.AddCatalogResponse{}, nil
}

func (c *CatalogService) GetCatalog(ctx context.Context, request *catalogservice.GetCatalogRequest) (*catalogservice.GetCatalogResponse, error) {
	panic("implement me")
}

func (c *CatalogService) ListCatalogs(ctx context.Context, request *catalogservice.ListCatalogsRequest) (*catalogservice.ListCatalogsResponse, error) {
	catalogs, err := c.Store.ListCatalogs(ctx)
	if err != nil {
		return nil, err
	}
	return &catalogservice.ListCatalogsResponse{
		Catalogs: catalogs,
	}, nil
}

func (c *CatalogService) DelCatalog(ctx context.Context, request *catalogservice.DelCatalogRequest) (*catalogservice.DelCatalogResponse, error) {
	panic("implement me")
}
