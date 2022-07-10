package doctorrepo

import (
	"context"
	"errors"
	"time"

	"github.com/praveennagaraj97/online-consultation/api"
	doctordto "github.com/praveennagaraj97/online-consultation/dto/doctor"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	doctormodel "github.com/praveennagaraj97/online-consultation/models/doctor"
	"github.com/praveennagaraj97/online-consultation/pkg/env"
	"github.com/praveennagaraj97/online-consultation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DoctorRepository struct {
	colln         *mongo.Collection
	imageBasePath string
}

func (r *DoctorRepository) Initialize(colln *mongo.Collection) {
	r.colln = colln

	r.imageBasePath = env.GetEnvVariable("S3_ACCESS_BASEURL")

	utils.CreateIndex(colln, bson.D{
		{Key: "phone.number", Value: 1},
		{Key: "phone.code", Value: 1}}, "PhoneIndex", true)

	utils.CreateIndex(colln, bson.D{{Key: "email", Value: 1}}, "EmailIndex", true)

	utils.CreateIndex(colln, bson.D{{Key: "consultation_type_id", Value: 1}}, "ConsultationTypeIndex", false)
	utils.CreateIndex(colln, bson.D{{Key: "speciality_id", Value: 1}}, "SpecialityIndex", false)
	utils.CreateIndex(colln, bson.D{{Key: "hospital_id", Value: 1}}, "HospitalIndex", false)
	utils.CreateIndex(colln, bson.D{{Key: "experience", Value: 1}}, "ExperienceIndex", false)
	utils.CreateIndex(colln, bson.D{{Key: "is_active", Value: 1}}, "AccountActiveStatusIndex", false)
	utils.CreateIndex(colln, bson.D{{Key: "languages_ids", Value: 1}}, "SpokenLanguagesIndex", false)
}

func (r *DoctorRepository) CreateOne(doc *doctormodel.DoctorEntity) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil

}

func (r *DoctorRepository) CheckIfDoctorExistsByEmailOrPhone(email string, phone interfaces.PhoneType) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	emailFilter := bson.M{"email": email}
	phoneFilter := bson.M{"$and": bson.A{bson.M{"phone.code": phone.Code}, bson.M{"phone.number": phone.Number}}}

	filter := bson.M{"$or": bson.A{
		emailFilter,
		phoneFilter,
	}}

	count, err := r.colln.CountDocuments(ctx, filter)
	if err != nil {
		return false
	}

	return count > 0
}

func (r *DoctorRepository) FindOne(id *primitive.ObjectID,
	email string,
	phone *interfaces.PhoneType,
	showInActive bool) (*doctormodel.DoctorEntity, error) {

	var filterPipe bson.D = make(bson.D, 0)

	if id != nil {
		if showInActive {
			filterPipe = bson.D{{Key: "$match", Value: bson.M{"$and": bson.A{bson.M{"_id": id}, bson.M{"is_active": showInActive}}}}}
		} else {
			filterPipe = bson.D{{Key: "$match", Value: bson.M{"_id": id}}}
		}
	}

	if email != "" {
		if showInActive {
			filterPipe = bson.D{{Key: "$match", Value: bson.M{"$and": bson.A{bson.M{"email": email}, bson.M{"is_active": showInActive}}}}}
		} else {
			filterPipe = bson.D{{Key: "$match", Value: bson.M{"email": email}}}
		}
	}

	if phone != nil {
		if showInActive {
			filterPipe = bson.D{{Key: "$match", Value: bson.M{"$and": bson.A{bson.M{"phone.code": phone.Code}, bson.M{"phone.number": phone.Number}, bson.M{"is_active": showInActive}}}}}
		} else {
			filterPipe = bson.D{{Key: "$match", Value: bson.M{"$and": bson.A{bson.M{"phone.code": phone.Code}, bson.M{"phone.number": phone.Number}}}}}
		}
	}

	// Consultation ID Populate
	typeMatchPipe := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "consultation",
		"localField":   "consultation_type_id",
		"foreignField": "_id",
		"as":           "consultation_type",
	}}}
	unwindTypePipe := bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$consultation_type",
		"preserveNullAndEmptyArrays": true,
	}}}
	setTypePipe := bson.D{{Key: "$set", Value: bson.M{"consultation_type": "$consultation_type.type"}}}

	pipeLine := mongo.Pipeline{
		filterPipe,
		typeMatchPipe,
		unwindTypePipe,
		setTypePipe,
	}

	// Add Prefix to image
	setImagePrefixPipe := bson.D{{Key: "$set",
		Value: bson.M{"profile_pic.image_src": bson.M{"$cond": bson.D{
			{Key: "if", Value: bson.M{"$eq": bson.A{"$profile_pic", nil}}},
			{Key: "then", Value: nil},
			{Key: "else", Value: bson.M{"$concat": bson.A{r.imageBasePath, "/", "$profile_pic.original_image_path"}}},
		}}}}}
	setBlurImagePrefixPipe := bson.D{{Key: "$set",
		Value: bson.M{"profile_pic.blur_data_url": bson.M{"$cond": bson.D{
			{Key: "if", Value: bson.M{"$eq": bson.A{"$profile_pic", nil}}},
			{Key: "then", Value: nil},
			{Key: "else", Value: bson.M{"$concat": bson.A{r.imageBasePath, "/", "$profile_pic.blur_image_path"}}},
		}}}}}
	resetNullImagePipe := bson.D{{Key: "$set", Value: bson.M{
		"profile_pic": bson.M{"$cond": bson.D{
			{Key: "if", Value: bson.M{"$eq": bson.A{"$profile_pic.image_src", nil}}},
			{Key: "then", Value: nil},
			{Key: "else", Value: "$profile_pic"},
		}},
	}}}
	pipeLine = append(pipeLine, setImagePrefixPipe, setBlurImagePrefixPipe, resetNullImagePipe)

	// Populate hospital
	lookUpHospital := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "hospital",
		"localField":   "hospital_id",
		"foreignField": "_id",
		"as":           "hospital",
	}}}

	unwindHospiatl := bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$hospital",
		"preserveNullAndEmptyArrays": true,
	}}}

	pipeLine = append(pipeLine, lookUpHospital, unwindHospiatl)

	// Populate Speciality
	lookupSpeciality := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "speciality",
		"as":           "speciality",
		"localField":   "speciality_id",
		"foreignField": "_id",
	}}}

	unwindSpeciality := bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$speciality",
		"preserveNullAndEmptyArrays": true,
	}}}

	setSpecialityTitle := bson.D{{Key: "$set", Value: bson.M{
		"speciality": "$speciality.title",
	}}}
	pipeLine = append(pipeLine, lookupSpeciality, unwindSpeciality, setSpecialityTitle)

	// Populate Languages
	languagesLookUp := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "language",
		"localField":   "languages_ids",
		"foreignField": "_id",
		"as":           "spoken_languages",
	}}}

	pipeLine = append(pipeLine, languagesLookUp)

	// Next Available Slot
	nextAvailableLookUp := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "appointment_slot",
		"localField":   "_id",
		"foreignField": "doctor_id",
		"as":           "next_available_slot",
		"pipeline": bson.A{
			bson.M{"$match": bson.M{"$and": bson.A{
				bson.M{"is_available": true},
				bson.M{"start": bson.M{"$gt": primitive.NewDateTimeFromTime(time.Now())}},
			}}},
			bson.M{"$limit": 1},
		},
	}}}

	unwindNextAvailableSlot := bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$next_available_slot",
		"preserveNullAndEmptyArrays": true,
	}}}

	pipeLine = append(pipeLine, nextAvailableLookUp, unwindNextAvailableSlot)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cur, err := r.colln.Aggregate(ctx, pipeLine)
	if err != nil {
		return nil, err
	}

	var result []doctormodel.DoctorEntity

	defer cur.Close(context.TODO())

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	if len(result) == 1 {
		return &result[0], nil
	}

	return nil, errors.New("couldn't find any doctor")

}

func (r *DoctorRepository) UpdateDoctorStatus(id *primitive.ObjectID, state bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := r.colln.UpdateByID(ctx, id, bson.M{"$set": bson.M{"is_active": state}})

	return err

}

func (r *DoctorRepository) FindAll(pgOpts *api.PaginationOptions,
	fltrOpts *map[string]primitive.M,
	srtOpts *map[string]int8,
	keySortBy string,
	searchOpts *bson.M,
	showInActive bool,
	slotsExistsOn *primitive.DateTime,
) ([]doctormodel.DoctorEntity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	pipeline := mongo.Pipeline{}

	// Filter Doctor By Appointment Slot Availability
	if slotsExistsOn != nil {
		apptAvailLookup := bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "appointment_slot",
			"localField":   "_id",
			"foreignField": "doctor_id",
			"as":           "slot",
			"pipeline": bson.A{
				bson.M{"$match": bson.M{
					"$and": bson.A{
						bson.M{"is_available": true},
						bson.M{"date": slotsExistsOn},
					},
				}},
				bson.M{"$limit": 1},
			},
		}}}

		// Filter doctors having slots
		apptAvailMatch := bson.D{{Key: "$match", Value: bson.M{
			"$expr": bson.M{
				"$gt": bson.A{bson.M{"$size": "$slot"}, 0},
			},
		}}}

		// Map the available slot to next available field
		apptAvailAdddFields := bson.D{{Key: "$addFields", Value: bson.M{
			"next_available_slot": "$slot",
		}}}

		pipeline = append(pipeline, apptAvailLookup, apptAvailMatch, apptAvailAdddFields)
	}

	// Search Match
	if searchOpts != nil {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: *searchOpts}})
	}

	// Filter Options
	if len(*fltrOpts) != 0 {
		fltr := bson.D{{Key: "$match", Value: *fltrOpts}}
		pipeline = append(pipeline, fltr)
	}

	if !showInActive {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{"is_active": true}}})
	}

	// Sort Options
	sortBy := bson.D{{Key: "$sort", Value: *srtOpts}}
	pipeline = append(pipeline, sortBy)

	// Pagination Options
	if pgOpts.PaginateId != nil {
		filter := bson.D{{Key: "$match", Value: bson.M{"_id": bson.M{keySortBy: pgOpts.PaginateId}}}}
		pipeline = append(pipeline, filter)
	} else if pgOpts != nil {
		skip := bson.D{{Key: "$skip", Value: (pgOpts.PerPage * (pgOpts.PageNum - 1))}}
		pipeline = append(pipeline, skip)
	}
	// Limit
	limit := bson.D{{Key: "$limit", Value: pgOpts.PerPage}}
	pipeline = append(pipeline, limit)

	// Populate Languages
	languageLookUp := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "language",
		"localField":   "languages_ids",
		"foreignField": "_id",
		"as":           "spoken_languages",
	}}}
	pipeline = append(pipeline, languageLookUp)

	// Hospital Populate
	hospitalLookUp := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "hospital",
		"localField":   "hospital_id",
		"foreignField": "_id",
		"as":           "hospital",
	}}}
	unwindHospital := bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$hospital",
		"preserveNullAndEmptyArrays": true,
	}}}
	pipeline = append(pipeline, hospitalLookUp, unwindHospital)

	// Add Prefix to image
	setImagePrefixPipe := bson.D{{Key: "$set",
		Value: bson.M{"profile_pic.image_src": bson.M{"$cond": bson.D{
			{Key: "if", Value: bson.M{"$eq": bson.A{"$profile_pic", nil}}},
			{Key: "then", Value: nil},
			{Key: "else", Value: bson.M{"$concat": bson.A{r.imageBasePath, "/", "$profile_pic.original_image_path"}}},
		}}}}}
	setBlurImagePrefixPipe := bson.D{{Key: "$set",
		Value: bson.M{"profile_pic.blur_data_url": bson.M{"$cond": bson.D{
			{Key: "if", Value: bson.M{"$eq": bson.A{"$profile_pic", nil}}},
			{Key: "then", Value: nil},
			{Key: "else", Value: bson.M{"$concat": bson.A{r.imageBasePath, "/", "$profile_pic.blur_image_path"}}},
		}}}}}
	resetNullImagePipe := bson.D{{Key: "$set", Value: bson.M{
		"profile_pic": bson.M{"$cond": bson.D{
			{Key: "if", Value: bson.M{"$eq": bson.A{"$profile_pic.image_src", nil}}},
			{Key: "then", Value: nil},
			{Key: "else", Value: "$profile_pic"},
		}},
	}}}
	pipeline = append(pipeline, setImagePrefixPipe, setBlurImagePrefixPipe, resetNullImagePipe)

	if slotsExistsOn == nil {
		// Next Available Slot
		nextAvailableLookUp := bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "appointment_slot",
			"localField":   "_id",
			"foreignField": "doctor_id",
			"as":           "next_available_slot",
			"pipeline": bson.A{
				bson.M{"$match": bson.M{"$and": bson.A{
					bson.M{"is_available": true},
					bson.M{"start": bson.M{"$gt": primitive.NewDateTimeFromTime(time.Now())}},
				}}},
				bson.M{"$limit": 1},
			},
		}}}

		pipeline = append(pipeline, nextAvailableLookUp)
	}

	unwindNextAvailableSlot := bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$next_available_slot",
		"preserveNullAndEmptyArrays": true,
	}}}

	pipeline = append(pipeline, unwindNextAvailableSlot)

	cur, err := r.colln.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}

	var res []doctormodel.DoctorEntity

	if err := cur.All(ctx, &res); err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())

	return res, nil

}

// Is Called only first first page.
func (r *DoctorRepository) GetDocumentsCount(
	fltrOpts *map[string]primitive.M,
	searchOpts *bson.M,
	showInActive bool,
	slotsExistsOn *primitive.DateTime,
) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	pipeline := mongo.Pipeline{}

	// Filter Doctor By Appointment Slot Availability
	if slotsExistsOn != nil {
		apptAvailLookup := bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "appointment_slot",
			"localField":   "_id",
			"foreignField": "doctor_id",
			"as":           "slot",
			"pipeline": bson.A{
				bson.M{"$match": bson.M{
					"$and": bson.A{
						bson.M{"is_available": true},
						bson.M{"date": slotsExistsOn},
					},
				}},
				bson.M{"$limit": 1},
			},
		}}}

		// Filter doctors having slots
		apptAvailMatch := bson.D{{Key: "$match", Value: bson.M{
			"$expr": bson.M{
				"$gt": bson.A{bson.M{"$size": "$slot"}, 0},
			},
		}}}

		// Map the available slot to next available field
		apptAvailAdddFields := bson.D{{Key: "$addFields", Value: bson.M{
			"next_available_slot": "$slot",
		}}}

		pipeline = append(pipeline, apptAvailLookup, apptAvailMatch, apptAvailAdddFields)
	}

	// Search Match
	if searchOpts != nil {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: *searchOpts}})
	}

	// Filter Options
	if len(*fltrOpts) != 0 {
		fltr := bson.D{{Key: "$match", Value: *fltrOpts}}
		pipeline = append(pipeline, fltr)
	}

	if !showInActive {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{"is_active": true}}})
	}

	countPipe := bson.D{{Key: "$count", Value: "total"}}
	pipeline = append(pipeline, countPipe)

	cur, err := r.colln.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var results []struct {
		Total int64 `bson:"total"`
	}

	cur.All(context.TODO(), &results)
	defer cur.Close(context.TODO())
	if len(results) > 0 {
		return results[0].Total, nil
	}

	return 0, nil

}

func (r *DoctorRepository) UpdateRefreshToken(id *primitive.ObjectID, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.UpdateByID(ctx, id, bson.M{"$set": bson.M{"refresh_token": token}})

	return err

}

func (r *DoctorRepository) UpdateById(id *primitive.ObjectID, payload *doctordto.EditDoctorDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.colln.UpdateByID(ctx, id, bson.M{"$set": payload})

	return err
}
