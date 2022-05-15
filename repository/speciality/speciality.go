package specialityrepo

import "go.mongodb.org/mongo-driver/mongo"

type SpecialityRepository struct {
	colln *mongo.Collection
}

func (c *SpecialityRepository) Init(colln *mongo.Collection) {
	c.colln = colln
}
