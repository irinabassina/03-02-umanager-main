package v1

import (
	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/pkg/api/apiv1"
	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/pkg/pb"
	"net/http"
)

func newUsersHandler(usersClient usersClient) *usersHandler {
	return &usersHandler{client: usersClient}
}

type usersHandler struct {
	client usersClient
}

func (h *usersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// implemented
	ctx := r.Context()
	resp, err := h.client.ListUsers(ctx, &pb.Empty{})
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	userList := make([]apiv1.User, 0, len(resp.Users))
	for _, u := range resp.Users {
		userList = append(
			userList, apiv1.User{
				CreatedAt: u.CreatedAt,
				Id:        u.Id,
				Password:  u.Password,
				UpdatedAt: u.UpdatedAt,
				Username:  u.Username,
			},
		)
	}

	MarshalResponse(w, http.StatusOK, userList)
}

func (h *usersHandler) PostUsers(w http.ResponseWriter, r *http.Request) {
	// implemented
	ctx := r.Context()

	var u apiv1.UserCreate
	code, err := Unmarshal(w, r, &u)
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

	if _, err := h.client.CreateUser(
		ctx, &pb.CreateUserRequest{
			Id:       u.Id,
			Username: u.Username,
			Password: u.Password,
		},
	); err != nil {
		handleGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *usersHandler) DeleteUsersId(w http.ResponseWriter, r *http.Request, id string) {
	// implemented
	ctx := r.Context()

	if _, err := h.client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id}); err != nil {
		handleGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *usersHandler) GetUsersId(w http.ResponseWriter, r *http.Request, id string) {
	// implemented
	ctx := r.Context()

	u, err := h.client.GetUser(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	MarshalResponse(
		w, http.StatusOK, apiv1.User{
			CreatedAt: u.CreatedAt,
			Id:        u.Id,
			Password:  u.Password,
			UpdatedAt: u.UpdatedAt,
			Username:  u.Username,
		},
	)
}

func (h *usersHandler) PutUsersId(w http.ResponseWriter, r *http.Request, id string) {
	// implemented
	ctx := r.Context()

	var u apiv1.UserCreate
	code, err := Unmarshal(w, r, &u)
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

	if _, err := h.client.UpdateUser(
		ctx, &pb.UpdateUserRequest{
			Id:       u.Id,
			Username: u.Username,
			Password: u.Password,
		},
	); err != nil {
		handleGRPCError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
