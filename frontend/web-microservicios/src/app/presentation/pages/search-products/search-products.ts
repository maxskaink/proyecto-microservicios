import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { ListProductTenantPreview } from '../../templates/list-product-tenant-preview/list-product-tenant-preview';
import { Header } from '../../templates/header/header';
import { IsLoading } from '../../components/is-loading/is-loading';
import { LoadingService } from '../../../service/loading-service';
import { ProductService } from '../../../service/ProductService';
import { ActivatedRoute, Router } from '@angular/router';
import { Product } from '../../../Models/Product';
import { FormsModule } from '@angular/forms';
import { TenantService } from '../../../service/TenantService';

@Component({
  selector: 'app-search-products',
  imports: [CommonModule, ListProductTenantPreview, Header, IsLoading, FormsModule],
  templateUrl: './search-products.html',
  styleUrl: './search-products.css',
})
export class SearchProducts implements OnInit{
  public searchTerm: string = '';
  public allProducts: Product[] = [];
  public filteredProducts: Product[] = [];
  public categories: string[] = [];
  public selectedCategory: string = 'todos';
  public maxPrice: number = 0;
  public maxProductPrice: number = 0;
  public searchName: string = '';

  constructor(
    private isLoading: LoadingService,
    private productService: ProductService,
    private cdr: ChangeDetectorRef,
    private route: ActivatedRoute,
    private tenantService: TenantService,
    private router: Router,
  ) { }

  ngOnInit(): void {
    this.loadSearchTerm();
    this.loadProductsBySearchTerm(this.searchTerm);
  }

  loadSearchTerm() {
    this.route.params.subscribe(params => {
      const term = params['searchTerm'];
      this.searchTerm = term;
      this.searchName = term;
      this.loadProductsBySearchTerm(term);
    });
  }

  loadProductsBySearchTerm(term: string) {
    this.isLoading.show("Buscando productos...");
    this.productService.getProducts(1, 100).subscribe({
      next: (products) => {
        this.searchTerm = term;
        
        // Filtrar productos que coincidan con el término de búsqueda
        const filteredBySearch = products.filter(p => 
          p.name.toLowerCase().includes(term.toLowerCase()) ||
          p.description.toLowerCase().includes(term.toLowerCase())
        );
        
        this.allProducts = filteredBySearch;
        this.filteredProducts = [...this.allProducts];
        
        // Extraer categorías únicas
        this.categories = [...new Set(this.allProducts.map(p => p.category))].sort();
        
        // Obtener precio máximo
        this.maxProductPrice = Math.max(...this.allProducts.map(p => Number(p.price) || 0), 0);
        this.maxPrice = this.maxProductPrice;
        
        this.isLoading.hide();
        this.cdr.detectChanges();
      },
      error: (error) => {
        console.error("Error al buscar productos:", error);
        this.isLoading.hide();
      }
    });
  }

  /**
   * Filtra productos por nombre en tiempo real
   */
  onSearchChange(searchValue: string) {
    this.searchName = searchValue.toLowerCase();
    this.applyFilters();
  }

  /**
   * Filtra productos por categoría
   */
  onCategoryChange(category: string) {
    this.selectedCategory = category;
    this.applyFilters();
  }

  /**
   * Filtra productos por precio
   */
  onPriceChange(price: number) {
    this.maxPrice = price;
    this.applyFilters();
  }

  /**
   * Aplica todos los filtros
   */
  applyFilters() {
    this.filteredProducts = this.allProducts.filter(product => {
      const nameMatch = product.name.toLowerCase().includes(this.searchName);
      const categoryMatch = this.selectedCategory === 'todos' || product.category === this.selectedCategory;
      const priceMatch = Number(product.price) <= this.maxPrice;
      return nameMatch && categoryMatch && priceMatch;
    });
    this.cdr.detectChanges();
  }

  /**
   * Limpia los filtros
   */
  clearFilters() {
    this.searchName = this.searchTerm;
    this.selectedCategory = 'todos';
    this.maxPrice = this.maxProductPrice;
    this.filteredProducts = [...this.allProducts];
    this.cdr.detectChanges();
  }
    /**
   * Maneja el click en un producto para navegar a su vista de detalles
   */
  onProductClick(product: Product): void {
    this.tenantService.getCurrentUserTenant().subscribe((tenant) => {
      if (!tenant) {
        console.error('No se encontró el tenant actual desde home.');
        return;
      }

      const route = ['product', tenant.tenant_id, product.id];

      this.router.navigate(route).then(
        (success) => console.log(' Navegación exitosa:', success),
        (error) => console.error(' Error en navegación:', error),
      );
    });
  }
}
