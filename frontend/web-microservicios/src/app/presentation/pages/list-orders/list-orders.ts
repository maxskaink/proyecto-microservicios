import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { ShoppingCartService } from '../../../service/ShoppinCartService';
import { Header } from '../../templates/header/header';
import { ArrowLeft } from '../../components/arrow-left/arrow-left';
import { catchError, of } from 'rxjs';
import { ItemOrder } from '../../components/item-order/item-order';
import { Order } from '../../../Models/OrderPeticion';
import { updateStatus } from '../../../Models/updateStatus';
import { LoadingService } from '../../../service/loading-service';
import { IsLoading } from '../../components/is-loading/is-loading';


@Component({
  selector: 'app-list-orders',
  imports: [CommonModule, RouterModule, Header, ItemOrder, IsLoading],
  templateUrl: './list-orders.html',
  styleUrl: './list-orders.css',
})
export class ListOrders implements OnInit {
  orders: Order[] = [];
  isLoading = true;
  error: string | null = null;
  updateStatus: updateStatus = {status: ''};
  constructor(
    private router: Router,
    private shoppingService: ShoppingCartService,
    private cdr: ChangeDetectorRef,
    private loadingService: LoadingService,
  ) {}

  ngOnInit() {
    console.log("iniciando componente");
    this.loadOrders();
  }

  /**
   * Carga las órdenes del usuario
   */
  loadOrders() {
    this.loadingService.show("cargando ordenes...");
    console.log("iniciando carga de ordenes");
    this.isLoading = true;
    this.error = null;

    this.shoppingService.getUserOrdersProducer().pipe(
      catchError(error => {
        console.error('Error al cargar las órdenes:', error);
        this.error = 'Error al cargar las órdenes';
        this.isLoading = false;
        // Forzar detección de cambios
        setTimeout(() => this.isLoading = false, 0);
        this.loadingService.hide();
        this.cdr.detectChanges();
        return of([]);
      })
    ).subscribe(orders => {

      this.orders = orders;
      this.isLoading = false;
      console.log("Órdenes cargadas:", orders);
      this.loadingService.hide();
      setTimeout(() => {
        this.isLoading = false;
        this.orders = [...(orders || [])];
        console.log("Órdenes actualizadas en el estado:", this.orders);
        this.cdr.detectChanges();
      }, 0);
    });
  }



  /**
   * Obtiene la clase CSS para el estado
   */
  getStatusClass(status: string): string {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'status-pending';
      case 'confirmed':
        return 'status-confirmed';
      case 'shipped':
        return 'status-shipped';
      case 'delivered':
        return 'status-delivered';
      case 'cancelled':
        return 'status-cancelled';
      default:
        return 'status-default';
    }
  }

  /**
   * Traduce el estado al español
   */
  translateStatus(status: string): string {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'Pendiente';
      case 'confirmed':
        return 'Confirmada';
      case 'shipped':
        return 'Enviada';
      case 'delivered':
        return 'Entregada';
      case 'cancelled':
        return 'Cancelada';
      default:
        return status;
    }
  }

  /**
   * Recarga las órdenes
   */
  refreshOrders() {
    this.loadOrders();
  }

  /**
   * Ver detalles de una orden específica
   */
  viewOrderDetails(orderId: string) {
    // Por ahora solo mostramos un alert, después se puede navegar a una página de detalles
    console.log('Ver detalles de la orden:', orderId);
    this.router.navigate(['list-order/view-order', orderId]);
  }

  handlerOrderAction(event: { status: string; id: string }) {
    this.updateStatus.status = event.status;
    this.shoppingService.updateStateOrder(event.status, event.id).subscribe({
      next: () => {
        this.loadOrders();
      },
      error: (error) => {
        console.error('❌ Error al completar la orden:', error);
      }
    });
  }
  
  
}
