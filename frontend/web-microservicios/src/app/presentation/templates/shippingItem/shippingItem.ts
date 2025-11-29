import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, Input } from '@angular/core';
import { ShoppingCartService } from '../../../service/ShoppinCartService';
import { Shipping } from '../../../Models/OrderPeticion';

@Component({
  selector: 'app-shipping',
  imports: [CommonModule],
  templateUrl: './shippingItem.html',
  styleUrl: './shippingItem.css',
})
export class shippingItem {
  @Input() showActions: boolean = true;
  @Input() shipping!: Shipping;
  constructor(
    private cdr: ChangeDetectorRef,
    
  ) {

  }
  updateSatus(status: string){

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
