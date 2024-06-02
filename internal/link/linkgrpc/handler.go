package linkgrpc

import (
	"context"
	"errors"
	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/internal/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/pkg/pb"
)

var _ pb.LinkServiceServer = (*Handler)(nil)

func New(linksRepository linksRepository, timeout time.Duration) *Handler {
	return &Handler{linksRepository: linksRepository, timeout: timeout}
}

type Handler struct {
	pb.UnimplementedLinkServiceServer
	linksRepository linksRepository
	timeout         time.Duration
}

func (h Handler) GetLinkByUserID(ctx context.Context, id *pb.GetLinksByUserId) (*pb.ListLinkResponse, error) {
	// implemented
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	list, err := h.linksRepository.FindByUserID(ctx, id.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := make([]*pb.Link, 0, len(list))

	for _, l := range list {
		response = append(
			response, &pb.Link{
				Id:        l.ID.Hex(),
				Title:     l.Title,
				Url:       l.URL,
				Images:    l.Images,
				Tags:      l.Tags,
				UserId:    l.UserID,
				CreatedAt: l.CreatedAt.Format(time.RFC3339),
				UpdatedAt: l.UpdatedAt.Format(time.RFC3339),
			},
		)
	}

	return &pb.ListLinkResponse{Links: response}, nil
}

func (h Handler) mustEmbedUnimplementedLinkServiceServer() {
	// nothing to implement here
}

func (h Handler) CreateLink(ctx context.Context, request *pb.CreateLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	// implemented
	objectID, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err := h.linksRepository.Create(
		ctx, database.CreateLinkReq{
			ID:     objectID,
			URL:    request.Url,
			Title:  request.Title,
			Tags:   request.Tags,
			Images: request.Images,
			UserID: request.UserId,
		},
	); err != nil {
		if errors.Is(err, errors.New("conflict")) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Empty{}, nil
}

func (h Handler) GetLink(ctx context.Context, request *pb.GetLinkRequest) (*pb.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	// implemented
	objectID, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	l, err := h.linksRepository.FindByID(ctx, objectID)
	if err != nil {
		if errors.Is(err, errors.New("not found")) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Link{
		Id:        l.ID.Hex(),
		Title:     l.Title,
		Url:       l.URL,
		Images:    l.Images,
		Tags:      l.Tags,
		UserId:    l.UserID,
		CreatedAt: l.CreatedAt.Format(time.RFC3339),
		UpdatedAt: l.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (h Handler) UpdateLink(ctx context.Context, request *pb.UpdateLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	// implemented
	objectID, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err := h.linksRepository.Update(
		ctx, database.UpdateLinkReq{
			ID:     objectID,
			URL:    request.Url,
			Title:  request.Title,
			Tags:   request.Tags,
			Images: request.Images,
			UserID: request.UserId,
		},
	); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Empty{}, nil
}

func (h Handler) DeleteLink(ctx context.Context, request *pb.DeleteLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	// implemented
	objectID, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.linksRepository.Delete(ctx, objectID); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Empty{}, nil
}

func (h Handler) ListLinks(ctx context.Context, request *pb.Empty) (*pb.ListLinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	// implemented
	list, err := h.linksRepository.FindAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := make([]*pb.Link, 0, len(list))

	for _, l := range list {
		response = append(
			response, &pb.Link{
				Id:        l.ID.Hex(),
				Title:     l.Title,
				Url:       l.URL,
				Images:    l.Images,
				Tags:      l.Tags,
				UserId:    l.UserID,
				CreatedAt: l.CreatedAt.Format(time.RFC3339),
				UpdatedAt: l.UpdatedAt.Format(time.RFC3339),
			},
		)
	}

	return &pb.ListLinkResponse{Links: response}, nil
}
