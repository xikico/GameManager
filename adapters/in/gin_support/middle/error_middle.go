package middle

import (
	"GameManager/adapters/in/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func GlobalErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 捕获 panic（最高优先级）
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered",
					zap.Any("panic", r),
					zap.String("path", c.Request.URL.Path),
				)
				// 使用 Failure 返回 500
				c.JSON(http.StatusInternalServerError, utils.Failure("服务器内部错误"))
				c.Abort()
			}
		}()

		// 执行后续 handler 和中间件
		c.Next()

		// 没有错误，且不是 404 → 直接返回
		if len(c.Errors) == 0 {
			if c.Writer.Status() == http.StatusNotFound {
				c.JSON(http.StatusNotFound, utils.Failure(c.FullPath(), "接口不存在"))
			}
			return
		}

		// 取最后一个错误（通常是最重要的那个）
		err := c.Errors.Last().Err

		// 情况1：参数校验失败（validator）
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			var errs []string
			for _, e := range ve {
				errs = append(errs, e.Error())
			}
			// 可以直接用 Failure，也可以自定义更详细的
			c.JSON(http.StatusBadRequest, utils.Failure(errs, "参数校验失败"))
			// 如果想把详细错误列表也返回，可以扩展 Failure 函数支持 data
			// 例如： c.JSON(400, Response[[]string]{Status:400, Message:"参数校验失败", Data:errs})
			return
		}

		// 情况2：JSON 格式错误（空 body、非法字符、EOF 等）
		if strings.Contains(err.Error(), "EOF") ||
			strings.Contains(err.Error(), "invalid character") ||
			strings.Contains(err.Error(), "cannot unmarshal") {
			c.JSON(http.StatusBadRequest, utils.Failure(err.Error(), "无效的 JSON 格式"))
			return
		}

		// 情况3：其他所有未预期错误 → 500
		logger.Error("unhandled error",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, utils.Failure(err.Error(), "服务器繁忙，请稍后再试"))
	}
}
