package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	httpNet "net/http"
	v1 "service/api/gobang/v1"
	"service/internal/conf"
	"service/internal/pkg/auth"
	"service/internal/service"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, jwtc *conf.JWT, greeter *service.GobangService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.ResponseEncoder(responseEncoder),
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(),
			selector.Server(auth.JWTAuth(jwtc.Secret)).Match(NewSkipRoutersMatcher()).Build(),
			logging.Server(logger),
		),
		http.Filter(
			handlers.CORS(
				handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "X-Session-Id", "Access-Control-Allow-Origin:http://127.0.0.1:8081/"}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
				handlers.AllowedOrigins([]string{"http://127.0.0.1:8081"}),
				handlers.AllowCredentials(),
			),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)
	v1.RegisterGobangHTTPServer(srv, greeter)
	return srv
}

func responseEncoder(w httpNet.ResponseWriter, r *httpNet.Request, i interface{}) error {
	type response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	reply := &response{
		Code:    0,
		Message: "success",
		Data:    i,
	}
	codec := encoding.GetCodec("json")
	data, err := codec.Marshal(reply)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	return nil
}

func NewSkipRoutersMatcher() selector.MatchFunc {
	skipRouters := map[string]struct{}{
		"/gobang.v1.Gobang/Login":        {},
		"/gobang.v1.Gobang/MemberStatus": {},
	}

	return func(ctx context.Context, operation string) bool {
		if _, ok := skipRouters[operation]; ok {
			return false
		}
		return true
	}
}
