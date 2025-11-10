import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Tenant } from '../../../Models/Tenant';
import { ProductBox } from '../../components/product-box/product-box';
import { Product } from '../../../Models/Product';
import { Route, Router } from '@angular/router';
import { ProductService } from '../../../service /ProductService';
import { TenantService } from '../../../service /TenantService';

@Component({
  selector: 'app-list-product-tenant-preview',
  imports: [CommonModule, ProductBox],
  templateUrl: './list-product-tenant-preview.html',
  styleUrl: './list-product-tenant-preview.css',
})
export class ListProductTenantPreview {
  public products: Product[] = [];
  public tenantid: string = '';
  @Output() productClick = new EventEmitter<Product>();
  constructor(private router: Router,
     private productService: ProductService,
    ) {}

  loadProductsForTenant(): void {
    this.productService.getProducts(1,10).subscribe((products: Product[]) => {
      this.products = products;
    });
  }
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
