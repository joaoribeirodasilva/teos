package models

/* type iModel interface {
	GetID() primitive.ObjectID
	SetID(id primitive.ObjectID)
	GetCreatedBy() primitive.ObjectID
	SetCreatedBy(createdBy primitive.ObjectID)
	GetCreatedAt() time.Time
	SetCreatedAt(createdAt time.Time)
	GetUpdatedBy() primitive.ObjectID
	SetUpdatedBy(updatedBy primitive.ObjectID)
	GetUpdatedAt() time.Time
	SetUpdatedAt(updatedAt time.Time)
	GetDeletedBy() *primitive.ObjectID
	SetDeletedBy(deletedBy *primitive.ObjectID)
	GetDeletedAt() *time.Time
	SetDeletedAt(deletedAt *time.Time)
	GetValues() *structures.RequestValues
	SetValues(values *structures.RequestValues)
	Validate(validateMeta bool) *service_errors.Error
	GetCollection() *mongo.Collection
	GetLocation() string
	SetLocation(location string)
	FillMeta(create bool, delete bool) *service_errors.Error
	SetFields(from interface{}, to interface{})
} */

/* type Model struct {
	ID             primitive.ObjectID  `json:"_id" bson:"_id"`
	CreatedBy      primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt      time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy      primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt      time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy      *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt      *time.Time          `json:"deletedAt" bson:"deletedAt"`
	values         *structures.RequestValues
	location       string
	collectionName string
	collection     *mongo.Collection
} */

/*
	 func (m *Model) GetID() primitive.ObjectID {
		return primitive.NewObjectID()
	}

func (m *Model) SetID(id primitive.ObjectID) {}

	func (m *Model) GetCreatedBy() primitive.ObjectID {
		return primitive.NewObjectID()
	}

func (m *Model) SetCreatedBy(createdBy primitive.ObjectID) {}

	func (m *Model) GetCreatedAt() time.Time {
		return time.Now()
	}

func (m *Model) SetCreatedAt(createdAt time.Time) {}

	func (m *Model) GetUpdatedBy() primitive.ObjectID {
		return primitive.NewObjectID()
	}

func (m *Model) SetUpdatedBy(updatedBy primitive.ObjectID)

	func (m *Model) GetUpdatedAt() time.Time {
		return time.Now()
	}

func (m *AppApp) SetUpdatedAt(updatedAt time.Time) {}

	func (m *AppApp) GetDeletedBy() *primitive.ObjectID {
		id := primitive.NewObjectID()
		return &id
	}

func (m *AppApp) SetDeletedBy(deletedBy *primitive.ObjectID) {}
func (m *AppApp) GetDeletedAt() *time.Time
func (m *AppApp) SetDeletedAt(deletedAt *time.Time) {}
func (m *AppApp) GetValues() *structures.RequestValues
func (m *AppApp) SetValues(values *structures.RequestValues) {}
*/
/* func (m *Model) Validate(validateMeta bool) *service_errors.Error
func (m *Model) GetCollection() *mongo.Collection */

/* func (m *AppApp) GetLocation() string {}
func (m *AppApp) SetLocation(location string)
func (m *AppApp) FillMeta(create bool, delete bool) *service_errors.Error {}
func (m *AppApp) SetFields(from interface{}, to interface{}) */

/* func (m *Model) FindByID(st iModel, id primitive.ObjectID, opts ...*options.FindOneOptions) *service_errors.Error {

	return m.First(st, bson.D{{Key: "_id", Value: id}}, opts...)
} */

/* func (m *Model) First(st iModel, filter interface{}, opts ...*options.FindOneOptions) *service_errors.Error {

	if err := m.collection.FindOne(context.TODO(), filter, opts...).Decode(st); err != nil {
		if err != mongo.ErrNoDocuments {
			return service_log.Error(0, http.StatusNotFound, st.GetLocation(), "", "document not found")
		}
		return service_log.Error(0, http.StatusInternalServerError, st.GetLocation(), "", "failed to query database. ERR: %s", err.Error())
	}
	return nil
}

func (m *Model) Exists(st iModel, filter interface{}) *service_errors.Error {

	if appErr := m.First(st, filter); appErr != nil {
		return appErr
	}
	return nil
} */

/* func (m *Model) Create(st iModel, query *bson.D, unique string, opts ...*options.InsertOneOptions) *service_errors.Error {

	// Checks if uniques are violated
	existsDoc := Model{}
	if appErr := existsDoc.Exists(st, query); appErr != nil {
		if appErr.HttpCode != http.StatusNotFound {
			return appErr
		}
	} else {
		return service_log.Error(0, http.StatusConflict, st.GetLocation(), unique, "document already exists with id: %s", existsDoc.GetID().Hex())
	}

	// Fill metadata fields
	if appErr := st.FillMeta(true, false); appErr != nil {
		return appErr
	}

	// Validate full document
	if appErr := st.Validate(true); appErr != nil {
		return appErr
	}

	// Insert document into the database
	if _, err := st.GetCollection().InsertOne(context.TODO(), st, opts...); err != nil {
		return service_log.Error(0, http.StatusInternalServerError, st.GetLocation(), "", "failed to insert document into database. ERR: %s", err.Error())
	}

	return nil
} */

/* func (m *Model) Update(st iModel, id primitive.ObjectID, query bson.D, opts ...*options.FindOneOptions) *service_errors.Error {

	doc := AppApp{}

	if appErr := m.Exists(st, query); appErr != nil {
		return appErr
	} else {
		return service_log.Error(0, http.StatusConflict, "MODELS::AppApp::Update", "appKey", "document already exists with id: %s", doc.ID.Hex())
	}

	//oldDoc := *m

	doc.Name = st.Name
	doc.Description = m.Description
	doc.AppKey = m.AppKey
	doc.Active = m.Active

	// Fill metadata fields
	if appErr := m.FillMeta(false, false); appErr != nil {
		return appErr
	}

	// Update the record
	if _, err := doc.GetCollection().UpdateOne(context.TODO(), bson.D{
		{Key: "_id"},
	}, bson.D{
		{Key: "$set", Value: doc},
	}); err != nil {
		return service_log.Error(0, http.StatusInternalServerError, "MODELS::AppApp::Update", "", "failed to update into database. ERR: %s", err.Error())
	}

	// Save old record in history
	if appErr := oldDoc.SaveHistory(); appErr != nil {
		return appErr
	}

	return nil
} */

/*
 */
