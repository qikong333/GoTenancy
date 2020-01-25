package GoTenancy

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Respond 返回一个带有特殊状态的 strruct JSON。
//
// 如果数据是错误它将会被打包成一个常见的 JSON 对象：
//
// 	{
// 		"status": 401,
// 		"error": "the result of data.Error()"
// 	}
//
// 使用实例:
//
// 	func handler(w http.ResponseWriter, r *http.Request) {
// 		task := Task{ID: 123, Name: "My Task", Done: false}
// 		GoTenancy.Respond(w, r, http.StatusOK, task)
// 	}
func Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) error {
	// 转换错误到一个真实的 JSON 序列化对象中
	if e, ok := data.(error); ok {
		var tmp = new(struct {
			Status string `json:"status"`
			Error  string `json:"error"`
		})
		tmp.Status = "error"
		tmp.Error = e.Error()
		data = tmp

		log.Println("error: ", e.Error())
	}

	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// 获取请求 ID
	reqID, ok := r.Context().Value(ContextRequestID).(string)
	if ok {
		w.Header().Set("X-Request-ID", reqID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	logRequest(r, status)

	return nil
}

// ParseBody 解析请求的 JSON 主体到一个 struct.ParseBody 中。
//
// 	func handler(w http.ResponseWriter, r *http.Request) {
// 		var task Task
// 		if err := GoTenancy.ParseBody(r.Body, &task); err != nil {
// 			GoTenancy.Respond(w, r, http.StatusBadRequest, err)
// 			return
// 		}
// 	}
func ParseBody(body io.ReadCloser, result interface{}) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(result)
}
