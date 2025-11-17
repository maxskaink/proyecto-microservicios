import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Product } from '../../../Models/Product';
import { ProductService } from '../../../service/ProductService';
import { finalize, catchError, of } from 'rxjs';
import { Header } from '../../templates/header/header';
import { ArrowLeft } from '../../components/arrow-left/arrow-left';
import { ShoppingCart } from '../shopping-cart/shopping-cart';
import { ShoppingCartService } from '../../../service/ShoppinCartService';
import { ShoppingPeticion } from '../../../Models/ShoppingPeticion';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-view-product',
  imports: [CommonModule, Header, ArrowLeft, FormsModule],
  templateUrl: './view-product.html',
  styleUrl: './view-product.css',
})
export class ViewProduct implements OnInit {
  product?: Product;

  isLoading: boolean = false;
  productId: string = '';
  isAddingToCart: boolean = false;
  addToCartMessage: string = '';
  showSuccessMessage: boolean = false;
  public quantity: number= 1;


  constructor(
    private productService: ProductService,
    private cdr: ChangeDetectorRef,
    private route: ActivatedRoute,
    private shoppingService: ShoppingCartService
  ) {}

  ngOnInit(): void {
    // Obtener el ID del producto y tenant desde los parámetros de la ruta
    this.route.params.subscribe(params => {
      this.productId = params['id'];
      
      if (this.productId) {
        this.loadProduct();
      } else {
        console.warn('No se proporcionaron productId');
      }
    });
  }
  
  loadProduct(): void {
    if (!this.productId ) {
      console.warn('No se proporcionó un productId ');
      return;
    }

    this.isLoading = true;

    this.productService.getProductById(this.productId)
      .pipe(
        finalize(() => {
          this.isLoading = false;
          this.cdr.markForCheck(); 
        })
      )
      .subscribe({
        next: (product) => {
          this.product = product;
          console.log('Producto cargado:', product);
        },
        error: (err) => {
          console.error('Error al cargar el producto:', err);
        }
      });
  }
  increaseQty() {
    this.quantity++;
  }

  decreaseQty() {
    if (this.quantity > 1) {
      this.quantity--;
    }
  }

  /**
   * Agrega el producto actual al carrito de compras
   */
  addToCart(quantity: number): void {
    if (!this.product) {
      console.warn('No hay producto para agregar al carrito');
      return;
    }

    this.isAddingToCart = true;
    this.addToCartMessage = '';

    // Crear el objeto ShoppingPeticion
    const shoppingItem: ShoppingPeticion = {
      product_id: this.product.id,
      quantity: quantity
    };

    this.shoppingService.addProductToCart(shoppingItem).pipe(
      catchError(error => {
        console.error('Error al agregar producto al carrito:', error);
        this.addToCartMessage = 'Error al agregar el producto al carrito';
        return of(null);
      }),
      finalize(() => {
        this.isAddingToCart = false;
        this.cdr.markForCheck();
      })
    ).subscribe({
      next: (cartItem) => {
        if (cartItem) {
          console.log('Producto agregado al carrito:', cartItem);
          this.addToCartMessage = '¡Producto agregado al carrito exitosamente!';
          this.showSuccessMessage = true;
          
          // Limpiar el mensaje después de 4 segundos
          setTimeout(() => {
            this.showSuccessMessage = false;
            this.addToCartMessage = '';
            this.cdr.markForCheck();
          }, 4000);
        }
      }
    });
  }

  /**
   * Cierra la notificación de éxito manualmente
   */
  closeSuccessMessage(): void {
    this.showSuccessMessage = false;
    this.addToCartMessage = '';
  }

  /**
   * Formatea el precio del producto
   */
  formatPrice(price: number): string {
    return price.toLocaleString('es-CO', {
      style: 'currency',
      currency: 'COP'
    });
  }

decreaseQuantity(input: HTMLInputElement): void {
  const currentValue = +input.value;
  if (currentValue > 1) {
    input.value = (currentValue - 1).toString();
  }
}

increaseQuantity(input: HTMLInputElement, maxStock: number): void {
  const currentValue = +input.value;
  if (currentValue < maxStock) {
    input.value = (currentValue + 1).toString();
  }
}
}
