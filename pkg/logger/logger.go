package logger

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/contexthandler"
	"github.com/sirupsen/logrus"
)

func Logger(ctx context.Context) *logrus.Entry {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	gc, err := contexthandler.RetrieveGinContext(ctx, "ContextKey")
	if err != nil {
		gc = &gin.Context{}
	}
	requestLogger := logger.WithFields(logrus.Fields{
		"requestID": gc.Value("requestID"),
	})
	return requestLogger
}
