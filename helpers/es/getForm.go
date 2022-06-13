package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"piepay/services/es"
	"piepay/services/logger"
	"piepay/structs"
	"piepay/structs/requests"
	"piepay/structs/response"

	"github.com/getsentry/sentry-go"
	jsoniter "github.com/json-iterator/go"
	"github.com/olivere/elastic/v7"
)

func GetSearchVideo(ctx context.Context, request *requests.SearchVideo, sentryCtx context.Context) (response.GetVideo, error) {
	defer sentry.Recover()
	span := sentry.StartSpan(sentryCtx, "[DAO] GetAny")
	defer span.Finish()
	if request.Page == 0 {
		request.Size = 10
	}
	dbSpan1 := sentry.StartSpan(span.Context(), "[DB] Get from videos")
	res, err := es.Client().Search().Index("video").SearchSource(elastic.NewSearchSource().Query(QueryDetails("description", "a")).SortBy(SortDetails("publishedAt")).From(request.Page).Size(request.Size)).Do(ctx)
	rescfg, _ := json.Marshal(elastic.NewSearchSource().Query(QueryDetails("description", "a")).SortBy(SortDetails("publishedAt")))
	fmt.Println(string(rescfg))

	fmt.Println(string(rescfg))
	dbSpan1.Finish()

	if err != nil {
		sentry.CaptureException(err)
		logger.Client().Error(err.Error())
		return response.GetVideo{}, err
	}
	var data1 structs.Video
	var dataRes []structs.Video
	if res != nil {
		for _, s := range res.Hits.Hits {
			jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(s.Source, &data1)
			dataRes = append(dataRes, data1)
		}
	}

	getRes := response.GetVideo{
		VideoDetails: dataRes,
		Page:         request.Page + 1,
		Size:         request.Size,
	}
	return getRes, nil
}
