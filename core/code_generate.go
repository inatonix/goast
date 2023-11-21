package core

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type CodeGenerator interface {
	LoadContext(ctx context.Context) error
	Predict(ctx context.Context, q string) ([]string, error)
}

type InitializeParameters struct {
	Model  string
	Client *openai.Client
}

type codeGeneratorModel struct {
	InitializeParameters
}

func NewCodeGenerator(p InitializeParameters) (CodeGenerator, error) {
	m := &codeGeneratorModel{p}
	return m, nil
}

const currentProjects = `
	This is current project directories by programming language Go.

	├── config
	├── controller
	│   ├── common
	│   │   └── error
	│   ├── customer_segment_group
	│   ├── customer_segment_group_parameter
	│   ├── data_connector
	│   ├── data_source
	│   ├── dataset
	│   ├── digging_analytics
	│   ├── engagement_parameter
	│   ├── external_connection
	│   ├── extractresult
	│   ├── oldcustomersegmenttemplate
	│   ├── optmps
	│   ├── overall_analytics
	│   ├── statistics
	│   └── task
	├── dao
	│   └── listutil
	├── infrastructure
	│   └── apiclient
	│       └── mock
	├── logger
	├── middleware
	├── model
	│   └── aimstarjava
	├── reference
	├── settings
	│   ├── aimstarjava
	│   └── oapi-codegen
	│       └── gin
	├── usecase
	│   └── repository
	│       └── mock
	└── util

	and here is explanation of the directory structure.

	"There are four layers. In the controller layer, there are processes written to handle API requests received via the web. In the model layer, we place structures for processing within the application. Additionally, methods containing domain logic are attached to these structures. The DAO (Data Access Object) layer contains processes for connecting to the DB (Database) and persisting the model. In the use case layer, we write the use cases that the application aims to realize, effectively utilizing the model and DAO."


	Remember this structure For generating the new API code.
`

const controllerFileExplanation = `
	For example, /Users/inada/src/github.com/inatonix/goast/controller/overallanalytics.go file is just like this.


	package overallanalytics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supsysjp/insight/controller/common"
	"github.com/supsysjp/insight/logger"
	"github.com/supsysjp/insight/model"
	"github.com/supsysjp/insight/usecase"
)

type overallAnalyticsController struct {
	u usecase.OverallAnalyticsUsecase
	ServerInterface
}

func NewOverallAnalyticsController(u usecase.OverallAnalyticsUsecase) ServerInterface {
	return &overallAnalyticsController{u: u}
}

// (GET /overall_analytics)
func (c *overallAnalyticsController) GetInsightV1OverallAnalytics(ctx *gin.Context, params GetInsightV1OverallAnalyticsParams) {
	var req model.SearchQueryRequestParameters
	if err := ctx.BindQuery(&req); err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	sqm, err := common.ToSearchQueryModel(req)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, common.NewBadRequest(err.Error()))
		return
	}

	l, s, canDownload, err := c.u.List(ctx, *sqm)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	meta, err := common.ToSearchMetaData(req, s, canDownload)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var ll []OverallAnalyticsListItem
	if s == 0 {
		ll = []OverallAnalyticsListItem{}
	} else {
		ll, err = toListResponse(l)
		if err != nil {
			logger.Error(err)
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	res := OverallAnalyticsList{
		Meta: meta,
		Data: &ll,
	}
	ctx.JSON(http.StatusOK, res)
}

// (POST /overall_analytics)
func (c *overallAnalyticsController) PostInsightV1OverallAnalytics(ctx *gin.Context) {
	var req PostInsightV1OverallAnalyticsJSONRequestBody
	if err := ctx.BindJSON(&req); err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	a, err := toPostModel(req)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	a, err = c.u.Create(ctx, a)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, a)
}

// (GET /overall_analytics/{id})
func (c *overallAnalyticsController) GetInsightV1OverallAnalyticsId(ctx *gin.Context, id string) {
	a, err := c.u.Read(ctx, id)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, a)
}

// (DELETE /overall_analytics/{id})
func (c *overallAnalyticsController) DeleteInsightV1OverallAnalyticsId(ctx *gin.Context, id string) {

	_, err := c.u.Delete(ctx, id)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// (PUT /overall_analytics/{id})
func (c *overallAnalyticsController) PutInsightV1OverallAnalyticsId(ctx *gin.Context, id string) {
	var req PutInsightV1OverallAnalyticsIdJSONRequestBody
	if err := ctx.BindJSON(&req); err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	a, err := toModel(req)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	a, err = c.u.Update(ctx, id, a)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, a)
}

// (POST /overall_analytics/{id}/execute)
func (c *overallAnalyticsController) PostInsightV1OverallAnalyticsIdExecute(ctx *gin.Context, id string) {
	// var req PutInsightV1OverallAnalyticsIdJSONRequestBody
	var req PostInsightV1OverallAnalyticsIdExecuteJSONRequestBody
	if err := ctx.BindJSON(&req); err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var a *model.OverallAnalyticsSetting
	var err error
	if req.AdhocParameter != nil {
		a, err = toExecuteModel(*req.AdhocParameter)
		if err != nil {
			logger.Error(err)
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}
	p := req.ParentFilePath

	taskId, err := c.u.Execute(ctx, id, a, p)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	res := struct {
		TaskId string json:"task_id"
	}{
		TaskId: taskId,
	}
	ctx.JSON(http.StatusOK, res)
}

// (GET /overall_analytics/_/status/{task_id})
func (c *overallAnalyticsController) GetInsightV1OverallAnalyticsStatusTaskId(ctx *gin.Context, taskId string) {
	a, err := c.u.GetResult(ctx, taskId)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, a)
}


`

const usecaseExplanation = `
For example, usecase file is just like this.


package usecase

import (
	"context"
	"log"

	"github.com/supsysjp/insight/model"
	"github.com/supsysjp/insight/usecase/repository"
	"github.com/supsysjp/insight/util"
	"golang.org/x/xerrors"
)

type OverallAnalyticsUsecase interface {
	Create(ctx context.Context, m *model.OverallAnalytics) (*model.OverallAnalytics, error)
	Read(ctx context.Context, id string) (*model.OverallAnalytics, error)
	Update(ctx context.Context, id string, m *model.OverallAnalytics) (*model.OverallAnalytics, error)
	Delete(ctx context.Context, id string) (*string, error)
	List(ctx context.Context, params model.SearchQueryParameters) (model.OverallAnalyticsList, int64, bool, error)
	Execute(ctx context.Context, id string, m *model.OverallAnalyticsSetting, p string) (string, error)
	GetResult(ctx context.Context, taskId string) (*model.OverallAnalyticsTaskStatus, error)
}
type overallAnalyticsUsecase struct {
	r  repository.OverallAnalyticsRepository
	re repository.OverallAnalyticsExecuteRepository
	rr repository.OverallAnalyticsResultRepository
}

func NewOverallAnalyticsUsecase(r repository.OverallAnalyticsRepository, re repository.OverallAnalyticsExecuteRepository, rr repository.OverallAnalyticsResultRepository) OverallAnalyticsUsecase {
	return &overallAnalyticsUsecase{
		r:  r,
		re: re,
		rr: rr,
	}
}

func (u *overallAnalyticsUsecase) Create(ctx context.Context, m *model.OverallAnalytics) (*model.OverallAnalytics, error) {

	m, err := u.r.Create(ctx, m)
	if err != nil {
		return nil, decodeError(err)
	}
	return m, nil
}
func (u *overallAnalyticsUsecase) Update(ctx context.Context, id string, m *model.OverallAnalytics) (*model.OverallAnalytics, error) {

	m, err := u.r.Update(ctx, id, m)
	if err != nil {
		return nil, decodeError(err)
	}
	return m, nil
}
func (u *overallAnalyticsUsecase) Read(ctx context.Context, id string) (*model.OverallAnalytics, error) {
	m, err := u.r.Read(ctx, id)
	if err != nil {
		return nil, decodeError(err)
	}
	return m, nil
}
func (u *overallAnalyticsUsecase) Delete(ctx context.Context, id string) (*string, error) {
	m, err := u.r.Delete(ctx, id)
	if err != nil {
		return nil, decodeError(err)
	}
	return m, nil

}
func (u *overallAnalyticsUsecase) List(ctx context.Context, params model.SearchQueryParameters) (model.OverallAnalyticsList, int64, bool, error) {
	// tenantId := au.TenantId
	l, total, err := u.r.List(ctx, params)
	if err != nil {
		return nil, 0, false, decodeError(err)
	}
	return l, total, false, nil

}
func (u *overallAnalyticsUsecase) Execute(ctx context.Context, id string, m *model.OverallAnalyticsSetting, p string) (string, error) {

	var err error
	if m == nil {
		// パラメータの指定がない場合はDBに保存された全体分析を使用する
		ms, err := u.r.Read(ctx, id)
		if err != nil {
			return "", decodeError(err)
		}
		m = &model.OverallAnalyticsSetting{
			AnalyticsType: ms.AnalyticsType,
			Parameters:    ms.Parameters,
			Segments:      ms.Segments,
		}
	}

	taskId := util.GetId(ctx)
	err = u.rr.Create(ctx, taskId) // 実行開始を登録
	if err != nil {
		return "", decodeError(err)
	}

	// TODO: 実際の実行部分はこれから実装
	go func() {
		sql, err := m.GenerateSQLQuery(ctx)
		if err != nil {
			u.rr.Update(ctx, taskId, nil, err)
			log.Println(err)
			return
		}
		result, err := u.re.ExecuteQuery(ctx, sql)
		if err != nil {
			u.rr.Update(ctx, taskId, nil, err)
			log.Println(err)
			return
		}

		// 実行結果を登録
		err = u.rr.Update(ctx, taskId, result, nil)
		if err != nil {
			// backgroundで実行しているのでエラーはログ出力のみ
			u.rr.Update(ctx, taskId, nil, err)
			log.Printf("failed to update task status: %v", err)
			return
		}
	}()
	return taskId, nil
}

func (u *overallAnalyticsUsecase) GetResult(ctx context.Context, taskId string) (*model.OverallAnalyticsTaskStatus, error) {
	return u.rr.Read(ctx, taskId)
}

// Errorの内容によりエラーコードを含むエラーに変換する
func decodeError(e error) error {
	if e == nil {
		return nil
	}
	// switch e.Error() {
	// case "mongo: no documents in result":
	// 	return e
	// }
	return xerrors.Errorf("unknown error: %w", e)
}



`

const daoExplanation = `
For example, dao file is just like this.

package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/supsysjp/insight/dao/listutil"
	"github.com/supsysjp/insight/infrastructure"
	"github.com/supsysjp/insight/model"
	"github.com/supsysjp/insight/usecase/repository"
	"github.com/supsysjp/insight/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/xerrors"
)

type overallAnalyticsDao struct {
	c         infrastructure.MongoClientWrapper
	tableName string
}

func NewOverallAnalyticsDao(client infrastructure.MongoClientWrapper) repository.OverallAnalyticsRepository {
	return &overallAnalyticsDao{
		c:         client,
		tableName: "overall_analytics",
	}
}

func (d *overallAnalyticsDao) Read(ctx context.Context, id string) (*model.OverallAnalytics, error) {
	tenantId, err := infrastructure.GetTenantID(ctx)
	if err != nil {
		return nil, err
	}
	var cs model.OverallAnalytics
	err = d.c.Database(ctx).Collection(d.tableName).FindOne(ctx, bson.M{"_id": id, "tenant_id": tenantId}).Decode(&cs)
	if err != nil {
		return nil, err
	}

	return &cs, nil
}
func (d *overallAnalyticsDao) Delete(ctx context.Context, id string) (*string, error) {
	tenantId, err := infrastructure.GetTenantID(ctx)
	if err != nil {
		return nil, err
	}
	result, err := d.c.Database(ctx).Collection(d.tableName).DeleteOne(ctx, bson.M{"_id": id, "tenant_id": tenantId})
	if err != nil {
		return nil, err
	}
	if result.DeletedCount > 0 {
		return &id, nil
	}
	return nil, nil
}

func (d *overallAnalyticsDao) Create(ctx context.Context, m *model.OverallAnalytics) (*model.OverallAnalytics, error) {
	// User情報をcontrollerから移動
	au, err := util.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to get authenticated user: %w", err)
	}
	now := time.Now()
	m.Id = util.GetId(ctx)
	m.TenantId = au.TenantId
	m.CreatedAt = now
	m.UpdatedAt = now
	m.CreatedBy = au.Id
	m.UpdatedBy = au.Id

	result, err := d.c.Database(ctx).Collection(d.tableName).InsertOne(ctx, m)
	if err != nil {
		return nil, err
	}
	v := result.InsertedID
	id, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("failed to convert to string")
	}
	m.Id = id
	return m, nil
}
func (d *overallAnalyticsDao) Update(ctx context.Context, id string, m *model.OverallAnalytics) (*model.OverallAnalytics, error) {
	// User情報をcontrollerから移動
	au, err := util.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to get authenticated user: %w", err)
	}
	now := time.Now()
	m.TenantId = au.TenantId
	m.UpdatedAt = now
	m.UpdatedBy = au.Id

	filter := bson.M{
		"_id":        m.Id,
		"tenant_id":  au.TenantId,
		"deleted_at": nil,
	}
	update := bson.M{
		"$set": bson.M{
			"title":          m.Title,
			"analytics_type": m.AnalyticsType,
			"updated_at":     m.UpdatedAt,
			"updated_by":     m.UpdatedBy,
			"parameters":     m.Parameters,
			"segments":       m.Segments,
		},
	}

	result, err := d.c.Database(ctx).Collection(d.tableName).UpdateOne(ctx, filter, update)
	fmt.Println(result, err)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (d *overallAnalyticsDao) List(ctx context.Context, params model.SearchQueryParameters) (model.OverallAnalyticsList, int64, error) {
	tenantID, err := infrastructure.GetTenantID(ctx)
	if err != nil {
		return nil, 0, err
	}

	findOption, err := listutil.FindOption(params.Orderby, params.StartRow, params.Rows)
	if err != nil {
		return nil, 0, err
	}

	filter, err := listutil.FilterOptions(params.Filter)
	if err != nil {
		return nil, 0, err
	}
	//FilterOptionsにはstring以外を入れることができない
	var options []bson.M
	if filter != nil && (*filter)["$and"] != nil {
		options = (*filter)["$and"].([]bson.M)
	} else {
		options = []bson.M{}
	}
	options = append(options, bson.M{"tenant_id": bson.M{"$eq": tenantID}})
	options = append(options, bson.M{"deleted_at": bson.M{"$eq": nil}})
	if params.Keywords != "" {
		options = append(options, bson.M{"title": bson.M{
			"$regex": primitive.Regex{
				Pattern: "^(?=.*" + params.Keywords + ").*$",
				Options: "im",
			},
		}})
	}

	// options = append(options, bson.M{"temporary": bson.M{"$eq": false}})
	filterAddDeletedAt := bson.M{"$and": options}

	findOption.SetProjection(bson.D{
		{Key: "_id", Value: 1},
		{Key: "title", Value: 1},
		{Key: "tenant_id", Value: 1},
		{Key: "analytics_type", Value: 1},
		{Key: "created_at", Value: 1},
		{Key: "created_by", Value: 1},
		{Key: "updated_at", Value: 1},
		{Key: "updated_by", Value: 1},
		{Key: "tags", Value: 1},
	})
	colls := d.c.Database(ctx).Collection(d.tableName)

	cursor, err := colls.Find(ctx, filterAddDeletedAt, findOption)
	if err != nil {
		return nil, 0, xerrors.Errorf("list overall analytics error: %w", err)
	}
	defer cursor.Close(ctx)

	var results model.OverallAnalyticsList
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, xerrors.Errorf("cursor overall analytics error: %w", err)
	}

	size, err := colls.CountDocuments(ctx, filterAddDeletedAt)
	if err != nil {
		return nil, 0, xerrors.Errorf("count overall analytics error: %w", err)
	}

	return results, size, nil
}



`

const modelExplanation = `
For example, model file is just like this.

package model

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// OverallAnalytics defines model for OverallAnalytics.
type OverallAnalytics struct {
	// Id 保存された分析ID
	Id string "json:"id" bson:"_id""

	// テナントID
	TenantId string "json:"tenant_id" bson:"tenant_id""

	// Title 保存名
	Title string "json:"title" bson:"title""

	// AnalyticsType 分析手法(RFM, ファネルなど)
	AnalyticsType OverallAnalyticsAnalyticsType "json:"analytics_type" bson:"analytics_type""

	// CreatedAt 作成日時
	CreatedAt time.Time "json:"created_at" bson:"created_at""

	// CreatedBy 作成ユーザID
	CreatedBy string "json:"created_by" bson:"created_by""

	// UpdatedAt 更新日時
	UpdatedAt time.Time "json:"updated_at,omitempty" bson:"updated_at,omitempty""

	// UpdatedBy 更新ユーザID
	UpdatedBy string "json:"updated_by,omitempty" bson:"updated_by,omitempty""

	// DeletedAt 削除日時
	DeletedAt time.Time "json:"deleted_at,omitempty" bson:"deleted_at,omitempty""

	// Parameters 集計のためのパラメータ
	Parameters OverallAnalyticsAnalyticsParameters "json:"parameters" bson:"parameters""

	// Segments 結果をセグメント分けする条件
	Segments []OverallAnalyticsSegment "json:"segments" bson:"segments""
}
`

const APIFile = `
This is API file defined by OpenAPI.


openapi: 3.0.0
info:
  title: digging_analytics
  version: '1.0'
  description: 深掘り分析用のAPI
servers: []
paths:
  /insight/v1/digging_analytics:
    get:
      summary: 保存された深掘り分析設定の一覧
      tags: []
      responses:
        '200':
          $ref: '#/components/responses/DiggingAnalyticsList'
      operationId: get-insight-v1-digging-analytics
      description: 保存された分析設定の一覧を取得する
      parameters:
        - $ref: ./common.yaml#/components/parameters/BrowserKeywords
        - $ref: ./common.yaml#/components/parameters/BrowserFilter
        - $ref: ./common.yaml#/components/parameters/BrowserOrderBy
        - $ref: ./common.yaml#/components/parameters/BrowserRows
        - $ref: ./common.yaml#/components/parameters/BrowserStartRow
    post:
      summary: 深掘り分析のパラメータ保存
      operationId: post-insight-v1-digging-analytics
      responses:
        '201':
          $ref: '#/components/responses/DiggingAnalyticsResponse'
      description: 分析を新規作成する
      requestBody:
        $ref: '#/components/requestBodies/DiggingAnalyticsPostRequest'
  '/insight/v1/digging_analytics/{id}':
    get:
      summary: 保存された深掘り分析設定の取得
      tags: []
      responses:
        '200':
          $ref: '#/components/responses/DiggingAnalyticsResponse'
      operationId: get-insight-v1-digging-analytics-id
      description: 保存された分析設定を取得する
    parameters:
      - schema:
          type: string
        name: id
        in: path
        required: true
    put:
      summary: 深掘り分析設定の更新
      operationId: put-insight-v1-digging-analytics-id
      responses:
        '200':
          $ref: '#/components/responses/DiggingAnalyticsResponse'
      description: 保存された分析用設定を更新する
      requestBody:
        $ref: '#/components/requestBodies/DiggingAnalyticsRequest'
    delete:
      summary: 深掘り分析設定の削除
      operationId: delete-insight-v1-digging-analytics-id
      responses:
        '204':
          description: No Content
      description: 指定された分析を削除する
components:
  schemas:
    DiggingAnalytics:
      allOf:
        - $ref: '#/components/schemas/DiggingAnalyticsMeta'
        - $ref: '#/components/schemas/DiggingAnalyticsSetting'
      x-examples:
        Example 1:
          id: string
          analytics_type: rfm
          title: string
          created_at: '2019-08-24T14:15:22Z'
          updated_at: '2019-08-24T14:15:22Z'
          deleted_at: '2019-08-24T14:15:22Z'
          created_by: string
          updated_by: string
          parameters:
            type: rf
            recency:
              - 1
              - 3
              - 6
              - 12
            recency_unit: month
            frequency:
              - 0
              - 1
              - 3
              - 6
          segments:
            - id: string
              title: string
              parameters:
                - key: string
                  value: string
      title: ''
      description: 深掘り分析の情報
    DiggingAnalyticsSetting:
      type: object
      description: |-
        深掘り分析を実行するための情報
        ・集計のためのparameters
        ・セグメントを分けのためのsegments
        を持つ
      properties:
        analytics_type:
          type: string
          description: '分析手法(RFM, ファネルなど)'
          enum:
            - rfm
            - funnel
        parameters:
          oneOf:
            - $ref: '#/components/schemas/DiggingAnalyticsParameterRFM'
            - $ref: '#/components/schemas/DiggingAnalyticsParameterFunnel'
        segments:
          type: array
          description: 結果をセグメント分けする条件
          items:
            $ref: '#/components/schemas/DiggingAnalyticsSegment'
      required:
        - analytics_type
        - parameters
        - segments
      x-examples:
        Example 1:
          analytics_type: rfm
          parameters:
            type: rf
            recency:
              - 0
            recency_unit: day
            frequency:
              - 0
            monetary:
              - 0
          segments:
            - id: string
              title: string
              parameters:
                - key: string
                  value: string
    DiggingAnalyticsMeta:
      title: DiggingAnalyticsMeta
      type: object
      description: 保存された分析条件の情報
      properties:
        id:
          type: string
          description: 保存された分析ID
        title:
          type: string
          description: 保存名
        created_at:
          type: string
          description: 作成日時
          format: date-time
        updated_at:
          type: string
          description: 更新日時
          format: date-time
        deleted_at:
          type: string
          description: 削除日時
          format: date-time
        created_by:
          type: string
          description: 作成ユーザID
        updated_by:
          type: string
          description: 更新ユーザID
      required:
        - id
        - title
        - created_at
        - created_by
    DiggingAnalyticsListItem:
      type: object
      properties:
        id:
          type: string
          description: 保存された分析ID
        title:
          type: string
          description: 保存名
        analytics_type:
          type: string
          description: '分析手法(RFM, ファネルなど)'
          enum:
            - rfm
            - funnel
        created_at:
          type: string
          description: 作成日時
          format: date-time
        updated_at:
          type: string
          description: 更新日時
          format: date-time
        deleted_at:
          type: string
          description: 削除日時
          format: date-time
        created_by:
          type: string
          description: 作成ユーザID
        updated_by:
          type: string
          description: 更新ユーザID
      required:
        - id
        - analytics_type
        - title
        - created_at
        - created_by
    DiggingAnalyticsParameter:
      title: AnalyticsParameter
      type: object
      properties:
        key:
          type: string
          description: パラメータのキー
        value:
          type: string
          description: パラメータの値
      required:
        - key
        - value
    DiggingAnalyticsSegment:
      title: DiggingAnalyticsSegmentSetting
      type: object
      description: 深掘り分析で表す軸の情報(RFMのセグメントの範囲など)を持つ
      x-examples:
        Example 1:
          id: '02'
          title: 育成現役
          parameters:
            - key: r_min
              value: '0'
            - key: r_max
              value: '6'
            - key: f_min
              value: '1'
            - key: f_max
              value: '3'
      properties:
        id:
          type: string
          description: SegmentのID(絞り込みなどで使用する)
        title:
          type: string
          description: セグメントを表す名前
        parameters:
          type: array
          description: 軸を特定する値(RFMのF値のような値)
          items:
            $ref: '#/components/schemas/DiggingAnalyticsParameter'
      required:
        - id
        - title
        - parameters
    DiggingAnalyticsSegmentAggregations:
      title: DiggingAnalyticsSegmentResult
      type: object
      description: 深掘り分析で集計した１つのセグメントの情報
      x-examples:
        Example 1:
          id: '01'
          title: 新規現役
          values:
            - key: customers
              value: 100000
            - key: sales_p
              value: 1900
      properties:
        id:
          type: string
          description: セグメントID
        title:
          type: string
          description: セグメント名
        aggregations:
          type: array
          items:
            type: object
            properties:
              key:
                type: string
                description: 'Valueの種類(customers, salesなど)'
              value:
                type: number
                description: 集計値(顧客数、購買金額など)
            required:
              - key
              - value
      required:
        - id
        - title
        - aggregations
    DiggingAnalyticsSegmentAggregationsList:
      title: DiggingAnalyticsSegmentAggregationsList
      type: array
      items:
        $ref: '#/components/schemas/DiggingAnalyticsSegmentAggregations'
      description: |
        深掘り分析で集計されたセグメント毎の集計値のリスト
      x-examples:
        Example 1:
          - id: '01'
            title: 新規現役
            values:
              - key: customers
                value: 100000
              - key: sales_p
                value: 1900
          - id: '02'
            Title: 育成現役
            values:
              - key: customers
                value: 33444
              - key: sales_p
                value: 16894
          - id: '03'
            title: 急進現役
            values:
              - key: customers
                value: 340
              - key: sales_p
                value: 70000
    DiggingAnalyticsParameterRFM:
      title: DiggingAnalyticsParameterRFM
      type: object
      description: RFM分析で使用するパラメータ
      properties:
        type:
          type: string
          description: 'RFMの種類(RF,RM,FM) 将来は3軸(RFM)もある'
          enum:
            - rf
            - rm
            - fm
            - rfm
        recency:
          type: array
          description: R値の分割値(この数値以下で分割)
          items:
            type: integer
        recency_unit:
          type: string
          description: |-
            Rの単位(day, week, month, year)
            ※ monthは月数ではなく期間を固定とするため30日相当としたい
          enum:
            - day
            - week
            - month
            - year
        frequency:
          type: array
          description: Fの区切り(この数値以下で分ける)
          items:
            type: integer
        monetary:
          type: array
          description: Mの区切り(この数値以下で分ける)
          items:
            type: integer
      required:
        - type
      x-examples:
        Example 1:
          type: rf
          recency:
            - 1
            - 3
            - 6
            - 12
          recency_unit: month
          frequency:
            - 0
            - 1
            - 3
            - 6
    DiggingAnalyticsParameterFunnel:
      title: DiggingAnalyticsParameterFunnel
      type: object
      description: |
        ファネル分析用の集計
        (現時点ではファネルの具体的なパラメータは仮)

        event_type : ["web"],
        step_event : [
          ["サイト訪問"],
          ["商品詳細閲覧"],
          ["カート投入", "購買条件確認"],
          ["購入"],
          ["レビュー投稿"]
        ]
        みたいな感じ？
      properties:
        event_type:
          type: array
          description: ファネルを進めるイベントの種類
          items:
            type: string
        step_events:
          type: array
          description: 各ステップ毎の次に進むイベントのリスト
          items:
            type: array
            items:
              type: string
  responses:
    DiggingAnalyticsList:
      description: Example response
      content:
        application/json:
          schema:
            type: object
            properties:
              data:
                type: array
                items:
                  $ref: '#/components/schemas/DiggingAnalyticsListItem'
              meta:
                $ref: ./common.yaml#/components/schemas/SearchMetadata
          examples:
            Example 1:
              value:
                data:
                  - id: '12345'
                    analytics_type: rfm
                    title: RFMセグメント
                    created_at: '2019-08-24T14:15:22Z'
                    updated_at: '2019-08-24T14:15:22Z'
                    created_by: user1
                    updated_by: user1
                  - id: '9999'
                    analytics_type: funnel
                    title: 購買ファネル
                    created_at: '2019-08-24T14:15:22Z'
                    updated_at: '2019-08-24T14:15:22Z'
                    created_by: user1
                    updated_by: user2
                meta:
                  total_object_count: 2
                  per_page: 30
                  current_page: 1
                  last_page: 1
                  from: 1
                  to: 2
                  can_bulk_download: true
    DiggingAnalyticsResponse:
      description: Example response
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/DiggingAnalytics'
          examples:
            Example 1:
              value:
                id: '12345'
                analytics_type: rfm
                title: RFMセグメント分析
                created_at: '2019-08-24T14:15:22Z'
                updated_at: '2019-08-24T14:15:22Z'
                created_by: user1
                updated_by: user2
                parameters:
                  type: rf
                  recency:
                    - 1
                    - 3
                    - 6
                    - 12
                  recency_unit: month
                  frequency:
                    - 0
                    - 1
                    - 3
                    - 6
                  monetary:
                    - 0
                segments:
                  - id: '01'
                    title: 新規
                    parameters:
                      - key: r_min
                        value: '0'
                      - key: r_max
                        value: '3'
                      - key: f_min
                        value: '1'
                      - key: f_max
                        value: '1'
                  - id: '02'
                    title: 継続
                    parameters:
                      - key: r_min
                        value: '0'
                      - key: r_max
                        value: '3'
                      - key: f_min
                        value: '2'
                      - key: f_max
                        value: '3'
  requestBodies:
    DiggingAnalyticsRequest:
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/DiggingAnalytics'
          examples:
            Example 1:
              value:
                id: '12345'
                analytics_type: rfm
                title: RFMセグメント分析
                created_at: '2019-08-24T14:15:22Z'
                updated_at: '2019-08-24T14:15:22Z'
                parameters:
                  type: rf
                  recency:
                    - 0
                    - 3
                    - 12
                  recency_type: month
                  frequency:
                    - 1
                    - 3
                    - 6
                  monetary:
                    - 0
                segments:
                  - id: '01'
                    title: 新規
                    parameters:
                      - key: r_min
                        value: '0'
                      - key: r_max
                        value: '3'
                      - key: f_min
                        value: '1'
                      - key: f_max
                        value: '1'
                  - id: '02'
                    title: 継続
                    parameters:
                      - key: r_min
                        value: '0'
                      - key: r_max
                        value: '3'
                      - key: f_min
                        value: '2'
                      - key: f_max
                        value: '3'
      description: 深掘り分析の登録用
    DiggingAnalyticsPostRequest:
      content:
        application/json:
          schema:
            type: object
            properties:
              analytics_type:
                type: string
                description: '分析手法(RFM, ファネルなど)'
              title:
                type: string
                description: 保存名
              parameters:
                type: array
                description: 集計のためのパラメータ
                items:
                  $ref: '#/components/schemas/DiggingAnalyticsParameter'
              segments:
                type: array
                description: 結果をセグメント分けする条件
                items:
                  $ref: '#/components/schemas/DiggingAnalyticsSegment'

`

func (m *codeGeneratorModel) LoadContext(ctx context.Context) error {
	// _, err := m.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
	// 	Model: m.Model,
	// 	Messages: []openai.ChatCompletionMessage{
	// 		{
	// 			Role:    "user",
	// 			Content: currentProjects,
	// 		},
	// 	}},
	// )
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (m *codeGeneratorModel) Predict(ctx context.Context, q string) ([]string, error) {
	resp, err := m.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: m.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are a talented software engineer. The API file will now be given, and the objective is to implement all functions needed according to the specifications in the API file.",
			},
			{
				Role:    "system",
				Content: "These are the code which already exist in the directory and have relation with the file you are going to generate. \n",
			},
			{
				Role:    "user",
				Content: currentProjects,
			},
			// {
			// 	Role:    "user",
			// 	Content: controllerFileExplanation,
			// },
			{
				Role:    "user",
				Content: modelExplanation,
			},
			// {
			// 	Role:    "user",
			// 	Content: usecaseExplanation,
			// },
			{
				Role:    "user",
				Content: daoExplanation,
			},
			{
				Role:    "user",
				Content: APIFile,
			},
			// {
			// 	Role:    "assistant",
			// 	Content: "OK, I will implement the dao file according to the given files",
			// },
			// {
			// 	Role:    "assistant",
			// 	Content: "/Users/inada/src/github.com/inatonix/goast/dao/digging_analytics.go",
			// },

			{
				Role:    "user",
				Content: "Please generate digging_analytics.go DAO file, according to the given files.",
			},
		},
	})
	if err != nil {
		return []string{""}, err
	}

	var messages []string
	for _, choice := range resp.Choices {
		messages = append(messages, choice.Message.Content)
	}

	return messages, nil
}
