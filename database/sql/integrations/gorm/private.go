package leafGorm

import (
	"errors"
	leafSql "github.com/enricodg/leaf-utilities/database/sql/sql"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

func (i *Impl) newImpl(db *gorm.DB) *Impl {
	return &Impl{
		GormDB:       db,
		GormDBDryRun: db.Session(&gorm.Session{DryRun: true}),
		Log:          i.Log,
		DatabaseName: i.DatabaseName,
		//DataStoreProduct: i.DataStoreProduct,
	}
}

func (i Impl) sort(sort []string, options leafSql.PaginateOptions) (string, []string) {
	sortQuery := ""

	arrSort := make([]string, 0)
	for _, s := range sort {
		s = strings.TrimSpace(s)
		if len(s) < 1 {
			continue
		}

		sortField := s
		sortDirection := ""
		if sortField[0] == '-' {
			sortDirection = " desc"
			sortField = sortField[1:]
		}

		sortField, ok := i.getMappedField(sortField, options)
		if len(sortField) < 1 {
			continue
		}

		if ok {
			arrSort = append(arrSort, s)
		}

		if len(sortQuery) != 0 {
			sortQuery += ","
		}

		sortQuery += sortField + sortDirection
	}

	return sortQuery, arrSort
}

func (i Impl) getMappedField(s string, paginateOptions leafSql.PaginateOptions) (string, bool) {
	if paginateOptions.FieldMap == nil {
		return s, false
	}

	if paginateOptions.MapOrDefault {
		if mapped := paginateOptions.FieldMap[s]; len(mapped) > 0 {
			return mapped, true
		}

		return s, false
	}

	if mapped := paginateOptions.FieldMap[s]; len(mapped) > 0 {
		return mapped, true
	}

	return "", false
}

func (i Impl) interfaceSlice(items interface{}) ([]interface{}, error) {
	switch reflect.TypeOf(items).Kind() {
	case reflect.Ptr:
		slice := reflect.ValueOf(items).Elem()
		result := make([]interface{}, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			result[i] = slice.Index(i).Interface()
		}
		return result, nil

	case reflect.Slice:
		slice := reflect.ValueOf(items)
		result := make([]interface{}, slice.Len())

		for i := 0; i < slice.Len(); i++ {
			result[i] = slice.Index(i).Interface()
		}
		return result, nil

	default:
		return nil, errors.New("can not proceed with non collection data")
	}
}

//func (i Impl) startDatastoreSegment(ctx *context.Context, operation string, statement *gorm.Statement) taniTracer.Span {
//	var span taniTracer.Span
//	span, *ctx = tracer.StartSpanFromContext(*ctx, operation,
//		taniNewRelicTracer.WithSpanType(taniNewRelicSpanType.DataStore),
//		taniNewRelicTracer.WithDataStore(taniNewRelicTracer.DataStoreOption{
//			Collection:         statement.Table,
//			DatabaseName:       i.DatabaseName,
//			Operation:          operation,
//			ParameterizedQuery: statement.SQL.String(),
//			QueryParameters:    statement.Vars,
//			DatastoreProduct:   i.DataStoreProduct,
//		}),
//		taniSentryTracer.WithSpanType(taniSentrySpanType.DataStore),
//		taniSentryTracer.WithDataStore(taniSentryTracer.DataStoreOption{
//			Collection:         statement.Table,
//			DatabaseName:       i.DatabaseName,
//			Operation:          operation,
//			ParameterizedQuery: statement.SQL.String(),
//			QueryParameters:    statement.Vars,
//			DatastoreProduct:   string(i.DataStoreProduct),
//		}),
//	)
//	return span
//}
