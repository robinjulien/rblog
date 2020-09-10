package server

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/golang-commonmark/markdown"
)

const (
	templatesFolder = "assets/"
)

// HeaderData contains all the informations that goes into the header section of the template
type HeaderData struct {
	BlogTitle string
	PageTitle string
}

// FooterData contains all the informations that goes into the footer section of the template
type FooterData struct {
}

// MenuData contains all the informations that goes into the menu section of the template
type MenuData struct {
}

// PageContent is all content a page need to be displayed
type PageContent struct {
	HeaderData
	FooterData
	MenuData
	Content template.HTML
}

// PostView is the summary of a post to display on the list posts page
type PostView struct {
	Title string
	Desc  string
	Date  time.Time
}

// PageListPosts is the content of the list posts page
type PageListPosts struct {
	PageContent
	PageNumber int
	MaxPage    int
	NumPosts   int
	Posts      []PostView
}

// function that respond when the home page is requested
func serveHomePage(w http.ResponseWriter, r *http.Request) {
	pageMD, err := ioutil.ReadFile("./pages/home.md")
	check(err)

	md := markdown.New()

	data := PageContent{
		HeaderData: HeaderData{
			BlogTitle: cfg.BlogName,
			PageTitle: cfg.HomePageTitle + cfg.PageTitleSuffix,
		},
		FooterData: FooterData{},
		MenuData:   MenuData{},
		Content:    template.HTML(md.RenderToString(pageMD)),
	}
	tmpl, err := template.ParseFiles(templatesFolder+"home.tpl", templatesFolder+"header.tpl", templatesFolder+"footer.tpl", templatesFolder+"menu.tpl")
	check(err)
	tmpl.ExecuteTemplate(w, "home.tpl", data)
	check(err)
}

// function that respond when a specific page is requested at /pages/*
func servePage(w http.ResponseWriter, r *http.Request) {
	fileName := strings.Replace(r.URL.Path, "/pages/", "", 1)

	pageMD, err := ioutil.ReadFile("./pages/" + fileName + ".md")

	if os.IsNotExist(err) { // If the specified page not found then 404 error
		serve404(w, r)
		return
	}
	check(err) // In case of another error

	md := markdown.New()

	data := PageContent{
		HeaderData: HeaderData{
			BlogTitle: cfg.BlogName,
			PageTitle: FormatFileNameToTitle(fileName) + cfg.PageTitleSuffix, // TODO
		},
		FooterData: FooterData{},
		MenuData:   MenuData{},
		Content:    template.HTML(md.RenderToString(pageMD)),
	}
	tmpl, err := template.ParseFiles(templatesFolder+"view_page.tpl", templatesFolder+"header.tpl", templatesFolder+"footer.tpl", templatesFolder+"menu.tpl")
	check(err)
	tmpl.ExecuteTemplate(w, "view_page.tpl", data)
	check(err)
}

func serveAssets(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, RewritePath(r.URL.Path, "/assets/", "./assets/public/"))
}

func serveListPosts(w http.ResponseWriter, r *http.Request) {
	files, err := ReadDirModTime("./posts") // custom ReadDir that sort by modtime
	check(err)

	data := PageListPosts{}
	data.HeaderData = HeaderData{
		BlogTitle: cfg.BlogName,
		PageTitle: cfg.ListPostsPageTitle + cfg.PageTitleSuffix, // TODO
	}
	data.FooterData = FooterData{}
	data.MenuData = MenuData{}
	data.NumPosts = len(files)
	data.MaxPage = data.NumPosts/cfg.MaxPostsOnListPage + 1

	queryPage := r.URL.Query().Get("p")
	var queryPageNumber int

	if queryPage == "" {
		queryPageNumber = 1
	} else {
		queryPageNumber, _ = strconv.Atoi(queryPage)
		// No need to check err, if the query param cannot be parsed, then queryPageNumber will be zero which is an invalid value
	}
	data.PageNumber = queryPageNumber

	if queryPageNumber <= 0 || queryPageNumber > data.MaxPage {
		data.Content = template.HTML("erreur, page non valide") // TODO
	} else {
		data.Posts = make([]PostView, 0)
		// No check for md file, we consider here that only md files are located into the posts folder
		// for easier handle of pagination // TODO change that to check if .md

		for i := (queryPageNumber - 1) * 10; (i < queryPageNumber*10) && i < len(files); i++ {
			data.Posts = append(data.Posts, PostView{
				Title: FormatFileNameToTitle(files[i].Name()),
				Desc:  "", // TODO
				Date:  files[i].ModTime(),
			})
		}
	}

	_tmpl := template.New("posts") // Create new empty template here because .Funcs must be called before parsing files
	_tmpl.Funcs(template.FuncMap{  // Callable functions from template to format the date the way you want
		"formatDate": func(date time.Time) string {
			return date.Format("02/01/2006 - 15:04")
		},
	})

	tmpl, err := _tmpl.ParseFiles(templatesFolder+"posts.tpl", templatesFolder+"header.tpl", templatesFolder+"footer.tpl", templatesFolder+"menu.tpl")
	check(err)

	err = tmpl.ExecuteTemplate(w, "posts.tpl", data)
	check(err)
}

func servePost(w http.ResponseWriter, r *http.Request) {
	fileName := strings.Replace(r.URL.Path, "/posts/", "", 1)

	pageMD, err := ioutil.ReadFile("./posts/" + fileName + ".md")

	if os.IsNotExist(err) { // If the specified page not found then 404 error
		serve404(w, r)
		return
	}
	check(err) // In case of another error

	md := markdown.New()

	data := PageContent{
		HeaderData: HeaderData{
			BlogTitle: cfg.BlogName,
			PageTitle: FormatFileNameToTitle(fileName) + cfg.PageTitleSuffix, //TODO
		},
		FooterData: FooterData{},
		MenuData:   MenuData{},
		Content:    template.HTML(md.RenderToString(pageMD)),
	}

	tmpl, err := template.ParseFiles(templatesFolder+"view_post.tpl", templatesFolder+"header.tpl", templatesFolder+"footer.tpl", templatesFolder+"menu.tpl")
	check(err)
	err = tmpl.ExecuteTemplate(w, "view_post.tpl", data)
	check(err)
}
