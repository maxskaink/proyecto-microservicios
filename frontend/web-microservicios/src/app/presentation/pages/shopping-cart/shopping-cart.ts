import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { ShoppingCartService } from '../../../service/ShoppinCartService';
import { CartItem } from '../../../Models/Cart';
import { OrderPeticion, OrderResponse } from '../../../Models/OrderPeticion';
import { catchError, of } from 'rxjs';
import { Header } from '../../templates/header/header';

@Component({
  selector: 'app-shopping-cart',
  imports: [CommonModule, RouterModule, FormsModule, Header],
  templateUrl: './shopping-cart.html',
  styleUrl: './shopping-cart.css',
})
export class ShoppingCart implements OnInit {

  cartItems: CartItem[] = [];
  isLoading = true;
  error: string | null = null;

  // Checkout
  isCreatingOrder = false;
  shippingAddress = '';
  orderSuccess = false;
  createdOrderResponse: OrderResponse | null = null;

  constructor(private shoppingCartService: ShoppingCartService, private cdr: ChangeDetectorRef) {}

  ngOnInit() {
    this.loadCart();
  }

  loadCart() {
  this.isLoading = true;
  this.error = null;

  this.shoppingCartService.getShoppingCart().pipe(
    catchError(error => {
      console.error('Error al cargar el carrito:', error);
      this.error = 'Error al cargar el carrito de compras';
      this.isLoading = false;
      return of([]);
    })
  ).subscribe(cartItems => {

    // 🔥 Asignamos dentro del siguiente microtask
    Promise.resolve().then(() => {
      this.cartItems = cartItems || [];
      this.isLoading = false;
      this.cdr.detectChanges(); 
    });

  });
}

  /** Actualizar cantidad */
  updateQuantity(item: CartItem, newQuantity: number) {
    if (newQuantity <= 0) {
      this.removeItem(item);
      return;
    }

    this.shoppingCartService.updateCartItem(item.id, newQuantity)
      .pipe(
        catchError(err => {
          console.error(err);
          return of(null);
        })
      )
      .subscribe(result => {
        if (result) this.loadCart();
      });
  }

  incrementQuantity(item: CartItem) {
    this.updateQuantity(item, item.quantity + 1);
  }

  decrementQuantity(item: CartItem) {
    this.updateQuantity(item, item.quantity - 1);
  }

  /** Eliminar un item */
  removeItem(item: CartItem) {
    if (!confirm('¿Eliminar este producto del carrito?')) return;

    this.shoppingCartService.removeFromCart(item.id)
      .pipe(catchError(err => {
        console.error(err);
        return of(null);
      }))
      .subscribe(() => this.loadCart());
  }

  /** Vaciar carrito */
  clearCart() {
    if (!confirm('¿Vaciar todo el carrito?')) return;

    this.shoppingCartService.clearCart()
      .pipe(catchError(err => {
        console.error(err);
        return of(null);
      }))
      .subscribe(() => this.loadCart());
  }

  /** Carrito vacío */
  get isCartEmpty(): boolean {
    return !this.cartItems.length;
  }

  /** Total formateado */
  get formattedTotal(): string {
    const total = this.cartItems.reduce(
      (sum, item) => sum + item.product.price * item.quantity,
      0
    );

    return total.toLocaleString('es-CO', {
      style: 'currency',
      currency: 'COP'
    });
  }

  formatPrice(price: number): string {
    return price.toLocaleString('es-CO', {
      style: 'currency',
      currency: 'COP'
    });
  }

  /** Crear una orden */
  createOrder() {
    if (!this.shippingAddress.trim()) {
      alert('Por favor ingresa una dirección de envío');
      return;
    }

    if (this.isCartEmpty) {
      alert('No puedes crear una orden con el carrito vacío');
      return;
    }

    this.isCreatingOrder = true;
    this.error = null;

    const data: OrderPeticion = {
      shipping_address: this.shippingAddress.trim()
    };

    this.shoppingCartService.createOrder(data)
      .pipe(
        catchError(err => {
          console.error(err);
          this.error = 'Error al crear la orden, inténtalo de nuevo.';
          this.isCreatingOrder = false;
          return of(null);
        })
      )
      .subscribe(resp => {
        this.isCreatingOrder = false;

        if (resp) {
          this.createdOrderResponse = resp;
          this.orderSuccess = true;

          this.loadCart();
          this.shippingAddress = '';
        }
      });
  }

  /** Cerrar modal */
  closeOrderSuccess() {
    this.orderSuccess = false;
    this.createdOrderResponse = null;
  }

  /** Validar checkout */
  get canCheckout(): boolean {
    return !this.isCartEmpty && this.shippingAddress.trim().length > 0;
  }
}
