package apiserver

import (
	"ShortURL/internal/app/store"
	"ShortURL/pkg/api"
	"ShortURL/pkg/grpcserver"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func getConn() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	store := store.NewStoreRedis(client)
	server := grpc.NewServer()
	srv := &grpcserver.GRPCServer{Store: store}
	api.RegisterShortlinkServer(server, srv)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestServer_HandleCreate(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(getConn()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := api.NewShortlinkClient(conn)
	s := newServer(c)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"url": "google.com",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid url",
			payload: map[string]interface{}{
				"url": "googlecom",
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/create", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleGet(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(getConn()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := api.NewShortlinkClient(conn)
	s := newServer(c)

	b := &bytes.Buffer{}
	payload := map[string]interface{}{
		"url": "google.com",
	}

	json.NewEncoder(b).Encode(payload)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/XXXaaa123_", nil)
	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}