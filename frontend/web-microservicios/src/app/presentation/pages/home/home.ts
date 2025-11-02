import { CommonModule } from '@angular/common';
import { Component, OnInit, ChangeDetectorRef, OnDestroy } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { Header } from '../../templates/header/header';
import { ProductBox } from '../../components/product-box/product-box';
import { ProductService } from '../../../service /ProductService';
import { finalize, Subscription } from 'rxjs';
import { Product } from '../../../Models/Product';

@Component({
  selector: 'app-home',
  imports: [CommonModule, Header, ProductBox],
  templateUrl: './home.html',
  styleUrl: './home.css',
})
export class Home implements OnInit, OnDestroy {
  allProducts: Product[] = [];
  filteredProducts: Product[] = [];
  isLoading: boolean = true;
  searchTerm: string = '';
  errorMessage: string = '';
  
  private subscriptions = new Subscription();

  constructor(
    private productService: ProductService,
    private cdr: ChangeDetectorRef,
    private router: Router,
    private route: ActivatedRoute
  ) {}

  ngOnInit(): void {
    this.loadProducts();
    this.setupSearchListener();
  }

  ngOnDestroy(): void {
    this.subscriptions.unsubscribe();
  }

  /**
   * Configure el listener para escuchar los parámetros de búsqueda de la URL
   */
  private setupSearchListener(): void {
    const queryParamsSub = this.route.queryParams.subscribe(params => {
      const searchParam = params['search'];
      if (searchParam !== this.searchTerm) {
        this.searchTerm = searchParam || '';
        this.filterProducts();
      }
    });
    this.subscriptions.add(queryParamsSub);
  }

  /**
   * Carga los productos desde el servicio
   */
  loadProducts(): void {
    this.isLoading = true;
    this.errorMessage = '';
    
    this.productService.getProducts()
      .pipe(finalize(() => {
        this.isLoading = false;
        this.cdr.markForCheck();
      }))
      .subscribe({
        next: (products) => {
          this.allProducts = products || [];
          this.filterProducts(); // Aplicar filtros después de cargar
          console.log('Productos cargados:', this.allProducts.length);
        },
        error: (err) => {
          console.error('Error al cargar productos:', err);
          this.errorMessage = 'Error al cargar los productos. Intenta de nuevo.';
          this.allProducts = [];
          this.filteredProducts = [];
        }
      });
  }

  /**
   * Filtra los productos basado en el término de búsqueda
   */
  private filterProducts(): void {
    if (!this.searchTerm || this.searchTerm.trim() === '') {
      this.filteredProducts = [...this.allProducts];
    } else {
      const searchLower = this.searchTerm.toLowerCase().trim();
      this.filteredProducts = this.allProducts.filter(product => 
        product.description?.toLowerCase().includes(searchLower) ||
        product.category?.toLowerCase().includes(searchLower) ||
        product.unit?.toLowerCase().includes(searchLower)
      );
    }
    
    console.log(`Productos filtrados: ${this.filteredProducts.length} de ${this.allProducts.length}`);
  }

  /**
   * Maneja el click en un producto para navegar a sus detalles
   */
  onProductClick(product: Product): void {
    console.log('Producto clickeado:', product);
    this.router.navigate(['/product', product.id]);
  }

  /**
   * Filtra productos por categoría específica
   */
  filterByCategory(category: string): void {
    if (!category || category.trim() === '') {
      this.filteredProducts = [...this.allProducts];
    } else {
      this.filteredProducts = this.allProducts.filter(product => 
        product.category?.toLowerCase() === category.toLowerCase()
      );
    }
  }

  /**
   * Limpia todos los filtros y muestra todos los productos
   */
  clearFilters(): void {
    this.searchTerm = '';
    this.filteredProducts = [...this.allProducts];
    // Limpiar query params
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: {},
      replaceUrl: true
    });
  }

  /**
   * Retorna el array de productos a mostrar en la vista
   */
  get productsToShow(): Product[] {
    return this.filteredProducts;
  }

  /**
   * Retorna si hay productos para mostrar
   */
  get hasProducts(): boolean {
    return this.productsToShow.length > 0;
  }

  /**
   * Retorna si se está mostrando una búsqueda filtrada
   */
  get isSearchActive(): boolean {
    return this.searchTerm.trim() !== '';
  }
}
