import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { ShoppingCartService } from '../../../service /ShoppinCartService';
import { CartItem } from '../../../Models/Cart';
import { OrderPeticion, OrderResponse } from '../../../Models/OrderPeticion';
import { catchError, of } from 'rxjs';
import { ArrowLeft } from '../../components/arrow-left/arrow-left';
import { Header } from '../../templates/header/header';

@Component({
  selector: 'app-shopping-cart',
  imports: [CommonModule, RouterModule, FormsModule, ArrowLeft, Header],
  templateUrl: './shopping-cart.html',
  styleUrl: './shopping-cart.css',
})
export class ShoppingCart implements OnInit {
  cartItems: CartItem[] = [];
  isLoading = true;
  error: string | null = null;
  
  // Propiedades para el checkout
  isCreatingOrder = false;
  shippingAddress = '';
  orderSuccess = false;
  createdOrderResponse: OrderResponse | null = null;

  constructor(private shoppingCartService: ShoppingCartService) {}

  ngOnInit() {
    this.loadCart();
  }

  /**
   * Carga el carrito de compras del usuario
   */
  loadCart() {
    this.isLoading = true;
    this.error = null;

    this.shoppingCartService.getShoppingCart().pipe(
      catchError(error => {
        console.error('Error al cargar el carrito:', error);
        this.error = 'Error al cargar el carrito de compras';
        this.isLoading = false;
        // Forzar detección de cambios
        setTimeout(() => this.isLoading = false, 0);
        return of([]);
      })
    ).subscribe(cartItems => {
      this.cartItems = cartItems || [];
      this.isLoading = false;
      // Forzar detección de cambios
      setTimeout(() => {
        this.isLoading = false;
        this.cartItems = [...(cartItems || [])];
      }, 0);
    });
  }

  /**
   * Actualiza la cantidad de un producto en el carrito
   */
  updateQuantity(item: CartItem, newQuantity: number) {
    if (newQuantity <= 0) {
      this.removeItem(item);
      return;
    }

    this.shoppingCartService.updateCartItem(item.id, newQuantity).pipe(
      catchError(error => {
        console.error('Error al actualizar cantidad:', error);
        return of(null);
      })
    ).subscribe(updatedItem => {
      if (updatedItem) {
        this.loadCart();
        // Forzar actualización adicional
        setTimeout(() => this.loadCart(), 100);
      }
    });
  }

  /**
   * Incrementa la cantidad de un producto
   */
  incrementQuantity(item: CartItem) {
    this.updateQuantity(item, item.quantity + 1);
  }

  /**
   * Decrementa la cantidad de un producto
   */
  decrementQuantity(item: CartItem) {
    this.updateQuantity(item, item.quantity - 1);
  }

  /**
   * Elimina un producto del carrito
   */
  removeItem(item: CartItem) {
    if (confirm('¿Estás seguro de que quieres eliminar este producto del carrito?')) {
      this.shoppingCartService.removeFromCart(item.id).pipe(
        catchError(error => {
          console.error('Error al eliminar producto:', error);
          return of(null);
        })
      ).subscribe(() => {
        this.loadCart();
      });
    }
  }

  /**
   * Vacía todo el carrito
   */
  clearCart() {
    if (confirm('¿Estás seguro de que quieres vaciar todo el carrito?')) {
      this.shoppingCartService.clearCart().pipe(
        catchError(error => {
          console.error('Error al vaciar carrito:', error);
          return of(null);
        })
      ).subscribe(() => {
        this.loadCart();
      });
    }
  }

  /**
   * Verifica si el carrito está vacío
   */
  get isCartEmpty(): boolean {
    return !this.cartItems || this.cartItems.length === 0;
  }

  /**
   * Obtiene el precio total formateado
   */
  get formattedTotal(): string {
    const total = this.cartItems.reduce((sum, item) => sum + (item.product.price * item.quantity), 0);
    return total.toLocaleString('es-CO', {
      style: 'currency',
      currency: 'COP'
    });
  }

  /**
   * Formatea el precio de un producto
   */
  formatPrice(price: number): string {
    return price.toLocaleString('es-CO', {
      style: 'currency',
      currency: 'COP'
    });
  }

  /**
   * Crea una nueva orden con los productos del carrito
   */
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

    const orderData: OrderPeticion = {
      shipping_address: this.shippingAddress.trim()
    };

    this.shoppingCartService.createOrder(orderData).pipe(
      catchError(error => {
        console.error('Error al crear la orden:', error);
        this.error = 'Error al crear la orden. Inténtalo de nuevo.';
        this.isCreatingOrder = false;
        // Forzar detección de cambios múltiples veces
        setTimeout(() => this.isCreatingOrder = false, 0);
        setTimeout(() => this.isCreatingOrder = false, 50);
        setTimeout(() => this.isCreatingOrder = false, 100);
        return of(null);
      })
    ).subscribe(orderResponse => {
      this.isCreatingOrder = false;
      
      if (orderResponse) {
        this.createdOrderResponse = orderResponse;
        this.orderSuccess = true;
        this.loadCart();
        this.shippingAddress = '';
        // Forzar detección de cambios múltiples veces
        setTimeout(() => {
          this.orderSuccess = true;
          this.isCreatingOrder = false;
        }, 0);
        setTimeout(() => {
          this.orderSuccess = true;
        }, 50);
        setTimeout(() => {
          this.orderSuccess = true;
        }, 100);
      }
    });
  }

  /**
   * Cierra el modal de éxito de la orden
   */
  closeOrderSuccess() {
    this.orderSuccess = false;
    this.createdOrderResponse = null;
    // Forzar detección de cambios múltiples veces
    setTimeout(() => {
      this.orderSuccess = false;
      this.createdOrderResponse = null;
    }, 0);
    setTimeout(() => {
      this.orderSuccess = false;
    }, 50);
    setTimeout(() => {
      this.orderSuccess = false;
    }, 100);
  }

  /**
   * Valida si se puede proceder al checkout
   */
  get canCheckout(): boolean {
    return !this.isCartEmpty && this.shippingAddress.trim().length > 0;
  }
}