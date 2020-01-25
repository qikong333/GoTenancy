package GoTenancy

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/snowlyg/GoTenancy/model"
)

var (
	pageTemplates *template.Template
	languagePacks map[string]map[string]string
)

func init() {
	loadTemplates()
	loadLanguagePacks()
}

func loadTemplates() {
	var tmpl []string

	files, err := ioutil.ReadDir("./templates")
	if err != nil {
		if os.IsNotExist(err) == false {
			log.Fatal("unable to load templates", err)
		}
		return
	}

	for _, f := range files {
		tmpl = append(tmpl, path.Join("./templates", f.Name()))
	}

	t, err := template.New("").Funcs(template.FuncMap{
		"translate":  Translate,
		"translatef": Translatef,
		"money": func(amount int) string {
			m := float64(amount) / 100.0
			return fmt.Sprintf("%.2f $", m)
		},
	}).ParseFiles(tmpl...)

	if err != nil {
		log.Fatal("error while parsing templates", err)
	}

	pageTemplates = t
}

// HTML 模版将被保存到一个被 templates.ServePage 命名的目录中。
// 	func handler(w http.ResponseWriter, r *http.Request) {
// 		data := HomePage{Title: "Hello world!"}
// 		GoTenancy.ServePage(w, r, "index.html", data)
// 	}
func ServePage(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	t := pageTemplates.Lookup(name)

	if err := t.Execute(w, data); err != nil {
		fmt.Println("error while rendering the template ", err)
	}

	logRequest(r, http.StatusOK)
}

func loadLanguagePacks() {
	languagePacks = make(map[string]map[string]string)

	files, err := ioutil.ReadDir("./languagepacks")
	if err != nil {
		log.Println("unable to load language packs: ", err)
		return
	}

	var pack = new(struct {
		Language string `json:"lang"`
		Keys     []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"keys"`
	})

	for _, f := range files {
		b, err := ioutil.ReadFile(path.Join("./languagepacks", f.Name()))
		if err != nil {
			log.Fatal("unable to read language pack: ", f.Name(), ": ", err)
		}

		if err := json.Unmarshal(b, &pack); err != nil {
			log.Fatal("unable to parse language pack: ", f.Name(), ": ", err)
		}

		values := make(map[string]string)
		for _, k := range pack.Keys {
			values[k.Key] = k.Value
		}

		languagePacks[pack.Language] = values
	}
}

// Translate 在语言包文件中找到一个键（保存在名为languagepack的目录中）
// 并且返回一个 template.HTML 对像，所以可以安全地在语言包文件中使用 HTML 。
// 语言包文件都是简单以 lng.json 命名的 JSON 文件，例如 en.json：
// 	{
// 		"lang": "en",
// 		"keys": [
// 			{"key": "landing-title", "value": "Welcome to my site"}
// 		]
// 	}
func Translate(lng, key string) template.HTML {
	if s, ok := languagePacks[lng][key]; ok {
		return template.HTML(s)
	}
	return template.HTML(fmt.Sprintf("key %s not found", key))
}

// Translatef 匹配一个翻译键并且替代格式化参数。
func Translatef(lng, key string, a ...interface{}) string {
	if s, ok := languagePacks[lng][key]; ok {
		return fmt.Sprintf(s, a...)
	}
	return fmt.Sprintf("key %s not found", key)
}

// BUG(dom): 这里需要更多的思考
func ExtractLimitAndOffset(r *http.Request) (limit int, offset int) {
	limit = 50
	offset = 0

	p := r.URL.Query().Get("limit")
	if len(p) > 0 {
		i, err := strconv.Atoi(p)
		if err == nil {
			limit = i
		}
	}

	p = r.URL.Query().Get("offset")
	if len(p) > 0 {
		i, err := strconv.Atoi(p)
		if err == nil {
			offset = i
		}
	}

	return
}

// ViewData 是需要在所有页面中渲染的基础数据。
// 它将自动获取用户的语言，角色以及是否有警报显示。
// 你可以将它看作一个需要发送到页面渲染数据的包装器
type ViewData struct {
	Language string
	Role     model.Roles
	Alert    *Notification
	Data     interface{}
}

// Notification 在一个 HTML 模版中显示一个弹窗给用户。
type Notification struct {
	Title     template.HTML
	Message   template.HTML
	IsSuccess bool
	IsError   bool
	IsWarning bool
}

func getLanguage(ctx context.Context) string {
	lng, ok := ctx.Value(ContextLanguage).(string)
	if !ok {
		lng = "en"
	}
	return lng
}

func getRole(ctx context.Context) model.Roles {
	auth, ok := ctx.Value(ContextAuth).(Auth)
	if !ok {
		return model.RolePublic
	}
	return auth.Role
}

// CreateViewData 将数据包装为ViewData类型，
// 在该类型中，语言，角色和通知将与数据一起自动添加。
func CreateViewData(ctx context.Context, alert *Notification, data interface{}) ViewData {
	return ViewData{
		Alert:    alert,
		Data:     data,
		Language: getLanguage(ctx),
		Role:     getRole(ctx),
	}
}
