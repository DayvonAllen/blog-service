package repo

import (
	"com.aharakitchen/app/domain"
	"fmt"
)

func ProcessMessage(message domain.Message) error {
	if message.ResourceType == "post" {
		// 201 is the created messageType
		if message.MessageType == 201 {
			posts := message.Post
			err := PostRepoImpl{}.Create(posts)

			if err != nil {
				return err
			}
			return nil
		}

		//200 is the updated messageType
		if message.MessageType == 200 {
			posts := message.Post
			err := PostRepoImpl{}.UpdateByTitle(posts)

			if err != nil {
				return err
			}

			return nil
		}

		// 204 is the deleted messageType
		//if message.MessageType == 204 {
		//	user := message.User
		//
		//	err := UserRepoImpl{}.DeleteByID(&user)
		//
		//	if err != nil {
		//		return err
		//	}
		//	return nil
		//}
	}

	if message.ResourceType == "tag" {
		// 201 is the created messageType
		if message.MessageType == 201 {
			tags := message.Tag
			err := TagRepoImpl{}.Create(tags)

			if err != nil {
				return err
			}
			return nil
		}

		// 200 is the updated messageType
		if message.MessageType == 200 {
			tags := message.Tag
			err := TagRepoImpl{}.UpdateTag(tags)

			if err != nil {
				return err
			}
			return nil
		}

		// 204 is the deleted messageType
		//if message.MessageType == 204 {
		//	user := message.User
		//
		//	err := UserRepoImpl{}.DeleteByID(&user)
		//
		//	if err != nil {
		//		return err
		//	}
		//	return nil
		//}
	}

	return fmt.Errorf("cannot process this message")
}

//func PushUserToQueue(message []byte, topic string) error {
//	producer := events.GetInstance()
//
//	msg := &sarama.ProducerMessage{
//		Topic: topic,
//		Value: sarama.StringEncoder(message),
//	}
//
//	partition, offset, err := producer.SendMessage(msg)
//	if err != nil {
//		fmt.Println(fmt.Errorf("%v", err))
//		err = producer.Close()
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println("Failed to send message to the queue")
//	}
//
//	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", "user", partition, offset)
//	return nil
//}
//
//func SendKafkaMessage(story *domain.Story, eventType int) error {
//	um := new(domain.Message)
//	um.Story = *story
//
//	// user created/updated event
//	um.MessageType = eventType
//	um.ResourceType = "story"
//
//	//turn user struct into a byte array
//	b, err := msgpack.Marshal(um)
//
//	if err != nil {
//		return err
//	}
//
//	err = PushUserToQueue(b, config.Config("PRODUCER_TOPIC"))
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func SendEventMessage(event *domain.Event, eventType int) error {
//	um := new(domain.Message)
//	um.Event = *event
//
//	// user created/updated event
//	um.MessageType = eventType
//	um.ResourceType = "event"
//
//	fmt.Println(um.Event)
//	//turn user struct into a byte array
//	b, err := msgpack.Marshal(um)
//
//	if err != nil {
//		return err
//	}
//
//	err = PushUserToQueue(b, "event")
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

