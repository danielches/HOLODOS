package server

import (
	"HOLODOS/internal/models"
	"HOLODOS/internal/service"
	fridgev1 "HOLODOS/proto"
	"context"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FridgeServer struct {
	fridgev1.UnimplementedFridgeServiceServer
	service *service.FridgeService
}

func NewFridgeServer(svc *service.FridgeService) *FridgeServer {
	return &FridgeServer{service: svc}
}

func (s *FridgeServer) OpenFridge(ctx context.Context, req *fridgev1.OpenRequest) (*fridgev1.OpenResponse, error) {
	return &fridgev1.OpenResponse{IsOpen: req.ToOpen}, nil
}

func (s *FridgeServer) CloseFridge(ctx context.Context, req *fridgev1.CloseRequest) (*fridgev1.CloseResponse, error) {
	return &fridgev1.CloseResponse{IsClosed: req.ToClose}, nil
}

func (s *FridgeServer) AddProduct(ctx context.Context, req *fridgev1.AddProductRequest) (*fridgev1.Product, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "product name is required")
	}

	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be positive")
	}

	product := models.Product{
		Name:       req.Name,
		Quantity:   int(req.Quantity),
		Category:   req.Category,
		ExpiryDate: req.ExpiryDate.AsTime(),
		DateAdded:  time.Now(),
	}

	addedProduct, err := s.service.AddProduct(product)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertToProtoProduct(addedProduct), nil
}

func (s *FridgeServer) GetProduct(ctx context.Context, req *fridgev1.GetProductRequest) (*fridgev1.Product, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "product ID is required")
	}

	product, err := s.service.GetProduct(req.Id)
	if err != nil {
		log.Println("Error getting product: product not found")
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return convertToProtoProduct(product), nil
}

func (s *FridgeServer) ListProducts(_ *emptypb.Empty, stream fridgev1.FridgeService_ListProductsServer) error {
	products, err := s.service.ListProducts()
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, product := range products {
		if err := stream.Send(convertToProtoProduct(product)); err != nil {
			return err
		}
	}

	return nil
}

func (s *FridgeServer) RemoveProduct(ctx context.Context, req *fridgev1.RemoveProductRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "product ID is required")
	}

	if err := s.service.RemoveProduct(req.Id); err != nil {
		log.Println("Error removing product: product not found")
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return &emptypb.Empty{}, nil
}

func (s *FridgeServer) IsExpiredProduct(ctx context.Context, req *fridgev1.IsExpiredProductRequest) (*fridgev1.IsExpiredProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "product ID is required")
	}

	product, isExpired, daysRemaining, err := s.service.CheckProductExpiry(req.Id)
	if err != nil {
		log.Println("Error getting product status: product not found")
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return &fridgev1.IsExpiredProductResponse{
		Product:       convertToProtoProduct(product),
		IsExpired:     isExpired,
		DaysRemaining: int64(daysRemaining),
	}, nil
}

func (s *FridgeServer) GetExpiringProducts(req *fridgev1.ExpiringProductsRequest, stream fridgev1.FridgeService_GetExpiringProductsServer) error {
	if req.DaysThreshold <= 0 {
		return status.Error(codes.InvalidArgument, "days threshold must be positive")
	}

	products, err := s.service.GetExpiringProducts(int(req.DaysThreshold))
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, product := range products {
		if err := stream.Send(convertToProtoProduct(product)); err != nil {
			return err
		}
	}

	return nil
}

func convertToProtoProduct(p models.Product) *fridgev1.Product {
	return &fridgev1.Product{
		Id:         p.ID,
		Name:       p.Name,
		Quantity:   int32(p.Quantity),
		Category:   p.Category,
		ExpiryDate: timestamppb.New(p.ExpiryDate),
		DateAdded:  timestamppb.New(p.DateAdded),
	}
}
