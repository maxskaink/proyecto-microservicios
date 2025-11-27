import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, OnChanges, Output, SimpleChanges } from '@angular/core';
import { CartItem } from '../../../Models/Cart';

@Component({
  selector: 'app-product-view-buy',
  imports: [CommonModule],
  templateUrl: './product-view-buy.html',
  styleUrl: './product-view-buy.css',
})
export class ProductViewBuy implements OnChanges{
  ngOnChanges(changes: SimpleChanges): void {
    if (changes['item']) {
    console.log('Item recibido:', this.item);
    console.log('URL de la imagen:', this.item.product.photo_url);
  }
  }

  // Lista de items recibidos desde el padre
  @Input() item!: CartItem;

  // Eventos que el padre escuchará
  @Output() increment = new EventEmitter<any>();
  @Output() decrement = new EventEmitter<any>();
  @Output() remove = new EventEmitter<any>();

  // Solo formatea, no modifica nada
  formatPrice(price: number): string {
    return price.toLocaleString('es-CO', {
      style: 'currency',
      currency: 'COP',
      minimumFractionDigits: 0
    });
  }

  // ✅ Emiten eventos con logs de debug
  onIncrement(item: any) {
    console.log('🔼 Incrementando cantidad para item:', item.id, 'cantidad actual:', item.quantity);
    this.increment.emit(item);
  }

  onDecrement(item: any) {
    console.log('🔽 Decrementando cantidad para item:', item.id, 'cantidad actual:', item.quantity);
    this.decrement.emit(item);
  }

  onRemove(item: any) {
    console.log('🗑️ Eliminando item:', item.id);
    this.remove.emit(item);
  }
}