package main

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/olvrng/rbot/be/cmd/rbot-server/config"
	"github.com/olvrng/rbot/be/com/flowdef/service"
	flowdefstore "github.com/olvrng/rbot/be/com/flowdef/store"
	flowexecservice "github.com/olvrng/rbot/be/com/flowexec/service"
	"github.com/olvrng/rbot/be/com/flowexec/store"
	"github.com/olvrng/rbot/be/com/integration/fbmsg"
	"github.com/olvrng/rbot/be/com/integration/webhook"
	"github.com/olvrng/rbot/be/pkg/httprpc"
	"github.com/olvrng/rbot/be/pkg/l"
	"github.com/olvrng/rbot/be/pkg/lifecycle"
)

var ll = l.New()
var ls = ll.Sugar()

func main() {
	ll.Debug("enable debug log")

	// load config
	initFlags()
	cfg, err := config.Load(flConfigFile)
	ll.Must("can not load config", err)

	// shutdown
	ctx, ctxCancel := context.WithCancel(context.Background())
	lifecycle.ListenForSignal(ctxCancel, 5*time.Second)

	// messenger client, webhook
	msgClient, err := fbmsg.NewClient(fbmsg.Config(cfg.Messenger))
	ll.Must("can not create messenger client", err)

	// build server
	mux := chi.NewMux()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	msgWebhook := buildAPIServer(cfg, mux, msgClient)
	mux.Get("/api/webhook/messenger", msgWebhook.HandleVerification)
	mux.Post("/api/webhook/messenger", msgWebhook.HandleWebhook)

	httpMux := http.NewServeMux()
	httpMux.Handle("/api/", mux)

	// static file
	fileHandler := http.FileServer(http.Dir(cfg.StaticPath))
	httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.Contains(path, ".") {
			fileHandler.ServeHTTP(w, r)
		} else {
			http.ServeFile(w, r, filepath.Join(cfg.StaticPath, "index.html"))
		}
	})

	// start http server
	httpServer := &http.Server{
		Addr:    cfg.HTTP.ListeningAddress(),
		Handler: httpMux,
	}
	ll.Info("server is listening at " + cfg.HTTP.ListeningAddress())
	go func() {
		err2 := httpServer.ListenAndServe()
		if err2 != http.ErrServerClosed {
			ll.Error("serving http", l.Error(err2))
		}
	}()

	<-ctx.Done()
	ll.Info("server is shutting down...")
}

func buildAPIServer(cfg config.Config, m *chi.Mux, msgClient *fbmsg.Client) *webhook.WebhookService {
	flowStore, err := flowdefstore.NewFlowFileStore(flFlowFile)
	ll.Must("can not open flow data file", err)
	stateStore, err := store.NewFlowStateStore(flStateFile)
	ll.Must("can not open state data file", err)

	actionExec := flowexecservice.NewActionExecutor(msgClient)
	flowService := service.NewFlowEditorService(flowStore)
	flowQuery := service.NewFlowQueryService(flowStore)
	orderService := flowexecservice.NewOrderService(flowQuery, stateStore, actionExec)
	messengerService := flowexecservice.NewMessengerService(flowQuery, stateStore, actionExec)
	msgWebhook := webhook.NewWebhookService(msgClient, cfg.Messenger.VerifyToken, messengerService)

	servers := httprpc.MustNewServers(flowService, orderService, messengerService)
	for _, s := range servers {
		m.Handle(s.PathPrefix()+"*", s)
	}
	return msgWebhook
}
