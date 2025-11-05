import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Tenant } from '../../../Models/Tenant';
import { ProductBox } from '../../components/product-box/product-box';
import { Product } from '../../../Models/Product';
import { Route, Router } from '@angular/router';

@Component({
  selector: 'app-list-product-tenant-preview',
  imports: [CommonModule, ProductBox],
  templateUrl: './list-product-tenant-preview.html',
  styleUrl: './list-product-tenant-preview.css',
})
export class ListProductTenantPreview {
  @Input() tenantid: string = "";
  @Input() products: Product[] = [];
  @Input() tenantName: string = "";
  @Output() productClick = new EventEmitter<Product>();

  constructor(private router: Router) {}

  /**
   * Retorna los productos a mostrar (todos los productos del tenant)
   */
  get productsToShow(): Product[] {
    return this.products || [];
  }

  /**
   * Retorna si hay productos para mostrar
   */
  get hasProducts(): boolean {
    return this.productsToShow.length > 0;
  }

  /**
   * Maneja el click en un producto para navegar a sus detalles
   */
  onProductClick(product: Product): void {
    console.log('Producto clickeado:', product);
    // Emitir el evento para el componente padre
    this.productClick.emit(product);
    // Navegar a la página de detalles del producto
    this.router.navigate(['/product', this.tenantid, product.id]);
  }

}
