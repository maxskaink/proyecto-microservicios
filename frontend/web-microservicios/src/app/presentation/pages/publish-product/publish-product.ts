import { CommonModule, isPlatformBrowser } from '@angular/common';
import { ChangeDetectorRef, Component, ElementRef, Inject, OnInit, PLATFORM_ID, ViewChild } from '@angular/core';
import { FormsModule, NgForm } from '@angular/forms';
import { Header } from '../../templates/header/header';
import { ProductPeticion } from '../../../Models/PrdocutPeticion';
import { ProductService } from '../../../service/ProductService';
import { switchMap, catchError } from 'rxjs/operators';
import { of } from 'rxjs';
import { FormProduct } from '../../templates/from-product/form-product';
import Swal from 'sweetalert2';



@Component({
  selector: 'app-publish-product',
  imports: [CommonModule, FormsModule, Header, FormProduct],
  templateUrl: './publish-product.html',
  styleUrl: './publish-product.css',
  standalone: true,
})// ✅ Verificar que estamos en el navegador
export class PublishProduct implements OnInit {
  product?: ProductPeticion;
    
  isSubmitting: boolean = false;
  imagePreview: string | null = null;
  isImageInvalid: boolean = false;
  imageError: string = '';
  selectedFile: File | null = null;
  @ViewChild('fileInput') fileInputRef?: ElementRef<HTMLInputElement>;
  categories: string[] =[];
  public units: string[] = ['kg', 'atado', 'libra'];

  constructor(
    private productService: ProductService,
    private cdr: ChangeDetectorRef,
    @Inject(PLATFORM_ID) private platformId: Object
  ) { }

  ngOnInit(): void {
    this.loadCategories();
  }

  handleProduct(event: { product: ProductPeticion; action: string; selectFIle: File }): void {
    console.log('Producto recibido:', event.product);
    console.log('Acción:', event.action);
    console.log('Archivo seleccionado:', event.selectFIle);
    this.product = event.product;
    this.selectedFile = event.selectFIle;
    this.isSubmitting = true;
    this.publishProduct(event.product);
  }
  
  publishProduct( productData: ProductPeticion ): void {
    if (!this.selectedFile) {
      alert('No hay imagen seleccionada.');
      this.isSubmitting = false;
      return;
    }

    let createdProductId: string;

    // PASO 1: Crear el producto sin la foto
    this.productService.postProduct(productData).pipe(
      switchMap((createdProduct) => {
        console.log('✅ Producto creado:', createdProduct);
        createdProductId = createdProduct.id;

        // PASO 2: Obtener URL de carga
        return this.productService.getUploadUrl(
          this.selectedFile!.name,
          this.selectedFile!.type
        );
      }),
      switchMap((uploadData) => {
        console.log('✅ URL de carga obtenida:', uploadData);

        // PASO 3: Subir la imagen a la URL pre-firmada
        return this.productService.uploadImageToUrl(
          uploadData.upload_url,
          this.selectedFile!
        ).pipe(
          switchMap(() => {
            console.log('✅ Imagen subida exitosamente');

            // PASO 4: Actualizar el producto con el object_key
            return this.productService.updateProductPhoto(
              createdProductId,
              uploadData.object_key
            );
          })
        );
      }),
      catchError((error) => {
        console.error('❌ Error en el proceso de publicación:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Error al publicar el producto. ${error.message}',
          buttonsStyling: false,
          customClass: {
            confirmButton: 'btn btn-danger'
          }
        });
        this.isSubmitting = false;
        return of(null);
      })
    ).subscribe({
      next: (result) => {
        if (result !== null) {
          console.log('✅ Producto publicado con foto exitosamente');
          alert('¡Producto publicado exitosamente con imagen!');
          this.imagePreview = null;
          this.selectedFile = null;
          this.cdr.detectChanges();
        }
        this.isSubmitting = false;
      },
      error: (err) => {
        console.error('❌ Error final:', err);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Error al publicar el producto.',
          buttonsStyling: false,
          customClass: {
            confirmButton: 'btn btn-danger'
          }
        });
        this.isSubmitting = false;
      }
    });
  }

  private loadCategories(): void { 
    this.productService.getCategories().subscribe({
      next: (categories: string[]) => {
        this.categories = categories;
        console.log('✅ Categorías cargadas:', categories);
        this.cdr.detectChanges();
      },
      error: (error) => {
        console.error('❌ Error al cargar categorías:', error);
      }
    });
  }
}