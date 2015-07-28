package dbutils

import (
	"reflect"

	"github.com/ilgooz/gorm"
	"github.com/ilgooz/sqlstruct"
	"github.com/lann/squirrel"
)

func RemoveM2M(table, col string, ids []int64, db *gorm.DB) error {
	sql, args, err := squirrel.Delete(table).Where(squirrel.Eq{col: ids}).ToSql()
	if err != nil {
		return err
	}
	if err := db.Exec(sql, args...).Error; err != nil {
		return err
	}
	return nil
}

func UpdateM2M(table, col1, col2 string, id int64, ids []int64, db *gorm.DB) error {
	if len(ids) == 0 {
		return nil
	}
	sql, args, err := squirrel.Delete(table).Where(squirrel.Eq{col1: id}).ToSql()
	if err != nil {
		return err
	}
	if err := db.Exec(sql, args...).Error; err != nil {
		return err
	}
	if ids[0] == 0 {
		return nil
	}
	query := squirrel.Insert(table).Columns(col1, col2)
	for _, id2 := range ids {
		query = query.Values(id, id2)
	}
	sql, args, err = query.ToSql()
	if err != nil {
		return err
	}
	return db.Exec(sql, args...).Error
}

//fix: id alanı gönderilince ne oluyor?
func MultipleInsert(db *gorm.DB, table string, st interface{}, cs interface{}) error {
	cats, _ := TakeSliceArg(cs)

	if len(cats) == 0 {
		return nil
	}

	sq := squirrel.
		Insert(table).
		Columns(sqlstruct.Columns(st))

	for _, cat := range cats {
		sq = sq.Values(sqlstruct.Values(cat)...)
	}

	sql, args, err := sq.ToSql()
	if err != nil {
		return err
	}

	return db.Exec(sql, args...).Error
}

func Count(db *gorm.DB, table string) (count int64, err error) {
	err = db.Table(table).Count(&count).Error
	return
}

func FoundErr(err error) error {
	if err == gorm.RecordNotFound {
		return nil
	}
	return err
}

func TakeSliceArg(arg interface{}) (out []interface{}, ok bool) {
	slice, success := takeArg(arg, reflect.Slice)
	if !success {
		ok = false
		return
	}
	c := slice.Len()
	out = make([]interface{}, c)
	for i := 0; i < c; i++ {
		out[i] = slice.Index(i).Interface()
	}
	return out, true
}

func takeArg(arg interface{}, kind reflect.Kind) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)
	if val.Kind() == kind {
		ok = true
	}
	return
}
