package repo

import (
	"com.aharakitchen/app/database"
	"com.aharakitchen/app/domain"
	"context"
	"fmt"
	cache2 "github.com/go-redis/cache/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"sync"
	"time"
)

type PostRepoImpl struct {
	postPreviews   []domain.PostPreviewDto
	postList   domain.PostList
	postDto    domain.PostDto
}

func (p PostRepoImpl) FindAllPosts(page string, newPosts bool) (*domain.PostList, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

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

	go func() {
		defer wg.Done()
		cur, err := conn.PostCollection.Find(context.TODO(), query, &findOptions)

		if err != nil {
			panic(err)
		}

		if err = cur.All(context.TODO(), &p.postPreviews); err != nil {
			log.Fatal(err)
		}

	}()

	go func() {
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
	}()

	wg.Wait()

	p.postList.Posts = p.postPreviews
	p.postList.CurrentPage = 1

	return &p.postList, nil
}

func (p PostRepoImpl) Create(post domain.Post) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	_, err := conn.PostCollection.InsertOne(context.TODO(), &post)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (p PostRepoImpl) UpdateByTitle(post domain.Post) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

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

	return nil
}

func (p PostRepoImpl) FindPostById(id primitive.ObjectID) (*domain.PostDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	query := bson.D{{"_id", id}}

	err := conn.PostCollection.FindOne(context.TODO(), query).Decode(&p.postDto)

	if err != nil {
		return nil, err
	}

	return &p.postDto, nil
}

func (p PostRepoImpl) FeaturedPosts(rdb *cache2.Cache) (*domain.PostList, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

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

	go func() {
		fmt.Println("set cache")
		err = rdb.Set(&cache2.Item{
			Ctx:   context.TODO(),
			Key:   "featuredstories",
			Value: p.postList,
			TTL:   time.Hour,
		})

		if err != nil {
			fmt.Println("Found in cache in find by ID...")
			panic(err)
		}
		fmt.Println("Cached in find by ID...")
		return
	}()

	return &p.postList, nil
}

func NewPostRepoImpl() PostRepoImpl {
	var postRepoImpl PostRepoImpl

	return postRepoImpl
}
