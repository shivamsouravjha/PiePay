package GET

import (
	"context"
	"fmt"
	"net/http"
	helpers "piepay/helpers/es"
	"piepay/structs/requests"
	"piepay/structs/response"
	"piepay/utils"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func GetVideo(c *gin.Context) {
	defer sentry.Recover()
	span := sentry.StartSpan(context.TODO(), "[GIN] GetAnyHandler", sentry.TransactionName("Get Any Handler"))
	defer span.Finish()

	formRequest := requests.GetVideo{}

	if err := c.ShouldBind(&formRequest); err != nil {
		span.Status = sentry.SpanStatusFailedPrecondition
		sentry.CaptureException(err)
		c.JSON(422, utils.SendErrorResponse(err))
		return
	}
	ctx := c.Request.Context()
	resp := response.VideoResponse{}
	fmt.Println(formRequest)
	response, err := helpers.GetAny(ctx, &formRequest, span.Context())
	if err != nil {
		resp.Status = "Failed"
		resp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Status = "Success"
	resp.Message = "Creator updated successfully"
	resp.Data = response
	span.Status = sentry.SpanStatusOK

	c.JSON(http.StatusOK, resp)

}
