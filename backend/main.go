package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Models
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Token    string `json:"token,omitempty" gorm:"default:null"`
}

type Item struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description"`
	Price       float64 `json:"price" gorm:"not null"`
}

type Cart struct {
	ID     uint        `json:"id" gorm:"primaryKey"`
	UserID uint        `json:"user_id" gorm:"not null"`
	User   User        `json:"user" gorm:"foreignKey:UserID"`
	Items  []CartItem  `json:"items" gorm:"foreignKey:CartID"`
}

type CartItem struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	CartID uint `json:"cart_id" gorm:"not null"`
	ItemID uint `json:"item_id" gorm:"not null"`
	Item   Item `json:"item" gorm:"foreignKey:ItemID"`
}

type Order struct {
	ID     uint         `json:"id" gorm:"primaryKey"`
	UserID uint         `json:"user_id" gorm:"not null"`
	User   User         `json:"user" gorm:"foreignKey:UserID"`
	Items  []OrderItem  `json:"items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID      uint `json:"id" gorm:"primaryKey"`
	OrderID uint `json:"order_id" gorm:"not null"`
	ItemID  uint `json:"item_id" gorm:"not null"`
	Item    Item `json:"item" gorm:"foreignKey:ItemID"`
}

// Request/Response structs
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type AddToCartRequest struct {
	ItemID uint `json:"item_id" binding:"required"`
}

type CreateOrderRequest struct {
	CartID uint `json:"cart_id" binding:"required"`
}

var db *gorm.DB
var jwtSecret = []byte("your-secret-key")

func main() {
	// Initialize database
	var err error
	db, err = gorm.Open(sqlite.Open("shopping_cart.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto migrate
	db.AutoMigrate(&User{}, &Item{}, &Cart{}, &CartItem{}, &Order{}, &OrderItem{})

	// Seed some items for testing
	seedItems()

	// Setup Gin
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes
	r.POST("/users", createUser)
	r.GET("/users", getUsers)
	r.POST("/users/login", loginUser)
	r.POST("/items", createItem)
	r.GET("/items", getItems)
	r.POST("/carts", authMiddleware(), addToCart)
	r.GET("/carts", authMiddleware(), getCarts)
	r.POST("/orders", authMiddleware(), createOrder)
	r.GET("/orders", authMiddleware(), getOrders)

	r.Run(":8000")
}

func seedItems() {
	var count int64
	db.Model(&Item{}).Count(&count)
	if count == 0 {
		items := []Item{
			{Name: "Laptop", Description: "High-performance laptop", Price: 999.99},
			{Name: "Mouse", Description: "Wireless mouse", Price: 29.99},
			{Name: "Keyboard", Description: "Mechanical keyboard", Price: 89.99},
			{Name: "Monitor", Description: "4K monitor", Price: 299.99},
			{Name: "Headphones", Description: "Noise-cancelling headphones", Price: 199.99},
		}
		db.Create(&items)
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["user_id"].(float64))
		c.Set("user_id", userID)
		c.Next()
	}
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	user.Password = "" // Don't return password
	c.JSON(http.StatusCreated, user)
}

func getUsers(c *gin.Context) {
	var users []User
	db.Select("id, username").Find(&users)
	c.JSON(http.StatusOK, users)
}

func loginUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/password"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Update user token in database
	db.Model(&user).Update("token", tokenString)

	user.Password = "" // Don't return password
	c.JSON(http.StatusOK, LoginResponse{Token: tokenString, User: user})
}

func createItem(c *gin.Context) {
	var item Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func getItems(c *gin.Context) {
	var items []Item
	db.Find(&items)
	c.JSON(http.StatusOK, items)
}

func addToCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if item exists
	var item Item
	if err := db.First(&item, req.ItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Find or create cart for user
	var cart Cart
	if err := db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		// Create new cart
		cart = Cart{UserID: userID}
		db.Create(&cart)
	}

	// Add item to cart
	cartItem := CartItem{CartID: cart.ID, ItemID: req.ItemID}
	if err := db.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Item added to cart"})
}

func getCarts(c *gin.Context) {
	userID := c.GetUint("user_id")
	var cart Cart
	if err := db.Preload("Items.Item").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func createOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find cart
	var cart Cart
	if err := db.Preload("Items").Where("id = ? AND user_id = ?", req.CartID, userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	if len(cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Create order
	order := Order{UserID: userID}
	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Copy cart items to order items
	for _, cartItem := range cart.Items {
		orderItem := OrderItem{OrderID: order.ID, ItemID: cartItem.ItemID}
		db.Create(&orderItem)
	}

	// Clear cart
	db.Where("cart_id = ?", cart.ID).Delete(&CartItem{})

	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully", "order_id": order.ID})
}

func getOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	var orders []Order
	db.Preload("Items.Item").Where("user_id = ?", userID).Find(&orders)
	c.JSON(http.StatusOK, orders)
}