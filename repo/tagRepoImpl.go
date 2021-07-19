package repo

import (
	"com.aharakitchen/app/database"
	"com.aharakitchen/app/domain"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"sync"
)

type TagRepoImpl struct {
	postPreviews   []domain.PostPreviewDto
	postList   domain.PostList
	tags []domain.TagDto
	tag domain.Tag
}

func (t TagRepoImpl) FindAllTags() (*[]domain.TagDto, error) {
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

	return &t.tags, nil
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

	_, err := conn.TagCollection.InsertOne(context.TODO(), &tag)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (t TagRepoImpl) UpdateTag(tag domain.Tag) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	_, err := conn.TagCollection.UpdateOne(context.TODO(), bson.M{"_id": tag.Id}, bson.M{"$set": bson.M{"associatedPosts": tag.AssociatedPosts}})

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func NewTagRepoImpl() TagRepoImpl {
	var tagRepoImpl TagRepoImpl

	return tagRepoImpl
}
