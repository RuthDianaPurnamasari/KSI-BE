package module

import (
	"context"
	"errors"
	"fmt"
	// "image"
	// "image/jpeg"
	// "image/png"
	// "net/http"
	// "os"
	// "path/filepath"
	"time"

	"github.com/Nidasakinaa/be_KaloriKu/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnect(dbname string) (db *mongo.Database) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MongoString))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

func InsertOneDoc(db string, collection string, doc interface{}) (insertedID interface{}) {
	insertResult, err := MongoConnect(db).Collection(collection).InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	return insertResult.InsertedID
}

func StaticAdminLogin(db *mongo.Database, col string, username, password string) (bool, error) {
	// Validasi input kosong
	if username == "" || password == "" {
		return false, errors.New("username and password cannot be empty")
	}

	// Placeholder untuk pengembangan ke database (jika diperlukan)
	// Implementasi statis sementara hanya memeriksa kecocokan dengan parameter
	return false, errors.New("invalid admin credentials")
}

// FUNCTION MENU ITEM
// // Helper function to download image from URL
// func downloadImage(url string) (image.Image, string, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	defer resp.Body.Close()

// 	img, format, err := image.Decode(resp.Body)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	if format != "jpeg" && format != "png" {
// 		return nil, "", errors.New("unsupported image format")
// 	}

// 	return img, format, nil
// }

// // Helper function to save image locally
// func saveImage(img image.Image, path string, format string) error {
// 	// Create the directory if it doesn't exist
// 	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
// 	if err != nil {
// 		return err
// 	}

// 	file, err := os.Create(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	switch format {
// 	case "jpeg":
// 		err = jpeg.Encode(file, img, nil)
// 	case "png":
// 		err = png.Encode(file, img)
// 	default:
// 		return errors.New("unsupported image format")
// 	}

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Helper function to display image in terminal
// func displayImage(img image.Image) {
// 	fmt.Println("Displaying image is not supported in this terminal.")
// }

// // Helper function to load image from local path
// func loadImage(path string) (image.Image, error) {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	img, _, err := image.Decode(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return img, nil
// }

// // GetMenuItemByIDAndDisplayImage retrieves a menu item by its ID and displays the image
// func GetMenuItemByIDAndDisplayImage(_id primitive.ObjectID, db *mongo.Database, col string) (model.MenuItem, error) {
// 	menu, err := GetMenuItemByID(_id, db, col)
// 	if err != nil {
// 		return menu, err
// 	}

// 	// Load and display the image
// 	img, err := loadImage(menu.Image)
// 	if err != nil {
// 		return menu, fmt.Errorf("error loading image: %v", err)
// 	}
// 	displayImage(img)

// 	return menu, nil
// }

// GetMenuItemByID retrieves a menu item from the database by its ID
func GetMenuItemByID(_id primitive.ObjectID, db *mongo.Database, col string) (model.MenuItem, error) {
	var menu model.MenuItem
	collection := db.Collection("Menu")
	filter := bson.M{"_id": _id}
	err := collection.FindOne(context.TODO(), filter).Decode(&menu)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return menu, fmt.Errorf("GetMenuItemByID: menu item dengan ID %s tidak ditemukan", _id.Hex())
		}
		return menu, fmt.Errorf("GetMenuItemByID: gagal mendapatkan menu item: %w", err)
	}
	return menu, nil
}

// GetAllMenuItem retrieves all menu items from the database
func GetAllMenuItem(db *mongo.Database, col string) (data []model.MenuItem) {
	menu := db.Collection(col)
	filter := bson.M{}
	cursor, err := menu.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("GetAllMenuItem :", err)
	}
	err = cursor.All(context.TODO(), &data)
	if err != nil {
		fmt.Println(err)
	}
	return
}

// GetMenuByCategory retrieves all menu items from the database by its category
func GetMenuItemByCategory(category string, db *mongo.Database, col string) ([]model.MenuItem, error) {
	var menus []model.MenuItem
	collection := db.Collection("Menu")
	filter := bson.M{"category": category}

	cursor, err := collection.Find(context.TODO(), filter, options.Find())
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan menu items: %w", err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var menu model.MenuItem
		if err := cursor.Decode(&menu); err != nil {
			continue
		}
		menus = append(menus, menu)
	}

	if len(menus) == 0 {
		return nil, fmt.Errorf("menu item dengan category %s tidak ditemukan", category)
	}

	return menus, nil
}

// InsertMenuItem creates a new menu item in the database
func InsertMenuItem(db *mongo.Database, col string, name string, ingredients string, description string, calories float64, category string, image string) (insertedID primitive.ObjectID, err error) {

	menu := bson.M{
		"name":        name,
		"ingredients": ingredients,
		"description": description,
		"calories":    calories,
		"category":    category,
		"image":       image,
	}
	result, err := db.Collection(col).InsertOne(context.Background(), menu)
	if err != nil {
		fmt.Printf("InsertMenuItem: %v\n", err)
		return
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

// UpdateMenuItem updates an existing menu item in the database
func UpdateMenuItem(ctx context.Context, db *mongo.Database, col string, _id primitive.ObjectID, name string, ingredients string, description string, calories float64, category string, image string) (err error) {
	filter := bson.M{"_id": _id}
	update := bson.M{
		"$set": bson.M{
			"name":        name,
			"ingredients": ingredients,
			"description": description,
			"calories":    calories,
			"category":    category,
			"image":       image,
		},
	}
	result, err := db.Collection(col).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("UpdateMenuItem: gagal memperbarui Menu Item: %w", err)
	}
	if result.MatchedCount == 0 {
		return errors.New("UpdateMenuItem: tidak ada data yang diubah dengan ID yang ditentukan")
	}
	return nil
}

// DeleteMenuItemByID deletes a menu item from the database by its ID
func DeleteMenuItemByID(_id primitive.ObjectID, db *mongo.Database, col string) error {
	menu := db.Collection(col)
	filter := bson.M{"_id": _id}

	result, err := menu.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}

	return nil
}

// FUNCTION USER
// GetUserByID retrieves a user from the database by its ID
func GetUserByID(_id primitive.ObjectID, db *mongo.Database, col string) (model.User, error) {
	var user model.User
	collection := db.Collection("User")
	filter := bson.M{"_id": _id}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, fmt.Errorf("GetUserByID: user dengan ID %s tidak ditemukan", _id.Hex())
		}
		return user, fmt.Errorf("GetUserByID: gagal mendapatkan data user: %w", err)
	}
	return user, nil
}

func GetRoleByAdmin(db *mongo.Database, collection string, role string) (*model.User, error) {
	var user model.User
	filter := bson.M{"role": role}
	opts := options.FindOne()

	err := db.Collection(collection).FindOne(context.Background(), filter, opts).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func InsertUsers(db *mongo.Database, col string, fullname string, phonenumber string, username string, password string, role string) (insertedID primitive.ObjectID, err error) {
	users := bson.M{
		"fullname":    fullname,
		"phone": phonenumber,
		"username":    username,
		"password":    password,
		"role":        role,
	}
	result, err := db.Collection(col).InsertOne(context.Background(), users)
	if err != nil {
		fmt.Printf("InsertUser: %v\n", err)
		return
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func GetByUsername(db *mongo.Database, col string, username string) (*model.User, error) {
	var admin model.User
	err := db.Collection(col).FindOne(context.Background(), bson.M{"username": username}).Decode(&admin)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func DeleteTokenFromMongoDB(db *mongo.Database, col string, token string) error {
	collection := db.Collection(col)
	filter := bson.M{"token": token}

	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}

// GetAllUser retrieves all users from the database
func GetAllUser(db *mongo.Database, col string) (data []model.User) {
	user := db.Collection(col)
	filter := bson.M{}
	cursor, err := user.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("GetAllUser :", err)
	}
	err = cursor.All(context.TODO(), &data)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func SaveTokenToDatabase(db *mongo.Database, col string, adminID string, token string) error {
	collection := db.Collection(col)
	filter := bson.M{"admin_id": adminID}
	update := bson.M{
		"$set": bson.M{
			"token":      token,
			"updated_at": time.Now(),
		},
	}
	_, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

// InsertUser creates a new order in the database
func InsertUser(db *mongo.Database, col string, name string, phone string, username string, password string, role string) (insertedID primitive.ObjectID, err error) {
	user := bson.M{
		"name":     name,
		"phone":    phone,
		"username": username,
		"password": password,
		"role":     role,
	}
	result, err := db.Collection(col).InsertOne(context.Background(), user)
	if err != nil {
		fmt.Printf("InsertUser: %v\n", err)
		return
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

// UpdateUser updates an existing user in the database
func UpdateUser(ctx context.Context, db *mongo.Database, col string, _id primitive.ObjectID, name string, phone string, username string, password string, role string) (err error) {
	filter := bson.M{"_id": _id}
	update := bson.M{
		"$set": bson.M{
			"name":     name,
			"phone":    phone,
			"username": username,
			"password": password,
			"role":     role,
		},
	}
	result, err := db.Collection(col).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("UpdateUser: gagal memperbarui User: %w", err)
	}
	if result.MatchedCount == 0 {
		return errors.New("UpdateUser: tidak ada data yang diubah dengan ID yang ditentukan")
	}
	return nil
}

// DeleteUserByID deletes a menu item from the database by its ID
func DeleteUserByID(_id primitive.ObjectID, db *mongo.Database, col string) error {
	user := db.Collection(col)
	filter := bson.M{"_id": _id}

	result, err := user.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}

	return nil
}
