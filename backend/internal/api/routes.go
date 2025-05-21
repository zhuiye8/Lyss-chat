package api

import (
	"github.com/gorilla/mux"
	"github.com/your-org/lyss-chat-2.0/backend/internal/api/auth"
	"github.com/your-org/lyss-chat-2.0/backend/internal/api/chat"
	"github.com/your-org/lyss-chat-2.0/backend/internal/api/model"
	"github.com/your-org/lyss-chat-2.0/backend/internal/api/user"
	"github.com/your-org/lyss-chat-2.0/backend/internal/middleware"
	"github.com/your-org/lyss-chat-2.0/backend/pkg/config"
	"github.com/your-org/lyss-chat-2.0/backend/pkg/db"
	"github.com/your-org/lyss-chat-2.0/backend/pkg/logger"
)

// RegisterRoutes 注册所有 API 路由
func RegisterRoutes(
	r *mux.Router,
	db *db.Postgres,
	redis *db.Redis,
	minio *db.MinIO,
	cfg *config.Config,
	logger *logger.Logger,
) {
	// API 版本前缀
	api := r.PathPrefix("/v1").Subrouter()

	// 认证路由
	authHandler := auth.NewHandler(db, redis, cfg, logger)
	authRoutes := api.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/login", authHandler.Login).Methods("POST")
	authRoutes.HandleFunc("/register", authHandler.Register).Methods("POST")
	authRoutes.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")

	// 需要认证的路由
	authenticated := api.NewRoute().Subrouter()
	authenticated.Use(middleware.Auth(cfg))

	// 用户路由
	userHandler := user.NewHandler(db, minio, cfg, logger)
	userRoutes := authenticated.PathPrefix("/users").Subrouter()
	userRoutes.HandleFunc("/me", userHandler.GetCurrentUser).Methods("GET")
	userRoutes.HandleFunc("", userHandler.ListUsers).Methods("GET")
	userRoutes.HandleFunc("/{id}", userHandler.GetUser).Methods("GET")
	userRoutes.HandleFunc("", userHandler.CreateUser).Methods("POST")
	userRoutes.HandleFunc("/{id}", userHandler.UpdateUser).Methods("PUT")
	userRoutes.HandleFunc("/{id}", userHandler.DeleteUser).Methods("DELETE")

	// 租户路由
	tenantRoutes := authenticated.PathPrefix("/tenants").Subrouter()
	tenantRoutes.HandleFunc("", userHandler.ListTenants).Methods("GET")
	tenantRoutes.HandleFunc("/{id}", userHandler.GetTenant).Methods("GET")
	tenantRoutes.HandleFunc("", userHandler.CreateTenant).Methods("POST")
	tenantRoutes.HandleFunc("/{id}", userHandler.UpdateTenant).Methods("PUT")
	tenantRoutes.HandleFunc("/{id}", userHandler.DeleteTenant).Methods("DELETE")

	// 画布路由
	canvasHandler := chat.NewCanvasHandler(db, cfg, logger)
	canvasRoutes := authenticated.PathPrefix("/canvases").Subrouter()
	canvasRoutes.HandleFunc("", canvasHandler.ListCanvases).Methods("GET")
	canvasRoutes.HandleFunc("/{id}", canvasHandler.GetCanvas).Methods("GET")
	canvasRoutes.HandleFunc("", canvasHandler.CreateCanvas).Methods("POST")
	canvasRoutes.HandleFunc("/{id}", canvasHandler.UpdateCanvas).Methods("PUT")
	canvasRoutes.HandleFunc("/{id}", canvasHandler.DeleteCanvas).Methods("DELETE")

	// 消息路由
	messageHandler := chat.NewMessageHandler(db, cfg, logger)
	messageRoutes := authenticated.PathPrefix("/canvases/{id}/messages").Subrouter()
	messageRoutes.HandleFunc("", messageHandler.ListMessages).Methods("GET")
	messageRoutes.HandleFunc("", messageHandler.SendMessage).Methods("POST")
	messageRoutes.HandleFunc("/stream", messageHandler.StreamMessage).Methods("POST")

	// 模型路由
	modelHandler := model.NewModelHandler(db, cfg, logger)
	modelRoutes := authenticated.PathPrefix("/models").Subrouter()
	modelRoutes.HandleFunc("", modelHandler.ListModels).Methods("GET")
	modelRoutes.HandleFunc("/{id}", modelHandler.GetModel).Methods("GET")
	modelRoutes.HandleFunc("", modelHandler.CreateModel).Methods("POST")
	modelRoutes.HandleFunc("/{id}", modelHandler.UpdateModel).Methods("PUT")
	modelRoutes.HandleFunc("/{id}", modelHandler.DeleteModel).Methods("DELETE")

	// 提供商路由
	providerRoutes := authenticated.PathPrefix("/providers").Subrouter()
	providerRoutes.HandleFunc("", modelHandler.ListProviders).Methods("GET")
	providerRoutes.HandleFunc("/{id}", modelHandler.GetProvider).Methods("GET")
	providerRoutes.HandleFunc("", modelHandler.CreateProvider).Methods("POST")
	providerRoutes.HandleFunc("/{id}", modelHandler.UpdateProvider).Methods("PUT")
	providerRoutes.HandleFunc("/{id}", modelHandler.DeleteProvider).Methods("DELETE")

	// API 密钥路由
	apiKeyRoutes := authenticated.PathPrefix("/providers/{id}/api-keys").Subrouter()
	apiKeyRoutes.HandleFunc("", modelHandler.ListAPIKeys).Methods("GET")
	apiKeyRoutes.HandleFunc("/{key_id}", modelHandler.GetAPIKey).Methods("GET")
	apiKeyRoutes.HandleFunc("", modelHandler.CreateAPIKey).Methods("POST")
	apiKeyRoutes.HandleFunc("/{key_id}", modelHandler.UpdateAPIKey).Methods("PUT")
	apiKeyRoutes.HandleFunc("/{key_id}", modelHandler.DeleteAPIKey).Methods("DELETE")
}
