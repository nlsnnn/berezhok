package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type webhookMiddleware struct {
	allowedCIDRs []string
	log          *slog.Logger
}

func NewWebhookMiddleware(allowedCIDRs []string, log *slog.Logger) *webhookMiddleware {
	return &webhookMiddleware{
		allowedCIDRs: allowedCIDRs,
		log:          log,
	}
}

// Middleware для проверки разрешенных IP-адресов.
func (wm *webhookMiddleware) IPFilterMiddleware(next http.Handler) http.Handler {
	log := wm.log.With(slog.String("op", "webhookMiddleware.IPFilterMiddleware"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteIP := r.RemoteAddr
		log.Info("Initial remote IP", slog.String("ip", remoteIP))

		// Проверяем X-Real-IP, если доступен.
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			log.Info("Using X-Real-IP header", slog.String("ip", realIP))
			remoteIP = realIP
		}

		// Разделяем адрес на хост и порт.
		var host string
		if strings.Contains(remoteIP, ":") {
			var err error
			host, _, err = net.SplitHostPort(remoteIP)
			if err != nil {
				response.BadRequest(w, "Invalid remote IP address")
				return
			}
		} else {
			host = remoteIP
		}

		// Проверяем, разрешен ли IP-адрес.
		if !IsIPAllowed(host, wm.allowedCIDRs) {
			response.Forbidden(w, "Access denied: IP not allowed")
			return
		}

		// Передаем управление дальше.
		next.ServeHTTP(w, r)
	})
}

// Проверяет, входит ли IP-адрес в разрешенные диапазоны.
func IsIPAllowed(ip string, allowedCIDRs []string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, cidr := range allowedCIDRs {
		_, allowedNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if allowedNet.Contains(parsedIP) {
			return true
		}
	}
	return false
}
