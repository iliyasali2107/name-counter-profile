package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"name-counter-profile/pkg/pb"

	"github.com/gin-gonic/gin"
)

func GetURL(ctx *gin.Context, client pb.ProfileServiceClient) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("bad request"))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}

	request := &pb.GetURLRequest{Id: id}

	res, err := client.GetURL(context.Background(), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}

	ctx.Redirect(http.StatusSeeOther, res.Url)
}
