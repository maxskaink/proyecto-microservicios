import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Header } from '../../templates/header/header';
import { Product } from '../../../Models/Product';
import { ListProductTenantPreview } from '../../templates/list-product-tenant-preview/list-product-tenant-preview';
import { ProductService } from '../../../service/ProductService';
import { TenantService } from '../../../service/TenantService';
import { LoadingService } from '../../../service/loading-service';
import { IsLoading } from '../../components/is-loading/is-loading';

@Component({
  selector: 'app-home',
  imports: [CommonModule, Header, ListProductTenantPreview, IsLoading],
  templateUrl: './home.html',
  styleUrl: './home.css',
})
export class Home implements OnInit {
  categories = [
    {
      name: 'tuberculo',
      img: 'https://blog.disfrutaverdura.com/wp-content/uploads/2018/12/tuberculos.jpg',
    },
    {
      name: 'medicional',
      img: 'https://www.cocinavital.mx/wp-content/uploads/2024/01/plantas-buenas-para-la-salud.jpg',
    },
    {
      name: 'fruta',
      img: 'https://www.shaio.org/_next/image?url=https%3A%2F%2Fbackend.shaio.org%2Fsites%2Fdefault%2Ffiles%2Fblog%2Ffrutas-saludables.jpg&w=640&q=75',
    },
    {
      name: 'verdura',
      img: 'https://media.scoolinary.app/blog/images/2021/02/hortalizas-portada.jpg',
    },
    {
      name: 'hortaliza',
      img: 'https://www.naturalcastello.com/wp-content/uploads/2019/08/hortalizas.jpg',
    },
  ];

  productsByCategory: { [key: string]: any[] } = {};
  searchTerm: string = '';
  allProducts: Product[] = [];
  public products: Product[] = [];
  constructor(
    private productService: ProductService,
    private cdr: ChangeDetectorRef,
    private router: Router,
    private tenantService: TenantService,
    private loadingService: LoadingService,
  ) {}

  ngOnInit(): void {
    console.log('Llamando los productos desde  el home');
    this.loadProductsForTenant();
  }

  /**
   * Maneja el click en un producto para navegar a su vista de detalles
   */
  onProductClick(product: Product): void {
    console.log('🔥 Click en producto desde Home:', product.id, product.name);

    this.tenantService.getCurrentUserTenant().subscribe((tenant) => {
      if (!tenant) {
        console.error('No se encontró el tenant actual desde home.');
        return;
      }

      const route = ['product', tenant.tenant_id, product.id];

      this.router.navigate(route).then(
        (success) => console.log(' Navegación exitosa:', success),
        (error) => console.error(' Error en navegación:', error),
      );
    });
  }
  /**
   * Agrupa los productos por categoría
   * @param products Lista de productos a agrupar
   */
  private groupProductsByCategory(products: Product[]): void {
    this.productsByCategory = products.reduce(
      (groups, product) => {
        const category = product.category;

        if (!groups[category]) {
          groups[category] = [];
          console.log('✨ Creando nueva categoría:', category);
        }
        groups[category].push(product);

        return groups;
      },
      {} as { [category: string]: Product[] },
    );
  }
  /**
   * Carga los productos para el tenant actual
   */
  loadProductsForTenant(): void {
    this.loadingService.show('Cargando productos...');

    this.productService.getProducts(1, 100).subscribe({
      next: (products: Product[]) => {

        this.products = products;
        this.allProducts = products; // Guarda todos los productos
        this.groupProductsByCategory(products);

        this.cdr.detectChanges(); // 🔥 Fuerza actualización de la vista
        this.loadingService.hide();
      },
      error: (error) => {
        console.error('Error al cargar productos:', error);
        this.products = [];
        this.allProducts = [];
        this.loadingService.hide();
        this.cdr.detectChanges(); // Actualiza la UI incluso si hay error
      }
    });
  }


  scrollToSection(sectionId: string, event?: Event) {
    if (event) {
      event.preventDefault();
      event.stopPropagation();
    }

    const element = document.getElementById(sectionId);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  }

  onCategoryClick(categoryName: string): void {
    console.log('🔥 Click en categoría desde Home:', categoryName);
    this.router.navigate(['home', categoryName]).then(
      (success) => console.log(' Navegación exitosa a categoría:', success),
      (error) => console.error(' Error en navegación a categoría:', error),
    );
  }
}
