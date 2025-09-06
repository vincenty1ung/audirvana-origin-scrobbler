package api

import (
	"context"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/vincenty1ung/lastfm-scrobbler/config"
	"github.com/vincenty1ung/lastfm-scrobbler/core/log"
	"github.com/vincenty1ung/lastfm-scrobbler/internal/logic/analysis"
	"github.com/vincenty1ung/lastfm-scrobbler/internal/logic/track"
	"github.com/vincenty1ung/lastfm-scrobbler/internal/model"
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
	trackService := track.NewTrackService()
	r.GET(
		"/api/track-play-counts", func(c *gin.Context) {
			limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
			offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

			if limit > 100 {
				limit = 100 // Limit max records per page
			}

			records, err := trackService.GetTrackPlayCounts(c.Request.Context(), limit, offset)
			log.Info(
				c.Request.Context(), "Fetched track play counts", zap.Int("count", len(records)),
				zap.Int("limit", limit), zap.Int("offset", offset),
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Check if client expects HTML response
			acceptHeader := c.GetHeader("Accept")
			if strings.Contains(acceptHeader, "text/html") || c.Query("format") == "html" {
				// Load HTML template
				tmplPath := filepath.Join("templates", "track_play_counts.html")
				tmpl, err := template.New("track_play_counts.html").Funcs(
					template.FuncMap{
						"addOne": func(i int) int {
							return i + 1
						},
						"add": func(a, b int) int {
							return a + b
						},
						"subtract": func(a, b int) int {
							return a - b
						},
					},
				).ParseFiles(tmplPath)
				if err != nil {
					log.Error(c.Request.Context(), "Failed to parse template", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load template"})
					return
				}

				// Execute template with records data
				data := struct {
					Records     []*model.TrackPlayCount
					Limit       int
					Offset      int
					RecordCount int
				}{
					Records:     records,
					Limit:       limit,
					Offset:      offset,
					RecordCount: len(records),
				}

				// Set content type and write HTML response
				c.Header("Content-Type", "text/html; charset=utf-8")
				if err := tmpl.Execute(c.Writer, data); err != nil {
					log.Error(c.Request.Context(), "Failed to execute template", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template"})
					return
				}
			} else {
				// Return JSON response for API clients
				c.JSON(http.StatusOK, records)
			}
		},
	)

	// Get play count for a specific track
	r.GET(
		"/api/track-play-count", func(c *gin.Context) {
			artist := c.Query("artist")
			album := c.Query("album")
			trackName := c.Query("trackName")

			if artist == "" || album == "" || trackName == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "artist, album, and trackName are required"})
				return
			}

			record, err := trackService.GetTrackPlayCount(c.Request.Context(), artist, album, trackName)
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

	// Generate music preference report
	musicAnalysisService := analysis.NewMusicAnalysisService()
	r.GET(
		"/api/music-analysis/report", func(c *gin.Context) {
			// Create a background context for the report generation
			ctx := c.Request.Context()

			// Generate the report data
			reportData, err := musicAnalysisService.GenerateMusicPreferenceReport(ctx)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Load HTML template
			tmplPath := filepath.Join("templates", "report.html")
			tmpl, err := template.New("report.html").Funcs(
				template.FuncMap{
					"addOne": func(i int) int {
						return i + 1
					},
				},
			).ParseFiles(tmplPath)
			if err != nil {
				log.Error(ctx, "Failed to parse template", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load template"})
				return
			}

			// Set content type and write HTML response
			c.Header("Content-Type", "text/html; charset=utf-8")
			if err := tmpl.Execute(c.Writer, reportData); err != nil {
				log.Error(ctx, "Failed to execute template", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template"})
				return
			}
		},
	)

	// Generate music recommendations
	r.GET(
		"/api/music-analysis/recommendations", func(c *gin.Context) {
			// Create a background context for the recommendation generation
			ctx := c.Request.Context()

			// Generate recommendations
			recommendations, err := musicAnalysisService.GenerateRecommendations(ctx)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Load HTML template
			tmplPath := filepath.Join("templates", "recommendations.html")
			tmpl, err := template.New("recommendations.html").Funcs(
				template.FuncMap{
					"addOne": func(i int) int {
						return i + 1
					},
				},
			).ParseFiles(tmplPath)
			if err != nil {
				log.Error(ctx, "Failed to parse template", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load template"})
				return
			}

			// Execute template with recommendations data
			data := struct {
				Recommendations []analysis.MusicRecommendation
			}{
				Recommendations: recommendations,
			}

			// Set content type and write HTML response
			c.Header("Content-Type", "text/html; charset=utf-8")
			if err := tmpl.Execute(c.Writer, data); err != nil {
				log.Error(ctx, "Failed to execute template", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template"})
				return
			}
		},
	)

	// Health check endpoint
	r.GET(
		"/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		},
	)

	// 首页
	r.GET(
		"/", func(c *gin.Context) {
			// Load HTML template
			tmplPath := filepath.Join("templates", "index.html")
			tmpl, err := template.New("index.html").ParseFiles(tmplPath)
			if err != nil {
				log.Error(c.Request.Context(), "Failed to parse template", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load template"})
				return
			}

			// Set content type and write HTML response
			c.Header("Content-Type", "text/html; charset=utf-8")
			if err := tmpl.Execute(c.Writer, nil); err != nil {
				log.Error(c.Request.Context(), "Failed to execute template", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template"})
				return
			}
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
	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
