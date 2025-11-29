import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { ProductBox } from '../../components/product-box/product-box';
import { Product } from '../../../Models/Product';


@Component({
  selector: 'app-list-product-tenant-preview',
  imports: [CommonModule, ProductBox],
  templateUrl: './list-product-tenant-preview.html',
  styleUrl: './list-product-tenant-preview.css',
})
export class ListProductTenantPreview {
  @Input() public products: Product[] = [];
  @Input() showActions: boolean = false;
  @Output() productClick = new EventEmitter<Product>();
  @Output() action = new EventEmitter<{state:string, id:string}>();
  
  constructor() {}
  
  /**
   * Solo emite el evento al componente padre
   */
  onProductClick(product: Product): void {
    console.log('📤 Emitiendo evento desde list-product-tenant-preview:', product.id);
    this.productClick.emit(product);
  }
  onAction(action: {state:string, id:string}): void {
    console.log('📤 Emitiendo acción desde list-product-tenant-preview:', action.state, action.id);
    this.action.emit(action);
  }
}
