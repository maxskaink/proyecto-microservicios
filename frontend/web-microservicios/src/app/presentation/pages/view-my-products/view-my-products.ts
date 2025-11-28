import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { Header } from '../../templates/header/header';
import { Product } from '../../../Models/Product';
import { ProductService } from '../../../service/ProductService';
import { AuthService } from '../../../service/Authser.vice';
import { ArrowLeft } from '../../components/arrow-left/arrow-left';
import { ListProductTenantPreview } from '../../templates/list-product-tenant-preview/list-product-tenant-preview';

@Component({
  selector: 'app-view-my-products',
  imports: [CommonModule, Header, ArrowLeft, ListProductTenantPreview],
  templateUrl: './view-my-products.html',
  styleUrl: './view-my-products.css',
})
export class ViewMyProducts implements OnInit {
  public allProducts: Product[] = [];
  public isLoading: boolean = true;
  
  constructor(
    private serviceProduct: ProductService,
    private authService: AuthService,
    private cdr: ChangeDetectorRef,
  ) { }
  ngOnInit(): void {
    this.loadMyProducts();
  }
loadMyProducts() {
  const userId = this.authService.userCurrentData?.id;
  if (!userId) {
    console.error('No hay usuario actual');
    return;
  }

  this.serviceProduct.getProducts(1, 100).subscribe({
    next: (products) => {
      console.log('Todos los productos:', products);
      this.allProducts = products.filter(p => p.producer_id === userId);
      console.log('Productos del usuario:', this.allProducts);
      this.isLoading = false;
      this.cdr.detectChanges();
    },
    error: (error) => {
      console.error('Error al cargar los productos del usuario:', error);
      this.isLoading = false;
      this.cdr.detectChanges();
    }
  });
}
onProductClick( product: Product ) {
  // Lógica para manejar el clic en un producto
}

}
