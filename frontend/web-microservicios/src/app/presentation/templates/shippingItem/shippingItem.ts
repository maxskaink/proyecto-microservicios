import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, Input, OnChanges, OnInit, SimpleChanges } from '@angular/core';
import { ShoppingCartService } from '../../../service/ShoppinCartService';
import { Shipping } from '../../../Models/OrderPeticion';
import { ShippingService } from '../../../service/shipping.service';
import { AuthService } from '../../../service/Authser.vice';

@Component({
  selector: 'app-shipping',
  imports: [CommonModule],
  templateUrl: './shippingItem.html',
  styleUrl: './shippingItem.css',
})
export class shippingItem implements OnChanges {
   showActions: boolean = false;
  @Input() shipping!: Shipping;
  constructor(
    private cdr: ChangeDetectorRef,
    private shippingService: ShippingService,
    private authService: AuthService
  ) {
    
  }
  ngOnChanges(changes: SimpleChanges): void {
    this.getRoleUser();
  }
  
  getRoleUser(){
    this.authService.getUserRole();
    if(this.authService.getUserRole() === 'admin'){
      this.showActions = true;
    }
  }
  updateSatus(status: string){
    this.shippingService.updateStatusShipping(this.shipping.id!, status).subscribe({
      next: (res) => {
        this.shipping.status = status;
        this.cdr.detectChanges();
      },
      error: (error) => {
        console.error('❌ Error actualizando el estado del envío:', error);
      }
    });
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
  /*
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
}
