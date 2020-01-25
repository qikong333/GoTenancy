package queue

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// TaskID 是代表特定排队任务的ID。
type TaskID int

const (
	// TaskEmail 发送邮件
	TaskEmail TaskID = iota
	// TaskCreateInvoice 创建条码发票
	TaskCreateInvoice
)

// QueueTask 表示一个任务队列
//
// 数据字段包含必要的数据，任务才能正确执行。
type QueueTask struct {
	ID      TaskID      `json:"id"`
	Data    interface{} `json:"data"`
	Created time.Time   `json:"created"`
}

// TaskExecutor 是一个基于任务 id 去执行任务的接口
type TaskExecutor interface {
	Run(t QueueTask) error
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		invalidTypeError := errors.New("provided value type didn't match obj field type")
		return invalidTypeError
	}

	structFieldValue.Set(val)
	return nil
}

func fillStruct(s interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := setField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
