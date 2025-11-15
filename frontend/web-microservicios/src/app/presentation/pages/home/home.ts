import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit} from '@angular/core';
import { Header } from '../../templates/header/header';
import { Product } from '../../../Models/Product';
import { ListProductTenantPreview } from '../../templates/list-product-tenant-preview/list-product-tenant-preview';
import { ProductService } from '../../../service/ProductService';


@Component({
  selector: 'app-home',
  imports: [CommonModule, Header, ListProductTenantPreview],
  templateUrl: './home.html',
  styleUrl: './home.css',
})
export class Home  implements OnInit {
  /**
   * Productos organizados por categoría
   */
  public productsByCategory: { [category: string]: Product[] } = {};
  /**
   * Lista de categorías disponibles
   */
  public categories: string[] = [];


  constructor(private productService: ProductService,   private cdr: ChangeDetectorRef) {}

  ngOnInit(): void {
    console.log('Llamando los productos desde  el home');
    this.loadProductsForTenant();
  }

  public products: Product[] = [];
  onProductClick(product: any) {
    console.log('Producto clickeado:', product);
    // Aquí puedes navegar, abrir un modal, etc.
  }
  /**
   * Agrupa los productos por categoría
   * @param products Lista de productos a agrupar
   */
private groupProductsByCategory(products: Product[]): void {

  this.productsByCategory = products.reduce((groups, product) => {
    const category = product.category || 'Sin categoría';
    console.log('📂 Categoría del producto:', category);
    console.log('📦 Producto completo:', product);

    if (!groups[category]) {
      groups[category] = [];
      console.log('✨ Creando nueva categoría:', category);
    }
    groups[category].push(product);

    return groups;
  }, {} as { [category: string]: Product[] });
  
  this.categories = Object.keys(this.productsByCategory);
  
  console.log('✅ productsByCategory FINAL:', this.productsByCategory);
  console.log('✅ categories FINAL:', this.categories);
  console.log('🔢 Cantidad de categorías:', this.categories.length);
}
  /**
   * Carga los productos para el tenant actual
   */
loadProductsForTenant(): void {
  this.productService.getProducts(1, 10).subscribe((products: Product[]) => {
    console.log('📦 Productos recibidos:', products);

    this.products = products;
    this.groupProductsByCategory(products);

    this.cdr.detectChanges();   // 🔥 Fuerza actualización de la vista
  });
}


}