import { CommonModule } from '@angular/common';
import { Component, OnInit, OnDestroy, ChangeDetectorRef, NgZone } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { finalize } from 'rxjs/operators';

import { Header } from '../../templates/header/header';
import { ProductService } from '../../../service /ProductService';
import { Product } from '../../../Models/Product';
import { ListProductTenantPreview } from '../../templates/list-product-tenant-preview/list-product-tenant-preview';

@Component({
  selector: 'app-home',
  imports: [CommonModule, Header, ListProductTenantPreview],
  templateUrl: './home.html',
  styleUrl: './home.css',
})
export class Home implements OnInit, OnDestroy {
  // ==================== PROPIEDADES ====================
  products: Product[] = [];
  filteredProducts: Product[] = [];
  
  // Estados de la UI
  isLoading: boolean = true;
  searchTerm: string = '';
  errorMessage: string = '';
  
  private subscriptions = new Subscription();

  // ==================== CONSTRUCTOR ====================
  constructor(
    private productService: ProductService,
    private router: Router,
    private route: ActivatedRoute,
    private cdr: ChangeDetectorRef, // ✅ Agregar ChangeDetectorRef
    private zone: NgZone // ✅ Agregar NgZone
  ) {}

  // ==================== CICLO DE VIDA ====================
  ngOnInit(): void {
    this.initializeComponent();
  }

  ngOnDestroy(): void {
    this.subscriptions.unsubscribe();
  }

  // ==================== INICIALIZACIÓN ====================
  private initializeComponent(): void {
    this.setupSearchListener();
    this.loadProducts();
  }

  /**
   * Configura el listener para escuchar los parámetros de búsqueda de la URL
   */
  private setupSearchListener(): void {
    const queryParamsSub = this.route.queryParams.subscribe(params => {
      const searchParam = params['search'];
      if (searchParam !== this.searchTerm) {
        this.searchTerm = searchParam || '';
        this.filterProducts();
        this.cdr.detectChanges(); // ✅ Forzar detección de cambios
      }
    });
    this.subscriptions.add(queryParamsSub);
  }

  // ==================== CARGA DE DATOS ====================
  /**
   * Carga los productos del tenant actual
   */
  loadProducts(): void {
    console.log('🚀 [HOME] Iniciando carga - isLoading:', this.isLoading);
    
    this.zone.run(() => {
      this.isLoading = true;
      this.errorMessage = '';
      this.cdr.detectChanges(); // ✅ Forzar detección inmediata
    });

    const loadSub = this.productService.getProducts(1, 50)
      .pipe(
        finalize(() => {
          this.zone.run(() => {
            this.isLoading = false;
            this.cdr.detectChanges();
          });
        })
      )
      .subscribe({
        next: (products) => {
          this.zone.run(() => {
            this.handleProductsLoaded(products);
          });
        },
        error: (error) => {
          this.zone.run(() => {
            this.handleLoadError(error);
          });
        }
      });

    this.subscriptions.add(loadSub);
  }

  /**
   * Maneja la carga exitosa de productos
   */
  private handleProductsLoaded(products: Product[]): void {
    this.products = products;
    this.filterProducts();
    this.cdr.detectChanges(); // ✅ Forzar detección después de actualizar datos
    
    // Debug: verificar estado después de un momento
    setTimeout(() => {
      console.log('⏰ [HOME] Estado después de procesar:', {
        isLoading: this.isLoading,
        productsLength: this.products.length,
        hasProducts: this.hasProducts
      });
    }, 100);
  }

  /**
   * Maneja errores en la carga
   */
  private handleLoadError(error: any): void {
    this.errorMessage = 'Error al cargar los productos. Intenta de nuevo.';
    this.products = [];
    this.filteredProducts = [];
    this.cdr.detectChanges(); 
  }

  // ==================== FILTROS Y BÚSQUEDA ====================
  /**
   * Filtra productos basado en el término de búsqueda
   */
  private filterProducts(): void {
    if (!this.searchTerm.trim()) {
      this.filteredProducts = [...this.products];
    } else {
      const term = this.searchTerm.toLowerCase().trim();
      this.filteredProducts = this.products.filter(product => 
        product.description?.toLowerCase().includes(term) ||
        product.category?.toLowerCase().includes(term)
      );
    }
    this.cdr.detectChanges(); // ✅ Forzar detección después de filtrar
  }

  /**
   * Realiza una búsqueda
   */
  performSearch(searchTerm: string): void {
    this.searchTerm = searchTerm;
    this.filterProducts();
    this.updateUrlWithSearch();
  }

  /**
   * Limpia la búsqueda
   */
  clearSearch(): void {
    this.searchTerm = '';
    this.filterProducts();
    this.clearUrlParams();
  }

  // ==================== GETTERS ====================
  /**
   * Productos a mostrar en la vista
   */
  get productsToShow(): Product[] {
    const products = this.isSearchActive ? this.filteredProducts : this.products;
    return products;
  }

  /**
   * Si hay productos para mostrar
   */
  get hasProducts(): boolean {
    const hasProducts = this.productsToShow.length > 0;
    return hasProducts;
  }

  /**
   * Si se está mostrando una búsqueda filtrada
   */
  get isSearchActive(): boolean {
    return this.searchTerm.trim() !== '';
  }

  /**
   * Mensaje a mostrar cuando no hay productos
   */
  get noProductsMessage(): string {
    if (this.isLoading) return '';
    if (this.errorMessage) return this.errorMessage;
    if (this.isSearchActive) return `No se encontraron productos para "${this.searchTerm}"`;
    return 'No hay productos disponibles en tu tienda';
  }

  // ==================== NAVEGACIÓN ====================
  /**
   * Maneja el click en un producto
   */
  onProductClick(product: Product): void {
    this.router.navigate(['/product', product.id]);
  }

  /**
   * Recarga los productos
   */
  refreshProducts(): void {
    this.loadProducts();
  }

  // ==================== UTILIDADES PRIVADAS ====================
  /**
   * Actualiza la URL con el término de búsqueda
   */
  private updateUrlWithSearch(): void {
    const queryParams = this.searchTerm ? { search: this.searchTerm } : {};
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams,
      replaceUrl: true
    });
  }

  /**
   * Limpia los parámetros de la URL
   */
  private clearUrlParams(): void {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: {},
      replaceUrl: true
    });
  }

  // ==================== DEBUG ====================
  /**
   * Estado actual del componente para debugging
   */
  get debugState(): any {
    return {
      isLoading: this.isLoading,
      productsLength: this.products.length,
      filteredProductsLength: this.filteredProducts.length,
      hasProducts: this.hasProducts,
      searchTerm: this.searchTerm,
      errorMessage: this.errorMessage
    };
  }

  /**
   * Método para debugging manual
   */
  logState(): void {
    console.log('🐛 [HOME] Estado actual:', this.debugState);
  }
}