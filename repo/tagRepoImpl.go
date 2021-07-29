package repo

import (
	"com.aharakitchen/app/database"
	"com.aharakitchen/app/domain"
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"sync"
)

type TagRepoImpl struct {
	postPreviews []domain.PostPreviewDto
	postList     domain.PostList
	tags         []domain.TagDto
	tag          domain.Tag
	tagList      domain.TagList
}

func (t TagRepoImpl) FindAllTags() (*domain.TagList, error) {
	conn, _ := database.ConnectToDB()

	defer func(conn *database.Connection, ctx context.Context) {
		err := conn.Disconnect(ctx)
		if err != nil {

		}
	}(conn, context.TODO())
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

	rdb := database.Conn.Get()

	rt := new(domain.RedisTagList)

	b, err := msgpack.Marshal(&t.tagList.Tags)


	rt.NumberOfCategories = t.tagList.NumberOfCategories
	rt.Tags = b

	_, err = rdb.Do("HMSET", redis.Args{}.Add("tags").AddFlat(rt)...)

	if err != nil {
		return nil, err
	}

	fmt.Println("set tags in the cache, find all tags")

	return &t.tagList, nil
}

func (t TagRepoImpl) FindAllPostsByCategory(category, page string) (*domain.PostList, error) {
	conn, _ := database.ConnectToDB()

	defer func(conn *database.Connection, ctx context.Context) {
		err := conn.Disconnect(ctx)
		if err != nil {

		}
	}(conn, context.TODO())

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
		count, err := conn.PostCollection.CountDocuments(context.TODO(), query)

		if err != nil {
			panic(err)
		}

		t.postList.NumberOfPosts = count

		if t.postList.NumberOfPosts < 10 {
			t.postList.NumberOfPages = 1
		} else {
			t.postList.NumberOfPages = int(count/10) + 1
		}
	}()

	wg.Wait()

	t.postList.Posts = t.postPreviews
	t.postList.CurrentPage = 1

	return &t.postList, nil
}

func (t TagRepoImpl) Create(tag domain.Tag) error {
	conn, _ := database.ConnectToDB()

	defer func(conn *database.Connection, ctx context.Context) {
		err := conn.Disconnect(ctx)
		if err != nil {

		}
	}(conn, context.TODO())

	err := conn.TagCollection.FindOneAndReplace(context.TODO(), bson.D{{"value", tag.Value}}, &tag).Decode(&t.tag)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err := conn.TagCollection.InsertOne(context.TODO(), &tag)

			if err != nil {
				return fmt.Errorf("error processing data")
			}

		}
		return err
	}

	rdb := database.Conn.Get()

	_, err = rdb.Do("HDEL", redis.Args{}.Add("tags").AddFlat(new(domain.RedisTagList))...)

	if err != nil {
		return err
	}

	fmt.Println("updated cache, create current tag")

	return nil
}

func (t TagRepoImpl) UpdateTag(tag domain.Tag) error {
	conn, _ := database.ConnectToDB()

	defer func(conn *database.Connection, ctx context.Context) {
		err := conn.Disconnect(ctx)
		if err != nil {

		}
	}(conn, context.TODO())

	err := conn.TagCollection.FindOneAndReplace(context.TODO(), bson.D{{"value", tag.Value}}, &tag).Decode(&t.tag)

	if err != nil {
		return err
	}

	rdb := database.Conn.Get()

	_, err = rdb.Do("HDEL", redis.Args{}.Add("tags").AddFlat(new(domain.RedisTagList))...)

	if err != nil {
		return err
	}

	return nil
}

func NewTagRepoImpl() TagRepoImpl {
	var tagRepoImpl TagRepoImpl

	return tagRepoImpl
}
