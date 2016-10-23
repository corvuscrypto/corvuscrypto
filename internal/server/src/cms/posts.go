package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Post represents a blog post
type Post struct {
	Date    time.Time `bson:"date"`
	Title   string    `bson:"title"`
	Summary string    `bson:"summary"`
	Body    string    `bson:"body"`
	Tags    []string  `bson:"tags"`
	URL     string    `bson:"url"`
	Publish bool      `bson:"publish"`
}

//PostsDB is the main post db
var PostsDB *mgo.Database

func initDBSession() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	PostsDB = session.DB("blog")
	postCollection := PostsDB.C("posts")

	//do a few asserts just to make sure the indexes are there
	postCollection.EnsureIndex(mgo.Index{
		Key:    []string{"url"},
		Unique: true,
	})
	postCollection.EnsureIndex(mgo.Index{
		Key: []string{"date"},
	})
	postCollection.EnsureIndex(mgo.Index{
		Key: []string{"tags"},
	})
}

//UpdatePost performs a whole document update
func UpdatePost(url string, p *Post) (err error) {
	err = PostsDB.C("posts").Update(bson.M{"url": url}, p)
	return
}

//InsertNewPost inserts a new post
func InsertNewPost(p *Post) (err error) {

	err = PostsDB.C("posts").Insert(p)
	if err != nil {
		return
	}

	return
}

//GetPostByURL retrieves a post by url
func GetPostByURL(url string) (post *Post, err error) {
	post = &Post{}
	err = PostsDB.C("posts").Find(bson.M{"url": url}).One(post)
	return
}
