package httprpc

import (
	"net/http"
	"path"
	"strings"
)

var _ Server = &stripPrefix{}

type stripPrefix struct {
	http.Handler
	prefix  string
	servers []Server
}

func StripPrefix(prefix string, servers ...Server) Server {
	mux := http.NewServeMux()
	for _, s := range servers {
		mux.Handle(
			path.Join(prefix+s.PathPrefix()),
			http.StripPrefix(prefix, s),
		)
	}
	w := &stripPrefix{
		Handler: mux,
		prefix:  prefix,
		servers: servers,
	}
	return w
}

func (s *stripPrefix) PathPrefix() string {
	return s.prefix
}

func (s *stripPrefix) WithHooks(builder HooksBuilder) Server {
	servers := WithHooks(s.servers, builder)
	return StripPrefix(s.prefix, servers...)
}

var _ Server = &withPrefix{}

type withPrefix struct {
	http.Handler
	prefix string
	server Server
}

func WithPrefix(prefix string, servers []Server) []Server {
	result := make([]Server, len(servers))
	for i, s := range servers {
		urlPrefix := path.Join(prefix, s.PathPrefix())
		if !strings.HasSuffix(urlPrefix, "/") {
			urlPrefix = urlPrefix + "/"
		}
		result[i] = &withPrefix{
			Handler: http.StripPrefix(strings.TrimSuffix(prefix, "/"), s),
			prefix:  urlPrefix,
			server:  s,
		}
	}
	return result
}

func (w *withPrefix) PathPrefix() string {
	return w.prefix
}

func (w *withPrefix) WithHooks(builder HooksBuilder) Server {
	w2 := *w
	w2.server = w2.server.WithHooks(builder)
	return &w2
}
