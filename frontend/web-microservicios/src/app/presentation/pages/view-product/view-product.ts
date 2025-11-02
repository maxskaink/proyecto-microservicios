import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Product } from '../../../Models/Product';
import { ProductService } from '../../../service /ProductService';
import { finalize } from 'rxjs';
import { Header } from '../../templates/header/header';

@Component({
  selector: 'app-view-product',
  imports: [CommonModule, Header],
  templateUrl: './view-product.html',
  styleUrl: './view-product.css',
})
export class ViewProduct implements OnInit {
  product?: Product;
  isLoading: boolean = false;
  productId: string = '';
  
  constructor(
    private productService: ProductService,
    private cdr: ChangeDetectorRef,
    private route: ActivatedRoute
  ) {}

  ngOnInit(): void {
    // Obtener el ID del producto desde los parámetros de la ruta
    this.route.params.subscribe(params => {
      this.productId = params['id'];
      if (this.productId) {
        this.loadProduct();
      }
    });
  }

  loadProduct(): void {
    if (!this.productId) {
      console.warn('No se proporcionó un productId');
      return;
    }

    this.isLoading = true;

    this.productService.getProductById(this.productId)
      .pipe(
        finalize(() => {
          this.isLoading = false;
          this.cdr.markForCheck(); 
        })
      )
      .subscribe({
        next: (product) => {
          this.product = product;
          console.log('Producto cargado:', product);
        },
        error: (err) => {
          console.error('Error al cargar el producto:', err);
        }
      });
  }
}
