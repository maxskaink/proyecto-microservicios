import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { ListProductTenantPreview } from '../../templates/list-product-tenant-preview/list-product-tenant-preview';
import { Header } from '../../templates/header/header';
import { IsLoading } from '../../components/is-loading/is-loading';
import { LoadingService } from '../../../service/loading-service';
import { ProductService } from '../../../service/ProductService';
import { ActivatedRoute } from '@angular/router';
import { Product } from '../../../Models/Product';

@Component({
  selector: 'app-search-products',
  imports: [CommonModule, ListProductTenantPreview, Header, IsLoading],
  templateUrl: './search-products.html',
  styleUrl: './search-products.css',
})
export class SearchProducts implements OnInit{
  public searchTerm: string = '';
  public products: Product[] = [];
  constructor(
    private isLoading: LoadingService,
    private productService: ProductService,
    private cdr: ChangeDetectorRef,
    private route: ActivatedRoute
  ) { }
  ngOnInit(): void {
    this.loadSearchTerm();
    this.loadProductsBySearchTerm(this.searchTerm);
  }
  loadSearchTerm() {
    this.route.params.subscribe(params => {
      const term = params['searchTerm'];
      this.searchTerm = term;
    });
  }
  ngAfterViewInit(): void {
    this.loadSearchTerm();
  }

  loadProductsBySearchTerm(term: string) {
    this.isLoading.show("Buscando productos...");
    this.productService.getProducts(1, 100).subscribe({
      next: (products) => {
        this.searchTerm = term;
        // Aquí filtras los productos según el término de búsqueda
        const filteredProducts = products.filter(p => 
          p.name.toLowerCase().includes(term.toLowerCase()) ||
          p.description.toLowerCase().includes(term.toLowerCase())
        );
        this.products = filteredProducts;
        this.isLoading.hide();
        this.cdr.detectChanges();
      },
      error: (error) => {
        console.error("Error al buscar productos:", error);
        this.isLoading.hide();
      }
    });
  }
  

}
