import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Product } from '../../../Models/Product';

@Component({
  selector: 'app-product-box',
  imports: [CommonModule],
  templateUrl: './product-box.html',
  styleUrl: './product-box.css',
})
export class ProductBox {
  @Input() product!: Product;
  @Output() productClick = new EventEmitter<Product>();

  onProductClick(): void {
    console.log('🎯 Click en ProductBox para producto:', this.product.id, this.product.name);
    this.productClick.emit(this.product);
  }
}
