package main

import (
  "os"
  "fmt"
  "strings"
  "strconv"
  "html/template"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "github.com/PuerkitoBio/goquery"
)

type Tv_name struct {
  Num int
  Name string
}

type Tv_magnet struct {
  Num int
  Magnet string
}

func store_data() {
  session, _ := mgo.Dial("127.0.0.1")
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)

  n := session.DB("korea").C("tv_name")
  m := session.DB("korea").C("tv_magnet")

  doc, _ := goquery.NewDocument("http://www.torrentbest.net/bbs/board.php?bo_table=torrent_kortv_ent")
  doc.Find("td.subject").Each(func(i int, s *goquery.Selection) {
    subject := s.Find("a").Text()
    val, _ := s.Find("a").Attr("href")
    n.Insert(&Tv_name{ i, subject})

    str := "http://www.torrentbest.net"
    substr := string([]byte(val[2:]))
    url := str + substr

    doc2, _ := goquery.NewDocument(url)
    doc2.Find("td.view_file").Each(func(j int, s2 *goquery.Selection) {
      magnet, _ := s2.Find("a").Attr("href")
      if strings.Contains(magnet, "magnet") {
        m.Insert(&Tv_magnet{i, magnet})
      }
    })
  })
}

func load_data(r render.Render) {
  data := map[string]interface{}{}
  
  session, _ := mgo.Dial("127.0.0.1")
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)

  n := session.DB("korea").C("tv_name")

  result_sub := Tv_name{}

  for i := 0 ; i < 30 ; i++ {
    n.Find(bson.M{"num" : i}).One(&result_sub)

    value := strconv.Itoa(i)
    safe_url := template.URL("../data/id=" + value)
    data["name" + value] = result_sub.Name
    data["url" + value] = safe_url
  }
  
  r.HTML(200, "test", data)
//  n.RemoveAll(bson.M{})
//  m.RemoveAll(bson.M{})
}

func mag_data (r2 render.Render, p martini.Params) {
  data2 := map[string]interface{}{}

  session, _:= mgo.Dial("127.0.0.1")
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)

  n2 := session.DB("korea").C("tv_name")
  m2 := session.DB("korea").C("tv_magnet")

  result_sub := Tv_name{}
  result_mag := Tv_magnet{}
  id, _ := strconv.ParseInt(p["id_num"], 0, 32)
  n2.Find(bson.M{"num": id}).One(&result_sub)
  m2.Find(bson.M{"num": id}).One(&result_mag)
  
  safe_mag := template.URL(result_mag.Magnet)
  data2["name"] = result_sub.Name
  data2["url"] = safe_mag

  r2.HTML(400, "second", data2)
}

func main() {  
  root := os.Getenv("GOPATH") + "/src/test/torrent_site"

  mar := martini.Classic()
  mar.Use(martini.Static(root))
      mar.Use(render.Renderer(render.Options{
        Directory: root,
        Layout: "layout",
        Extensions: []string{".html"},
        Charset: "UTF-8",
        IndentJSON: true,
      }))
  
  store_data()

  mar.Get("/", load_data)
  mar.Get("/data/id=:id_num", mag_data)
  mar.RunOnAddr(":8888")
}
