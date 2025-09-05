package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/lastfm-scrobbler/config"
	"github.com/lastfm-scrobbler/log"
	"github.com/lastfm-scrobbler/model"
)

func setupRouter(name string) *gin.Engine {
	r := gin.Default()

	// Add OpenTelemetry middleware
	r.Use(
		otelgin.Middleware(name),
		func(c *gin.Context) {
			traceID := trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()
			c.Header("Trace-Id", traceID)
			c.Next()
		},
	)

	// Get track play counts with pagination
	r.GET(
		"/api/track-play-counts", func(c *gin.Context) {
			limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
			offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

			if limit > 100 {
				limit = 100 // Limit max records per page
			}

			records, err := model.GetTrackPlayCounts(c.Request.Context(), limit, offset)
			log.Info(
				c.Request.Context(), "Fetched track play counts", zap.Int("count", len(records)),
				zap.Int("limit", limit), zap.Int("offset", offset),
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, records)
		},
	)

	// Get play count for a specific track
	r.GET(
		"/api/track-play-count", func(c *gin.Context) {
			artist := c.Query("artist")
			album := c.Query("album")
			track := c.Query("track")

			if artist == "" || album == "" || track == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "artist, album, and track are required"})
				return
			}

			record, err := model.GetTrackPlayCount(c.Request.Context(), artist, album, track)
			if err != nil {
				if err.Error() == "record not found" {
					c.JSON(http.StatusOK, gin.H{"play_count": 0})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, record)
		},
	)

	// Health check endpoint
	r.GET(
		"/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		},
	)

	return r
}

func StartHTTPServer(ctx context.Context, name string) {
	r := setupRouter(name)
	port := config.ConfigObj.HTTP.Port
	if port == "" {
		port = "8080" // Default port
	}
	log.Info(ctx, "Starting HTTP server on port", zap.String("port", port))
	r.Run(":" + port)
}
