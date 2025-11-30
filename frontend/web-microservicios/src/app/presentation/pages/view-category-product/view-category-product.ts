import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, Input } from '@angular/core';
import { ListProductTenantPreview } from '../../templates/list-product-tenant-preview/list-product-tenant-preview';
import { IsLoading } from '../../components/is-loading/is-loading';
import { ProductService } from '../../../service/ProductService';
import { Product } from '../../../Models/Product';
import { LoadingService } from '../../../service/loading-service';
import { Header } from '../../templates/header/header';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-view-category-product',
  imports: [CommonModule, ListProductTenantPreview, Header, IsLoading],
  templateUrl: './view-category-product.html',
  styleUrl: './view-category-product.css',
})
export class ViewCategoryProduct {
   categoryName: string = '';
  products: Product[] = [];
  constructor( private productService: ProductService,
    private isLoading: LoadingService,
    private router: ActivatedRoute,
    private cdr: ChangeDetectorRef
   ) {
    this.router.params.subscribe(params => {
      this.categoryName = params['categoryName'];
      this.cdr.detectChanges();
    });
    this.loadProductsByCategory();
  }  
  loadProductsByCategory() {
    this.isLoading.show("Cargando productos...");
    this.productService.getProducts(1,100).subscribe((data: Product[]) => {
      this.products = data.filter(product => product.category.toLowerCase() === this.categoryName.toLowerCase());
      console.log(this.products);
      this.cdr.detectChanges();
      this.isLoading.hide();
    });
  }
}
