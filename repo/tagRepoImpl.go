package repo

import (
	"com.aharakitchen/app/cache"
	"com.aharakitchen/app/database"
	"com.aharakitchen/app/domain"
	"context"
	"fmt"
	cache2 "github.com/go-redis/cache/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"sync"
	"time"
)

type TagRepoImpl struct {
	postPreviews   []domain.PostPreviewDto
	postList   domain.PostList
	tags []domain.TagDto
	tag domain.Tag
	tagList domain.TagList
}

func (t TagRepoImpl) FindAllTags(rdb *cache2.Cache) (*domain.TagList, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	cur, err := conn.TagCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &t.tags); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	t.tagList.Tags = t.tags
	t.tagList.NumberOfCategories = len(t.tags)

	go func() {
		fmt.Println("set cache")
		err = rdb.Set(&cache2.Item{
			Ctx:   context.TODO(),
			Key:   "tagList",
			Value: &t.tagList,
			TTL:   time.Hour,
		})

		if err != nil {
			fmt.Println("Found in cache in find by ID...")
			panic(err)
		}
		fmt.Println("Cached in find by ID...")
		return
	}()

	return &t.tagList, nil
}

func (t TagRepoImpl) FindAllPostsByCategory(category, page string) (*domain.PostList, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	err := conn.TagCollection.FindOne(context.TODO(), bson.D{{"value", category}}).Decode(&t.tag)

	if err != nil {
		return nil, err
	}

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)
	findOptions.SetSort(bson.D{{"score", -1}})

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}
	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	query := bson.M{"_id": bson.M{"$in": t.tag.AssociatedPosts}}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		cur, err := conn.PostCollection.Find(context.TODO(), query, &findOptions)

		if err != nil {
			panic(err)
		}

		if err = cur.All(context.TODO(), &t.postPreviews); err != nil {
			log.Fatal(err)
		}

	}()

	go func() {
		defer wg.Done()
		count, err:= conn.PostCollection.CountDocuments(context.TODO(),query)

		if err != nil {
			panic(err)
		}

		t.postList.NumberOfPosts = count

		if t.postList.NumberOfPosts < 10 {
			t.postList.NumberOfPages = 1
		} else {
			t.postList.NumberOfPages = int(count / 10) + 1
		}
	}()

	wg.Wait()

	t.postList.Posts = t.postPreviews
	t.postList.CurrentPage = 1

	return &t.postList, nil
}

func (t TagRepoImpl) Create(tag domain.Tag) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)
	rdb := cache.RedisCachePool.Get().(*cache2.Cache)
	defer cache.RedisCachePool.Put(rdb)

	_, err := conn.TagCollection.InsertOne(context.TODO(), &tag)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	go func() {
		err := rdb.Delete(context.TODO(), "tags")

		if err != nil {
			panic(err)
		}

		fmt.Println("Removed from cache, update current tag")

		return
	}()

	return nil
}

func (t TagRepoImpl) UpdateTag(tag domain.Tag) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)
	rdb := cache.RedisCachePool.Get().(*cache2.Cache)
	defer cache.RedisCachePool.Put(rdb)

	_, err := conn.TagCollection.UpdateOne(context.TODO(), bson.M{"_id": tag.Id}, bson.M{"$set": bson.M{"associatedPosts": tag.AssociatedPosts}})

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	go func() {
		err := rdb.Delete(context.TODO(), "tags")

		if err != nil {
			panic(err)
		}

		fmt.Println("Removed from cache, update current tag")

		return
	}()

	return nil
}

func NewTagRepoImpl() TagRepoImpl {
	var tagRepoImpl TagRepoImpl

	return tagRepoImpl
}
