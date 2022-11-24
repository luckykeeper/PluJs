// pluJs ，生成 Plumemo 前端文件，批量添加轮播图
// Powered By Luckykeeper <luckykeeper@luckykeeper.site | https://luckykeeper.site> 2022/11/19
// 考虑 memo 酱换成 momo 酱（MomoTalk）

package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xuri/excelize/v2"
	"gopkg.in/ini.v1"
)

// 显示中文
// 设置环境变量   通过go-findfont 寻找simkai.ttf 字体
func init() {
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "simkai.ttf") {
			fmt.Println(path)
			os.Setenv("FYNE_FONT", path) // 设置环境变量  // 取消环境变量 os.Unsetenv("FYNE_FONT")
			break
		}
	}
}

var (
	Content              *fyne.Container
	ContentTitle         *widget.Label
	ContentTitleBox      *fyne.Container
	contentBox           *fyne.Container
	RunTime              = 0
	LayoutChanged        = false
	sql_initialize_table = `CREATE TABLE file (url TEXT PRIMARY KEY NOT NULL,fileName TEXT NOT NULL,type TEXT NOT NULL,md5 TEXT NOT NULL);`
)

func main() {
	// 先初始化数据库
	if exists, _ := PathExists("./pluJs.db"); exists {
		log.Println("数据库存在")
	} else {
		log.Println("创建数据库！")
		db, _ := sql.Open("sqlite3", "./pluJs.db")
		defer db.Close()
		db.Exec(sql_initialize_table) //初始化
	}

	// App 基本信息
	a := app.NewWithID("plujs.luckykeeper.site")
	logo, _ := fyne.LoadResourceFromPath("plujs_icon.ico")
	a.SetIcon(logo)
	makeTray(a)
	logLifecycle(a)
	w := a.NewWindow("PluJs, A software to generate plumemo frontend JavaScript File | Powered by Luckykeeper | Build 20221119 | Ver 1.0.0")
	w.SetMainMenu(makeMenu(a, w))

	// 左侧菜单
	menu := container.NewVBox(
		widget.NewButtonWithIcon("Welcome PluJs!",
			theme.HomeIcon(),
			welcomeScreen),
		widget.NewButtonWithIcon("给momo酱喂食！",
			theme.DeleteIcon(),
			addImage),
		widget.NewButtonWithIcon("momo酱今天吃什么？",
			theme.SearchReplaceIcon(),
			database),
		// 在这里添加换旧、删除功能【给momo酱换道菜】
		widget.NewButtonWithIcon("连接momo酱到心智云图数据库！",
			theme.StorageIcon(),
			couldDatabase),
		widget.NewButtonWithIcon("再见，momo酱！",
			theme.VisibilityOffIcon(),
			func() { fyne.App.Quit(a) }),
	)

	left := container.New(layout.NewHBoxLayout(), menu, widget.NewSeparator())

	go func() {
		for {
			time.Sleep(time.Millisecond * 500)
			if RunTime == 0 {
				ContentTitle = widget.NewLabel("雷猴哇~<(￣︶￣)↗[GO!]")
				ContentTitleBox = container.New(layout.NewVBoxLayout(), ContentTitle, widget.NewSeparator())

				addImageIcon := canvas.NewImageFromFile("./icon/top.jpg")
				addImageIcon.FillMode = canvas.ImageFillContain
				addImageIcon.SetMinSize(fyne.NewSize(290.75, 500))

				Content = container.NewCenter(container.NewVBox(
					widget.NewLabelWithStyle("↑Shirasu Azusa, A student from Blue Archive", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
					widget.NewLabelWithStyle("Welcome to PluJs, A software to generate plumemo frontend JavaScript File", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
					container.NewHBox(
						widget.NewHyperlink("Powered By Luckykeeper", parseURL("https://luckykeeper.site/")),
						widget.NewLabel("-"),
						widget.NewHyperlink("Github", parseURL("https://github.com/luckykeeper/")),
						widget.NewLabel("-"),
						widget.NewHyperlink("Blog", parseURL("https://luckykeeper.site/")),
					),
				))
				Content = container.New(layout.NewVBoxLayout(), addImageIcon, Content)

			}
			contentBox = container.New(layout.NewBorderLayout(ContentTitleBox, nil, nil, nil), ContentTitleBox, Content)
			if RunTime == 0 || LayoutChanged {
				contentBox.Refresh()
				RunTime = 1
				LayoutChanged = false
			}
			// 显示主界面，分别：适应宽度，左侧菜单，分割线，右侧内容
			w.SetContent(container.New(layout.NewBorderLayout(nil, nil, left, nil), left, contentBox))
		}
	}()

	// 设置窗口大小
	w.Resize(fyne.NewSize(1280, 720))
	w.SetFixedSize(true)
	// 润！
	w.ShowAndRun()

}

// 欢迎界面
func welcomeScreen() {
	ContentTitle = widget.NewLabel("雷猴哇~<(￣︶￣)↗[GO!]")
	ContentTitleBox = container.New(layout.NewVBoxLayout(), ContentTitle, widget.NewSeparator())

	addImageIcon := canvas.NewImageFromFile("./icon/top.jpg")
	addImageIcon.FillMode = canvas.ImageFillContain
	addImageIcon.SetMinSize(fyne.NewSize(290.75, 500))

	Content = container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("↑Shirasu Azusa, A student from Blue Archive", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Welcome to PluJs, A software to generate plumemo frontend JavaScript File", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewHBox(
			widget.NewHyperlink("Powered By Luckykeeper", parseURL("https://luckykeeper.site/")),
			widget.NewLabel("-"),
			widget.NewHyperlink("Github", parseURL("https://github.com/luckykeeper/")),
			widget.NewLabel("-"),
			widget.NewHyperlink("Blog", parseURL("https://luckykeeper.site/")),
		),
	))

	Content = container.New(layout.NewVBoxLayout(), addImageIcon, Content)

	contentBox = container.New(layout.NewBorderLayout(ContentTitleBox, nil, nil, nil), ContentTitleBox, Content)
	LayoutChanged = true
}

// 添加界面
func addImage() {
	ContentTitle = widget.NewLabel("momo酱今天吃嘉然！ヾ(≧▽≦*)o")
	ContentTitleBox = container.New(layout.NewVBoxLayout(), ContentTitle, widget.NewSeparator())

	// 添加图片小技巧：container 套 container ，注意 layout.NewVBoxLayout （或者 HBox ）的时候会按照组件最小大小来排列
	// 下面带字的会自动计算最小大小，但是你加的图片不会，所以你需要手动给它一个大小，不然就会被压成一个 1x1 的像素点（乐）
	addImageIcon := canvas.NewImageFromFile("./icon/eat.png")
	addImageIcon.FillMode = canvas.ImageFillContain
	addImageIcon.SetMinSize(fyne.NewSize(232.5, 350.8))

	Content = container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle("momo酱今天吃点什么好呢？", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewButton("momo酱今天恰_____",
				addImageFromTable,
			),
			widget.NewLabelWithStyle("提示，在 addImage.xlsx 文件中添加外链图片链接", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})),
	)

	Content = container.New(layout.NewVBoxLayout(), addImageIcon, Content)

	contentBox = container.New(layout.NewBorderLayout(ContentTitleBox, nil, nil, nil), ContentTitleBox, Content)

	LayoutChanged = true
}

// 添加界面->添加
// 说明， addImage.xlsx 文件中图片分三类，一类是只在轮播图中使用的，一类是只在文章图使用的，一类是两边都使用的（共通图）
func addImageFromTable() {
	addImageTable, _ := excelize.OpenFile("./addImage.xlsx")
	if exists, _ := PathExists("./conflict.xlsx"); exists {
		os.Remove("./conflict.xlsx")
	}
	f := excelize.NewFile()
	f.SaveAs("./conflict.xlsx")

	fileType := "banner"
	checkAndWrite(addImageTable, fileType, "轮播图")

	fileType = "article"
	checkAndWrite(addImageTable, fileType, "文章图")

	fileType = "all"
	checkAndWrite(addImageTable, fileType, "共通图")
	os.Remove("./addImage.xlsx")
	copy("addImage_Template.xlsx", "./addImage.xlsx", 64)
	// temp 文件夹一锅端
	os.RemoveAll("./temp")
	// insertSql := "insert into file (id, fileName,type) values ('" + id + "','" + upFileName + "','" + fileType + "');"
	// log.Println("insertSql:", insertSql)
	// db.Exec(insertSql)
}

// [DB 路径可写死]

// 检查并写入
func checkAndWrite(addImageTable *excelize.File, fileType, sheet string) {
	for line := 2; line >= 2; line++ {
		urlRaw1, _ := addImageTable.GetCellValue(sheet, "A"+strconv.Itoa(line))

		if urlRaw1 == "" {
			line = 0
		} else if strings.Contains(urlRaw1, ": http") {
			log.Println("Cloudreve!")
			fileName := strings.Split(strings.Split(urlRaw1, ": ")[0], ".")[0]
			url := strings.Split(urlRaw1, ": ")[1]
			log.Println("fileName:", fileName)
			log.Println("url:", url)
			fileHash := getFileMd5(url, fileName)
			log.Println("MD5:", fileHash)
			// MD5 存在，导出冲突表
			log.Println("_____________")
			if DataExists(url, fileName, fileHash) {
				exportConflict(url, fileName, fileHash)
			} else { // MD5 不存在，写入数据库
				writeToSqlite(url, fileName, fileType, fileHash)
			}
		} else {
			log.Println("Direct Link!")
			fileName, _ := addImageTable.GetCellValue(sheet, "B"+strconv.Itoa(line))
			url := urlRaw1
			log.Println("fileName:", fileName)
			log.Println("url:", url)
			fileHash := getFileMd5(url, fileName)
			log.Println("MD5:", fileHash)
			log.Println("_____________")
			// MD5 存在，导出冲突表
			if DataExists(url, fileName, fileHash) {
				exportConflict(url, fileName, fileHash)
			} else { // MD5 不存在，写入数据库
				writeToSqlite(url, fileName, fileType, fileHash)
			}
		}

	}
}

// 获取在线文件的 MD5 值，以便比较
func getFileMd5(url, name string) (hash string) {
	if exists, _ := PathExists("./temp"); !exists {
		os.Mkdir("./temp", os.ModePerm)
	}

	// 下载指定 URL 文件，准备计算 MD5
	resp, err := http.Get(url)
	if err != nil {
		log.Println("respErr: ", err)
	}
	defer resp.Body.Close()

	// 创建一个文件用于保存
	saveTempFile, _ := os.Create("./temp/" + name)
	defer saveTempFile.Close()

	// 然后将响应流和文件流对接起来
	io.Copy(saveTempFile, resp.Body)

	f, _ := os.Open("./temp/" + name)
	defer f.Close()

	h := md5.New()

	io.Copy(h, f)
	hash = fmt.Sprintf("%x", h.Sum(nil))
	return hash
}

// 判断 MD5 在文件中是否存在，存在输出同 MD5 信息（Excel）
func DataExists(urlNew, fileNameNew, fileHash string) bool {
	db, _ := sql.Open("sqlite3", "./pluJs.db")
	defer db.Close()
	querySql := "select url,filename from file where md5='" + fileHash + "';"
	var data, data1 string
	queryResult := db.QueryRow(querySql).Scan(&data, &data1)
	if queryResult == sql.ErrNoRows {
		return false
	} else { // 数据库存在同 MD5 文件
		return true
	}
}

// 冲突文件导出
func exportConflict(urlNew, fileNameNew, fileHash string) {
	f, _ := excelize.OpenFile("./conflict.xlsx")
	f.NewSheet("冲突表")
	// 设置单元格的值
	f.SetCellValue("冲突表", "A1", "已存在文件URL")
	f.SetCellValue("冲突表", "B1", "已存在文件名称")
	f.SetCellValue("冲突表", "C1", "冲突文件URL")
	f.SetCellValue("冲突表", "D1", "冲突文件名称")
	f.SetColWidth("冲突表", "B", "B", 80)
	f.SetColWidth("冲突表", "D", "D", 80)
	f.SetColWidth("冲突表", "A", "A", 20)
	f.SetColWidth("冲突表", "C", "C", 20)
	// 数据库读数据
	db, _ := sql.Open("sqlite3", "./pluJs.db")
	defer db.Close()
	rows, _ := db.Query("select url,fileName from file where md5='" + fileHash + "';")

	// 寻找最下一行
	Line := 2
	for Line >= 2 {
		excelContent, _ := f.GetCellValue("冲突表", "A"+strconv.Itoa(Line))

		if excelContent == "" {
			break
		} else {
			Line++
		}
	}

	for rows.Next() {
		var urlExists string
		var fileNameExists string
		rows.Scan(&urlExists, &fileNameExists)
		f.SetCellValue("冲突表", fmt.Sprintf("A%d", Line), urlExists)
		f.SetCellValue("冲突表", fmt.Sprintf("B%d", Line), fileNameExists)
		f.SetCellValue("冲突表", fmt.Sprintf("C%d", Line), urlNew)
		f.SetCellValue("冲突表", fmt.Sprintf("D%d", Line), fileNameNew)
		Line++
	}
	f.DeleteSheet("Sheet1")
	f.SaveAs("./conflict.xlsx")
}

// 写入不冲突数据到数据库
func writeToSqlite(url, fileName, fileType, fileHash string) {
	db, _ := sql.Open("sqlite3", "./pluJs.db")
	defer db.Close()
	insertSql := "insert into file (url, fileName,type,md5) values ('" + url + "','" + fileName + "','" + fileType + "','" + fileHash + "');"
	// log.Println("insertSql:", insertSql)
	db.Exec(insertSql)
}

// 判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 库存界面
func database() {
	ContentTitle = widget.NewLabel("momo酱吃不下了啦！o(*≧▽≦)ツ┏━┓")
	ContentTitleBox = container.New(layout.NewVBoxLayout(), ContentTitle, widget.NewSeparator())

	addImageIcon := canvas.NewImageFromFile("./icon/change.jpg")
	addImageIcon.FillMode = canvas.ImageFillContain
	addImageIcon.SetMinSize(fyne.NewSize(265.5, 375.5))

	Content = container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("momo酱吃不下了啦！", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewHBox(
			widget.NewButtonWithIcon("给我恰！",
				theme.DownloadIcon(),
				generateFrontendMenu,
			),
			widget.NewButtonWithIcon("塞回去！",
				theme.UploadIcon(),
				updateFromExcel,
			),
		)))
	Content = container.New(layout.NewVBoxLayout(), addImageIcon, Content)

	contentBox = container.New(layout.NewBorderLayout(ContentTitleBox, nil, nil, nil), ContentTitleBox, Content)
	LayoutChanged = true
}

// 导出前端文件及Excel
func generateFrontendMenu() {
	ContentTitle = widget.NewLabel("呜呜呜o(一︿一+)o")
	ContentTitleBox = container.New(layout.NewVBoxLayout(), ContentTitle, widget.NewSeparator())

	addImageIcon := canvas.NewImageFromFile("./icon/change.jpg")
	addImageIcon.FillMode = canvas.ImageFillContain
	addImageIcon.SetMinSize(fyne.NewSize(265.5, 375.5))

	Content = container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("呜呜呜o(一︿一+)o", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("已导出前端文件及Excel！", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	))
	Content = container.New(layout.NewVBoxLayout(), addImageIcon, Content)
	contentBox = container.New(layout.NewBorderLayout(ContentTitleBox, nil, nil, nil), ContentTitleBox, Content)
	LayoutChanged = true

	if exists, _ := PathExists("./fileList.xlsx"); exists {
		os.Remove("./fileList.xlsx")
	}
	generateFrontend()
}

// 实行生成前端文件及Excel
func generateFrontend() {

	if exists, _ := PathExists("./conflict.xlsx"); exists {
		os.Remove("./conflict.xlsx")
	} else if exists, _ := PathExists("./fileList.xlsx"); exists {
		os.Remove("./fileList.xlsx")
	}

	// 生成前端文件
	frontendBeforeRaw, _ := ioutil.ReadFile("frontendBefore.txt")
	frontendBefore := string(frontendBeforeRaw)
	frontendAfterRaw, _ := ioutil.ReadFile("frontendAfter.txt")
	frontendAfter := string(frontendAfterRaw)

	db, _ := sql.Open("sqlite3", "./pluJs.db")
	defer db.Close()

	// // 查询需要写入的条数
	// queryBannerCountSql := "SELECT COUNT(*) FROM (SELECT type FROM file WHERE type='banner');"
	// queryArticleCountSql := "SELECT COUNT(*) FROM (SELECT type FROM file WHERE type='article');"
	// queryAllCountSql := "SELECT COUNT(*) FROM (SELECT type FROM file WHERE type='all');"

	// var queryBannerResult, queryArticleResult, queryAllCountResult int
	// db.QueryRow(queryBannerCountSql).Scan(&queryBannerResult)
	// db.QueryRow(queryArticleCountSql).Scan(&queryArticleResult)
	// db.QueryRow(queryAllCountSql).Scan(&queryAllCountResult)

	// bannerNumber := queryBannerResult + queryAllCountResult
	// articleNumber := queryArticleResult + queryAllCountResult

	// ListImg: 在前，文章随机图
	listImgReturnBefore := "ListImg:["
	listImgReturnAfter := "],"
	ListImgReturnInside := ""

	// bannerList: 在后，轮播随机图
	bannerImgReturnBefore := "bannerList: ["
	bannerImgReturnAfter := "]"
	BannerImgReturnInside := ""

	// 查询语句
	queryAllUrlRows, _ := db.Query("SELECT url FROM file WHERE type='all';")
	queryArticleUrlRows, _ := db.Query("SELECT url FROM file WHERE type='article';")
	queryBannerUrlRows, _ := db.Query("SELECT url FROM file WHERE type='banner';")

	// 生成 ListImg
	for queryAllUrlRows.Next() {
		var url string
		queryAllUrlRows.Scan(&url)
		ListImgReturnInside = ListImgReturnInside + "{img:\"" + url + "\"},"
		BannerImgReturnInside = BannerImgReturnInside + "{img:\"" + url + "\"},"
	}

	for queryArticleUrlRows.Next() {
		var url string
		queryArticleUrlRows.Scan(&url)
		ListImgReturnInside = ListImgReturnInside + "{img:\"" + url + "\"},"
	}

	listImgReturn := listImgReturnBefore + strings.TrimRight(ListImgReturnInside, ",") + listImgReturnAfter

	// 生成 bannerList
	for queryBannerUrlRows.Next() {
		var url string
		queryBannerUrlRows.Scan(&url)
		BannerImgReturnInside = BannerImgReturnInside + "{img:\"" + url + "\"},"
	}

	bannerListReturn := bannerImgReturnBefore + strings.TrimRight(BannerImgReturnInside, ",") + bannerImgReturnAfter

	// 随机图部分
	randomImgList := listImgReturn + bannerListReturn

	// 组织输出文本
	frontendStrings := frontendBefore + randomImgList + frontendAfter

	// 输出到 JavaScript 文件
	frontendFile, _ := os.Create("./main.3a574d82.chunk.js")
	frontendFile.WriteString(frontendStrings)

	log.Println("前端 JS 文件输出完成!")

	// 输出到 Excel 文件
	exportExcel()

	log.Println("Excel 文件输出完成")
}

// 导出到 Excel
func exportExcel() {
	f := excelize.NewFile()
	db, _ := sql.Open("sqlite3", "./pluJs.db")
	defer db.Close()
	queryFileByTypeAndWriteToExcel(f, db, "banner", "轮播图")
	queryFileByTypeAndWriteToExcel(f, db, "article", "文章图")
	queryFileByTypeAndWriteToExcel(f, db, "all", "共通图")
	f.DeleteSheet("Sheet1")
	f.SaveAs("./fileList.xlsx")
}

// 根据文件类型查询数据
func queryFileByTypeAndWriteToExcel(f *excelize.File, db *sql.DB, fileType, sheetName string) {
	f.NewSheet(sheetName)
	// 设置单元格的值
	f.SetCellValue(sheetName, "A1", "文件名")
	f.SetCellValue(sheetName, "B1", "文件URL")
	f.SetCellValue(sheetName, "C1", "状态(dd-删除;cg-更改)")
	f.SetCellValue(sheetName, "D1", "新文件名")
	f.SetCellValue(sheetName, "E1", "新URL")
	f.SetColWidth(sheetName, "A", "E", 40)
	// 数据库读数据
	rows, _ := db.Query("SELECT url,fileName FROM file WHERE type = '" + fileType + "';")
	line := 2
	for rows.Next() {
		var url string
		var fileName string
		rows.Scan(&url, &fileName)

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", line), fileName)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", line), url)
		line++
	}
}

// 从输出的 Excel 文件更改内容
func updateFromExcel() {
	ContentTitle = widget.NewLabel("呜呜呜o(一︿一+)o")
	ContentTitleBox = container.New(layout.NewVBoxLayout(), ContentTitle, widget.NewSeparator())

	addImageIcon := canvas.NewImageFromFile("./icon/change.jpg")
	addImageIcon.FillMode = canvas.ImageFillContain
	addImageIcon.SetMinSize(fyne.NewSize(265.5, 375.5))

	Content = container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("呜呜呜o(一︿一+)o", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("已经塞回到数据库啦！", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	))
	Content = container.New(layout.NewVBoxLayout(), addImageIcon, Content)
	contentBox = container.New(layout.NewBorderLayout(ContentTitleBox, nil, nil, nil), ContentTitleBox, Content)
	LayoutChanged = true

	updateToSqlite("轮播图")
	updateToSqlite("文章图")
	updateToSqlite("共通图")
}

// 从表格导回数据库
func updateToSqlite(sheet string) {
	updateTable, _ := excelize.OpenFile("./fileList.xlsx")
	db, _ := sql.Open("sqlite3", "./pluJs.db")
	defer db.Close()

	for line := 2; line >= 2; line++ {
		oldUrl, _ := updateTable.GetCellValue(sheet, "B"+strconv.Itoa(line))
		status, _ := updateTable.GetCellValue(sheet, "C"+strconv.Itoa(line))

		if oldUrl == "" {
			line = 0
			// 删除一条记录
		} else if status == "dd" {
			log.Println("Delete!")
			deleteSql := "DELETE FROM file WHERE url='" + oldUrl + "';"
			db.Exec(deleteSql)

			log.Println("Delete complete!")
			log.Println("_____________")
			// 修改一条记录
		} else if status == "cg" {
			newfileName, _ := updateTable.GetCellValue(sheet, "D"+strconv.Itoa(line))
			newUrl, _ := updateTable.GetCellValue(sheet, "E"+strconv.Itoa(line))
			updateSql := "UPDATE file SET url='" + newUrl + "',fileName='" + newfileName + "' WHERE url='" + oldUrl + "';"
			db.Exec(updateSql)
			log.Println("Update complete!")
			log.Println("_____________")
		}
	}
}

// 备份云端
func couldDatabase() {
	ContentTitle = widget.NewLabel("心智云图数据库操作中心")
	ContentTitleBox = container.New(layout.NewVBoxLayout(), ContentTitle, widget.NewSeparator())

	addImageIcon := canvas.NewImageFromFile("./icon/sync.jpg")
	addImageIcon.FillMode = canvas.ImageFillContain
	addImageIcon.SetMinSize(fyne.NewSize(265.5, 375.5))

	Content = container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Sensei，请在config.ini文件配置心智云图数据库好相关参数，再使用云图数据库！", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewHBox(
			widget.NewButtonWithIcon("上传momo酱的心智数据到心智云图数据库",
				theme.UploadIcon(),
				uploadToCloud,
			),
			widget.NewButtonWithIcon("从心智云图数据库下载momo酱的心智数据",
				theme.DownloadIcon(),
				downloadToLocal,
			),
		),
	))
	Content = container.New(layout.NewVBoxLayout(), addImageIcon, Content)
	contentBox = container.New(layout.NewBorderLayout(ContentTitleBox, nil, nil, nil), ContentTitleBox, Content)
	LayoutChanged = true
}

// 从云端恢复momo酱的心智数据
func downloadToLocal() {
	dbAddress, dbPort, dbUsername, dbPassword, dbName := readIniConfig()
	// 删除已经存在的数据（当然也不应该存在）
	if exists, _ := PathExists("./pluJs.db"); exists {
		os.Remove("./pluJs.db")
	}
	// 创建 sqlite 并初始化
	sqliteDb, _ := sql.Open("sqlite3", "./pluJs.db")
	sqliteDb.Exec(sql_initialize_table)
	pgsqlDb, _ := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbAddress, dbPort, dbUsername, dbPassword, dbName))
	defer sqliteDb.Close()
	defer pgsqlDb.Close()
	rows, _ := pgsqlDb.Query("SELECT * FROM file;")
	for rows.Next() {
		var url string
		var fileName string
		var fileType string
		var md5 string
		rows.Scan(&url, &fileName, &fileType, &md5)
		log.Println("同步完成:", fileName)
		sqliteDb.Exec("insert into file (url, fileName,type,md5) values ('" + url + "','" + fileName + "','" + fileType + "','" + md5 + "');")
	}
	log.Println("心智云图恢复完成！")
}

// 上传momo酱的心智数据到云端
func uploadToCloud() {
	dbAddress, dbPort, dbUsername, dbPassword, dbName := readIniConfig()
	connectPgsqldb(dbAddress, dbPort, dbUsername, dbPassword, dbName)
	sqliteDb, _ := sql.Open("sqlite3", "./pluJs.db")
	pgsqlDb, _ := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbAddress, dbPort, dbUsername, dbPassword, dbName))
	defer sqliteDb.Close()
	defer pgsqlDb.Close()

	rows, _ := sqliteDb.Query("SELECT * FROM file;")
	for rows.Next() {
		var url string
		var fileName string
		var fileType string
		var md5 string
		rows.Scan(&url, &fileName, &fileType, &md5)
		log.Println("同步完成:", fileName)
		pgsqlDb.Exec("INSERT INTO file(url, fileName,type,md5) values ('" + url + "','" + fileName + "','" + fileType + "','" + md5 + "');")
	}
	log.Println("心智云图上传完成！")
}

// 数据库初始化（大象数据库）
func connectPgsqldb(dbAddress string, dbPort string, dbUsername string, dbPassword string, dbName string) {
	// 数据库初始连接
	db, _ := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		dbAddress, dbPort, dbUsername, dbPassword))
	defer db.Close()
	checkInitialize, _ := db.Query(fmt.Sprintf((`SELECT u.datname FROM pg_catalog.pg_database u where u.datname='%s';`),
		dbName))
	// 查询结果，存在 true ，不存在 false
	checkInitializeResult := checkInitialize.Next()

	// 不存在就创建，存在就先扬了云端数据库再创建
	if !checkInitializeResult {
		createPgsqlDatabase(db, dbAddress, dbPort, dbUsername, dbPassword, dbName)
	} else {
		db.Exec("DROP DATABASE IF EXISTS " + dbName + ";")
		// 修复不能删除旧表的问题
		dropFile, _ := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbAddress, dbPort, dbUsername, dbPassword, dbName))
		defer dropFile.Close()
		dropFile.Exec("DROP TABLE IF EXISTS file;")

		createPgsqlDatabase(db, dbAddress, dbPort, dbUsername, dbPassword, dbName)
	}
}

// pgsql 建表
func createPgsqlDatabase(db *sql.DB, dbAddress string, dbPort string, dbUsername string, dbPassword string, dbName string) {
	db.Exec(fmt.Sprintf((`CREATE DATABASE %s;`), dbName))
	pluJsDb, _ := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbAddress, dbPort, dbUsername, dbPassword, dbName))
	defer pluJsDb.Close()
	pluJsDb.Exec(`CREATE TABLE file(
			url TEXT NOT NULL,
			fileName TEXT NOT NULL,
			type TEXT NOT NULL,
			md5 TEXT NOT NULL
			)WITH (OIDS=FALSE);`)
}

// 读取心智云图数据库配置
func readIniConfig() (dbAddress string, dbPort string, dbUsername string, dbPassword string, dbName string) {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	dbAddress = cfg.Section("cloudDatabase").Key("address").String()
	dbPort = cfg.Section("cloudDatabase").Key("port").String()
	dbUsername = cfg.Section("cloudDatabase").Key("username").String()
	dbPassword = cfg.Section("cloudDatabase").Key("password").String()
	dbName = cfg.Section("cloudDatabase").Key("dbName").String()
	// 打印配置数据
	// log.Println("心智云图数据库配置：")
	// log.Println("dbAddress:", dbAddress)
	// log.Println("dbPort:", dbPort)
	// log.Println("dbUsername:", dbUsername)
	// log.Println("dbPassword:", dbPassword)
	// log.Println("dbName:", dbName)
	return
}

// 生命周期日志
func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}

// 顶部菜单
func makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	aboutMenu := fyne.NewMenu("关于",
		fyne.NewMenuItem("访问 Blog", func() {
			u, _ := url.Parse("https://luckykeeper.site")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("访问 Github", func() {
			u, _ := url.Parse("https://github.com/luckykeeper")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("访问 Plumemo 项目", func() {
			u, _ := url.Parse("https://github.com/byteblogs168/plumemo")
			_ = a.OpenURL(u)
		}))

	main := fyne.NewMainMenu(
		aboutMenu,
	)
	return main
}

// 任务栏托盘
func makeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem("PluJs By Luckykeeper", func() {})
		menu := fyne.NewMenu("Hello World", h)
		h.Action = func() {
			log.Println("Hi there!")
			h.Label = "PluJs By Luckykeeper"
			u, _ := url.Parse("https://github.com/luckykeeper")
			a.OpenURL(u)
			menu.Refresh()
		}
		desk.SetSystemTrayMenu(menu)
	}
}

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

// 复制文件
func copy(src, dst string, BUFFERSIZE int64) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		log.Println("src is not a regular file.", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		log.Println("dst is not a regular file.", src)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}

	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}
