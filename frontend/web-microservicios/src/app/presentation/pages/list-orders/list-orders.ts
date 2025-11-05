import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { ShoppingCartService } from '../../../service /ShoppinCartService';
import { Header } from '../../templates/header/header';
import { ArrowLeft } from '../../components/arrow-left/arrow-left';
import { catchError, of } from 'rxjs';

@Component({
  selector: 'app-list-orders',
  imports: [CommonModule, RouterModule, Header, ArrowLeft],
  templateUrl: './list-orders.html',
  styleUrl: './list-orders.css',
})
export class ListOrders implements OnInit {
  orders: any[] = [];
  isLoading = true;
  error: string | null = null;

  constructor(private shoppingService: ShoppingCartService) {}

  ngOnInit() {
    this.loadOrders();
  }

  /**
   * Carga las órdenes del usuario
   */
  loadOrders() {
    this.isLoading = true;
    this.error = null;

    this.shoppingService.getUserOrders().pipe(
      catchError(error => {
        console.error('Error al cargar las órdenes:', error);
        this.error = 'Error al cargar las órdenes';
        this.isLoading = false;
        // Forzar detección de cambios
        setTimeout(() => this.isLoading = false, 0);
        return of([]);
      })
    ).subscribe(orders => {
      this.orders = orders || [];
      this.isLoading = false;
      // Forzar detección de cambios
      setTimeout(() => {
        this.isLoading = false;
        this.orders = [...(orders || [])];
      }, 0);
    });
  }

  /**
   * Formatea el precio
   */
  formatPrice(price: number): string {
    return price.toLocaleString('es-CO', {
      style: 'currency',
      currency: 'COP'
    });
  }

  /**
   * Formatea la fecha
   */
  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString('es-CO', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
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
    alert(`Detalles de la orden: ${orderId}`);
  }
}
