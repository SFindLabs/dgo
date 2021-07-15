package utils

import (
	kcode "dgo/work/code"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
)

/**
 *  分页方法
 * 	param c *gin.Context
 * 	param count int	                    总记录数
 * 	param params map[string]interface{} 额外的查询参数
 *  param option int 				    显示当前页条数(1: 显示  0: 不显示)、初始每页条数、下拉初始条数和每次选择间隔条数(可选)
 *  return 	分页html字符串, 当前路径(包含页数和每页条数,用于其他查询请求拼接), 当前页数, 每页条数
 */
func Paginate(c *gin.Context, count int, params map[string]interface{}, option ...int) (string, string, int, int) {

	toPage, total, pageSize := 1, 0, kcode.PAGE_NUMBER

	countOption := len(option)
	tmpCount, tmpPageSize, optionCount, step := 0, kcode.PAGE_NUMBER, 5, 5
	switch countOption {
	case 1:
		tmpCount = option[0]
	case 2:
		tmpCount = option[0]
		tmpPageSize = option[1]
	case 3:
		tmpCount = option[0]
		tmpPageSize = option[1]
		optionCount = option[2]
	case 4:
		tmpCount = option[0]
		tmpPageSize = option[1]
		optionCount = option[2]
		step = option[3]
	}

	if pageSize != tmpPageSize {
		pageSize = tmpPageSize
	}

	page, _ := c.GetQuery("page")
	if num, err := strconv.Atoi(page); err == nil && num > 1 {
		toPage = num
	}

	pageNum, _ := c.GetQuery("pageSize")
	if size, err := strconv.Atoi(pageNum); err == nil && size > 0 {
		pageSize = size
	}

	if count > 0 {
		total = int(math.Ceil(float64(count) / float64(pageSize)))
		if toPage > total {
			toPage = total
		}
	}

	prePage := toPage - 1
	nextPage := toPage + 1

	url := c.Request.URL.Path
	addParams := ""
	if len(params) > 0 {
		for k, v := range params {
			addParams += fmt.Sprint("&", k, "=", v)
		}
	}

	pageSizeStr := strconv.Itoa(pageSize)
	prePageStr := strconv.Itoa(prePage)
	toPageStr := strconv.Itoa(toPage)
	totalStr := strconv.Itoa(total)
	nextPageStr := strconv.Itoa(nextPage)

	html := `<ul class="pagination">`

	if toPage > 1 {
		html += `<li><a href="` + url + `?page=1&pageSize=` + pageSizeStr + addParams + `" rel="prev">首页</a></li>
				<li><a href="` + url + `?page=` + prePageStr + `&pageSize=` + pageSizeStr + addParams + `" rel="prev">上一页</a></li>
				<li><a href="` + url + `?page=` + prePageStr + `&pageSize=` + pageSizeStr + addParams + `">` + prePageStr + `</a></li>`
	} else {
		html += `<li class="disabled"><span>首页</span></li>
				<li class="disabled"><span>上一页</span></li>`
	}

	html += `<li class="active"><span>` + toPageStr + `</span></li>`

	if toPage < total {
		html += `<li><a href="` + url + `?page=` + nextPageStr + `&pageSize=` + pageSizeStr + addParams + `">` + nextPageStr + `</a></li>
				 <li><a href="` + url + `?page=` + nextPageStr + `&pageSize=` + pageSizeStr + addParams + `" rel="next">下一页</a></li>
			     <li><a href="` + url + `?page=` + totalStr + `&pageSize=` + pageSizeStr + addParams + `">尾页</a></li>`
	} else {
		html += `<li class="disabled"><span>下一页</span></li>
				<li class="disabled"><span>尾页</span></li>`
	}

	html += `<li>
				<span data-toggle="tooltip" title="选择每页条数" data-placement="bottom">
				<select aria-label="" style="appearance:none; -moz-appearance:none; -webkit-appearance:none; border: 0;outline:none;" class="text-center no-padding"  name="pageSize" id="setPageSize">`

	for i := 0; i < optionCount; i++ {
		optionValue := tmpPageSize + i*step
		selected := ""
		if pageSize == optionValue {
			selected = "selected"
		}
		html += fmt.Sprintf(`<option value="%d" %s >%d条/页</option>`, optionValue, selected, optionValue)
	}

	html += `</select>
			</li> 
			<li>
			 <span data-toggle="tooltip" data-placement="bottom" title="输入页码，鼠标失去焦点时快速跳转">
			<input type="number" class="text-center no-padding href_to" value="` + toPageStr + `" style="width: 50px; border: 0;outline:none;" aria-label=""> 页
			</li>
			<li class="disabled"><span>共` + totalStr + `页</span></li>
			`

	//增加本页条数显示
	if tmpCount == 1 {
		html += `<li class="disabled"><span>本页%d条</span></li>`
	}

	html += `</ul>
			<script>
				$("#setPageSize").on('change', function(){
					var size = $("option:selected", this).val();
					window.location.href = "` + url + `?page=1&pageSize=" + size + "` + addParams + `";
				});
				$(".href_to").on('blur',function(){
					var page = $(this).val();
					window.location.href = "` + url + `?page=" + page + "&pageSize=` + pageSizeStr + addParams + `";
				});
			</script>`

	toUrl := url + `?page=` + toPageStr + `&pageSize=` + pageSizeStr

	return html, toUrl, toPage, pageSize
}
