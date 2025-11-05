import { CommonModule } from '@angular/common';
import { Component, NgModule } from '@angular/core';
import { FormsModule, NgForm } from '@angular/forms';
import { Header } from '../../templates/header/header';
import { ProductPeticion } from '../../../Models/PrdocutPeticion';
import { ProductService } from '../../../service /ProductService';

@Component({
  selector: 'app-publish-product',
  imports: [CommonModule, FormsModule, Header],
  templateUrl: './publish-product.html',
  styleUrl: './publish-product.css',
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
  
  constructor(private productService: ProductService) { 
  }

  onSubmit(form: NgForm): void {
    if (this.isSubmitting) return; 
    
    // Validar que el formulario sea válido
    if (!form.valid) {
      alert('Por favor, complete todos los campos del formulario.');
      return;
    }
    
    this.isSubmitting = true;
    this.publishProduct(form);
  }
  
  publishProduct(form: NgForm): void {
    this.productService.postProduct(this.product).subscribe({
      next: (response) => {
        console.log('Producto publicado con éxito:', response);
        alert('Producto publicado exitosamente!');
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
        this.isSubmitting = false;
      },
      error: (err) => {
        console.error('Error al publicar el producto:', err);
        alert('Error al publicar el producto. Intente nuevamente.');
        this.isSubmitting = false;
      },
    });
  }
}
