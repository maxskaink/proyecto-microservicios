import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Order } from '../../../Models/OrderPeticion';
import { catchError, forkJoin, of } from 'rxjs';
import { ProductService } from '../../../service/ProductService';
import { Product } from '../../../Models/Product';
import { AnyARecord } from 'dns';

@Component({
  selector: 'app-item-order',
  imports: [CommonModule],
  templateUrl: './item-order.html',
  styleUrl: './item-order.css',
})
export class ItemOrder implements OnInit{
  /**
   * Orden a mostrar
   * action: accion del boton indicado, manda estado a actualizar y el ide del itema a actualizar
   */
  @Input() order!: Order;
  @Output() action = new EventEmitter<{ status: string; id: string }>();

  public products: Product[] = [];
  public isLoadingProducts = false;
  
  ngOnInit(): void {
    this.loadProducts();
  }
  constructor(
    private productService: ProductService,
    private cdr: ChangeDetectorRef
  ) {}

  /**
   * 
   */
  private loadProducts(): void {
      if (!this.order?.items?.length) {
        console.warn('⚠️ No hay items en la orden');
        return;
      }
  
      this.isLoadingProducts = true;
      console.log('🔄 Cargando productos para los items:', this.order.items);
  
      const productRequests = this.order.items.map(item => {
        console.log(`📞 Solicitando producto con ID: ${item.product_id}`);
        return this.productService.getProductById(item.product_id).pipe(
          catchError(error => {
            console.error(`❌ Error al cargar producto ${item.product_id}:`, error);
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
}
