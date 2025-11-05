import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, from, switchMap, combineLatest } from 'rxjs';
import { map, filter, take } from 'rxjs/operators';
import { ShoppingPeticion } from '../Models/ShoppingPeticion';
import { AuthService } from './Authser.vice';
import { CartItem, ShoppingCart } from '../Models/Cart';
import { OrderResponse } from '../Models/OrderPeticion';

@Injectable({
  providedIn: 'root'
})
export class ShoppingCartService {
  private apiUrlShoppingCart = 'http://localhost:80/';

  constructor(
    private http: HttpClient,
    private authService: AuthService
  ) {}

  /**
   * Método privado para obtener headers con autenticación
   */
  private getAuthHeaders(): Observable<HttpHeaders> {
    return from(this.authService.getToken()).pipe(
      map(token => {
        if (!token) {
          throw new Error('No se pudo obtener el token de autenticación');
        }
        return new HttpHeaders({
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        });
      })
    );
  }

  /**
   * Combina tenant ID y headers de autenticación
   */
  private getTenantAndHeaders(): Observable<{ tenantId: string; headers: HttpHeaders }> {
    return combineLatest([
      this.authService.idTenant$.pipe(
        filter((tenantId): tenantId is string => !!tenantId),
        take(1)
      ),
      this.getAuthHeaders()
    ]).pipe(
      map(([tenantId, headers]) => ({ tenantId, headers }))
    );
  }

  /**
   * Obtiene el carrito de compras del usuario actual
   */
  getShoppingCart(): Observable<CartItem[]> {
    return this.getTenantAndHeaders().pipe(
      switchMap(({ tenantId, headers }) => {
        const url = `${this.apiUrlShoppingCart}${tenantId}/api/cart`;
        return this.http.get<CartItem[]>(url, { headers });
      })
    );
  }

  /**
   * Agrega un producto al carrito de compras
   */
  addProductToCart(shoppingItem: ShoppingPeticion): Observable<CartItem> {
    return this.getTenantAndHeaders().pipe(
      switchMap(({ tenantId, headers }) => {
        const url = `${this.apiUrlShoppingCart}${tenantId}/api/cart/items`;
        return this.http.post<CartItem>(url, shoppingItem, { headers });
      })
    );
  }

  /**
   * Actualiza la cantidad de un producto en el carrito
   */
  updateCartItem(itemId: string, quantity: number): Observable<CartItem> {
    return this.getTenantAndHeaders().pipe(
      switchMap(({ tenantId, headers }) => {
        const url = `${this.apiUrlShoppingCart}${tenantId}/api/users/me/cart/${itemId}`;
        const updateData = { quantity };
        return this.http.put<CartItem>(url, updateData, { headers });
      })
    );
  }

  /**
   * Elimina un producto del carrito
   */
  removeFromCart(itemId: string): Observable<any> {
    return this.getTenantAndHeaders().pipe(
      switchMap(({ tenantId, headers }) => {
        const url = `${this.apiUrlShoppingCart}${tenantId}/api/users/me/cart/${itemId}`;
        return this.http.delete<any>(url, { headers });
      })
    );
  }

  /**
   * Vacía completamente el carrito de compras
   */
  clearCart(): Observable<any> {
    return this.getTenantAndHeaders().pipe(
      switchMap(({ tenantId, headers }) => {
        const url = `${this.apiUrlShoppingCart}${tenantId}/api/users/me/cart`;
        return this.http.delete<any>(url, { headers });
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
    return this.getTenantAndHeaders().pipe(
      switchMap(({ tenantId, headers }) => {
        const url = `${this.apiUrlShoppingCart}${tenantId}/api/orders`;
        return this.http.post<OrderResponse>(url, orderData, { headers });
      })
    );
  }

  /**
   * Obtiene las órdenes del usuario actual
   */
  getUserOrders(): Observable<any[]> {
    return this.getTenantAndHeaders().pipe(
      switchMap(({ tenantId, headers }) => {
        const url = `${this.apiUrlShoppingCart}${tenantId}/api/orders`;
        return this.http.get<any[]>(url, { headers });
      })
    );
  }

  /**
   * Obtiene una orden específica por su ID
   */
  getOrderById(orderId: string): Observable<any> {
    return this.getTenantAndHeaders().pipe(
      switchMap(({ tenantId, headers }) => {
        const url = `${this.apiUrlShoppingCart}${tenantId}/api/orders/${orderId}`;
        return this.http.get<any>(url, { headers });
      })
    );
  }
}