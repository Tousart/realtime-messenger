package wsmw

import (
	"log/slog"
	"time"

	"github.com/tousart/messenger/pkg/types/wstypes"
)

func Logging(logger *slog.Logger) func(wstypes.Method) wstypes.Method {
	return func(next wstypes.Method) wstypes.Method {
		return func(metadata *wstypes.Metadata, cw *wstypes.ConnWriter, req *wstypes.Request) {
			start := time.Now()

			next(metadata, cw, req)

			logger.Info("ws call",
				slog.String("method", req.Method),
				slog.Int64("user_id", metadata.UserID),
				slog.Duration("duration", time.Since(start)),
			)
		}
	}
}
