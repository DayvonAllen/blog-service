package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	Action string `bson:"action" json:"action"`
	Target string `bson:"target" json:"target"`
	ResourceId primitive.ObjectID `bson:"resourceId" json:"resourceId"`
	ActorIp string `bson:"actorUsername" json:"actorUsername"`
	Message string `bson:"message" json:"message"`
}

