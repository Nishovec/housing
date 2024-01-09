package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var client *mongo.Client

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	Email     string             `bson:"email,omitempty"`
	Password  string             `bson:"password,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/users", getUsersHandler)
	http.HandleFunc("/users/update", updateUserHandler)
	http.HandleFunc("/users/delete", deleteUserHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Housing</title>
			<style>
					body {
							font-family: Arial, sans-serif;
							background-color: #f4f4f4;
							margin: 0;
							padding: 0;
					}
	
					header {
							background-color: #333;
							color: #fff;
							padding: 10px;
							text-align: center;
					}
	
					header a {
							color: #fff;
							text-decoration: none;
							margin-right: 20px;
					}
	
					header a:hover {
							text-decoration: underline;
					}
	
					.logo {
							display: block;
							margin: 20px auto;
							max-width: 200px;
					}
	
					.main-content {
							max-width: 800px;
							margin: 20px auto;
							padding: 20px;
							background-color: #fff;
							border-radius: 10px;
							box-shadow: 0 0 20px rgba(0, 0, 0, 0.2);
					}
	
					.block {
							background-color: #e0e0e0;
							padding: 20px;
							margin: 20px 0;
							border-radius: 10px;
					}

					.block-1 {
						background-color: green;
					}
					
					.block-2 {
						background-color: orange;
					}

					.block-1 {
						background-color: yellow;
					}

					.logo {
						font-family: sans-serif;
						font-size: 30px;
					}
			</style>
	</head>
	<body>
			<header>
					<a href="/">Home</a>
					<a href="/register">Register</a>
					<p class="logo">Housing</p>
			</header>
	
			<div class="main-content">
					<h1>Welcome to Housing</h1>
	
					<div class="block-1 block">
							<h2>Search for Your Dream Home</h2>
							<p>Explore properties in different locations and find the perfect home for you.</p>
					</div>
	
					<div class="block block-2">
							<h2>Featured Properties</h2>
							<p>Discover our featured properties with beautiful images and detailed descriptions.</p>
					</div>
	
					<div class="block block-3">
							<h2>Get Started Today</h2>
							<p>Register now to access exclusive features and save your favorite properties.</p>
					</div>
			</div>
	</body>
	</html>
	
`)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl, err := template.New("register").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Register</title>
		<style>
    body {
        font-family: Arial, sans-serif;
        background-color: #f4f4f4;
        margin: 0;
        padding: 0;
    }

    h1 {
        text-align: center;
        color: #333;
    }

    form {
        max-width: 400px;
        margin: 20px auto;
        padding: 20px;
        background-color: #fff;
        border-radius: 5px;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    }

    label {
        display: block;
        margin-bottom: 8px;
    }

    input {
        width: 100%;
        padding: 8px;
        margin-bottom: 16px;
        box-sizing: border-box;
    }

    button {
        background-color: #4caf50;
        color: #fff;
        padding: 10px;
        border: none;
        border-radius: 5px;
        cursor: pointer;
    }

    button:hover {
        background-color: #45a049;
    }
</style>
</head>
<body>
    <h1>Register</h1>
    <form action="/register" method="post">
        <label for="name">Name:</label>
        <input type="text" id="name" name="name" required><br>

        <label for="email">Email:</label>
        <input type="email" id="email" name="email" required><br>

        <label for="password">Password:</label>
        <input type="password" id="password" name="password" required><br>

        <button type="submit">Register</button>
    </form>
</body>
</html>
`)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user := User{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	collection := client.Database("housing").Collection("users")
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Registration successful! Welcome, "+name+"!")
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	collection := client.Database("housing").Collection("users")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var users []User
	if err := cursor.All(context.Background(), &users); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jsonUsers, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUsers)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	collection := client.Database("housing").Collection("users")
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{"name": user.Name, "email": user.Email, "password": user.Password, "updated_at": time.Now()}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	collection := client.Database("housing").Collection("users")
	filter := bson.M{"_id": user.ID}

	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}