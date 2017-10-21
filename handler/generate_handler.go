package handler

import (
	"github.com/gogank/papillon/configuration"
	"github.com/gogank/papillon/utils"
	"errors"
	"github.com/gogank/papillon/render"
	"path"
	"fmt"
	"strings"
	"time"
	"strconv"
	"math/rand"
)

func Generate(conf_path string) error{
	config := config.NewConfig(conf_path)

	sourceDir := config.GetString(utils.DIR_SOURCE)
	postsDir  := config.GetString(utils.DIR_POSTS)
	publicDir := config.GetString(utils.DIR_PUBLIC)
	themeDir := config.GetString(utils.DIR_THEME)

	if !utils.ExistDir(sourceDir) {
		return errors.New(fmt.Sprintf("source directory '%s' doesn't exist, cann't generate", sourceDir))
	}

	if utils.ExistDir(publicDir) {
		if err := utils.RemoveDir(publicDir); err != nil {
			return err
		}
	}

	// create public dir
	if !utils.Mkdir(publicDir) {
		return errors.New(fmt.Sprintf("create directory %s failed", publicDir))
	}

	if utils.ExistDir(postsDir) {

		// create public/posts dir
		if !utils.Mkdir(path.Join(publicDir, "posts")) {
			return errors.New(fmt.Sprintf("create directory %s failed", path.Join(publicDir, "posts")))
		}

		// 遍历source/posts/ 目录中的所有的markdown文件
		files, err := utils.ListDir(postsDir, "md")
		if err != nil {
			return err
		}

		parse := render.New()

		for _, fname := range files {
			mdContent, err := utils.ReadFile(path.Join(postsDir, fname))
			if err != nil {
				return err
			}

			postsTpl, err := utils.ReadFile(path.Join(themeDir, "post.hbs"))
			if err != nil {
				return err
			}

			// 调用markdown－>html方法, 得到文章信息、文章内容
			fileInfo, htmlContent, err := parse.DoRender(mdContent, postsTpl)
			if err != nil {
				return err
			}

			year := strconv.Itoa(time.Now().Year())
			month := strconv.Itoa(int(time.Now().Month()))
			day := strconv.Itoa(time.Now().Day())
			title := "untitled"+ strconv.Itoa(rand.Int())

			// 根据文章信息创建文件夹
			for k, v := range fileInfo {

				// 确定日期文件夹目录
				if k == "date" {
					ds := strings.Split(v,"/")

					if len(ds) == 3 {
						year = ds[0]
						month = ds[1]
						day = ds[2]
					}
				}

				// 确定文章文件夹目录
				if k == "title" {
					title = v
				}
			}

			// 检查年份文件夹是否存在
			if !utils.ExistDir(path.Join(publicDir, "posts", year)) {
				if !utils.Mkdir(path.Join(publicDir, "posts", year)) {
					return errors.New(fmt.Sprintf("create directory %s failed", path.Join(publicDir, "posts", year)))
				}
			}

			// 检查月份文件夹是否存在
			if !utils.ExistDir(path.Join(publicDir, "posts", year, month)) {
				if !utils.Mkdir(path.Join(publicDir, "posts", year, month)) {
					return errors.New(fmt.Sprintf("create directory %s failed", path.Join(publicDir, "posts", year, month)))
				}
			}

			// 检查日期文件夹是否存在
			if !utils.ExistDir(path.Join(publicDir, "posts", year, month, day)) {
				if !utils.Mkdir(path.Join(publicDir, "posts", year, month, day)) {
					return errors.New(fmt.Sprintf("create directory %s failed", path.Join(publicDir, "posts", year, month, day)))
				}
			}

			newTitle := strings.Replace(title, " ", "_", -1)
			if !utils.Mkdir(path.Join(publicDir, "posts", year, month, day, newTitle)) {
				return errors.New(fmt.Sprintf("create directory %s failed",
					path.Join(publicDir, "posts", year, month, day, newTitle)))
			}

			// 根据文章内容创建html文件
			if !utils.Mkfile(path.Join(publicDir,"posts", year, month, day, newTitle, "index.html"), htmlContent) {
				return errors.New(fmt.Sprintf("create file %s failed",
					path.Join(publicDir,"posts", year, month, day, newTitle, "index.html")))
			}

		}
	}
	return nil
}
