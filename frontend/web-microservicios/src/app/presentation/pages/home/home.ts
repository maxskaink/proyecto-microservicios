import { CommonModule } from '@angular/common';
import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { Header } from '../../templates/header/header';
import { ProductBox } from '../../components/product-box/product-box';
import { ProductService } from '../../../service /ProductService';
import { finalize } from 'rxjs';
import { Product } from '../../../Models/Product';

@Component({
  selector: 'app-home',
  imports: [CommonModule, Header, ProductBox],
  templateUrl: './home.html',
  styleUrl: './home.css',
})
export class Home implements OnInit {
  allProducts: Product[] = [];
  isLoading: boolean = true; // Iniciar en true

  constructor(
    private productService: ProductService,
    private cdr: ChangeDetectorRef
  ) {}

  ngOnInit(): void {
    this.loadProducts();
  }

  
loadProducts(): void {
  this.isLoading = true;
  this.productService.getProducts()
    .pipe(finalize(() => {
      this.isLoading = false;
      this.cdr.markForCheck();
    }))
    .subscribe({
      next: (products) => this.allProducts = products || [],
      error: (err) => console.error(' Error:', err)
    });
}
}
