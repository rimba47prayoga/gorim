package pagination

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mcuadros/go-defaults"
	"github.com/rimba47prayoga/gorim.git/errors"
	"github.com/rimba47prayoga/gorim.git/models"
	"github.com/rimba47prayoga/gorim.git/utils"
	"gorm.io/gorm"
)

type Pagination struct {
	Db				*gorm.DB    `json:"-"`
    Page         	int         `json:"page" default:"1"`
	PageSize		int			`json:"page_size" default:"10"`
    Sort         	string      `json:"sort"`
    TotalRows    	int64       `json:"total_rows"`    
    TotalPages   	int         `json:"total_pages"`   
    Data         	interface{} `json:"-"`  
}

func (p *Pagination) GetOffset() int {  
    return (p.GetPage() - 1) * p.GetLimit() 
}   

func (p *Pagination) GetLimit() int {   
    if p.PageSize == 0 {   
        p.PageSize = 10    
    }   
    return p.PageSize  
}

func (p *Pagination) GetPage() int {    
    if p.Page == 0 {    
        p.Page = 1  
    }   
    return p.Page   
}

func (p *Pagination) GetSort() []string { 
    if p.Sort == "" {   
        p.Sort = "ID"  
    }
    // Split the sort string by commas
    sortFields := strings.Split(p.Sort, ",")
    var sortClauses []string

    for _, field := range sortFields {
        field = strings.TrimSpace(field)
        if field != "" {
            sortClause := field+" asc"
            if strings.HasPrefix(field, "-") {
                // If field starts with "-", it's descending
                sortClause = strings.TrimPrefix(field, "-")+" desc"
            }
            sortClauses = append(sortClauses, strings.ToLower(sortClause))
        }
    }
    return sortClauses
}

func InitPagination(ctx echo.Context, db *gorm.DB) *Pagination {
	pagination := Pagination{
		Db: db,
	}
	defaults.SetDefaults(&pagination)

	page, _ := strconv.Atoi(ctx.QueryParam("page"))
    pageSize, _ := strconv.Atoi(ctx.QueryParam("page_size"))
    sort := ctx.QueryParam("sort")

	if page > 0 {
        pagination.Page = page
    }
    if pageSize > 0 {
        pagination.PageSize = pageSize
    }
    if sort != "" {
        pagination.Sort = sort
    }
	return &pagination
}

func (p *Pagination) SortQuery(results interface{}) []string {
    // Extract the model type from the slice
    sliceType := reflect.TypeOf(results)
    
    // Handle pointer to slice
    if sliceType.Kind() == reflect.Ptr {
        sliceType = sliceType.Elem()
    }
    if sliceType.Kind() != reflect.Slice {
        errors.Raise(&errors.InternalServerError{
            Message: "not a slice or pointer to slice, cannot proceed",
        })
    }
    modelType := sliceType.Elem()
    if modelType.Kind() == reflect.Ptr {
        modelType = modelType.Elem() // Dereference the pointer to get the actual struct type
    }
     // Create a new instance of the model and pass it to GetModelFields
    modelInstance := reflect.New(modelType).Elem().Interface()
    // Get valid fields for sorting
    validFields, err := models.GetModelFields(modelInstance)
    if err != nil {
        errors.Raise(&errors.InternalServerError{
            Message: err.Error(),
        })
    }
    // Convert valid fields to lowercase
    for i, field := range validFields {
        validFields[i] = strings.ToLower(field)
    }
    // Validate and filter sort fields
    var validSortClauses []string
    for _, sortClause := range p.GetSort() {
        parts := strings.Fields(sortClause)
        if len(parts) != 2 {
            continue // Invalid format, skip
        }
        field, direction := parts[0], parts[1]
        field = strings.ToLower(field)
        if utils.Contains(validFields, field) {
            validSortClauses = append(validSortClauses, fmt.Sprintf("%s %s", field, direction))
        }
    }
    // Apply validated sort clauses
    for _, sortClause := range validSortClauses {
        p.Db = p.Db.Order(sortClause)
    }
    return validSortClauses
}

func (p *Pagination) PaginateQuery(results interface{}) {
    
    var totalRows int64
    p.Db.Model(results).Count(&totalRows)
    p.TotalRows = totalRows
    totalPages := int(totalRows/int64(p.PageSize)) + 1
    if totalRows <= 1 {
        totalPages = 1
    }
    p.TotalPages = totalPages
    offset := (p.Page - 1) * p.PageSize
    sortClauses := p.SortQuery(results)
    p.Sort = strings.Join(sortClauses, ",")
    p.Db.Offset(offset).Limit(p.PageSize).Find(results)
    p.Data = results
}

func (p *Pagination) GetPaginatedResponse() *PaginatedResponse {
	status := true
	message := "Data Found."
	if p.TotalRows == 0 {
		status = false
		message = "Data Not Found."
	}
	return &PaginatedResponse{
		Status:  status,
		Message: message,
		Data:    p.Data,
		Pagination: Pagination{
			Page:       p.Page,
			PageSize:   p.PageSize,
			Sort:       p.Sort,
			TotalRows:  p.TotalRows,
			TotalPages: p.TotalPages,
		},
	}
}


func Paginate(value interface{}, pagination *Pagination) func(db *gorm.DB) *gorm.DB {  
    var totalRows int64 
    pagination.Db.Model(value).Count(&totalRows)   
    pagination.TotalRows = totalRows 
    totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.PageSize)))
    pagination.TotalPages = totalPages  
    return func(db *gorm.DB) *gorm.DB { 
        return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())   
    }   
}   
