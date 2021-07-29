package repo

import (
	"com.aharakitchen/app/database"
	"com.aharakitchen/app/domain"
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"sync"
)

type PostRepoImpl struct {
	postPreviews   []domain.PostPreviewDto
	postList   domain.PostList
	postDto    domain.PostDto
}

func (p PostRepoImpl) FindAllPosts(page string, newPosts bool) (*domain.PostList, error) {
	conn := database.MongoConn

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}
	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	if newPosts {
		findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	query := bson.M{}

	var wg sync.WaitGroup
	wg.Add(2)

	go func(conn *database.Connection) {
		defer wg.Done()
		cur, err := conn.PostCollection.Find(context.TODO(), query, &findOptions)

		if err != nil {
			panic(err)
		}

		if err = cur.All(context.TODO(), &p.postPreviews); err != nil {
			log.Fatal(err)
		}
		return
	}(conn)

	go func(conn *database.Connection) {
		defer wg.Done()
		count, err:= conn.PostCollection.CountDocuments(context.TODO(),query)

		if err != nil {
			panic(err)
		}

		p.postList.NumberOfPosts = count

		if p.postList.NumberOfPosts < 10 {
			p.postList.NumberOfPages = 1
		} else {
			p.postList.NumberOfPages = int(count / 10) + 1
		}
	}(conn)

	wg.Wait()

	p.postList.Posts = p.postPreviews
	p.postList.CurrentPage = 1

	return &p.postList, nil
}

func (p PostRepoImpl) Create(post domain.Post) error {
	conn := database.MongoConn

	_, err := conn.PostCollection.InsertOne(context.TODO(), &post)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (p PostRepoImpl) UpdateByTitle(post domain.Post) error {
	conn := database.MongoConn

	rdb := database.Conn.Get()

	updatedPost := new(domain.Post)
	err := conn.PostCollection.FindOneAndUpdate(context.TODO(), bson.D{{"_id", post.Id}},
		bson.M{"$set": bson.M{
			"title":       post.Title,
			"content":     post.Content,
			"mainImage":   post.MainImage,
			"storyImages": post.StoryImages,
			"tag":         post.Tag,
			"visible":     post.Visible,
			"updated":      post.Updated,
			"updatedAt":   post.UpdatedAt,
		},
		}).Decode(&updatedPost)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	_, err = rdb.Do("HDEL", redis.Args{}.Add(post.Id.String() + "getbyID").AddFlat(new(domain.RedisPostDto))...)

	if err != nil {
		return err
	}

	fmt.Println("Removed from cache, update current tag")

	return nil
}

func (p PostRepoImpl) DeleteById(post domain.Post) error {
	conn := database.MongoConn

	rdb := database.Conn.Get()

	_, err := conn.PostCollection.DeleteOne(context.TODO(), bson.D{{"_id", post.Id}})

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	_, err = rdb.Do("HDEL", redis.Args{}.Add(post.Id.String() + "getbyID").AddFlat(new(domain.RedisPostDto))...)

	if err != nil {
		return err
	}

	fmt.Println("Removed from cache, delete current post")

	return nil
}

func (p PostRepoImpl) FindPostById(id primitive.ObjectID) (*domain.PostDto, error) {
	conn := database.MongoConn

	query := bson.D{{"_id", id}}

	err := conn.PostCollection.FindOne(context.TODO(), query).Decode(&p.postDto)

	if err != nil {
		return nil, err
	}

	rdb := database.Conn.Get()

	rp := new(domain.RedisPostDto)

	b, err := msgpack.Marshal(&p.postDto.StoryImages)
	crb, err := msgpack.Marshal(&p.postDto.CreatedAt)
	urb, err := msgpack.Marshal(&p.postDto.UpdatedAt)

	rp.Tag = p.postDto.Tag
	rp.Title = p.postDto.Title
	rp.Content = p.postDto.Content
	rp.Author = p.postDto.Author
	rp.MainImage = p.postDto.MainImage
	rp.UpdatedAt = urb
	rp.CreatedAt = crb
	rp.Updated = p.postDto.Updated
	rp.StoryImages = b

	_, err = rdb.Do("HMSET", redis.Args{}.Add(id.String() + "getbyID").AddFlat(rp)...)

	if err != nil {
		fmt.Println("Found in cache in find by ID...")
		return nil, err
	}
	fmt.Println("Cached in find by ID...")

	return &p.postDto, nil
}

func (p PostRepoImpl) FeaturedPosts() (*domain.PostList, error) {
	conn := database.MongoConn

	findOptions := options.FindOptions{}

	findOptions.SetLimit(10)
	findOptions.SetSort(bson.D{{"score", -1}})

	cur, err := conn.PostCollection.Find(context.TODO(), bson.M{}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &p.postPreviews); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	p.postList.Posts = p.postPreviews
	p.postList.CurrentPage = 1
	p.postList.NumberOfPages = 1
	p.postList.NumberOfPosts = int64(len(p.postPreviews))

	rdb := database.Conn.Get()

	rp := new(domain.RedisPostList)

	b, err := msgpack.Marshal(&p.postList.Posts)

	rp.NumberOfPosts = p.postList.NumberOfPosts
	rp.CurrentPage = p.postList.CurrentPage
	rp.NumberOfPages = p.postList.NumberOfPages
	rp.Posts = b

	_, err = rdb.Do("HMSET", redis.Args{}.Add("featuredstories").AddFlat(rp)...)

	if err != nil {
		fmt.Println("Found in cache in find by ID...")
		return nil, err
	}

	fmt.Println("Cached in find by ID...")
	return &p.postList, nil
}

func NewPostRepoImpl() PostRepoImpl {
	var postRepoImpl PostRepoImpl

	return postRepoImpl
}
