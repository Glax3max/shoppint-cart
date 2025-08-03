import React, { useState, useEffect } from 'react';

const API_BASE = 'http://localhost:8000';

function App() {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(localStorage.getItem('token'));
  const [items, setItems] = useState([]);
  const [cart, setCart] = useState(null);
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(false);

  // Login state
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  useEffect(() => {
    if (token) {
      fetchItems();
    }
  }, [token]);

  const login = async () => {
    setLoading(true);
    
    try {
      const response = await fetch(`${API_BASE}/users/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });

      if (response.ok) {
        const data = await response.json();
        setToken(data.token);
        setUser(data.user);
        localStorage.setItem('token', data.token);
      } else {
        window.alert('Invalid username/password');
      }
    } catch (error) {
      window.alert('Login failed. Please try again.');
    }
    
    setLoading(false);
  };

  const fetchItems = async () => {
    try {
      const response = await fetch(`${API_BASE}/items`);
      const data = await response.json();
      setItems(data);
    } catch (error) {
      console.error('Failed to fetch items:', error);
    }
  };

  const addToCart = async (itemId) => {
    try {
      const response = await fetch(`${API_BASE}/carts`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({ item_id: itemId }),
      });

      if (response.ok) {
        window.alert('Item added to cart!');
      } else {
        window.alert('Failed to add item to cart');
      }
    } catch (error) {
      window.alert('Failed to add item to cart');
    }
  };

  const viewCart = async () => {
    try {
      const response = await fetch(`${API_BASE}/carts`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.ok) {
        const cartData = await response.json();
        setCart(cartData);
        
        if (cartData.items && cartData.items.length > 0) {
          const cartItems = cartData.items.map(item => 
            `Cart ID: ${item.cart_id}, Item: ${item.item.name} (ID: ${item.item_id})`
          ).join('\n');
          window.alert(`Cart Items:\n${cartItems}`);
        } else {
          window.alert('Your cart is empty');
        }
      } else {
        window.alert('No cart found');
      }
    } catch (error) {
      window.alert('Failed to fetch cart');
    }
  };

  const viewOrderHistory = async () => {
    try {
      const response = await fetch(`${API_BASE}/orders`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.ok) {
        const ordersData = await response.json();
        setOrders(ordersData);
        
        if (ordersData && ordersData.length > 0) {
          const orderIds = ordersData.map(order => `Order ID: ${order.id}`).join('\n');
          window.alert(`Order History:\n${orderIds}`);
        } else {
          window.alert('No orders found');
        }
      } else {
        window.alert('Failed to fetch orders');
      }
    } catch (error) {
      window.alert('Failed to fetch order history');
    }
  };

  const checkout = async () => {
    if (!cart || !cart.items || cart.items.length === 0) {
      // Fetch current cart first
      try {
        const response = await fetch(`${API_BASE}/carts`, {
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        });

        if (response.ok) {
          const cartData = await response.json();
          if (!cartData.items || cartData.items.length === 0) {
            window.alert('Your cart is empty. Add some items first!');
            return;
          }
          setCart(cartData);
          
          // Proceed with checkout
          const orderResponse = await fetch(`${API_BASE}/orders`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify({ cart_id: cartData.id }),
          });

          if (orderResponse.ok) {
            window.alert('Order successful!');
            setCart(null); // Clear cart state
          } else {
            window.alert('Failed to create order');
          }
        } else {
          window.alert('Your cart is empty. Add some items first!');
        }
      } catch (error) {
        window.alert('Failed to process checkout');
      }
    } else {
      // Cart already loaded, proceed with checkout
      try {
        const response = await fetch(`${API_BASE}/orders`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
          },
          body: JSON.stringify({ cart_id: cart.id }),
        });

        if (response.ok) {
          window.alert('Order successful!');
          setCart(null); // Clear cart state
        } else {
          window.alert('Failed to create order');
        }
      } catch (error) {
        window.alert('Failed to process checkout');
      }
    }
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    setCart(null);
    setOrders([]);
    localStorage.removeItem('token');
    setUsername('');
    setPassword('');
  };

  if (!token) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="bg-white p-8 rounded-lg shadow-md w-96">
          <h1 className="text-2xl font-bold mb-6 text-center">Shopping Cart Login</h1>
          <div>
            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2">
                Username
              </label>
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-blue-500"
                required
              />
            </div>
            <div className="mb-6">
              <label className="block text-gray-700 text-sm font-bold mb-2">
                Password
              </label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-blue-500"
                required
              />
            </div>
            <button
              type="button"
              onClick={login}
              disabled={loading}
              className="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline disabled:opacity-50"
            >
              {loading ? 'Logging in...' : 'Login'}
            </button>
          </div>
          <div className="mt-4 text-sm text-gray-600">
            <p>Demo credentials:</p>
            <p>Username: testuser</p>
            <p>Password: password123</p>
            <p className="text-xs mt-2">Note: You'll need to create this user first using the API</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <h1 className="text-2xl font-bold text-gray-900">Shopping Portal</h1>
            <div className="flex items-center space-x-4">
              <span className="text-gray-700">Welcome, {user?.username}</span>
              <button
                onClick={logout}
                className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-xl font-semibold">Available Items</h2>
          <div className="space-x-2">
            <button
              onClick={checkout}
              className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
            >
              Checkout
            </button>
            <button
              onClick={viewCart}
              className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
            >
              Cart
            </button>
            <button
              onClick={viewOrderHistory}
              className="bg-purple-500 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded"
            >
              Order History
            </button>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {items.map((item) => (
            <div key={item.id} className="bg-white rounded-lg shadow-md p-6">
              <h3 className="text-lg font-semibold mb-2">{item.name}</h3>
              <p className="text-gray-600 mb-2">{item.description}</p>
              <p className="text-xl font-bold text-green-600 mb-4">${item.price}</p>
              <button
                onClick={() => addToCart(item.id)}
                className="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
              >
                Add to Cart
              </button>
            </div>
          ))}
        </div>

        {items.length === 0 && (
          <div className="text-center py-8">
            <p className="text-gray-500">No items available. Please check back later.</p>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;