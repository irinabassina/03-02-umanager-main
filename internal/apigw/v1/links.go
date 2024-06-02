package v1

import (
	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/pkg/api/apiv1"
	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/pkg/pb"
	"net/http"
)

func newLinksHandler(linksClient linksClient) *linksHandler {
	return &linksHandler{client: linksClient}
}

type linksHandler struct {
	client linksClient
}

func (h *linksHandler) GetLinks(w http.ResponseWriter, r *http.Request) {
	// implemented
	ctx := r.Context()

	resp, err := h.client.ListLinks(ctx, nil)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	linkList := make([]apiv1.Link, 0, len(resp.Links))
	for _, l := range resp.Links {
		linkList = append(
			linkList, apiv1.Link{
				CreatedAt: l.CreatedAt,
				Id:        l.Id,
				Images:    l.Images,
				Tags:      l.Tags,
				Title:     l.Title,
				UpdatedAt: l.UpdatedAt,
				Url:       l.Url,
				UserId:    l.UserId,
			},
		)
	}

	MarshalResponse(w, http.StatusOK, linkList)
}

func (h *linksHandler) PostLinks(w http.ResponseWriter, r *http.Request) {
	// implemented
	ctx := r.Context()

	var l apiv1.LinkCreate
	code, err := Unmarshal(w, r, &l)
	if err != nil {
		errStr := err.Error()
		MarshalResponse(
			w, code, apiv1.Error{
				Code:    ConvertHTTPToErrorCode(code),
				Message: &errStr,
			},
		)
		return
	}

	if _, err := h.client.CreateLink(
		ctx, &pb.CreateLinkRequest{
			Id:     l.Id,
			Title:  l.Title,
			Url:    l.Url,
			Images: l.Images,
			Tags:   l.Tags,
			UserId: l.UserId,
		},
	); err != nil {
		handleGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *linksHandler) DeleteLinksId(w http.ResponseWriter, r *http.Request, id string) {
	// implemented
	ctx := r.Context()

	if _, err := h.client.DeleteLink(ctx, &pb.DeleteLinkRequest{Id: id}); err != nil {
		handleGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *linksHandler) GetLinksId(w http.ResponseWriter, r *http.Request, id string) {
	// implemented
	ctx := r.Context()

	link, err := h.client.GetLink(ctx, &pb.GetLinkRequest{Id: id})
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	MarshalResponse(
		w, http.StatusOK, apiv1.Link{
			CreatedAt: link.CreatedAt,
			Id:        link.Id,
			Images:    link.Images,
			Tags:      link.Tags,
			Title:     link.Title,
			UpdatedAt: link.UpdatedAt,
			Url:       link.Url,
			UserId:    link.UserId,
		},
	)
}

func (h *linksHandler) PutLinksId(w http.ResponseWriter, r *http.Request, id string) {
	// implemented
	ctx := r.Context()
	var l apiv1.LinkCreate
	code, err := Unmarshal(w, r, &l)
	if err != nil {
		errStr := err.Error()
		MarshalResponse(
			w, code, apiv1.Error{
				Code:    ConvertHTTPToErrorCode(code),
				Message: &errStr,
			},
		)
		return
	}

	if _, err := h.client.UpdateLink(
		ctx, &pb.UpdateLinkRequest{
			Id:     l.Id,
			Title:  l.Title,
			Url:    l.Url,
			Images: l.Images,
			Tags:   l.Tags,
			UserId: l.UserId,
		},
	); err != nil {
		handleGRPCError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *linksHandler) GetLinksUserUserID(w http.ResponseWriter, r *http.Request, userID string) {
	// implemented
	ctx := r.Context()
	resp, err := h.client.GetLinkByUserID(ctx, &pb.GetLinksByUserId{UserId: userID})
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	linkList := make([]apiv1.Link, 0, len(resp.Links))
	for _, l := range resp.Links {
		linkList = append(
			linkList, apiv1.Link{
				CreatedAt: l.CreatedAt,
				Id:        l.Id,
				Images:    l.Images,
				Tags:      l.Tags,
				Title:     l.Title,
				UpdatedAt: l.UpdatedAt,
				Url:       l.Url,
				UserId:    l.UserId,
			},
		)
	}

	MarshalResponse(w, http.StatusOK, linkList)
}
