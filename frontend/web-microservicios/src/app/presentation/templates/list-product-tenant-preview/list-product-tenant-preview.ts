import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { ProductBox } from '../../components/product-box/product-box';
import { Product } from '../../../Models/Product';
import {  Router } from '@angular/router';
import { ProductService } from '../../../service/ProductService';
import { TenantService } from '../../../service/TenantService';
import { Tenant } from '../../../Models/Tenant';


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
  constructor(private router: Router, private  serviceTenant: TenantService) {}
  /**
   * Maneja el click en un producto para navegar a sus detalles
   */
  onProductClick(product: Product): void {
    this.productClick.emit(product);

    this.serviceTenant.getCurrentUserTenant().subscribe((tenant: Tenant | null) => {

      if (!tenant) {
        console.error('No se encontró el tenant actual');
        return;
      }

      this.router.navigate([
        'product',
        tenant.tenant_id,  // <-- tenant_id seguro aquí
        product.id
      ]);
    });
  }
}
