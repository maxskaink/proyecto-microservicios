import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component } from '@angular/core';
import { FormProduct } from '../../templates/from-product/form-product';
import { ActivatedRoute, Router } from '@angular/router';
import { ProductService } from '../../../service/ProductService';
import Swal from 'sweetalert2';
import { Header } from '../../templates/header/header';
import { ProductPeticion } from '../../../Models/PrdocutPeticion';

@Component({
  selector: 'app-edit-item',
  imports: [CommonModule, FormProduct, Header],
  templateUrl: './edit-item.html',
  styleUrl: './edit-item.css',
})
export class EditItem  {
  public idProduct: string = '';
  public categories: string[] = [];
  public productData!: ProductPeticion;
  public isLoading: boolean = false;
  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private cdr: ChangeDetectorRef,
    private serviceProduct: ProductService
  ) { 
  
  }
  ngOnInit() {
    this.route.paramMap.subscribe(params => {
      this.idProduct = params.get('id') || '';
      console.log('ID recibido por paramMap:', this.idProduct);
    });
    this.loadCategories();
    this.loadProduct();

  }
  loadCategories() {
    this.serviceProduct.getCategories().subscribe({
      next: (cats) => {
        console.log('Categorías cargadas:', cats);
        this.categories = cats;
        this.cdr.detectChanges();
      },
      error: (error) => {
        console.error('Error al cargar las categorías:', error);
      }
    });
  }
  loadProduct() {
    this.isLoading = true;
    if (!this.idProduct) {
      console.warn('No se tiene un ID de producto válido');
      return;
    }
    this.serviceProduct.getProductById(this.idProduct).subscribe({
      next: (product) => {
        console.log('Producto cargado para edición:', product);
        this.productData = product
        this.isLoading = false;
        this.cdr.detectChanges();
      },
      error: (error) => {
        console.error('Error al cargar el producto para edición:', error);
        Swal.fire({
              icon: 'error',
              title: 'Error',
              text: 'Error al traer el producto.',
              buttonsStyling: false,
              customClass: {
                confirmButton: 'btn btn-danger'
              }
            });
        this.isLoading = false;
        this.cdr.detectChanges();
      }
    });
  }

  /**
   * Maneja el evento de creación/edición del producto
   * @param event Datos del producto creado/ editado
   */
  handleProduct(event: { product: any; action: string; selectFIle: File }) {
    console.log('Producto editado recibido:', event);
    this.productData = event.product;
  }
  /**
   * Actualiza el producto existente al back
   */
  updateProduct() {
    if (!this.productData) {
      console.error('No hay datos de producto para actualizar.');
      return;
    }
    this.serviceProduct.updateProduct(this.idProduct, this.productData).subscribe({
      next: (response) => {
        console.log('Producto actualizado con éxito:', response);
        Swal.fire({
              icon: 'success',
              title: 'Éxito',
              text: 'Producto actualizado con éxito.',
              buttonsStyling: false,
              customClass: {
                confirmButton: 'btn btn-success'
              }
            });
        this.router.navigate(['/my-products']);
      },
      error: (error) => {
        console.error('Error al actualizar el producto:', error);
        Swal.fire({
              icon: 'error',
              title: 'Error',
              text: 'Error al actualizar el producto. ${error.message}',
              buttonsStyling: false,
              customClass: {
                confirmButton: 'btn btn-danger'
              }
            });
      }
    });
  }
}
