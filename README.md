# Shopping Cart Full-Stack Application

A complete e-commerce shopping cart application built with Go (Gin, GORM) backend and React frontend.

## Tech Stack

### Backend
- **Go** with Gin web framework
- **GORM** for database operations
- **SQLite** for database storage
- **JWT** for authentication
- **bcrypt** for password hashing

### Frontend
- **React** with hooks
- **Tailwind CSS** for styling
- **Local Storage** for token persistence

## Features

- User registration and authentication
- JWT-based session management
- Product catalog browsing
- Shopping cart functionality
- Order management
- Single cart per user
- Order history tracking

## Project Structure

```
shopping-cart/
├── backend/
│   ├── main.go
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── App.js
│   │   └── index.js
│   ├── public/
│   └── package.json
└── README.md
```

## API Endpoints

### User Management
- `POST /users` - Create a new user
- `GET /users` - List all users
- `POST /users/login` - User login

### Items
- `POST /items` - Create an item
- `GET /items` - List all items

### Cart (Requires Authentication)
- `POST /carts` - Add item to cart
- `GET /carts` - Get user's cart

### Orders (Requires Authentication)
- `POST /orders` - Create order from cart
- `GET /orders` - Get user's order history

## Setup Instructions

### Prerequisites
- Go 1.19+
- Node.js 16+
- npm or yarn

### Backend Setup

1. **Clone the repository**
```bash
git clone <your-repo-url>
cd shopping-cart/backend
```

2. **Initialize Go modules**
```bash
go mod init shopping-cart-backend
go mod tidy
```

3. **Install dependencies**
```bash
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get github.com/golang-jwt/jwt/v4
go get golang.org/x/crypto/bcrypt
```

4. **Run the backend**
```bash
go run main.go
```

The backend will start on `http://localhost:8080`

### Frontend Setup

1. **Navigate to frontend directory**
```bash
cd ../frontend
```

2. **Install dependencies**
```bash
npm install
```

3. **Add Tailwind CSS**
```bash
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

4. **Configure Tailwind** (tailwind.config.js)
```javascript
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
```

5. **Add Tailwind to CSS** (src/index.css)
```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

6. **Start the React app**
```bash
npm start
```

The frontend will start on `http://localhost:3000`

## Usage Guide

### 1. Create a Test User

First, create a user using the API:

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 2. Login to the Application

1. Open `http://localhost:3000` in your browser
2. Use the credentials:
   - Username: `testuser`
   - Password: `password123`

### 3. Shopping Flow

1. **Browse Items**: View all available items on the main page
2. **Add to Cart**: Click "Add to Cart" on any item
3. **View Cart**: Click the "Cart" button to see items in your cart
4. **Checkout**: Click "Checkout" to convert cart to order
5. **Order History**: Click "Order History" to view past orders

## API Testing with Postman

### Collection Setup

Create a Postman collection with the following requests:

#### 1. Create User
```
POST http://localhost:8080/users
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

#### 2. Login User
```
POST http://localhost:8080/users/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

#### 3. Get Items
```
GET http://localhost:8080/items
```

#### 4. Add to Cart (requires token)
```
POST http://localhost:8080/carts
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "item_id": 1
}
```

#### 5. Get Cart (requires token)
```
GET http://localhost:8080/carts
Authorization: Bearer <your-jwt-token>
```

#### 6. Create Order (requires token)
```
POST http://localhost:8080/orders
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "cart_id": 1
}
```

#### 7. Get Orders (requires token)
```
GET http://localhost:8080/orders
Authorization: Bearer <your-jwt-token>
```

## Database Schema

The application uses SQLite with the following tables:

- **users**: User accounts with authentication
- **items**: Product catalog
- **carts**: Shopping carts (one per user)
- **cart_items**: Items in carts
- **orders**: Completed orders
- **order_items**: Items in orders

## Development Notes

### Authentication
- JWT tokens expire after 24 hours
- Users can only be logged in from one device at a time
- Tokens are stored in localStorage on the frontend

### Cart Behavior
- Each user can have only one active cart
- Adding items to cart creates a new cart if none exists
- Cart is cleared when converted to order

### Default Data
The application seeds 5 sample items on startup:
- Laptop ($999.99)
- Mouse ($29.99)
- Keyboard ($89.99)
- Monitor ($299.99)
- Headphones ($199.99)

## Troubleshooting

### Common Issues

1. **CORS Errors**: Ensure the backend CORS configuration allows `http://localhost:3000`

2. **Authentication Failures**: Check that the JWT token is included in request headers

3. **Database Issues**: Delete `shopping_cart.db` file to reset the database

4. **Port Conflicts**: Ensure ports 8080 and 3000 are available

### Debug Tips

- Check browser console for frontend errors
- Check terminal output for backend logs
- Verify API responses using browser dev tools or Postman

## Future Enhancements

- User registration in frontend
- Item inventory management
- Multiple items quantity in cart
- Payment integration
- Order status tracking
- Admin panel for item management

## License

This project is created for educational purposes.
