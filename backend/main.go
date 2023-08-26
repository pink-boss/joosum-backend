package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
	"joosum-backend/pkg/config"
	"joosum-backend/pkg/middleware"
	"joosum-backend/pkg/routes"
	"joosum-backend/pkg/util"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

// @title Joosum App
// @description This is API Docs for Joosum App.
// @termsOfService http://swagger.io/terms/
// @contact.name Pink boss
// @contact.email pinkjoosum@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /api
func main() {
	config.EnvConfig()

	util.StartMongoDB()
	util.LoadApplePublicKeys()
	util.Validate = validator.New()

	router := config.GetRouter()
	if *config.Env == "prod" {
		router.Use(middleware.LoggingMiddleware()) // 커스텀 로깅 미들웨어 적용
	}

	routes.PublicRoutes(router)
	routes.SwaggerRoutes(router)

	// SwaggerRoutes 보다 위에 있으면 swagger 문서가 보이지 않음
	routes.PrivateRoutes(router)

	// http.Server 인스턴스 생성

	server := &http.Server{
		Addr:    ":5001",
		Handler: router,
	}

	// 서버를 고루틴에서 시작
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 서버를 꺼도 활성 연결을 중단시키지 않고 서버를 정상적으로 종료합니다.
	// Shutdown은 먼저 모든 오픈 리스너를 닫고, 그 다음 모든 유휴 연결을 닫고, 그리고 연결이 유휴 상태로 돌아올 때까지 무기한으로 대기한 다음 종료하도록 작동합니다.
	// 제공된 context가 종료가 완료되기 전에 만료되면, Shutdown은 context의 오류를 반환하고, 그렇지 않으면 서버의 기본 리스너를 닫는데서 발생한 오류를 반환합니다.

	// Shutdown이 호출되면, Serve, ListenAndServe, 그리고 ListenAndServeTLS는 즉시 ErrServerClosed를 반환합니다.
	// 프로그램이 종료되지 않고 대신 Shutdown이 반환될 때까지 기다리도록 합니다.

	// Shutdown은 WebSockets과 같은 인계된(hijacked) 연결을 닫거나 대기하려고 시도하지 않습니다.
	// Shutdown의 호출자는, 원한다면, 종료를 알리고 그것들이 닫힐 때까지 기다리기 위해 그러한 오래 지속되는 연결을 별도로 알려야 합니다.
	// 종료 알림 기능을 등록하는 방법으로 RegisterOnShutdown을 참조하십시오.

	if err := server.Shutdown(ctx); err != nil {
		// 예상치 못한 오류로 서버 셧다운 실패
		panic(err)
	}

	// mongoDB 종료
	util.CloseMongoDB()

}
