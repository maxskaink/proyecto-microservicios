import { CommonModule, isPlatformBrowser } from '@angular/common';
import { ChangeDetectorRef, Component, ElementRef, EventEmitter, Inject, Input, OnInit, Output, PLATFORM_ID, ViewChild } from '@angular/core';
import { FormBuilder, FormsModule, NgForm, NgModel, Validators } from '@angular/forms';
import { ProductService } from '../../../service/ProductService';
import { ProductPeticion } from '../../../Models/PrdocutPeticion';

@Component({
  selector: 'app-form-product',
  imports: [CommonModule, FormsModule],
  standalone: true,
  templateUrl: './form-product.html',
  styleUrl: './form-product.css',
})
export class FormProduct implements OnInit{

  @Output() productCreated = new EventEmitter< { product: ProductPeticion, action: string, selectFIle: File }>();
  @Input() categories: string[] = [];
  @Input()  action: string = 'create';
  @ViewChild('productForm') productForm?: NgForm;

  
  @Input() product: ProductPeticion  = {
    name: '',
    category: '',
    price: 100,
    description: '',
    stock: 1,
    unit: '',
    photo_url: ''
  };
  isSubmitting: boolean = false;
  isLoading: boolean = false;

  imagePreview: string | null = null;
  isImageInvalid: boolean = false;
  imageError: string = '';
  selectedFile: File | null = null;
  @ViewChild('fileInput') fileInputRef?: ElementRef<HTMLInputElement>;
  @ViewChild('categoryBtn') categoryControl?: NgModel;
  public units: string[] = ['kg', 'atado', 'libra'];
  
  constructor(
    private serviceProduct: ProductService,
    private cdr: ChangeDetectorRef,
    @Inject(PLATFORM_ID) private platformId: Object
  ) {
      
  }
  ngOnInit() {
    if (this.action === 'edit' ) {
      this.imagePreview = this.product.photo_url; // Muestra la imagen existente
    }
  }

  /**
   * funcion de apoyo para seleccionar la categoria
   * @param cat categoría seleccionada
   */
  selectCategory(cat: string) {
    this.product.category = cat;
    if (this.categoryControl) {
      this.categoryControl.control?.markAsTouched();
    }
  }

  /**
   * Emite los datos del producto creado al padre
   * @param productData datos del producto
   */
  onClick(productData: ProductPeticion): void {
    this.productCreated.emit({product: productData, action: this.action, selectFIle: this.selectedFile!});
  }

  /**
   * 
   * @param form formulario de producto
   * @returns si el formulario es valido, crea el producto
   */
  onSubmit(form: NgForm): void {

    if (this.isSubmitting) return; 
    
    if (!form.valid) {
      alert('Por favor, complete todos los campos del formulario.');
      return;
    }
    const hasExistingImage = !!this.product.photo_url; 
    if (!this.selectedFile && !hasExistingImage) {
      alert('Por favor, seleccione una imagen para el producto.');
      return;
    }
    
    this.isSubmitting = true;
    const productData: ProductPeticion = {
      ...this.product,
      photo_url: ''
    };
    this.onClick(productData);
    setTimeout(() => {
      this.cleanForm();
      this.isSubmitting = false;
    }, 500);
  }

  /** Limpia el formulario
   */
cleanForm(): void {
  // Resetea el modelo
  this.product = {
    name: '',
    category: '',
    price: 100,
    description: '',
    stock: 1,
    unit: '',
    photo_url: ''
  };

  // Resetea vista previa e imagen
  this.imagePreview = null;
  this.isImageInvalid = false;
  this.imageError = '';
  this.selectedFile = null;

  if (this.fileInputRef) {
    this.fileInputRef.nativeElement.value = '';
  }

  // Resetea el estado del formulario y errores
  if (this.productForm) {
    this.productForm.resetForm();
  }

  this.cdr.detectChanges();
}

  /**
   * 
   * @param event  
   * @returns sube la imagen y la valida
   */
  onImageChange(event: Event): void {
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
  
}
