package websocket

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/itohin/gophkeeper/internal/server/events"
	"log"
	"net/http"
)

type WSNotifier struct {
	srv          *http.Server
	certFilePath string
	keyFilePath  string
}

func NewWSNotifier(address, certFilePath, keyFilePath string, secretEventsCh chan *events.SecretEvent) *WSNotifier {
	tlsCfg := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		InsecureSkipVerify: true,
	}
	srv := &http.Server{
		Addr:      address,
		TLSConfig: tlsCfg,
		Handler:   NewRouter(NewHub(secretEventsCh)),
	}
	return &WSNotifier{
		srv:          srv,
		certFilePath: certFilePath,
		keyFilePath:  keyFilePath,
	}
}

func (ws *WSNotifier) Run() error {
	err := ws.srv.ListenAndServeTLS(ws.certFilePath, ws.keyFilePath)
	if err != nil {
		log.Printf("ws error: %v", err)
	}
	return err
}

func (ws *WSNotifier) Stop(ctx context.Context) {
	if err := ws.srv.Shutdown(ctx); err != nil {
		log.Printf("WS server Shutdown error: %v", err)
	}
}

type Router struct {
	*http.ServeMux
	hub *Hub
}

func NewRouter(hub *Hub) *Router {
	r := &Router{
		ServeMux: http.NewServeMux(),
		hub:      hub,
	}
	r.HandleFunc("/connect", r.connect)
	return r
}

func (rt *Router) connect(w http.ResponseWriter, r *http.Request) {
	if err := rt.handleConn(w, r); err != nil {
		log.Println(err)
	}
}

func (rt *Router) handleConn(w http.ResponseWriter, r *http.Request) error {
	params := r.URL.Query()
	clientId := params.Get("user_id")
	deviceId := params.Get("finger_print")
	if len(clientId) == 0 || len(deviceId) == 0 {
		return fmt.Errorf("no required params to connect to ws")
	}
	fmt.Println("params: ", clientId, deviceId)
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return fmt.Errorf("failed to upgrade HTTP connection: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("failed to close ws connection: %w", err)
		}
	}()
	done, err := rt.hub.Connect(conn, clientId, deviceId)
	if err != nil {
		return fmt.Errorf("failed to connect user: %v", err)
	}
	<-done
	return nil
}
