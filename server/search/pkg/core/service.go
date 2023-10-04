package core

import (
	"context"
	searchpb "search/gen/protos"

	"github.com/jmoiron/sqlx"
)

func CreateSearchService(db *sqlx.DB) *SearchService {
	return &SearchService{
		db: db,
	}
}

type SearchService struct {
	searchpb.UnimplementedSearchServiceServer

	db *sqlx.DB
}

func (ss *SearchService) FetchCategories(_ context.Context, _ *searchpb.Empty) (*searchpb.FetchCategoriesResponse, error) {
	categories, err := FetchCategories(ss.db)
	if err != nil {
		return nil, err
	}

	result := make([]*searchpb.Category, len(*categories))
	for i, category := range *categories {
		result[i] = &searchpb.Category{
			Id:       category.ID,
			Name:     category.Name,
			Slug:     category.Slug,
			ParentId: category.ParentID,
		}
	}

	return &searchpb.FetchCategoriesResponse{
		Categories: result,
	}, nil
}

func (ss *SearchService) CreateCategory(ctx context.Context, req *searchpb.CreateCategoryRequest) (*searchpb.Empty, error) {
	_, err := CreateCategory(ss.db, req)
	return &searchpb.Empty{}, err
}

func (ss *SearchService) DeleteCategory(_ context.Context, req *searchpb.DeleteCategoryRequest) (*searchpb.Empty, error) {
	err := DeleteCategory(ss.db, req)
	return &searchpb.Empty{}, err
}

func (ss *SearchService) GetCategory(_ context.Context, req *searchpb.GetCategoryRequest) (*searchpb.Category, error) {
	cat, err := GetCategory(ss.db, req.Id)

	if err != nil {
		return nil, err
	}

	return &searchpb.Category{
		Id:       cat.ID,
		Name:     cat.Name,
		Slug:     cat.Slug,
		ParentId: cat.ParentID,
	}, nil
}

func (ss *SearchService) UpdateCategory(_ context.Context, req *searchpb.UpdateCategoryRequest) (*searchpb.Empty, error) {
	err := UpdateCategory(ss.db, req)
	return &searchpb.Empty{}, err
}
