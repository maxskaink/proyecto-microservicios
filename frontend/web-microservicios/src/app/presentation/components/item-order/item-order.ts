import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges } from '@angular/core';
import { Order } from '../../../Models/OrderPeticion';
import { catchError, forkJoin, map, Observable, of, tap } from 'rxjs';
import { ProductService } from '../../../service/ProductService';
import { Product } from '../../../Models/Product';
import { AnyARecord } from 'dns';
import { ShippingService } from '../../../service/shipping.service';

@Component({
  selector: 'app-item-order',
  imports: [CommonModule],
  templateUrl: './item-order.html',
  styleUrl: './item-order.css',
})
export class ItemOrder implements OnChanges{
  /**
   * Orden a mostrar
   * action: accion del boton indicado, manda estado a actualizar y el ide del itema a actualizar
   */
  @Input() order!: Order;
  @Input() showActions: boolean = true;
  @Output() action = new EventEmitter<{ status: string; id: string }>();
  @Output() details = new EventEmitter<string>();

  public products: Product[] = [];
  public isLoadingProducts = false;
  public stateShipping: string = '';
  
  constructor(
    private productService: ProductService,
    private shippingService: ShippingService,
    private cdr: ChangeDetectorRef
  ) {}
  ngOnChanges(changes: SimpleChanges): void {
      if (changes['showActions']) {
    this.showActions = changes['showActions'].currentValue;
  }
  if (this.order && this.order.items?.length) {

    this.isLoadingProducts = true; // un solo loader general

    forkJoin({
      products: this.loadProducts(),
      shipping: this.loadShippingInfo()
    }).subscribe({
      next: ({ products, shipping }) => {
        console.log("📦 Todo cargado antes de pintar");

        
        this.isLoadingProducts = false;
        this.cdr.detectChanges();
      },
      error: (err) => {
        console.error("❌ Error cargando datos:", err);
        this.isLoadingProducts = false;
      }
    });
  }
}


private loadShippingInfo(): Observable<any> {
  return this.shippingService.getShippingByOrderId(this.order.id).pipe(
    tap((shipping) => {
      console.log('✅ Información de envío cargada:', shipping);
      this.stateShipping = shipping.status;
    }),
    catchError(error => {
      console.error('Error al cargar información de envío:', error);
      return of(null);
    })
  );
}


  /**
   * 
   */
  private loadProducts(): Observable<Product[]> {
  if (!this.order?.items?.length) {
    console.warn('⚠️ No hay items en la orden');
    return of([]);
  }

  const productRequests = this.order.items.map(item =>
    this.productService.getProductById(item.product_id).pipe(
      catchError(error => {
        console.error(`❌ Error al cargar producto ${item.product_id}:`, error);
        return of(null);
      })
    )
  );

  return forkJoin(productRequests).pipe(
    map((products: (Product | null)[]) =>
      products.filter((p): p is Product => p !== null) // ← cambia el tipo
    ),
    tap((filteredProducts: Product[]) => {
      this.products = filteredProducts;
      console.log('✅ Productos cargados:', this.products);
    })
  );
}

  formatState(state: string): string {
    switch (state) {
      case 'pending':
        return 'Pendiente';
      case 'in_transit':
        return 'En tramite';
      case 'shipped':
        return 'Enviada';
      case 'delivered':
        return 'Entregada';
      case 'cancelled':
        return 'Cancelada';
      default:
        return state;
    }
  }
  /**
   * Formatea la hora
   */
  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString('es-CO', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
  /**
   * Formatea el precio
   */
  formatPrice(price: number): string {
    return price.toLocaleString('es-CO', {
      style: 'currency',
      currency: 'COP',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    });
  }
  /**
   * emitir complete order
   */
  onCompleteOrder() {
    this.action.emit({
      status: "paid",
      id: this.order.id
    });
  }
  /**
   * emitir complete order
   */
  onCancelOrder() {
    this.action.emit({
      status: "cancelled",
      id: this.order.id
    });
  }
  onClickDetails() {
    this.details.emit(this.order.id);
  }
  get canShowActions(): boolean {
  return this.showActions && this.order?.status !== 'cancelled';
}

}
