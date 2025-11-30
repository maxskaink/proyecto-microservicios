import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnChanges, OnInit, SimpleChanges } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Route, Router, RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { ShoppingCartService } from '../../../service/ShoppinCartService';
import { CartItem } from '../../../Models/Cart';
import { OrderPeticion, OrderResponse } from '../../../Models/OrderPeticion';
import { catchError, map, of, switchMap } from 'rxjs';
import { Header } from '../../templates/header/header';
import { ProductViewBuy } from '../../templates/product-view-buy/product-view-buy';
import { Order } from '../../templates/order/order';
import Swal from 'sweetalert2';
import { LoadingService } from '../../../service/loading-service';
import { IsLoading } from '../../components/is-loading/is-loading';

@Component({
  selector: 'app-shopping-cart',
  imports: [CommonModule, RouterModule, FormsModule, Header, ProductViewBuy, Order, IsLoading],
  templateUrl: './shopping-cart.html',
  styleUrl: './shopping-cart.css',
})
export class ShoppingCart implements OnChanges {

  cartItems: CartItem[] = [];
  isLoading = true;
  error: string | null = null;

  // Checkout
  isCreatingOrder = false;
  shippingAddress = '';
  orderSuccess = false;
  createdOrderResponse!: OrderResponse ;

  constructor( 
    private shoppingCartService: ShoppingCartService, 
    private cdr: ChangeDetectorRef,
    private router: Router,
    private isLoadign: LoadingService
  ) {}
  ngOnChanges(changes: SimpleChanges): void {
    throw new Error('Method not implemented.');
  }

  ngOnInit() {
    this.loadCart();
  }

  loadCart() {
  this.isLoadign.show("Cargando carrito");
  console.log("iniciando carrito")
  this.isLoading = true;
  this.error = null;

  this.shoppingCartService.getShoppingCart().pipe(
    catchError(error => {
      this.isLoadign.hide();
      console.error('Error al cargar el carrito:', error);
      this.error = 'Error al cargar el carrito de compras';
      this.isLoading = false;
      return of([]);
    })
  ).subscribe(cartItems => {
    this.isLoadign.hide();
    console.log("Guardando informacion", cartItems);
    // 🔥 Asignamos dentro del siguiente microtask
    Promise.resolve().then(() => {
      this.cartItems = cartItems || [];
      this.isLoading = false;
      this.cdr.detectChanges(); 
    });

  });
}

/** Eliminar un item */
removeItem(item: CartItem) {
  
  Swal.fire({
    title: '¿Eliminar este producto del carrito?',
    icon: 'warning',
    showCancelButton: true,
    confirmButtonText: 'Sí, eliminar',
    cancelButtonText: 'Cancelar',
    buttonsStyling: false,
    customClass: {
      confirmButton: 'btn btn-danger mx-2',
      cancelButton: 'btn btn-secondary mx-2'
    }
  }).then(result => {
    if (result.isConfirmed) {
      this.isLoadign.show("Eliminando producto del carrito");
      this.shoppingCartService.removeFromCart(item.id)
        .pipe(
          catchError(err => {
            this.isLoadign.hide();
            console.error(err);
            Swal.fire('Error', 'No se pudo eliminar el producto', 'error');
            return of(null);
          })
        )
        .subscribe(() => {
          this.isLoadign.hide();
          this.loadCart();
          Swal.fire('¡Eliminado!', 'El producto ha sido eliminado del carrito', 'success');
        });
    }
  });
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

  this.shoppingCartService.createOrder(data).pipe(

    catchError(err => {
      console.error(err);
      this.error = 'Error al crear la orden, inténtalo de nuevo.';
      this.isCreatingOrder = false;
      return of(null);
    }),

    // Si creó la orden, encadenar al updateStateOrder
    switchMap(resp => {
      if (!resp) return of(null); // si hubo error

      this.createdOrderResponse = resp;

      return this.shoppingCartService.updateStateOrder("paid", resp.order.id).pipe(
        map(() => resp) // devolver resp para usar después
      );
    })

  ).subscribe({

    next: (resp) => {
      this.isCreatingOrder = false;

      if (!resp) return;

      console.log("Orden creada:", resp);
      console.log("Estado actualizado a paid");

      this.orderSuccess = true;
      this.loadCart();
      this.shippingAddress = '';

      this.cdr.detectChanges();
    },

    error: (error) => {
      console.error("Error inesperado:", error);
      this.error = "Ocurrió un error, inténtalo de nuevo.";
      this.isCreatingOrder = false;
    }

  });
}

closeOrderSuccess(){
  this.router.navigate(['/home']);
}

  /** Validar checkout */
  get canCheckout(): boolean {
    return !this.isCartEmpty && this.shippingAddress.trim().length > 0;
  }
}
