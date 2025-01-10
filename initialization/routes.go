package initialization

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/global"
	"watchAlert/internal/middleware"
	"watchAlert/internal/routers"
	"watchAlert/internal/routers/v1"
)

func InitRoute() {
	global.Logger.Sugar().Info("服务启动")

	mode := global.Config.Server.Mode
	if mode == "" {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	ginEngine := gin.New()

	ginEngine.Use(
		// 启用CORS中间件
		middleware.Cors(),
		// 自定义请求日志格式
		gin.LoggerWithFormatter(middleware.RequestLoggerFormatter),
		// 处理异常panic，go template解析错误的模板时有可能会出现
		gin.Recovery(),
	)
	allRouter(ginEngine)

	err := ginEngine.Run(":" + global.Config.Server.Port)
	if err != nil {
		global.Logger.Sugar().Error("服务启动失败:", err)
		return
	}
}

func allRouter(engine *gin.Engine) {

	routers.HealthCheck(engine)
	v1.Router(engine)

}
