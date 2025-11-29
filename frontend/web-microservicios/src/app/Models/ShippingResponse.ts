export interface ShippingResponse {
      id: string;
  order_id: string;
  tracking_number: string;
  shipping_address: string;
  status:  string; 
  created_at: string;   // ISO date string
  updated_at: string;   // ISO date string
}