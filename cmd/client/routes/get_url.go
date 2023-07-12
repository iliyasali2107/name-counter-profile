package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"url-redirecter-url/pkg/pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetURL(ctx *gin.Context, client pb.URLServiceClient) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, "incorrect id")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "incorrect id")
		return
	}

	request := &pb.GetURLRequest{Id: id}

	res, err := client.GetURL(context.Background(), request)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			ctx.JSON(http.StatusNotFound, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, res.Url.Url)
}

type SetActiveRequestBody struct {
	Id int64 `json:"id"`
}

func SetActiveURL(ctx *gin.Context, client pb.URLServiceClient) {
	var reqBody SetActiveRequestBody
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "id is incorrect form"+err.Error())
		return
	}

	userId, ok := ctx.Get("userId")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, nil)
		return
	}

	req := &pb.SetActiveUrlRequest{
		UrlId:  reqBody.Id,
		UserId: userId.(int64),
	}

	_, err = client.SetActiveUrl(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			ctx.JSON(http.StatusNotFound, err)
		case codes.PermissionDenied:
			ctx.JSON(http.StatusForbidden, err)
		case codes.AlreadyExists:
			ctx.JSON(http.StatusForbidden, err)
		case codes.Internal:
			ctx.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type AddURLRequestBody struct {
	Url string `json:"url"`
}

func AddURL(ctx *gin.Context, client pb.URLServiceClient) {
	var reqBody AddURLRequestBody

	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "url is not correct")
		return
	}

	userId, ok := ctx.Get("userId")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, nil)
		return
	}

	req := &pb.AddURLRequest{
		UserId: userId.(int64),
		Url:    reqBody.Url,
	}

	res, err := client.AddURL(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, fmt.Sprintf("%s was added to url list", res.Url.Url))
}

type GetUserUrlsRequestBody struct {
	UserId int64 `json:"user_id"`
}

func GetUserURLs(ctx *gin.Context, client pb.URLServiceClient) {
	var reqBody GetUserUrlsRequestBody
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "incorrect user id")
		return
	}

	req := &pb.GetUserURLsRequest{
		UserId: reqBody.UserId,
	}

	res, err := client.GetUserURLs(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			ctx.JSON(http.StatusNotFound, "you have no urls")
			return
		}

		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
