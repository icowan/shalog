package templates

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/flosch/pongo2"
	"github.com/icowan/shalog/src/repository"
	"github.com/russross/blackfriday"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var (
	tplExt         = ".html"
	templatesCache = make(map[string]*pongo2.Template)
)

func init() {
	if err := pongo2.RegisterFilter("markdown", func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
		return pongo2.AsSafeValue(string(blackfriday.Markdown([]byte(in.String()), blackfriday.HtmlRenderer(0, "", "markdown-body"), blackfriday.EXTENSION_NO_INTRA_EMPHASIS|
			blackfriday.EXTENSION_TABLES|
			blackfriday.EXTENSION_FENCED_CODE|
			blackfriday.EXTENSION_AUTOLINK|
			blackfriday.EXTENSION_STRIKETHROUGH|
			blackfriday.EXTENSION_SPACE_HEADERS|
			blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK))), nil
	}); err != nil {
		fmt.Println("err", err.Error())
	}

	if err := pongo2.RegisterFilter("toString", func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
		return pongo2.AsValue(strconv.Itoa(in.Integer())), nil
	}); err != nil {
		fmt.Println("err", err.Error())
	}

	if err := pongo2.RegisterFilter("str2html", func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
		return pongo2.AsValue(in), nil
	}); err != nil {
		fmt.Println("err", err.Error())
	}
}

func Render(data map[string]interface{}, body io.Writer, tplName string) error {

	tpl, ok := templatesCache[tplName]
	if !ok {
		tpl = pongo2.Must(pongo2.FromFile(tplName + tplExt))
		//templatesCache[tplName] = tpl
	}
	action := strings.Split(tplName, "/")
	if data == nil {
		data = map[string]interface{}{
			"action": action[2],
		}
	}

	b, _ := json.Marshal(data)

	var ctxData pongo2.Context
	if err := json.Unmarshal(b, &ctxData); err != nil {
		fmt.Println("err", err.Error())
	}

	if err := tpl.ExecuteWriter(ctxData, body); err != nil {
		return err
	}

	return nil
}

func RenderHtml(ctx context.Context, w http.ResponseWriter, response map[string]interface{}) error {
	name := ctx.Value("method").(string)
	var viewPath string
	if settings, ok := ctx.Value("settings").(map[string]string); ok {
		if response == nil {
			response = make(map[string]interface{})
		}
		for k, v := range settings {
			response[strings.ReplaceAll(k, "-", "_")] = v
		}
		viewPath = settings[repository.SettingViewTemplate.String()]
	}

	buf := new(bytes.Buffer)
	if err := Render(response, buf, viewPath+name); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(buf.Bytes())); err != nil {
		return err
	}

	return nil
}
