package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/pkg/api/apiv1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"strings"
)

const MaxBodyBytes = 64000

type serverInterface interface {
	apiv1.ServerInterface
}

var _ serverInterface = (*Handler)(nil)

func New(usersRepository usersClient, linksRepository linksClient) *Handler {
	return &Handler{usersHandler: newUsersHandler(usersRepository), linksHandler: newLinksHandler(linksRepository)}
}

type Handler struct {
	*usersHandler
	*linksHandler
}

func handleGRPCError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	st := status.Convert(err)
	code := st.Code()
	w.WriteHeader(ConvertGRPCCodeToHTTP(code))
	if err := json.NewEncoder(w).Encode(
		apiv1.Error{
			Code:    ConvertGRPCToErrorCode(code),
			Message: nil,
		},
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ConvertGRPCCodeToHTTP(grpcCode codes.Code) int {
	switch grpcCode {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func ConvertGRPCToErrorCode(grpcCode codes.Code) apiv1.ErrorCode {
	switch grpcCode {
	case codes.Internal, codes.Unknown, codes.DataLoss:
		return apiv1.InternalServerError
	case codes.NotFound:
		return apiv1.NotFound
	case codes.InvalidArgument, codes.FailedPrecondition, codes.OutOfRange:
		return apiv1.BadRequest
	case codes.Aborted, codes.AlreadyExists:
		return apiv1.Conflict
	}

	return apiv1.InternalServerError
}

func ConvertHTTPToErrorCode(code int) apiv1.ErrorCode {
	switch code {
	case http.StatusBadRequest:
		return apiv1.BadRequest
	case http.StatusInternalServerError:
		return apiv1.InternalServerError
	case http.StatusRequestEntityTooLarge:
		return apiv1.BadRequest
	case http.StatusUnsupportedMediaType:
		return apiv1.BadRequest
	case http.StatusConflict:
		return apiv1.Conflict
	}
	return apiv1.InternalServerError
}

func MarshalResponse(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, "%s", data)
}

func Unmarshal(w http.ResponseWriter, r *http.Request, data interface{}) (int, error) {
	if t := r.Header.Get("content-type"); len(t) < 16 || t[:16] != "application/json" {
		return http.StatusUnsupportedMediaType, fmt.Errorf("content-type is not application/json")
	}

	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	if err := d.Decode(&data); err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalError *json.UnmarshalTypeError
		switch {
		case errors.As(err, &syntaxErr):
			return http.StatusBadRequest, fmt.Errorf("malformed json at position %d", syntaxErr.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return http.StatusBadRequest, fmt.Errorf("malformed json")
		case errors.As(err, &unmarshalError):
			return http.StatusBadRequest, fmt.Errorf(
				"invalid value %q at position %d", unmarshalError.Field, unmarshalError.Offset,
			)
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return http.StatusBadRequest, fmt.Errorf("unknown field %s", fieldName)
		case errors.Is(err, io.EOF):
			return http.StatusBadRequest, fmt.Errorf("body must not be empty")
		case err.Error() == "http: request body too large":
			return http.StatusRequestEntityTooLarge, err
		default:
			return http.StatusInternalServerError, fmt.Errorf("failed to decode json: %w", err)
		}
	}

	if d.More() {
		return http.StatusBadRequest, fmt.Errorf("body must contain only one JSON object")
	}

	return http.StatusOK, nil
}
