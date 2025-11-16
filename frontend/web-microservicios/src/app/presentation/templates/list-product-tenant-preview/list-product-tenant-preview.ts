import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { ProductBox } from '../../components/product-box/product-box';
import { Product } from '../../../Models/Product';
import {  Router } from '@angular/router';
import { ProductService } from '../../../service/ProductService';


@Component({
  selector: 'app-list-product-tenant-preview',
  imports: [CommonModule, ProductBox],
  templateUrl: './list-product-tenant-preview.html',
  styleUrl: './list-product-tenant-preview.css',
})
export class ListProductTenantPreview {
  @Input() public category: string = '';
  @Input() public products: Product[] = [
  ];
  public tenantid: string = '';
  @Output() productClick = new EventEmitter<Product>();
  constructor(private router: Router) {}
  /**
   * Maneja el click en un producto para navegar a sus detalles
   */
  onProductClick(product: Product): void {
    console.log('Producto clickeado:', product);
    // Emitir el evento para el componente padre
    this.productClick.emit(product);
    // Navegar a la página de detalles del producto
  }

}
