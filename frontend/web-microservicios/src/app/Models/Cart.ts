export interface CartItem {
  id: string;
  user_id: string;
  product_id: string;
  product: {
    id: string;
    name: string;
    description: string;
    price: number;
    stock: number;
    created_at: string;
    updated_at: string;
  };
  quantity: number;
  created_at: string;
  updated_at: string;
}

export interface ShoppingCart {
  user_id: string;
  items: CartItem[];
  total_items: number;
  total_price: number;
  created_at: string;
  updated_at: string;
}
