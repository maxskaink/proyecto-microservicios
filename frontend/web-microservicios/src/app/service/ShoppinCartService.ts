import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, from, switchMap, combineLatest, of } from 'rxjs';
import { map, filter, take } from 'rxjs/operators';
import { ShoppingPeticion } from '../Models/ShoppingPeticion';
import { AuthService } from './Authser.vice';
import { CartItem, ShoppingCart } from '../Models/Cart';
import { Order, OrderResponse } from '../Models/OrderPeticion';
import { TenantService } from './TenantService';

@Injectable({
  providedIn: 'root'
})
export class ShoppingCartService {
  private apiUrlShoppingCart = 'http://localhost:80/';

  constructor(
    private http: HttpClient,
    private authService: AuthService,
    private tenantService: TenantService
  ) {}




getShoppingCart(): Observable<CartItem[]> {
  return this.tenantService.getTenantId().pipe(
    switchMap((tenantId) => {
      console.log("Pidiendo carrito para tenant:", tenantId);
      const url = `${this.apiUrlShoppingCart}${tenantId}/api/cart`;
      return this.http.get<CartItem[]>(url);
    })
  );
}


  /**
   * Agrega un producto al carrito de compras
   */
  addProductToCart(shoppingItem: ShoppingPeticion): Observable<CartItem> {
    return this.tenantService.getTenantId().pipe(
      switchMap(( tenantId) => {
        const url = `${this.apiUrlShoppingCart}${tenantId}/api/cart/items`;
        return this.http.post<CartItem>(url, shoppingItem);
      })
    );
  }

  /**
   * Actualiza la cantidad de un producto en el carrito
   */
  updateCartItem(itemId: string, quantity: number): Observable<CartItem> {
    return this.tenantService.getTenantId().pipe(
      switchMap(( tenantId ) => {
        const url = `${this.apiUrlShoppingCart}${tenantId}/api/cart/items/${itemId}`;
        const updateData = { quantity };
        return this.http.put<CartItem>(url, updateData);
      })
    );
  }

  /**
   * Elimina un producto del carrito
   */
  removeFromCart(itemId: string): Observable<any> {
    return this.tenantService.getTenantId().pipe(
      switchMap(( tenantId ) => {
        const url = `${this.apiUrlShoppingCart}${tenantId}/api/cart/items/${itemId}`;
        return this.http.delete<any>(url);
      })
    );
  }

  /**
   * Vacía completamente el carrito de compras
   */
  clearCart(): Observable<any> {
    return this.tenantService.getTenantId().pipe(
      switchMap(( tenantId ) => {
        const url = `${this.apiUrlShoppingCart}${tenantId}/api/users/me/cart`;
        return this.http.delete<any>(url);
      })
    );
  }

  /**
   * Obtiene el número total de items en el carrito
   */
  getCartItemCount(): Observable<number> {
    return this.getShoppingCart().pipe(
      map((cartItems: CartItem[]) => {
        if (!cartItems || !Array.isArray(cartItems)) return 0;
        return cartItems.reduce((total, item) => total + item.quantity, 0);
      })
    );
  }

  /**
   * Obtiene el precio total del carrito
   */
  getCartTotalPrice(): Observable<number> {
    return this.getShoppingCart().pipe(
      map((cartItems: CartItem[]) => {
        if (!cartItems || !Array.isArray(cartItems)) return 0;
        return cartItems.reduce((total, item) => total + (item.product.price * item.quantity), 0);
      })
    );
  }

  /**
   * Crea una nueva orden con los productos del carrito actual
   * orderData is a plain object (e.g. { shipping_address: '...' })
   */
  createOrder(orderData: any): Observable<OrderResponse> {
    return this.tenantService.getTenantId().pipe(
      switchMap(( tenantId) => {
        const url = `${this.apiUrlShoppingCart}${tenantId}/api/orders`;
        return this.http.post<OrderResponse>(url, orderData);
      })
    );
  }

  /**
   * Obtiene las órdenes del usuario actual
   */
    getUserOrders(): Observable<Order[]> {
      return this.tenantService.getTenantId().pipe(
        switchMap((tenantId ) => {
          const url = `${this.apiUrlShoppingCart}${tenantId}/api/orders`;
          return this.http.get<Order[]>(url);
        })
      );
    }

  /**
   * Obtiene una orden específica por su ID
   */
  getOrderById(orderId: string): Observable<Order> {
    return this.tenantService.getTenantId().pipe(
      switchMap(( tenantId ) => {
        const url = `${this.apiUrlShoppingCart}${tenantId}/api/orders/${orderId}`;
        return this.http.get<Order>(url);
      })
    );
  }
}
