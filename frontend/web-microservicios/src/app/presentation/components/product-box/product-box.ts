import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Product } from '../../../Models/Product';
import { Console } from 'console';

@Component({
  selector: 'app-product-box',
  imports: [CommonModule],
  standalone: true,
  templateUrl: './product-box.html',
  styleUrl: './product-box.css',
})
export class ProductBox {
  @Input() product!: Product;
  @Output() productClick = new EventEmitter<Product>();
  constructor() {
    console.log("ProductBox creado para producto:", this.product);  
    
  }
  onProductClick(): void {
    console.log('🎯 Click en ProductBox para producto:', this.product.id, this.product.name);
    this.productClick.emit(this.product);
  }
}
