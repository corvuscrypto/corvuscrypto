package main

import (
	"errors"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Only error we really need
var ErrPostNotFound = errors.New("Post specified was not found!")

//Post represents a blog post
type Post struct {
	Number  int       `bson:"number"`
	Date    time.Time `bson:"date"`
	Title   string    `bson:"title"`
	Summary string    `bson:"summary"`
	Body    string    `bson:"body"`
	Tags    []string  `bson:"tags"`
	URL     string    `bson:"url"`
}

//PostsDB is the only db connection that will be used for this app
var PostsDB *mgo.Database

//initializeSession connects to mongo and sets the global DB variable to the blog_posts db.
func initializeDBSession() {
	//just use default. SOOOOO ORIGINAL
	dbSession, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	PostsDB = dbSession.DB("blog")
}

func getLatestPost() (*Post, error) {
	var err error
	latestPost := new(Post)
	err = PostsDB.C("posts").Find(nil).Sort("-_id").One(latestPost)
	return latestPost, err
}

func getPostByURL(url string) (*Post, error) {
	var err error
	post := new(Post)
	err = PostsDB.C("posts").Find(bson.M{
		"url": url,
	}).One(post)
	if post == nil {
		err = ErrPostNotFound
	}
	return post, err
}

func searchPosts(searchTerms []string, prevLast int) ([]*Post, error) {
	var err error
	var posts []*Post

	iter := PostsDB.C("posts").Find(bson.M{
		"searchTags": bson.M{
			"$all": searchTerms,
		},
		"number": bson.M{
			"$gt": prevLast,
		},
	}).Sort("-_id").Limit(10).Iter()

	defer iter.Close()

	tempPost := new(Post)

	for iter.Next(tempPost) {
		var newPost = new(Post)
		*newPost = *tempPost
		posts = append(posts, newPost)
		tempPost = nil
	}
	return posts, err
}

func getPostsByTags(tags []string, prevLast int) ([]*Post, error) {
	var err error
	var posts []*Post

	iter := PostsDB.C("posts").Find(bson.M{
		"tags": bson.M{
			"$all": tags,
		},
		"number": bson.M{
			"$gt": prevLast,
		},
	}).Sort("-_id").Limit(10).Iter()

	defer iter.Close()

	tempPost := new(Post)

	for iter.Next(tempPost) {
		var newPost *Post
		*newPost = *tempPost
		posts = append(posts, newPost)
		tempPost = nil
	}
	return posts, err
}

func getAllPosts(prevLast int) ([]*Post, error) {
	var err error
	var posts []*Post

	iter := PostsDB.C("posts").Find(bson.M{
		"number": bson.M{
			"$gt": prevLast,
		},
	}).Sort("-_id").Limit(10).Iter()

	defer iter.Close()

	tempPost := new(Post)

	for iter.Next(tempPost) {
		var newPost = new(Post)
		*newPost = *tempPost
		posts = append(posts, newPost)
		tempPost = nil
	}
	return posts, err
}
