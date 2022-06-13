package helpers

import (
	"context"
	"piepay/services/es"
	"piepay/services/logger"
	"piepay/structs"
	"piepay/structs/requests"
	"piepay/structs/response"

	"github.com/getsentry/sentry-go"
	jsoniter "github.com/json-iterator/go"
	"github.com/olivere/elastic/v7"
)

func GetLatestVideo(ctx context.Context, request *requests.GetVideo, sentryCtx context.Context) (response.GetVideo, error) {
	defer sentry.Recover()
	span := sentry.StartSpan(sentryCtx, "[DAO] GetLatestVideo") //sentry to log db calls
	defer span.Finish()

	if request.Size == 0 {
		request.Size = 10
	}

	dbSpan1 := sentry.StartSpan(span.Context(), "[DB] Get from videos index")
	res, err := es.Client().Search().Index("video").SearchSource(elastic.NewSearchSource().SortBy(SortDetails("publishedAt")).From(request.Page).Size(request.Size)).Do(ctx)

	dbSpan1.Finish() //noting time of query

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

func SortDetails(param string) *elastic.FieldSort {
	return elastic.NewFieldSort(param).Desc()
}
func FormDetails(id string) *elastic.TermQuery {
	return elastic.NewTermQuery("title", id)
}
func QueryDetails(param string, value string) *elastic.MatchQuery {
	return elastic.NewMatchQuery(param, value)
}