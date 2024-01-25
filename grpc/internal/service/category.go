package service

import (
	"context"
	"github.com/jardelkuhnen/fullcyclecourse/grpc/internal/database"
	"github.com/jardelkuhnen/fullcyclecourse/grpc/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(db database.Category) *CategoryService {
	return &CategoryService{CategoryDB: db}
}

func (c *CategoryService) CreateCategory(ctx context.Context, request *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	category, err := c.CategoryDB.Create(request.GetName(), request.GetDescription())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating new category")
	}

	categoryResp := &pb.CategoryResponse{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return categoryResp, nil
}

func (c *CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryList, error) {
	categories, err := c.CategoryDB.FindAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error listing categories")
	}

	var categoryList []*pb.Category
	for _, category := range categories {
		categoryResp := &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}
		categoryList = append(categoryList, categoryResp)
	}

	return &pb.CategoryList{Categories: categoryList}, nil
}

func (c *CategoryService) GetCategory(ctx context.Context, in *pb.CategoryGetRequest) (*pb.CategoryResponse, error) {
	category, err := c.CategoryDB.FindByCategoryID(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting category")
	}

	return &pb.CategoryResponse{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil

}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}

	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}
		if err != nil {
			return err
		}

		category, err := c.CategoryDB.Create(request.GetName(), request.GetDescription())
		if err != nil {
			return err
		}

		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})

	}
}

func (c *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		categoryCcreated, err := c.CategoryDB.Create(request.GetName(), request.GetDescription())
		if err != nil {
			return err
		}

		err = stream.Send(&pb.Category{
			Id:          categoryCcreated.ID,
			Name:        categoryCcreated.Name,
			Description: categoryCcreated.Description,
		})

		if err != nil {
			return err
		}
	}
}
