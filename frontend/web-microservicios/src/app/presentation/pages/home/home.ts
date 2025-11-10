import { CommonModule } from '@angular/common';
import { Component} from '@angular/core';
import { Header } from '../../templates/header/header';
import { Product } from '../../../Models/Product';
import { ListProductTenantPreview } from '../../templates/list-product-tenant-preview/list-product-tenant-preview';


@Component({
  selector: 'app-home',
  imports: [CommonModule, Header, ListProductTenantPreview],
  templateUrl: './home.html',
  styleUrl: './home.css',
})
export class Home  {
  
  onProductClick(product: any) {
    console.log('Producto clickeado:', product);
    // Aquí puedes navegar, abrir un modal, etc.
  }
}