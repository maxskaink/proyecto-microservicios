import { CommonModule, isPlatformBrowser } from '@angular/common';
import { ChangeDetectorRef, Component, ElementRef, Inject, PLATFORM_ID, ViewChild } from '@angular/core';
import { FormsModule, NgForm } from '@angular/forms';
import { Header } from '../../templates/header/header';
import { ProductPeticion } from '../../../Models/PrdocutPeticion';
import { ProductService } from '../../../service /ProductService';
import { switchMap, catchError } from 'rxjs/operators';
import { of } from 'rxjs';

@Component({
  selector: 'app-publish-product',
  imports: [CommonModule, FormsModule, Header],
  templateUrl: './publish-product.html',
  styleUrl: './publish-product.css',
  standalone: true,
})
export class PublishProduct {
  product: ProductPeticion = {
    name: '',
    category: '',
    price: 0,
    description: '',
    stock: 0,
    unit: '',
    photo_url: ''
  };
  isSubmitting: boolean = false;

  imagePreview: string | null = null;
  isImageInvalid: boolean = false;
  imageError: string = '';
  selectedFile: File | null = null;
  @ViewChild('fileInput') fileInputRef?: ElementRef<HTMLInputElement>;

  constructor(
    private productService: ProductService,
    private cdr: ChangeDetectorRef,
    @Inject(PLATFORM_ID) private platformId: Object
  ) { }

  onSubmit(form: NgForm): void {
    if (this.isSubmitting) return; 
    
    if (!form.valid) {
      alert('Por favor, complete todos los campos del formulario.');
      return;
    }

    if (!this.selectedFile) {
      alert('Por favor, seleccione una imagen para el producto.');
      return;
    }
    
    this.isSubmitting = true;
    this.publishProduct(form);
  }

  onImageChange(event: Event): void {
    // ✅ Verificar que estamos en el navegador
    if (!isPlatformBrowser(this.platformId)) {
      console.warn('⚠️ FileReader no disponible en SSR');
      return;
    }

    this.isImageInvalid = false;
    this.imageError = '';
    
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) {
      this.imagePreview = null;
      this.selectedFile = null;
      this.product.photo_url = '';
      this.cdr.detectChanges();
      return;
    }

    const file = input.files[0];
    console.log('📁 Archivo seleccionado:', file.name, file.type, file.size);

    // Validar tipo
    if (!file.type.startsWith('image/')) {
      this.isImageInvalid = true;
      this.imageError = 'El archivo no es una imagen válida.';
      this.imagePreview = null;
      this.selectedFile = null;
      this.cdr.detectChanges();
      return;
    }

    // Validar tamaño (5MB)
    const maxSizeMB = 5;
    if (file.size > maxSizeMB * 1024 * 1024) {
      this.isImageInvalid = true;
      this.imageError = `La imagen supera el tamaño máximo de ${maxSizeMB} MB.`;
      this.imagePreview = null;
      this.selectedFile = null;
      this.cdr.detectChanges();
      return;
    }

    this.selectedFile = file;

    // Convertir a base64 para vista previa
    const reader = new FileReader();
    
    reader.onload = () => {
      this.imagePreview = reader.result as string;
      console.log('✅ Vista previa cargada, longitud:', this.imagePreview?.length);
      this.cdr.detectChanges(); // ✅ Forzar detección de cambios
    };
    
    reader.onerror = (error) => {
      console.error('❌ Error al leer la imagen:', error);
      this.isImageInvalid = true;
      this.imageError = 'Error al leer la imagen.';
      this.imagePreview = null;
      this.selectedFile = null;
      this.cdr.detectChanges();
    };
    
    reader.readAsDataURL(file);
  }

  publishProduct(form: NgForm): void {
    if (!this.selectedFile) {
      alert('No hay imagen seleccionada.');
      this.isSubmitting = false;
      return;
    }

    // Preparar el producto sin la foto
    const productData: ProductPeticion = {
      ...this.product,
      photo_url: ''
    };

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
        alert('Error al publicar el producto. Intente nuevamente.');
        this.isSubmitting = false;
        return of(null);
      })
    ).subscribe({
      next: (result) => {
        if (result !== null) {
          console.log('✅ Producto publicado con foto exitosamente');
          alert('¡Producto publicado exitosamente con imagen!');
          
          // Resetear el formulario
          form.resetForm();
          this.product = {
            name: '',
            category: '',
            price: 0,
            description: '',
            stock: 0,
            unit: '',
            photo_url: ''
          };
          this.imagePreview = null;
          this.selectedFile = null;
          
          // Limpiar input file
          if (this.fileInputRef?.nativeElement) {
            this.fileInputRef.nativeElement.value = '';
          }
          
          this.cdr.detectChanges();
        }
        this.isSubmitting = false;
      },
      error: (err) => {
        console.error('❌ Error final:', err);
        alert('Error al publicar el producto.');
        this.isSubmitting = false;
      }
    });
  }
}