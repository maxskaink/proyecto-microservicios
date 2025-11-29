import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { Header } from '../../templates/header/header';
import { Product } from '../../../Models/Product';
import { ProductService } from '../../../service/ProductService';
import { AuthService } from '../../../service/Authser.vice';
import { ArrowLeft } from '../../components/arrow-left/arrow-left';
import { ListProductTenantPreview } from '../../templates/list-product-tenant-preview/list-product-tenant-preview';
import { Router } from '@angular/router';
import Swal from 'sweetalert2';

@Component({
  selector: 'app-view-my-products',
  imports: [CommonModule, Header, ArrowLeft, ListProductTenantPreview],
  templateUrl: './view-my-products.html',
  styleUrl: './view-my-products.css',
})
export class ViewMyProducts implements OnInit {
  public allProducts: Product[] = [];
  public isLoading: boolean = false;
  
  constructor(
    private serviceProduct: ProductService,
    private authService: AuthService,
    private cdr: ChangeDetectorRef,
    private router: Router
  ) { }
  ngOnInit(): void {
    this.authService.userData.subscribe(userData => {
      if (userData) {
        this.loadMyProducts();
      }
    });
  }

loadMyProducts() {
  const userId = this.authService.userCurrentData?.id;

  if (!userId) {
    console.warn("UserData aún no está listo. Reintentando...");
    setTimeout(() => this.loadMyProducts(), 150);
    return;
  }
  this.serviceProduct.getProducts(1, 100).subscribe({
    next: (products) => {
      console.log('Todos los productos:', products);
      this.allProducts = products.filter(p => String(p.producer_id) === String(userId));
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
/**
 * si la persona hace click en un producto, lo redirige a la pagina de detalle del producto
 * @param product Producto seleccionado
 */
  onProductClick( product: Product ) {
    console.log('Producto seleccionado en view-my-products:', product);
    this.router.navigate(['/edit-pruduct', product.id]);
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
          confirmButton: ' btn btn-danger mx-2',
          cancelButton: ' btn btn-secondary mx-2'
        },
        confirmButtonText: 'Sí, eliminarlo!'
      }).then((result) => {
        if (result.isConfirmed) {
          this.deleteProduct(action.id);
          Swal.fire(
            '¡Eliminado!',
            'Tu producto ha sido eliminado.',
            'success'
          );
        }
      });
    }else if (action.state === 'edit') {
      this.router.navigate(['/edit-product', action.id]);
    }
  }

  deleteProduct(productId: string) {
    this.serviceProduct.deleteProduct(productId).subscribe({
      next: () => {
        console.log('Producto eliminado con éxito:', productId);
        this.loadMyProducts();
      },
      error: (error) => {
        console.error('Error al eliminar el producto:', error);
      }
    });
  }


}
