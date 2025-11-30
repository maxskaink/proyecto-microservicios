import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, Input } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ListProductTenantPreview } from '../../templates/list-product-tenant-preview/list-product-tenant-preview';
import { IsLoading } from '../../components/is-loading/is-loading';
import { ProductService } from '../../../service/ProductService';
import { Product } from '../../../Models/Product';
import { LoadingService } from '../../../service/loading-service';
import { Header } from '../../templates/header/header';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-view-category-product',
  imports: [CommonModule, ListProductTenantPreview, Header, IsLoading, FormsModule],
  templateUrl: './view-category-product.html',
  styleUrl: './view-category-product.css',
})
export class ViewCategoryProduct {
  categoryName: string = '';
  products: Product[] = [];
  filteredProducts: Product[] = [];
  searchName: string = '';
  maxPrice: number = 0;
  allProducts: Product[] = [];

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
      this.allProducts = data.filter(product => product.category.toLowerCase() === this.categoryName.toLowerCase());
      this.products = [...this.allProducts];
      this.filteredProducts = [...this.allProducts];
      this.maxPrice = Math.max(...this.allProducts.map(p => p.price || 0), 0) || 1000;
      console.log(this.products);
      this.cdr.detectChanges();
      this.isLoading.hide();
    });
  }

  onSearchChange(searchValue: string) {
    this.searchName = searchValue.toLowerCase();
    this.applyFilters();
  }

  onPriceChange(priceValue: number) {
    this.maxPrice = priceValue;
    this.applyFilters();
  }

  applyFilters() {
    this.filteredProducts = this.allProducts.filter(product => {
      const nameMatch = product.name.toLowerCase().includes(this.searchName);
      const priceMatch = (product.price || 0) <= this.maxPrice;
      return nameMatch && priceMatch;
    });
    this.products = this.filteredProducts;
    this.cdr.detectChanges();
  }

  clearFilters() {
    this.searchName = '';
    this.maxPrice = Math.max(...this.allProducts.map(p => p.price || 0), 0) || 1000;
    this.products = [...this.allProducts];
    this.filteredProducts = [...this.allProducts];
    this.cdr.detectChanges();
  }
  get maxProductPrice(): number {
  if (!this.allProducts || this.allProducts.length === 0) {
    return 1000;
  }
  return Math.max(...this.allProducts.map(p => p.price ?? 0));
  }
  get shouldShowClear(): boolean {
    return !!this.searchName || this.maxPrice < this.maxProductPrice;
  }

}
