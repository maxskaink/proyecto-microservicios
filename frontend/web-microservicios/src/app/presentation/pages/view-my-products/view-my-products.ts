import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { Header } from '../../templates/header/header';
import { Product } from '../../../Models/Product';
import { ProductService } from '../../../service/ProductService';
import { AuthService } from '../../../service/Authser.vice';
import { ArrowLeft } from '../../components/arrow-left/arrow-left';
import { ListProductTenantPreview } from '../../templates/list-product-tenant-preview/list-product-tenant-preview';
import { Router, RouterLink } from '@angular/router';
import Swal from 'sweetalert2';
import { FormsModule } from '@angular/forms';
import { IsLoading } from '../../components/is-loading/is-loading';
import { LoadingService } from '../../../service/loading-service';

@Component({
  selector: 'app-view-my-products',
  imports: [CommonModule, Header, ListProductTenantPreview, FormsModule, RouterLink, IsLoading],
  templateUrl: './view-my-products.html',
  styleUrl: './view-my-products.css',
})
export class ViewMyProducts implements OnInit {
  public allProducts: Product[] = [];
  public filteredProducts: Product[] = [];
  public isLoading: boolean = false;
  public searchName: string = '';
  public selectedCategory: string = 'todos';
  public categories: string[] = [];
  
  constructor(
    private serviceProduct: ProductService,
    private authService: AuthService,
    private cdr: ChangeDetectorRef,
    private router: Router,
    private loadingService: LoadingService
  ) { }
  ngOnInit(): void {
      this.loadMyProducts();
  }
  
loadMyProducts() {
  const userId = localStorage.getItem('user_data') ? JSON.parse(localStorage.getItem('user_data')!).id : null;

  this.loadingService.show("Cargando productos...");
  this.isLoading = true;  // <<< ACTIVAS EL LOADING AQUÍ

  if (!userId && this.allProducts.length === 0) {
    console.warn("UserData aún no está listo. Reintentando...");
    setTimeout(() => this.loadMyProducts(), 150);
    return;
  }

  this.serviceProduct.getProducts(1, 100).subscribe({
    next: (products) => {
      this.allProducts = products.filter(p => String(p.producer_id) === String(userId));
      this.filteredProducts = [...this.allProducts];
      this.categories = [...new Set(this.allProducts.map(p => p.category))].sort();
      
      this.isLoading = false;        
      this.loadingService.hide();
      this.cdr.detectChanges();
    },
    error: (error) => {
      console.error('Error al cargar productos:', error);
      this.isLoading = false;        // <<< TAMBIÉN LO DETIENE EN ERROR
      this.loadingService.hide();
      this.cdr.detectChanges();
    }
  });
}

/**
 * Filtra productos por nombre
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
 * Aplica todos los filtros
 */
applyFilters() {
  this.filteredProducts = this.allProducts.filter(product => {
    const nameMatch = product.name.toLowerCase().includes(this.searchName);
    const categoryMatch = this.selectedCategory === 'todos' || product.category === this.selectedCategory;
    return nameMatch && categoryMatch;
  });
  this.cdr.detectChanges();
}

/**
 * Limpia los filtros
 */
clearFilters() {
  this.searchName = '';
  this.selectedCategory = 'todos';
  this.filteredProducts = [...this.allProducts];
  this.cdr.detectChanges();
}

/**
 * si la persona hace click en un producto, lo redirige a la pagina de detalle del producto
 * @param product Producto seleccionado
 */
  onProductClick( product: Product ) {
    console.log('Producto seleccionado en view-my-products:', product);
    this.router.navigate(['/edit-product', product.id]);
  }
 onAction(action: {state:string, id:string}) {
  console.log('Acción recibida en view-my-products:', action.state, action.id);

  if (action.state === 'delete') {
    Swal.fire({
      title: '¿Estás seguro?',
      text: "¡No podrás revertir esto!",
      icon: 'warning',
      showCancelButton: true,
      buttonsStyling: false,
      customClass: {
        confirmButton: 'btn btn-danger mx-2',
        cancelButton: 'btn btn-secondary mx-2'
      },
      confirmButtonText: 'Sí, eliminarlo!'
    }).then((result) => {
      if (result.isConfirmed) {
        // Llamamos al servicio de eliminación
        this.deleteProduct(action.id).subscribe({
          next: () => {
            Swal.fire(
              '¡Eliminado!',
              'Tu producto ha sido eliminado.',
              'success'
            ).then(() => {
              // Redirige después de confirmar el Swal de éxito
              this.router.navigate(['/admin-panel']);
            });
          },
          error: (err) => {
            Swal.fire(
              'Error',
              'No se pudo eliminar el producto. Intenta de nuevo.',
              'error'
            );
          }
        });
      }
    });
  } else if (action.state === 'edit') {
    this.router.navigate(['/edit-product', action.id]);
  }
}

// Ahora devuelve el observable sin subscribirse
deleteProduct(productId: string) {
  return this.serviceProduct.deleteProduct(productId);
}

  goToPublush() {
    this.router.navigate(['/publishProduct']);
  }
  

}
