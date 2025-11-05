export interface OrderPeticion {
    shipping_address: string;
}

export interface OrderItem {
    id: string;
    order_id: string;
    product_id: string;
    quantity: number;
    price: number;
    created_at: string;
}

export interface Order {
    id: string;
    user_id: string;
    total_price: number;
    status: string;
    items: OrderItem[];
    created_at: string;
    updated_at: string;
}

export interface Shipping {
    id: string;
    order_id: string;
    tracking_number: string;
    shipping_address: string;
    status: string;
    created_at: string;
    updated_at: string;
}

export interface OrderResponse {
    order: Order;
    shipping: Shipping;
}
