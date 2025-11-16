import { Component, OnInit, OnDestroy, NgZone, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { AuthService } from '../../../service/Authser.vice';
import { combineLatest, finalize, Subscription } from 'rxjs';
import { UserData } from '../../../Models/UserData';
interface MenuOption {
  id: string;
  title: string;
  description: string;
  icon: string;
  route: string;
  category: 'shopping' | 'profile' | 'business';
}


@Component({
  selector: 'app-user',
  imports: [CommonModule],
  templateUrl: './user.html',
  styleUrl: './user.css',
})
export class User implements OnInit, OnDestroy {
  
  // Información del usuario actual
  currentUser: UserData | null = null;
  isLoggedIn: boolean = false;
  isLoading: boolean = true;
  
  // Suscripciones para limpiar al destruir el componente
  private subscriptions: Subscription = new Subscription();
  
  menuOptions: MenuOption[] = [
    // Sección de Compras
    {
      id: 'catalog',
      title: 'Catálogo',
      description: 'Explora todos los productos disponibles',
      icon: 'bi-grid-3x3-gap',
      route: '/home',
      category: 'shopping'
    },
    {
      id: 'cart',
      title: 'Carrito de Compras',
      description: 'Ver productos en tu carrito',
      icon: 'bi-cart3',
      route: '/shopping-cart',
      category: 'shopping'
    },
    {
      id: 'purchase-history',
      title: 'Historial de Compras',
      description: 'Revisa tus compras anteriores',
      icon: 'bi-clock-history',
      route: '/purchase-history',
      category: 'shopping'
    },
    // Sección de Perfil
    {
      id: 'edit-profile',
      title: 'Editar Perfil',
      description: 'Actualiza tu información personal',
      icon: 'bi-person-gear',
      route: '/edit-profile',
      category: 'profile'
    },
    // Sección de Negocio
    {
      id: 'my-products',
      title: 'Mis Productos',
      description: 'Mira tus pedidos',
      icon: 'bi-box-seam',
      route: '/list-order',
      category: 'business'
    },
    {
      id: 'publish-product',
      title: 'Publicar Producto',
      description: 'Añade un nuevo producto a la venta',
      icon: 'bi-plus-circle',
      route: '/publishProduct',
      category: 'business'
    },
    {
      id: 'orders',
      title: 'Pedidos',
      description: 'Gestiona los pedidos de tus productos',
      icon: 'bi-clipboard-check',
      route: '/orders',
      category: 'business'
    },
    {
      id: 'register',
      title: 'Registrar zona veredal',
      description: 'Registra tu zona veredal',
      icon: 'bi-clipboard-check',
      route: '/register-tenant',
      category: 'business'
    }
  ];

  constructor(private router: Router, 
    private authService: AuthService,
    private cdr: ChangeDetectorRef,
    private zone: NgZone) {}

  ngOnInit(): void {
    this.loadUserData();
  }

  ngOnDestroy(): void {
    this.subscriptions.unsubscribe();
  }

private loadUserData(): void {
  this.isLoading = true;

  // Suscripción principal a los datos del usuario
  const userDataSub = this.authService.userData.subscribe({
    next: (userData) => {
      console.log('Datos del usuario recibidos en componente:', userData);
      this.currentUser = userData;
      this.isLoading = false;
      this.cdr.detectChanges();
    },
    error: (error) => {
      console.error('Error al cargar datos del usuario:', error);
      this.isLoading = false;
      this.cdr.detectChanges();
    }
  });

  // Suscripción al estado de login
  const loginSub = this.authService.isLoggedIn$.subscribe(isLoggedIn => {
    this.isLoggedIn = isLoggedIn;
    if (!isLoggedIn) {
      this.currentUser = null;
      this.isLoading = false;
    }
    this.cdr.detectChanges();
  });

  this.subscriptions.add(userDataSub);
  this.subscriptions.add(loginSub);
}
  // Cerrar sesión
  async onLogout(): Promise<void> {
    try {
      await this.authService.logout();
      this.router.navigate(['/login']);
    } catch (error) {
      console.error('Error al cerrar sesión:', error);
    }
  }

  // Iniciar sesión (redirigir al login)
  onLogin(): void {
    this.router.navigate(['/login']);
  }

  // Obtener opciones por categoría
  getOptionsByCategory(category: string): MenuOption[] {
    return this.menuOptions.filter(option => option.category === category);
  }

  // Navegar a una opción
  navigateToOption(option: MenuOption): void {
    this.router.navigate([option.route]);
  }

  // Obtener todas las categorías únicas
  getCategories(): string[] {
    return [...new Set(this.menuOptions.map(option => option.category))];
  }

  // Obtener nombre de categoría en español
  getCategoryName(category: string): string {
    const categoryNames: {[key: string]: string} = {
      'shopping': 'Compras',
      'profile': 'Mi Perfil',
      'business': 'Mi Negocio'
    };
    return categoryNames[category] || category;
  }

  // Obtener nombre del usuario para mostrar
  getUserDisplayName(): string {
    if (!this.currentUser) return 'Usuario';
    
    // Usar directamente el campo name del backend
    return this.currentUser.name || 
           this.currentUser.email?.split('@')[0] || 
           'Usuario';
  }

  // Obtener rol del usuario con fallback
  getUserRole(): string {
    if (!this.currentUser?.rol) return 'Usuario';
    
    // Mapear roles técnicos a nombres amigables
    const roleNames: {[key: string]: string} = {
      'admin': 'Administrador',
      'vendedor': 'Vendedor',
      'comprador': 'Comprador',
      'usuario': 'Usuario'
    };
    
    return roleNames[this.currentUser.rol] || this.currentUser.rol || 'Usuario';
  }

  // Verificar si el usuario es administrador
  isUserAdmin(): boolean {
    return this.currentUser?.rol === 'admin';
  }


  getUserAddress(): string {
    // Aquí deberías usar el campo correcto para la dirección del usuario
    // Por ahora retorno un placeholder
    return 'Dirección no disponible';
  }



}
