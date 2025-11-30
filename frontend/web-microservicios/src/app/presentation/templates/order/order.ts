import { Component, EventEmitter, Input, OnChanges, Output, SimpleChanges, ChangeDetectorRef } from '@angular/core';
import { OrderResponse } from '../../../Models/OrderPeticion';
import { forkJoin, of } from 'rxjs';
import { CommonModule } from '@angular/common';
import { ProductBox } from '../../components/product-box/product-box';
import { ProductService } from '../../../service/ProductService';
import { Product } from '../../../Models/Product';
import { catchError } from 'rxjs/operators';

@Component({
  selector: 'app-order',
  imports: [CommonModule],
  templateUrl: './order.html',
  standalone: true,
  styleUrl: './order.css',
})
export class Order implements OnChanges {
  
  @Output() close = new EventEmitter<boolean>();
  @Input() createdOrderResponse!: OrderResponse;
  
  public products: Product[] = [];
  public isLoadingProducts = false;

  constructor(
    private productService: ProductService,
    private cdr: ChangeDetectorRef
  ) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['createdOrderResponse'] && this.createdOrderResponse) {
      this.loadProducts();
    }
  }

  private loadProducts(): void {
    if (!this.createdOrderResponse?.order?.items?.length) {
      return;
    }

    this.isLoadingProducts = true;
    console.log('🔄 Cargando productos para los items:', this.createdOrderResponse.order.items);

    const productRequests = this.createdOrderResponse.order.items.map(item => {
      console.log(`📞 Solicitando producto con ID: ${item.product_id}`);
      return this.productService.getProductById(item.product_id).pipe(
        catchError(error => {
          return of(null); // Continúa con los demás productos aunque uno falle
        })
      );
    });

    forkJoin(productRequests).subscribe({
      next: (products: (Product | null)[]) => {
        this.products = products.filter(product => product !== null) as Product[];
        this.isLoadingProducts = false;
        console.log('✅ Productos cargados:', this.products);
        this.cdr.detectChanges();
      },
      error: (error) => {
        console.error('❌ Error al cargar productos:', error);
        this.isLoadingProducts = false;
        this.products = [];
        this.cdr.detectChanges();
      }
    });
  }

  closeOrder(): void {
    this.close.emit(true);
  }

  formatPrice(price: number): string {
    return price.toLocaleString('es-CO', {
      style: 'currency',
      currency: 'COP'
    });
  }
}
